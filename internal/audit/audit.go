package audit

import (
	"encoding/json"
	"time"

	"github.com/swlee3306/ai-company-os/internal/store"
)

type FileAudit struct {
	st *store.FileStore
}

func NewFileAudit(st *store.FileStore) *FileAudit {
	return &FileAudit{st: st}
}

func (a *FileAudit) Emit(actor, action string, fields map[string]any) {
	entry := map[string]any{
		"ts":     time.Now().UTC().Format(time.RFC3339Nano),
		"actor":  actor,
		"action": action,
		"fields": fields,
	}
	b, _ := json.Marshal(entry)
	_ = a.st.AppendAudit(string(b) + "\n")
}
