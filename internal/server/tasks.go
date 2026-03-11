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

func registerTaskRoutes(api *gin.RouterGroup, st *store.FileStore, au *audit.FileAudit) {
	api.GET("/tasks", func(c *gin.Context) {
		au.Emit("api", "task.list", nil)
		b, err := st.ReadTasks()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.Data(200, "application/json", b)
	})

	api.POST("/tasks", func(c *gin.Context) {
		var body struct {
			Title string `json:"title"`
			Desc  string `json:"desc"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || body.Title == "" {
			c.JSON(400, gin.H{"error": "title required"})
			return
		}
		now := time.Now().UTC()
		t := model.Task{
			ID:               fmtID("T"),
			Title:            body.Title,
			Desc:             body.Desc,
			State:            "draft",
			ReviewerRequired: true,
			CreatedAt:        now,
			UpdatedAt:        now,
		}
		au.Emit("api", "task.create", map[string]any{"id": t.ID})

		b, err := st.ReadTasks()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var arr []model.Task
		_ = json.Unmarshal(b, &arr)
		arr = append(arr, t)
		out, _ := json.MarshalIndent(arr, "", "  ")
		if err := st.WriteTasks(out); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, t)
	})

	api.GET("/tasks/:id", func(c *gin.Context) {
		id := c.Param("id")
		au.Emit("api", "task.get", map[string]any{"id": id})
		b, err := st.ReadTasks()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var arr []model.Task
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

	api.POST("/tasks/:id/transition", func(c *gin.Context) {
		id := c.Param("id")
		var body struct {
			To string `json:"to"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || body.To == "" {
			c.JSON(400, gin.H{"error": "to required"})
			return
		}
		b, err := st.ReadTasks()
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		var arr []model.Task
		if err := json.Unmarshal(b, &arr); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		for i := range arr {
			if arr[i].ID == id {
				from := arr[i].State
				to := body.To
				if !validTransition(from, to) {
					c.JSON(400, gin.H{"error": "invalid transition"})
					return
				}
				arr[i].State = to
				arr[i].UpdatedAt = time.Now().UTC()
				au.Emit("api", "task.transition", map[string]any{"id": id, "from": from, "to": to})
				out, _ := json.MarshalIndent(arr, "", "  ")
				if err := st.WriteTasks(out); err != nil {
					c.JSON(500, gin.H{"error": err.Error()})
					return
				}
				c.JSON(200, arr[i])
				return
			}
		}
		c.JSON(404, gin.H{"error": "not found"})
	})
}

func validTransition(from, to string) bool {
	if to == "blocked" {
		return true
	}
	if from == "blocked" && to == "planned" {
		return true
	}
	switch from {
	case "draft":
		return to == "planned"
	case "planned":
		return to == "assigned"
	case "assigned":
		return to == "running"
	case "running":
		return to == "reviewing"
	case "reviewing":
		return to == "done"
	default:
		return false
	}
}

func fmtID(prefix string) string {
	return prefix + "-" + time.Now().UTC().Format("20060102-150405")
}
