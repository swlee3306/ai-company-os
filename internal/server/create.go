package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/model"
	"github.com/swlee3306/ai-company-os/internal/store"
)

type createAgentBody struct {
	Name             string   `json:"name"`
	PersonaRole      string   `json:"persona_role"`
	OpsSpecialty     string   `json:"ops_specialty"`
	Status           string   `json:"status"`
	Scope            []string `json:"scope"`
	Version          string   `json:"version"`
	HeartbeatSeconds int      `json:"heartbeat_seconds"`
	ApprovalRequired bool     `json:"approval_required"`
	RiskScope        []string `json:"risk_scope"`
}

type createProjectBody struct {
	Name     string   `json:"name"`
	Summary  string   `json:"summary"`
	Status   string   `json:"status"`
	Phase    string   `json:"phase"`
	OwnerCEO string   `json:"owner_ceo"`
	TeamLead string   `json:"team_lead"`
	Due      string   `json:"due"`
	Agents   []string `json:"agents"`
	Evidence []string `json:"evidence_bundle"`
}

type createApprovalBody struct {
	Type      string `json:"type"`
	Requester string `json:"requester"`
	Target    string `json:"target"`
	Risk      string `json:"risk"`
	Action    string `json:"action"`
	TaskID    string `json:"task_id"`
}

func registerCreateRoutes(api *gin.RouterGroup, st *store.FileStore, au *audit.FileAudit) {
	api.POST("/agents", func(c *gin.Context) {
		var body createAgentBody
		if err := c.ShouldBindJSON(&body); err != nil || body.Name == "" || body.PersonaRole == "" {
			c.JSON(400, gin.H{"error": "name and persona_role required"})
			return
		}
		a := model.Agent{
			ID:               fmtID("AG"),
			Name:             body.Name,
			PersonaRole:      body.PersonaRole,
			OpsSpecialty:     body.OpsSpecialty,
			Status:           body.Status,
			Scope:            body.Scope,
			Version:          body.Version,
			HeartbeatSeconds: body.HeartbeatSeconds,
			ApprovalRequired: body.ApprovalRequired,
			RiskScope:        body.RiskScope,
		}
		if a.Status == "" {
			a.Status = "active"
		}
		if a.Version == "" {
			a.Version = "dev"
		}
		if a.HeartbeatSeconds == 0 {
			a.HeartbeatSeconds = 15
		}

		au.Emit("api", "agent.create", map[string]any{"id": a.ID, "name": a.Name, "persona_role": a.PersonaRole})

		b, err := st.ReadAgents()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var arr []model.Agent
		_ = json.Unmarshal(b, &arr)
		arr = append(arr, a)
		out, _ := json.MarshalIndent(arr, "", "  ")
		if err := st.WriteAgents(out); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, a)
	})

	api.POST("/projects", func(c *gin.Context) {
		var body createProjectBody
		if err := c.ShouldBindJSON(&body); err != nil || body.Name == "" || body.Summary == "" {
			c.JSON(400, gin.H{"error": "name and summary required"})
			return
		}
		p := model.Project{
			ID:       fmtID("P"),
			Name:     body.Name,
			Status:   body.Status,
			Phase:    body.Phase,
			OwnerCEO: body.OwnerCEO,
			TeamLead: body.TeamLead,
			Due:      body.Due,
			Summary:  body.Summary,
			Evidence: body.Evidence,
			Agents:   body.Agents,
		}
		if p.Status == "" {
			p.Status = "planned"
		}

		au.Emit("api", "project.create", map[string]any{"id": p.ID, "name": p.Name, "status": p.Status})

		b, err := st.ReadProjects()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var arr []model.Project
		_ = json.Unmarshal(b, &arr)
		arr = append(arr, p)
		out, _ := json.MarshalIndent(arr, "", "  ")
		if err := st.WriteProjects(out); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, p)
	})

	api.POST("/approvals", func(c *gin.Context) {
		var body createApprovalBody
		if err := c.ShouldBindJSON(&body); err != nil || body.Type == "" || body.Requester == "" || body.Target == "" || body.Risk == "" {
			c.JSON(400, gin.H{"error": "type, requester, target, risk required"})
			return
		}
		it := model.ApprovalItem{
			ID:        fmtID("appr"),
			Type:      body.Type,
			Requester: body.Requester,
			Target:    body.Target,
			Risk:      body.Risk,
			Action:    body.Action,
			Status:    "pending",
			TaskID:    body.TaskID,
		}
		if it.Action == "" {
			it.Action = "approve or reject"
		}

		au.Emit("api", "approval.create", map[string]any{"id": it.ID, "type": it.Type, "target": it.Target, "task_id": it.TaskID})

		b, err := st.ReadApprovals()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var arr []model.ApprovalItem
		_ = json.Unmarshal(b, &arr)
		arr = append(arr, it)
		out, _ := json.MarshalIndent(arr, "", "  ")
		if err := st.WriteApprovals(out); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, it)
	})

	// TODO: extend later with edit/delete
	_ = time.Second
}
