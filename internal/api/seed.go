package api

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/model"
	"github.com/swlee3306/ai-company-os/internal/store"
)

func init() {
	cmd := &cobra.Command{
		Use:   "seed",
		Short: "Seed local store with demo data (agents/projects/approvals)",
		RunE: func(cmd *cobra.Command, args []string) error {
			st := store.NewFileStore(defaultDataDir())
			au := audit.NewFileAudit(st)
			au.Emit("cli", "seed", nil)

			agents := []model.Agent{
				{ID: "planner-manager-01", Name: "planner-manager-01", PersonaRole: "PM", Status: "active", Scope: []string{"org:*", "task:assign"}, Version: "v1.8.2", HeartbeatSeconds: 8, ApprovalRequired: false},
				{ID: "backend-worker-03", Name: "backend-worker-03", PersonaRole: "BE", Status: "active", Scope: []string{"repo:billing-core"}, Version: "v2.4.1", HeartbeatSeconds: 15, ApprovalRequired: false},
				{ID: "devops-worker-01", Name: "devops-worker-01", PersonaRole: "BE", OpsSpecialty: "DevOps", Status: "blocked", Scope: []string{"deploy:staging"}, Version: "v1.3.0", HeartbeatSeconds: 42, ApprovalRequired: true, RiskScope: []string{"deploy:prod", "secret:read"}},
			}
			projects := []model.Project{
				{ID: "billing-core", Name: "billing-core migration", Status: "running", Phase: "schema + application cutover prep", OwnerCEO: "Flant", TeamLead: "openclaw", Due: "Mar 18", Summary: "Critical backend modernization program with staged rollout and approval gates.", Evidence: []string{"ADR-014", "schema-diff-v7", "benchmark-2026-03-09", "rollback-checklist.md"}, Agents: []string{"planner-manager-01", "backend-worker-03", "devops-worker-01"}},
				{ID: "agent-runtime-v2", Name: "agent-runtime v2", Status: "reviewing", Phase: "perf + stability", OwnerCEO: "Flant", TeamLead: "openclaw", Due: "", Summary: "Runtime upgrade with tighter audit + approval semantics."},
			}
			approvals := []model.ApprovalItem{
				{ID: "appr-1", Type: "production deploy", Requester: "PM (planner-manager-01)", Target: "billing-core", Risk: "HIGH", Action: "approve or reject", Status: "pending", TaskID: "T-184"},
				{ID: "appr-2", Type: "tool permission", Requester: "Team Lead (manager-ai)", Target: "devops-worker-01", Risk: "HIGH", Action: "approve or reject", Status: "pending"},
				{ID: "appr-3", Type: "agent activation", Requester: "Owner (CEO)", Target: "QA (qa-reviewer-02)", Risk: "MEDIUM", Action: "approve or reject", Status: "pending"},
			}

			ab, _ := json.MarshalIndent(agents, "", "  ")
			pb, _ := json.MarshalIndent(projects, "", "  ")
			qb, _ := json.MarshalIndent(approvals, "", "  ")
			if err := st.WriteAgents(ab); err != nil {
				return err
			}
			if err := st.WriteProjects(pb); err != nil {
				return err
			}
			if err := st.WriteApprovals(qb); err != nil {
				return err
			}
			if err := seedTasks(st); err != nil {
				return err
			}

			return nil
		},
	}
	rootCmd.AddCommand(cmd)
}
