package server

import (
	"encoding/json"
	"os"

	"github.com/swlee3306/ai-company-os/internal/store"
)

// applySecretsToEnv loads secrets.json and applies keys into process env.
// NOTE: This is an MVP implementation; secrets are local-only and not audited.
func applySecretsToEnv(st *store.FileStore) {
	b, err := st.ReadSecrets()
	if err != nil {
		return
	}
	var m map[string]string
	if json.Unmarshal(b, &m) != nil {
		return
	}
	for k, v := range m {
		if v == "" {
			continue
		}
		_ = os.Setenv(k, v)
	}
}
