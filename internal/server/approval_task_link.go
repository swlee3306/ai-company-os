package server

import (
	"encoding/json"
	"time"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/model"
	"github.com/swlee3306/ai-company-os/internal/store"
)

func applyApprovalTaskLink(st *store.FileStore, au *audit.FileAudit, approval map[string]any) {
	decision, _ := approval["status"].(string)
	taskID, _ := approval["task_id"].(string)
	if taskID == "" {
		return
	}
	if decision != "approve" {
		return
	}

	b, err := st.ReadTasks()
	if err != nil {
		return
	}
	var arr []model.Task
	if err := json.Unmarshal(b, &arr); err != nil {
		return
	}
	for i := range arr {
		if arr[i].ID == taskID {
			from := arr[i].State
			if from == "blocked" {
				arr[i].State = "planned"
				arr[i].UpdatedAt = time.Now().UTC()
				au.Emit("system", "task.transition", map[string]any{"id": taskID, "from": from, "to": "planned", "cause": "approval"})
				out, _ := json.MarshalIndent(arr, "", "  ")
				_ = st.WriteTasks(out)
			}
			return
		}
	}
}
