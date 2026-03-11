package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/swlee3306/ai-company-os/internal/driver"

	"github.com/spf13/cobra"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/store"
)

func init() {
	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Run diagnostics (MVP: mock checks, stores last run)",
		RunE: func(cmd *cobra.Command, args []string) error {
			st := store.NewFileStore(defaultDataDir())
			au := audit.NewFileAudit(st)
			au.Emit("cli", "doctor.run", nil)

			res := map[string]any{
				"ts":      time.Now().UTC().Format(time.RFC3339),
				"overall": "ok",
				"checks": []map[string]any{
					{"name": "filesystem", "status": "ok"},
					{"name": "network", "status": "ok"},
				},
				"driver": driver.CheckAll(readSelectedDriver(st)),
			}
			b, _ := json.Marshal(res)
			if err := st.WriteDoctor(b); err != nil {
				return err
			}
			out, _ := json.MarshalIndent(res, "", "  ")
			fmt.Println(string(out))
			return nil
		},
	}
	rootCmd.AddCommand(cmd)
}
