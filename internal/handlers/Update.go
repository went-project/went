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

var (
	latestReleaseURL = "https://api.github.com/repos/went-project/went/releases/latest"
	releasesURL      = "https://api.github.com/repos/went-project/went/releases?per_page=100"
	httpClient       = &http.Client{Timeout: 10 * time.Second}
)

const installScriptBase = "https://raw.githubusercontent.com/went-project/went/main"

type UpdateOptions struct {
	Channel        utils.Channel
	Force          bool
	CurrentVersion string
}

type VersionCheckResult struct {
	Current         string
	Latest          string
	UpdateAvailable bool
}

func getLatestStableReleaseTag(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, latestReleaseURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "went-update-checker")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := httpClient.Do(req)
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

func getLatestBetaReleaseTag(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, releasesURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "went-update-checker")
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch beta releases: %s", resp.Status)
	}

	var releases []struct {
		TagName    string `json:"tag_name"`
		Prerelease bool   `json:"prerelease"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", err
	}

	var bestTag string
	for _, release := range releases {
		tag := strings.TrimSpace(release.TagName)
		if tag == "" || !release.Prerelease || !strings.Contains(tag, "-beta") {
			continue
		}

		if bestTag == "" || utils.CompareVersions(bestTag, tag) < 0 {
			bestTag = tag
		}
	}

	if bestTag == "" {
		return "", errors.New("no beta release found")
	}

	return bestTag, nil
}

func getLatestReleaseTag(ctx context.Context, opts UpdateOptions) (string, error) {
	if opts.Channel == utils.ChannelBeta {
		return getLatestBetaReleaseTag(ctx)
	}
	return getLatestStableReleaseTag(ctx)
}

func CheckLatestVersion(opts UpdateOptions) (*VersionCheckResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	latest, err := getLatestReleaseTag(ctx, opts)
	if err != nil {
		return nil, err
	}

	current := opts.CurrentVersion
	if current == "" {
		current = utils.GetCurrentVersion()
	}

	compare := utils.CompareVersions(current, latest)

	return &VersionCheckResult{
		Current:         current,
		Latest:          latest,
		UpdateAvailable: compare < 0,
	}, nil
}

func PrintCreateVersionWarning() {
	result, err := CheckLatestVersion(UpdateOptions{Channel: utils.ChannelStable})
	if err != nil {
		output.PrintWarning("Unable to check for the latest version: %v", err)
		return
	}

	if result.UpdateAvailable {
		output.PrintWarning("A newer version is available: %s (currently %s). Run 'went update' to install it.", result.Latest, result.Current)
	}
}

func UpdateWent(opts UpdateOptions) error {
	if opts.Channel == "" {
		opts.Channel = utils.ChannelStable
	}

	result, err := CheckLatestVersion(opts)
	if err != nil {
		if isNetworkError(err) {
			return fmt.Errorf("Update failed: network is unavailable or GitHub cannot be reached: %w", err)
		}
		return fmt.Errorf("Update check failed: %w", err)
	}

	if !result.UpdateAvailable && !opts.Force {
		output.PrintSuccess("Went is already up to date. Current version: %s", result.Current)
		return nil
	}

	channelLabel := string(opts.Channel)
	if channelLabel == "" {
		channelLabel = string(utils.ChannelStable)
	}

	if !result.UpdateAvailable && opts.Force {
		output.PrintInfo("No newer version found, but force install requested for %s channel: %s.", channelLabel, result.Latest)
	} else {
		output.PrintInfo("A new version is available on %s channel: %s (currently %s). Starting update...", channelLabel, result.Latest, result.Current)
	}

	cmdName, args, err := installCommand()
	if err != nil {
		return err
	}

	cmd := exec.Command(cmdName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = mergeEnv(os.Environ(), "WENT_VERSION", result.Latest)
	cmd.Env = mergeEnv(cmd.Env, "WENT_CHANNEL", channelLabel)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("install script failed: %w", err)
	}

	output.PrintSuccess("Update completed successfully.")
	return nil
}

func mergeEnv(env []string, key, value string) []string {
	prefix := key + "="
	filtered := make([]string, 0, len(env)+1)
	for _, entry := range env {
		if !strings.HasPrefix(entry, prefix) {
			filtered = append(filtered, entry)
		}
	}
	filtered = append(filtered, fmt.Sprintf("%s=%s", key, value))
	return filtered
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
