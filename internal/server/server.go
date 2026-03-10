package server

import (
	"net/http"

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
	api.GET("/status", func(c *gin.Context) {
		au.Emit("api", "status.get", nil)
		b, err := st.ReadState()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Data(200, "application/json", b)
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

	return r.Run(addr)
}
