package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"went/internal/output"
	"went/internal/utils"
)

const (
	latestReleaseURL  = "https://api.github.com/repos/went-project/went/releases/latest"
	installScriptBase = "https://raw.githubusercontent.com/went-project/went/main"
)

type VersionCheckResult struct {
	Current         string
	Latest          string
	UpdateAvailable bool
}

func getLatestReleaseTag(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, latestReleaseURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "went-update-checker")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch latest release: %s", resp.Status)
	}

	var payload struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", err
	}

	if strings.TrimSpace(payload.TagName) == "" {
		return "", errors.New("GitHub release tag_name is empty")
	}

	return strings.TrimSpace(payload.TagName), nil
}

func CheckLatestVersion() (*VersionCheckResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	latest, err := getLatestReleaseTag(ctx)
	if err != nil {
		return nil, err
	}

	current := utils.GetCurrentVersion()
	compare := utils.CompareVersions(current, latest)

	return &VersionCheckResult{
		Current:         current,
		Latest:          latest,
		UpdateAvailable: compare < 0,
	}, nil
}

func PrintCreateVersionWarning() {
	result, err := CheckLatestVersion()
	if err != nil {
		output.PrintWarning("Unable to check for the latest version: %v", err)
		return
	}

	if result.UpdateAvailable {
		output.PrintWarning("A newer version is available: %s (currently %s). Run 'went update' to install it.", result.Latest, result.Current)
	}
}

func UpdateWent() error {
	result, err := CheckLatestVersion()
	if err != nil {
		if isNetworkError(err) {
			return fmt.Errorf("Update failed: network is unavailable or GitHub cannot be reached: %w", err)
		}
		return fmt.Errorf("Update check failed: %w", err)
	}

	if !result.UpdateAvailable {
		output.PrintSuccess("Went is already up to date. Current version: %s", result.Current)
		return nil
	}

	output.PrintInfo("A new version is available: %s (currently %s). Starting update...", result.Latest, result.Current)
	cmdName, args, err := installCommand()
	if err != nil {
		return err
	}

	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("install script failed: %w", err)
	}

	output.PrintSuccess("Update completed successfully.")
	return nil
}

func installCommand() (string, []string, error) {
	switch runtime.GOOS {
	case "darwin", "linux":
		command := "sh"
		script := fmt.Sprintf("curl -fsSL %s/install.sh | sh", installScriptBase)
		return command, []string{"-c", script}, nil
	case "windows":
		command := "powershell.exe"
		script := fmt.Sprintf("Invoke-WebRequest -Uri '%s/install.ps1' -UseBasicParsing -OutFile $env:TEMP\\went-update.ps1; & $env:TEMP\\went-update.ps1; Remove-Item -Force $env:TEMP\\went-update.ps1 -ErrorAction SilentlyContinue", installScriptBase)
		return command, []string{"-NoProfile", "-Command", script}, nil
	default:
		return "", nil, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

func isNetworkError(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return true
	}

	if utils.IsNetworkError(err) {
		return true
	}

	return false
}
