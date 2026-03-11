package driver

import (
	"fmt"
	"os/exec"
	"strings"
)

func K3DUp(cluster string) error {
	if _, err := exec.LookPath("k3d"); err != nil {
		return fmt.Errorf("k3d not found; run `company install k3d`")
	}
	if clusterExists(cluster) {
		return nil
	}
	cmd := exec.Command("k3d", "cluster", "create", cluster)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("k3d cluster create failed: %s", string(out))
	}
	return nil
}

func K3DDown(cluster string) error {
	if _, err := exec.LookPath("k3d"); err != nil {
		return nil
	}
	if !clusterExists(cluster) {
		return nil
	}
	cmd := exec.Command("k3d", "cluster", "delete", cluster)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("k3d cluster delete failed: %s", string(out))
	}
	return nil
}

func clusterExists(cluster string) bool {
	cmd := exec.Command("k3d", "cluster", "list")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(out), cluster)
}
