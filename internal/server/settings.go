package server

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/store"
)

type settingsObj struct {
	Driver struct {
		Selected string `json:"selected"`
	} `json:"driver"`
	Approval struct {
		PolicyText string `json:"policy_text"`
	} `json:"approval"`
	Runner struct {
		Backend string `json:"backend"`
		Type    string `json:"type"`
		Command string `json:"command"`
		Workdir string `json:"workdir"`
		Agents  struct {
			PM       string `json:"pm"`
			FE       string `json:"fe"`
			BE       string `json:"be"`
			QA       string `json:"qa"`
			Reviewer string `json:"reviewer"`
		} `json:"agents"`
	} `json:"runner"`
}

func registerSettingsRoutes(api *gin.RouterGroup, st *store.FileStore, au *audit.FileAudit) {
	api.GET("/settings", func(c *gin.Context) {
		au.Emit("api", "settings.get", nil)
		b, err := st.ReadSettings()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Data(200, "application/json", b)
	})

	api.POST("/settings", func(c *gin.Context) {
		var body settingsObj
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "invalid json"})
			return
		}
		if body.Driver.Selected != "k3d" && body.Driver.Selected != "k3s" {
			c.JSON(400, gin.H{"error": "driver.selected must be k3d|k3s"})
			return
		}
		switch body.Runner.Type {
		case "", "codex_cli", "claude_code", "gemini_cli", "custom":
		default:
			c.JSON(400, gin.H{"error": "runner.type must be codex_cli|claude_code|gemini_cli|custom"})
			return
		}
		switch body.Runner.Backend {
		case "", "local_placeholder", "local_cli", "openclaw_acp":
		default:
			c.JSON(400, gin.H{"error": "runner.backend must be local_placeholder|local_cli|openclaw_acp"})
			return
		}
		au.Emit("api", "settings.update", map[string]any{"driver": body.Driver.Selected, "runner": body.Runner.Type, "runner_backend": body.Runner.Backend})
		out, _ := json.MarshalIndent(body, "", "  ")
		if err := st.WriteSettings(out); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Data(200, "application/json", out)
	})
}
