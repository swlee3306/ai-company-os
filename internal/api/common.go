package api

import (
	"os"
	"path/filepath"
)

func defaultDataDir() string {
	if v := os.Getenv("AI_COMPANY_OS_HOME"); v != "" {
		return v
	}
	h, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(h, ".ai-company-os")
}
