package handlers

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"went/internal/output"
	"went/internal/utils"

	"github.com/fsnotify/fsnotify"
)

var (
	activeCmd     *exec.Cmd
	activeCmdDone chan struct{}
	activeCmdMu   sync.Mutex
)

func RunWithWatcher(root string) error {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("failed to resolve project root: %w", err)
	}

	ignorePatterns, ignoreSource, err := utils.LoadIgnorePatterns(".wentignore", ".devwatchignore")
	if err != nil {
		return fmt.Errorf("failed to load ignore patterns: %w", err)
	}

	port, portSource, err := utils.LoadPortFromDotEnv(".env.local", ".env")
	if err != nil {
		return fmt.Errorf("failed to read env file: %w", err)
	}

	printRunHeader(absRoot, port, portSource, ignoreSource)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to initialize file watcher: %w", err)
	}
	defer watcher.Close()

	if err := watchRecursive(absRoot, watcher, ignorePatterns); err != nil {
		return fmt.Errorf("failed to watch project directories: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), terminationSignals()...)
	defer cancel()

	if err := startApp(absRoot, port); err != nil {
		return fmt.Errorf("failed to start application: %w", err)
	}

	events := make(chan struct{}, 1)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if utils.ShouldIgnore(event.Name, ignorePatterns, absRoot) {
					continue
				}

				if event.Op&fsnotify.Create != 0 {
					info, err := os.Stat(event.Name)
					if err == nil && info.IsDir() && !utils.ShouldIgnore(event.Name, ignorePatterns, absRoot) {
						_ = watcher.Add(event.Name)
					}
				}

				if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
					output.PrintWarning("Change detected: %s", event.Name)
					select {
					case events <- struct{}{}:
					default:
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				output.PrintWarning("Watcher error: %v", err)
			}
		}
	}()

	var restartTimer *time.Timer
	var restartC <-chan time.Time
	const debounceDelay = 250 * time.Millisecond

	for {
		select {
		case <-ctx.Done():
			output.PrintInfo("Shutting down watcher...")
			topAndWait()
			return nil
		case <-events:
			if restartTimer == nil {
				restartTimer = time.NewTimer(debounceDelay)
				restartC = restartTimer.C
			} else {
				if !restartTimer.Stop() {
					<-restartTimer.C
				}
				restartTimer.Reset(debounceDelay)
			}
		case <-restartC:
			restartTimer = nil
			restartC = nil
			ignorePatterns, ignoreSource, err = utils.LoadIgnorePatterns(".wentignore", ".devwatchignore")
			if err != nil {
				output.PrintWarning("Failed to reload ignore rules: %v", err)
			}
			output.PrintInfo("Restarting application...")
			if err := restartApp(absRoot, port); err != nil {
				output.PrintError("Failed to restart application: %v", err)
			}
		}
	}
}

func printRunHeader(root, port, portSource, ignoreSource string) {
	output.PrintHeader("WENT hot reload started")
	output.PrintLine("Watching", root)
	defaultPort := "8080"
	if port == "" {
		port = defaultPort
		portSource = "default"
	}
	output.PrintLine("Port", fmt.Sprintf("%s (%s)", port, portSource))
	output.PrintLine("Local", fmt.Sprintf("http://127.0.0.1:%s", port))
	output.PrintLine("Swagger", fmt.Sprintf("http://127.0.0.1:%s/swagger/index.html", port))
	networkIP := getNetworkIP()
	if networkIP != "" {
		output.PrintLine("Network", fmt.Sprintf("http://%s:%s", networkIP, port))
		output.PrintLine("Swagger", fmt.Sprintf("http://%s:%s/swagger/index.html", networkIP, port))
	} else {
		output.PrintWarning("No network address found")
	}
	if ignoreSource != "" {
		output.PrintLine("Ignore file", ignoreSource)
	} else {
		output.PrintWarning("No ignore file found")
	}
	output.PrintSuccess("Hot reload is active.")
	output.PrintSpacing()
}

func getNetworkIP() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip := v.IP
				if ip.IsLoopback() {
					continue
				}
				if ipv4 := ip.To4(); ipv4 != nil {
					return ipv4.String()
				}
			case *net.IPAddr:
				ip := v.IP
				if ip.IsLoopback() {
					continue
				}
				if ipv4 := ip.To4(); ipv4 != nil {
					return ipv4.String()
				}
			}
		}
	}

	return ""
}

func watchRecursive(root string, watcher *fsnotify.Watcher, ignorePatterns []string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if utils.ShouldIgnore(path, ignorePatterns, root) {
				return filepath.SkipDir
			}
			return watcher.Add(path)
		}

		return nil
	})
}

func startApp(root, port string) error {
	activeCmdMu.Lock()
	defer activeCmdMu.Unlock()

	cmd := exec.Command("go", "run", ".")
	cmd.Dir = root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = buildAppEnv(port)

	setupProcessGroup(cmd)

	if err := cmd.Start(); err != nil {
		return err
	}

	activeCmd = cmd
	activeCmdDone = make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(activeCmdDone)
	}()

	output.PrintSuccess("Application started (PID: %d)", cmd.Process.Pid)
	return nil
}

func restartApp(root, port string) error {
	if err := terminateActiveProcess(5 * time.Second); err != nil {
		output.PrintWarning("Failed to stop existing process: %v", err)
	}
	return startApp(root, port)
}

func topAndWait() {
	if err := terminateActiveProcess(5 * time.Second); err != nil {
		output.PrintWarning("Failed to stop running process: %v", err)
	}
}

func terminateActiveProcess(timeout time.Duration) error {
	activeCmdMu.Lock()
	cmd := activeCmd
	done := activeCmdDone
	activeCmdMu.Unlock()

	if cmd == nil || cmd.Process == nil {
		return nil
	}

	_ = signalProcessGroup(cmd)

	finished := make(chan struct{})
	go func() {
		if done != nil {
			<-done
		}
		close(finished)
	}()

	select {
	case <-finished:
		return nil
	case <-time.After(timeout):
		_ = killProcessGroup(cmd)
		<-finished
	}

	activeCmdMu.Lock()
	if activeCmd == cmd {
		activeCmd = nil
		activeCmdDone = nil
	}
	activeCmdMu.Unlock()

	return nil
}

func buildAppEnv(port string) []string {
	env := os.Environ()
	if port == "" {
		return env
	}

	filtered := make([]string, 0, len(env)+1)
	for _, variable := range env {
		if strings.HasPrefix(variable, "PORT=") {
			continue
		}
		filtered = append(filtered, variable)
	}
	filtered = append(filtered, fmt.Sprintf("PORT=%s", port))
	return filtered
}
