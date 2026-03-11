package installer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func K3DPlan(goos string) string {
	switch goos {
	case "darwin":
		return "macOS: install via Homebrew: `brew install k3d` (requires brew)."
	case "linux":
		return "Linux: download the k3d release binary and install to ~/.local/bin (or /usr/local/bin if writable)."
	default:
		return fmt.Sprintf("Unsupported OS for auto-install: %s", goos)
	}
}

func InstallK3D(goos string) error {
	switch goos {
	case "darwin":
		return installK3DViaBrew()
	case "linux":
		return installK3DViaReleaseBinary()
	default:
		return fmt.Errorf("unsupported OS: %s", goos)
	}
}

func installK3DViaBrew() error {
	if _, err := exec.LookPath("brew"); err != nil {
		return errors.New("brew not found; install Homebrew or install k3d manually")
	}
	cmd := exec.Command("brew", "install", "k3d")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("brew install k3d failed: %s", string(out))
	}
	return nil
}

func installK3DViaReleaseBinary() error {
	assetName := assetFor(runtime.GOOS, runtime.GOARCH)
	if assetName == "" {
		return fmt.Errorf("unsupported platform for k3d binary: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	rel, err := fetchLatestRelease()
	if err != nil {
		return err
	}

	binURL := ""
	checksumsURL := ""
	for _, a := range rel.Assets {
		if a.Name == assetName {
			binURL = a.BrowserDownloadURL
		}
		if strings.Contains(a.Name, "checks") && strings.HasSuffix(a.Name, ".txt") {
			checksumsURL = a.BrowserDownloadURL
		}
	}
	if binURL == "" {
		return fmt.Errorf("k3d asset not found in release %s: %s", rel.TagName, assetName)
	}

	tmpDir, err := os.MkdirTemp("", "k3d-install-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	binPath := filepath.Join(tmpDir, "k3d")
	if err := downloadTo(binURL, binPath); err != nil {
		return err
	}
	if err := os.Chmod(binPath, 0o755); err != nil {
		return err
	}

	// optional checksum verification when available
	if checksumsURL != "" {
		chkPath := filepath.Join(tmpDir, "checksums.txt")
		if err := downloadTo(checksumsURL, chkPath); err == nil {
			_ = verifySHA256(chkPath, binPath, assetName)
		}
	}

	installDir := defaultInstallDir()
	if err := os.MkdirAll(installDir, 0o755); err != nil {
		return err
	}
	finalPath := filepath.Join(installDir, "k3d")

	// try direct install
	if err := copyFile(binPath, finalPath, 0o755); err == nil {
		return nil
	}

	// fallback: try sudo install to /usr/local/bin
	if _, err := exec.LookPath("sudo"); err == nil {
		cmd := exec.Command("sudo", "install", "-m", "0755", binPath, "/usr/local/bin/k3d")
		out, err := cmd.CombinedOutput()
		if err == nil {
			return nil
		}
		return fmt.Errorf("failed to install k3d (sudo install): %s", string(out))
	}

	return fmt.Errorf("failed to install k3d to %s (permission). Try: sudo install -m 0755 %s /usr/local/bin/k3d", finalPath, binPath)
}

func assetFor(goos, goarch string) string {
	if goos == "linux" && goarch == "amd64" {
		return "k3d-linux-amd64"
	}
	if goos == "linux" && goarch == "arm64" {
		return "k3d-linux-arm64"
	}
	return ""
}

func fetchLatestRelease() (*githubRelease, error) {
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/k3d-io/k3d/releases/latest", nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github api error: %s: %s", resp.Status, string(b))
	}
	var rel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}
	return &rel, nil
}

func downloadTo(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download failed: %s: %s", resp.Status, string(b))
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return os.Chmod(dst, mode)
}

func defaultInstallDir() string {
	h, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(h, ".local", "bin")
}
