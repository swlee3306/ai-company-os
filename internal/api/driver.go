package api

import (
	"encoding/json"

	"github.com/spf13/cobra"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/driver"
	"github.com/swlee3306/ai-company-os/internal/store"
)

func init() {
	// status is already implemented as reading state.json; enhance by writing computed driver status.
	cmd := &cobra.Command{
		Use:   "driver-check",
		Short: "Run driver checks (Docker/k3d) and persist to state.json (debug helper)",
		RunE: func(cmd *cobra.Command, args []string) error {
			st := store.NewFileStore(defaultDataDir())
			au := audit.NewFileAudit(st)
			au.Emit("cli", "driver.check", nil)

			res := driver.CheckAll()
			b, _ := json.MarshalIndent(map[string]any{"driver": res}, "", "  ")
			return st.WriteState(b)
		},
	}
	rootCmd.AddCommand(cmd)
}
