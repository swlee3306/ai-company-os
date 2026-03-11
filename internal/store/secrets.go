package store

import (
	"os"
)

func (s *FileStore) secretsPath() string {
	return s.path("secrets.json")
}

func (s *FileStore) ReadSecrets() ([]byte, error) {
	return s.readJSONOrDefault("secrets.json", "{}")
}

func (s *FileStore) WriteSecrets(b []byte) error {
	if err := s.ensureDir(); err != nil {
		return err
	}
	p := s.secretsPath()
	// Write with restrictive permissions.
	if err := os.WriteFile(p, b, 0o600); err != nil {
		return err
	}
	return nil
}
