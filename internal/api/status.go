package api

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/store"
)

func init() {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show current system status (reads local state store)",
		RunE: func(cmd *cobra.Command, args []string) error {
			st := store.NewFileStore(defaultDataDir())
			au := audit.NewFileAudit(st)
			au.Emit("cli", "status", nil)
			b, err := st.ReadState()
			if err != nil {
				return err
			}
			var obj any
			if err := json.Unmarshal(b, &obj); err != nil {
				return err
			}
			out, _ := json.MarshalIndent(obj, "", "  ")
			fmt.Println(string(out))
			return nil
		},
	}
	rootCmd.AddCommand(cmd)
}
