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

func registerArtifactRoutes(api *gin.RouterGroup, st *store.FileStore, au *audit.FileAudit) {
	api.GET("/artifacts", func(c *gin.Context) {
		au.Emit("api", "artifact.list", nil)
		b, err := st.ReadArtifacts()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Data(200, "application/json", b)
	})

	api.POST("/artifacts", func(c *gin.Context) {
		var body struct {
			Type      string         `json:"type"`
			Title     string         `json:"title"`
			ProjectID string         `json:"project_id"`
			TaskID    string         `json:"task_id"`
			URI       string         `json:"uri"`
			Meta      map[string]any `json:"meta"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || body.Title == "" || body.URI == "" {
			c.JSON(400, gin.H{"error": "title and uri required"})
			return
		}
		a := model.Artifact{
			ID:        fmtID("A"),
			Type:      body.Type,
			Title:     body.Title,
			ProjectID: body.ProjectID,
			TaskID:    body.TaskID,
			URI:       body.URI,
			CreatedAt: time.Now().UTC(),
			Meta:      body.Meta,
		}
		au.Emit("api", "artifact.create", map[string]any{"artifact_id": a.ID, "project_id": a.ProjectID, "task_id": a.TaskID})

		b, err := st.ReadArtifacts()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var arr []model.Artifact
		_ = json.Unmarshal(b, &arr)
		arr = append(arr, a)
		out, _ := json.MarshalIndent(arr, "", "  ")
		if err := st.WriteArtifacts(out); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, a)
	})

	api.GET("/artifacts/:id", func(c *gin.Context) {
		id := c.Param("id")
		au.Emit("api", "artifact.get", map[string]any{"id": id})
		b, err := st.ReadArtifacts()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var arr []model.Artifact
		if err := json.Unmarshal(b, &arr); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		for _, it := range arr {
			if it.ID == id {
				c.JSON(200, it)
				return
			}
		}
		c.JSON(404, gin.H{"error": "not found"})
	})
}
