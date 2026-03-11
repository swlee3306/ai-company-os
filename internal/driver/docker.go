package driver

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func DockerDaemonOK() (bool, string) {
	if _, err := exec.LookPath("docker"); err != nil {
		return false, "docker not found"
	}
	cmd := exec.Command("docker", "info")
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return false, msg
	}
	return true, ""
}

func DockerPreflightHint(goos string) string {
	if goos == "darwin" {
		return "Start Docker Desktop and wait until it is running, then retry."
	}
	return "Start the docker daemon/service (e.g., `sudo systemctl start docker`) then retry."
}

func FormatDockerPreflightError(goos string, detail string) error {
	return fmt.Errorf("docker daemon not reachable: %s\nHint: %s", detail, DockerPreflightHint(goos))
}
