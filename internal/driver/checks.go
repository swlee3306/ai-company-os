package driver

import (
	"bytes"
	"os/exec"
	"strings"
	"time"
)

type Check struct {
	Name   string `json:"name"`
	Status string `json:"status"` // ok|warn|error|unknown
	Detail string `json:"detail,omitempty"`
}

type DriverStatus struct {
	Selected        string  `json:"selected"`
	Docker          bool    `json:"docker"`
	K3D             bool    `json:"k3d"`
	K3DClusterCount int     `json:"k3d_cluster_count"`
	Checks          []Check `json:"checks"`
}

func CheckAll() DriverStatus {
	st := DriverStatus{Selected: "k3d"}

	dockerOK, dockerDetail := cmdOK("docker", "info")
	st.Docker = dockerOK
	checks := []Check{mk("docker.info", dockerOK, dockerDetail)}

	k3dOK, k3dDetail := cmdOK("k3d", "version")
	st.K3D = k3dOK
	checks = append(checks, mk("k3d.version", k3dOK, k3dDetail))

	if k3dOK {
		out, ok, detail := cmdOut("k3d", "cluster", "list")
		checks = append(checks, mk("k3d.cluster_list", ok, detail))
		if ok {
			st.K3DClusterCount = countClusters(out)
		}
	}

	st.Checks = checks
	return st
}

func mk(name string, ok bool, detail string) Check {
	c := Check{Name: name}
	if ok {
		c.Status = "ok"
	} else {
		c.Status = "warn"
		c.Detail = detail
	}
	return c
}

func cmdOK(bin string, args ...string) (bool, string) {
	_, ok, detail := cmdOut(bin, args...)
	return ok, detail
}

func cmdOut(bin string, args ...string) (string, bool, string) {
	start := time.Now()
	cmd := exec.Command(bin, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	_ = time.Since(start)
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" {
			msg = err.Error()
		}
		return stdout.String(), false, msg
	}
	return stdout.String(), true, ""
}

func countClusters(out string) int {
	lines := strings.Split(out, "\n")
	cnt := 0
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "NAME") {
			continue
		}
		cnt++
	}
	return cnt
}
