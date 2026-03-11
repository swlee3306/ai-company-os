package driver

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func K3SChecks() (bool, string) {
	if runtime.GOOS != "linux" {
		return false, "k3s is linux-only"
	}
	if _, err := exec.LookPath("k3s"); err != nil {
		return false, "k3s not found"
	}
	cmd := exec.Command("k3s", "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(out.String())
		if msg == "" {
			msg = err.Error()
		}
		return false, msg
	}
	return true, strings.TrimSpace(out.String())
}

func K3SServiceState() (string, string) {
	if runtime.GOOS != "linux" {
		return "unsupported", "k3s is linux-only"
	}
	if _, err := exec.LookPath("systemctl"); err != nil {
		return "unknown", "systemctl not available"
	}
	cmd := exec.Command("systemctl", "is-active", "k3s")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	s := strings.TrimSpace(out.String())
	if err != nil {
		if s == "" {
			s = "inactive"
		}
		return s, strings.TrimSpace(out.String())
	}
	return s, ""
}

func K3SUp() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("k3s is linux-only")
	}
	if _, err := exec.LookPath("systemctl"); err != nil {
		return fmt.Errorf("systemctl not available")
	}
	cmd := exec.Command("systemctl", "start", "k3s")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to start k3s (try with sudo or as root): %s", strings.TrimSpace(string(out)))
	}
	return nil
}

func K3SDown() error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("k3s is linux-only")
	}
	if _, err := exec.LookPath("systemctl"); err != nil {
		return fmt.Errorf("systemctl not available")
	}
	cmd := exec.Command("systemctl", "stop", "k3s")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to stop k3s (try with sudo or as root): %s", strings.TrimSpace(string(out)))
	}
	return nil
}
