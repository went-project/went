package handlers

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	"went/internal/utils"

	"github.com/fsnotify/fsnotify"
)

var (
	activeCmd     *exec.Cmd
	activeCmdDone chan struct{}
	activeCmdMu   sync.Mutex
)

const (
	runBannerColor  = "\033[36m"
	runSuccessColor = "\033[32m"
	runWarningColor = "\033[33m"
	resetColor      = "\033[0m"
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

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
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
					fmt.Printf("%s⟳ Change detected:%s %s\n", runWarningColor, resetColor, event.Name)
					select {
					case events <- struct{}{}:
					default:
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("%swatcher error:%s %v\n", runWarningColor, resetColor, err)
			}
		}
	}()

	var restartTimer *time.Timer
	var restartC <-chan time.Time
	const debounceDelay = 250 * time.Millisecond

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("\n%s⏹ Shutting down watcher...%s\n", runBannerColor, resetColor)
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
				fmt.Printf("%sWarning: failed to reload ignore rules:%s %v\n", runWarningColor, resetColor, err)
			}
			fmt.Printf("%s↺ Restarting application...%s\n", runBannerColor, resetColor)
			if err := restartApp(absRoot, port); err != nil {
				fmt.Printf("%sError: failed to restart application:%s %v\n", runWarningColor, resetColor, err)
			}
		}
	}
}

func printRunHeader(root, port, portSource, ignoreSource string) {
	fmt.Println(runBannerColor + "╭────────────────────────────────────────────╮" + resetColor)
	fmt.Println(runBannerColor + "│" + resetColor + "  WENT hot reload started                   " + runBannerColor + "│" + resetColor)
	fmt.Println(runBannerColor + "╰────────────────────────────────────────────╯" + resetColor)
	fmt.Printf("%sWatching:%s %s\n", runSuccessColor, resetColor, root)
	if port != "" {
		fmt.Printf("%sPORT:%s %s (%s)\n", runSuccessColor, resetColor, port, portSource)
		fmt.Printf("%sLocal:%s   http://127.0.0.1:%s\n", runSuccessColor, resetColor, port)
		fmt.Printf("%sSwagger:%s http://127.0.0.1:%s/swagger/index.html\n", runSuccessColor, resetColor, port)
		networkIP := getNetworkIP()
		if networkIP != "" {
			fmt.Printf("%sNetwork:%s http://%s:%s\n", runSuccessColor, resetColor, networkIP, port)
			fmt.Printf("%sSwagger:%s http://%s:%s/swagger/index.html\n", runSuccessColor, resetColor, networkIP, port)
		} else {
			fmt.Printf("%sNetwork:%s no network address found\n", runWarningColor, resetColor)
		}
	} else {
		fmt.Printf("%sPORT:%s not set\n", runWarningColor, resetColor)
	}
	if ignoreSource != "" {
		fmt.Printf("%sIgnore file:%s %s\n", runSuccessColor, resetColor, ignoreSource)
	} else {
		fmt.Printf("%sIgnore file:%s not found\n", runWarningColor, resetColor)
	}
	fmt.Printf("%sHot reload is active.%s\n\n", runSuccessColor, resetColor)
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

	if runtime.GOOS != "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	activeCmd = cmd
	activeCmdDone = make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(activeCmdDone)
	}()

	fmt.Printf("%s▶ Application started%s (PID: %d)\n\n", runSuccessColor, resetColor, cmd.Process.Pid)
	return nil
}

func restartApp(root, port string) error {
	if err := terminateActiveProcess(5 * time.Second); err != nil {
		fmt.Printf("%sWarning: failed to stop existing process:%s %v\n", runWarningColor, resetColor, err)
	}
	return startApp(root, port)
}

func topAndWait() {
	if err := terminateActiveProcess(5 * time.Second); err != nil {
		fmt.Printf("%sWarning: failed to stop running process:%s %v\n", runWarningColor, resetColor, err)
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

	if runtime.GOOS != "windows" {
		if pgid, err := syscall.Getpgid(cmd.Process.Pid); err == nil {
			_ = syscall.Kill(-pgid, syscall.SIGTERM)
		} else {
			_ = cmd.Process.Signal(syscall.SIGTERM)
		}
	} else {
		_ = cmd.Process.Kill()
	}

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
		if runtime.GOOS != "windows" {
			if pgid, err := syscall.Getpgid(cmd.Process.Pid); err == nil {
				_ = syscall.Kill(-pgid, syscall.SIGKILL)
			}
		} else {
			_ = cmd.Process.Kill()
		}
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
