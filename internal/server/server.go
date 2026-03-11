package server

import (
	"encoding/json"
	"net/http"

	"github.com/swlee3306/ai-company-os/internal/driver"

	"github.com/gin-gonic/gin"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/store"
)

func Run(addr string, st *store.FileStore, au *audit.FileAudit) error {
	r := gin.New()
	r.Use(gin.Recovery())

	// Basic CORS for local dashboard dev
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	registerEvidenceRoutes(api, st, au)
	registerTaskRoutes(api, st, au)
	registerArtifactRoutes(api, st, au)
	api.GET("/status", func(c *gin.Context) {
		au.Emit("api", "status.get", nil)
		b, err := st.ReadState()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		base := map[string]any{}
		_ = json.Unmarshal(b, &base)
		base["driver"] = driver.CheckAll()
		out, _ := json.MarshalIndent(base, "", "  ")
		c.Data(200, "application/json", out)
	})

	api.GET("/audit", func(c *gin.Context) {
		au.Emit("api", "audit.list", nil)
		b, err := st.ReadAudit()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Data(200, "text/plain", b)
	})

	api.GET("/doctor", func(c *gin.Context) {
		au.Emit("api", "doctor.get", nil)
		b, err := st.ReadDoctor()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Data(200, "application/json", b)
	})

	api.GET("/agents", func(c *gin.Context) {
		au.Emit("api", "agents.list", nil)
		b, err := st.ReadAgents()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Data(200, "application/json", b)
	})

	api.GET("/agents/:id", func(c *gin.Context) {
		au.Emit("api", "agents.get", map[string]any{"id": c.Param("id")})
		b, err := st.ReadAgents()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var arr []map[string]any
		if err := json.Unmarshal(b, &arr); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		for _, it := range arr {
			if it["id"] == c.Param("id") {
				c.JSON(200, it)
				return
			}
		}
		c.JSON(404, gin.H{"error": "not found"})
	})

	api.GET("/projects", func(c *gin.Context) {
		au.Emit("api", "projects.list", nil)
		b, err := st.ReadProjects()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Data(200, "application/json", b)
	})

	api.GET("/projects/:id", func(c *gin.Context) {
		au.Emit("api", "projects.get", map[string]any{"id": c.Param("id")})
		b, err := st.ReadProjects()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var arr []map[string]any
		if err := json.Unmarshal(b, &arr); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		for _, it := range arr {
			if it["id"] == c.Param("id") {
				c.JSON(200, it)
				return
			}
		}
		c.JSON(404, gin.H{"error": "not found"})
	})

	api.GET("/approvals", func(c *gin.Context) {
		au.Emit("api", "approvals.list", nil)
		b, err := st.ReadApprovals()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Data(200, "application/json", b)
	})

	api.GET("/approvals/:id", func(c *gin.Context) {
		au.Emit("api", "approvals.get", map[string]any{"id": c.Param("id")})
		b, err := st.ReadApprovals()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var arr []map[string]any
		if err := json.Unmarshal(b, &arr); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		for _, it := range arr {
			if it["id"] == c.Param("id") {
				c.JSON(200, it)
				return
			}
		}
		c.JSON(404, gin.H{"error": "not found"})
	})

	api.POST("/approvals/:id/decision", func(c *gin.Context) {
		id := c.Param("id")
		var body struct {
			Decision string `json:"decision"`
			Reason   string `json:"reason"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "invalid json"})
			return
		}
		if body.Decision != "approve" && body.Decision != "reject" {
			c.JSON(400, gin.H{"error": "decision must be approve|reject"})
			return
		}
		if body.Decision == "reject" && body.Reason == "" {
			c.JSON(400, gin.H{"error": "reason required for reject"})
			return
		}

		au.Emit("api", "approvals.decision", map[string]any{"id": id, "decision": body.Decision})
		b, err := st.ReadApprovals()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var arr []map[string]any
		if err := json.Unmarshal(b, &arr); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		for _, it := range arr {
			if it["id"] == id {
				it["status"] = body.Decision
				it["decision_reason"] = body.Reason

				// MVP side-effects: update related agent/project state for demo flows
				applyApprovalSideEffects(st, au, it)

				out, _ := json.MarshalIndent(arr, "", "  ")
				if err := st.WriteApprovals(out); err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				c.JSON(200, it)
				return
			}
		}
		c.JSON(404, gin.H{"error": "not found"})
	})

	return r.Run(addr)
}
