package store

import (
	"errors"
	"os"
	"path/filepath"
)

type FileStore struct {
	base string
}

func NewFileStore(base string) *FileStore {
	return &FileStore{base: base}
}

func (s *FileStore) ensureDir() error {
	return os.MkdirAll(s.base, 0o755)
}

func (s *FileStore) path(name string) string {
	return filepath.Join(s.base, name)
}

func (s *FileStore) ReadState() ([]byte, error) {
	if err := s.ensureDir(); err != nil {
		return nil, err
	}
	p := s.path("state.json")
	b, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		// minimal default state
		return []byte(`{"driver":"k3d","health":"unknown","note":"state not initialized"}`), nil
	}
	return b, err
}

func (s *FileStore) WriteState(b []byte) error {
	if err := s.ensureDir(); err != nil {
		return err
	}
	return os.WriteFile(s.path("state.json"), b, 0o644)
}

func (s *FileStore) AppendAudit(line string) error {
	if err := s.ensureDir(); err != nil {
		return err
	}
	f, err := os.OpenFile(s.path("audit.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(line)
	return err
}

func (s *FileStore) ReadAudit() ([]byte, error) {
	if err := s.ensureDir(); err != nil {
		return nil, err
	}
	p := s.path("audit.log")
	b, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return []byte(""), nil
	}
	return b, err
}

func (s *FileStore) WriteDoctor(b []byte) error {
	if err := s.ensureDir(); err != nil {
		return err
	}
	return os.WriteFile(s.path("doctor.json"), b, 0o644)
}

func (s *FileStore) ReadDoctor() ([]byte, error) {
	if err := s.ensureDir(); err != nil {
		return nil, err
	}
	p := s.path("doctor.json")
	b, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return []byte("{}"), nil
	}
	return b, err
}

func (s *FileStore) readJSONOrDefault(name string, def string) ([]byte, error) {
	if err := s.ensureDir(); err != nil {
		return nil, err
	}
	p := s.path(name)
	b, err := os.ReadFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return []byte(def), nil
	}
	return b, err
}

func (s *FileStore) writeJSON(name string, b []byte) error {
	if err := s.ensureDir(); err != nil {
		return err
	}
	return os.WriteFile(s.path(name), b, 0o644)
}

func (s *FileStore) ReadAgents() ([]byte, error) {
	return s.readJSONOrDefault("agents.json", "[]")
}

func (s *FileStore) WriteAgents(b []byte) error {
	return s.writeJSON("agents.json", b)
}

func (s *FileStore) ReadProjects() ([]byte, error) {
	return s.readJSONOrDefault("projects.json", "[]")
}

func (s *FileStore) WriteProjects(b []byte) error {
	return s.writeJSON("projects.json", b)
}

func (s *FileStore) ReadApprovals() ([]byte, error) {
	return s.readJSONOrDefault("approvals.json", "[]")
}

func (s *FileStore) WriteApprovals(b []byte) error {
	return s.writeJSON("approvals.json", b)
}
