package store

import (
	"fmt"
	"os"
	"path/filepath"
)

func (s *FileStore) RunsDir() string {
	return filepath.Join(s.base, "runs")
}

func (s *FileStore) EnsureRunsDir() error {
	return os.MkdirAll(s.RunsDir(), 0o755)
}

func (s *FileStore) RunDir(id string) string {
	return filepath.Join(s.RunsDir(), id)
}

func (s *FileStore) EnsureRunDir(id string) error {
	if id == "" {
		return fmt.Errorf("run id required")
	}
	return os.MkdirAll(s.RunDir(id), 0o755)
}
