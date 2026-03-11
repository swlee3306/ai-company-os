package api

import (
	"encoding/json"

	"github.com/swlee3306/ai-company-os/internal/store"
)

type settingsObj struct {
	Driver struct {
		Selected string `json:"selected"`
	} `json:"driver"`
}

func readSelectedDriver(st *store.FileStore) string {
	b, err := st.ReadSettings()
	if err != nil {
		return "k3d"
	}
	var s settingsObj
	if err := json.Unmarshal(b, &s); err != nil {
		return "k3d"
	}
	if s.Driver.Selected == "k3s" {
		return "k3s"
	}
	return "k3d"
}
