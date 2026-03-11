package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/store"
)

func registerWorkspaceRoutes(api *gin.RouterGroup, st *store.FileStore, au *audit.FileAudit) {
	api.POST("/workspace/validate", func(c *gin.Context) {
		var body struct {
			RepoPath string `json:"repo_path"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || body.RepoPath == "" {
			c.JSON(400, gin.H{"ok": false, "error": "repo_path required"})
			return
		}
		au.Emit("api", "workspace.validate", map[string]any{"repo_path": body.RepoPath})

		stt, err := os.Stat(body.RepoPath)
		if err != nil || !stt.IsDir() {
			c.JSON(200, gin.H{"ok": false, "status": "fail", "detail": "directory not found"})
			return
		}
		if _, err := os.Stat(filepath.Join(body.RepoPath, ".git")); err != nil {
			c.JSON(200, gin.H{"ok": false, "status": "warn", "detail": ".git not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true, "status": "ok"})
	})
}
