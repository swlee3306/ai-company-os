package api

import (
	"encoding/json"
	"time"

	"github.com/swlee3306/ai-company-os/internal/model"
	"github.com/swlee3306/ai-company-os/internal/store"
)

func seedTasks(st *store.FileStore) error {
	now := time.Now().UTC()
	tasks := []model.Task{
		{ID: "T-184", Title: "write migration rollback validator", Desc: "smoke seeded blocked task", State: "blocked", Assignee: "backend-worker-03", ReviewerRequired: true, CreatedAt: now, UpdatedAt: now},
	}
	b, _ := json.MarshalIndent(tasks, "", "  ")
	return st.WriteTasks(b)
}
