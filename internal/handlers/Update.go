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
		return "", errors.New("GitHub release tag_name boş")
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
		fmt.Printf("Warning: Güncel sürüm kontrolü yapılamadı: %v\n", err)
		return
	}

	if result.UpdateAvailable {
		fmt.Printf("Warning: Yeni sürüm bulundu: %s (şu an: %s). Güncellemek için 'went update' kullanabilirsiniz.\n", result.Latest, result.Current)
	}
}

func UpdateWent() error {
	result, err := CheckLatestVersion()
	if err != nil {
		if isNetworkError(err) {
			return fmt.Errorf("Güncelleme başarısız: internet bağlantısı yok veya GitHub erişilemiyor: %w", err)
		}
		return fmt.Errorf("Güncelleme kontrolü başarısız: %w", err)
	}

	if !result.UpdateAvailable {
		fmt.Printf("Went zaten güncel. Şu anki sürüm: %s\n", result.Current)
		return nil
	}

	fmt.Printf("Yeni sürüm bulundu: %s (şu an %s). Güncelleme başlatılıyor...\n", result.Latest, result.Current)
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

	fmt.Println("Güncelleme tamamlandı.")
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
		return "", nil, fmt.Errorf("desteklenmeyen işletim sistemi: %s", runtime.GOOS)
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
