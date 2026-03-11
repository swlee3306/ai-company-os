package api

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/driver"
	"github.com/swlee3306/ai-company-os/internal/store"
)

func init() {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current system status",
		RunE: func(cmd *cobra.Command, args []string) error {
			st := store.NewFileStore(defaultDataDir())
			au := audit.NewFileAudit(st)
			au.Emit("cli", "status", nil)

			base := map[string]any{}
			b, err := st.ReadState()
			if err == nil {
				_ = json.Unmarshal(b, &base)
			}

			base["driver"] = driver.CheckAll(readSelectedDriver(st))

			out, _ := json.MarshalIndent(base, "", "  ")
			fmt.Println(string(out))
			return nil
		},
	}
	rootCmd.AddCommand(cmd)
}
