package server

import (
	"encoding/json"

	"github.com/swlee3306/ai-company-os/internal/store"
)

// applyApprovalSideEffects is an MVP rule engine that updates related store objects
// to reflect approval decisions. This is intentionally simple and deterministic.
func applyApprovalSideEffects(st *store.FileStore, approval map[string]any) {
	decision, _ := approval["status"].(string)
	apprType, _ := approval["type"].(string)
	target, _ := approval["target"].(string)

	if target == "" {
		return
	}

	// tool permission: approving unblocks agent; rejecting keeps blocked
	if apprType == "tool permission" {
		b, err := st.ReadAgents()
		if err != nil {
			return
		}
		var agents []map[string]any
		if err := json.Unmarshal(b, &agents); err != nil {
			return
		}
		changed := false
		for _, a := range agents {
			if a["id"] == target || a["name"] == target {
				if decision == "approve" {
					a["status"] = "active"
					a["approval_required"] = false
					changed = true
				}
				if decision == "reject" {
					a["status"] = "blocked"
					a["approval_required"] = true
					changed = true
				}
			}
		}
		if changed {
			out, _ := json.MarshalIndent(agents, "", "  ")
			_ = st.WriteAgents(out)
		}
		return
	}

	// production deploy: set project status to running on approve; blocked on reject
	if apprType == "production deploy" {
		b, err := st.ReadProjects()
		if err != nil {
			return
		}
		var projects []map[string]any
		if err := json.Unmarshal(b, &projects); err != nil {
			return
		}
		changed := false
		for _, p := range projects {
			if p["id"] == target {
				if decision == "approve" {
					p["status"] = "running"
					changed = true
				}
				if decision == "reject" {
					p["status"] = "blocked"
					changed = true
				}
			}
		}
		if changed {
			out, _ := json.MarshalIndent(projects, "", "  ")
			_ = st.WriteProjects(out)
		}
		return
	}

	// agent activation: no-op for MVP beyond status
}
