package server

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/store"
)

type evidenceResponse struct {
	Approval map[string]any   `json:"approval"`
	Agent    map[string]any   `json:"agent,omitempty"`
	Project  map[string]any   `json:"project,omitempty"`
	Audit    []map[string]any `json:"audit_recent"`
}

func registerEvidenceRoutes(api *gin.RouterGroup, st *store.FileStore, au *audit.FileAudit) {
	api.GET("/approvals/:id/evidence", func(c *gin.Context) {
		id := c.Param("id")
		au.Emit("api", "approvals.evidence", map[string]any{"id": id})

		// approvals
		ab, err := st.ReadApprovals()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var approvals []map[string]any
		if err := json.Unmarshal(ab, &approvals); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var appr map[string]any
		for _, it := range approvals {
			if it["id"] == id {
				appr = it
				break
			}
		}
		if appr == nil {
			c.JSON(404, gin.H{"error": "not found"})
			return
		}

		resp := evidenceResponse{Approval: appr}

		// attempt link target to agent/project by id
		target, _ := appr["target"].(string)
		if target != "" {
			// agents
			gb, _ := st.ReadAgents()
			var agents []map[string]any
			_ = json.Unmarshal(gb, &agents)
			for _, a := range agents {
				if a["id"] == target || a["name"] == target {
					resp.Agent = a
					break
				}
			}
			// projects
			pb, _ := st.ReadProjects()
			var projects []map[string]any
			_ = json.Unmarshal(pb, &projects)
			for _, p := range projects {
				if p["id"] == target {
					resp.Project = p
					break
				}
			}
		}

		// audit recent: naive parse last N json lines
		aud, _ := st.ReadAudit()
		lines := bytesLines(string(aud))
		for i := len(lines) - 1; i >= 0 && len(resp.Audit) < 10; i-- {
			var e map[string]any
			if err := json.Unmarshal([]byte(lines[i]), &e); err == nil {
				resp.Audit = append(resp.Audit, e)
			}
		}
		c.JSON(200, resp)
	})
}

func bytesLines(s string) []string {
	out := []string{}
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			line := s[start:i]
			start = i + 1
			if len(line) > 0 {
				out = append(out, line)
			}
		}
	}
	if start < len(s) {
		line := s[start:]
		if len(line) > 0 {
			out = append(out, line)
		}
	}
	return out
}
