package api

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/server"
	"github.com/swlee3306/ai-company-os/internal/store"
)

func init() {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the local API server for the dashboard",
		RunE: func(cmd *cobra.Command, args []string) error {
			dataDir := defaultDataDir()
			st := store.NewFileStore(dataDir)
			au := audit.NewFileAudit(st)

			addr, _ := cmd.Flags().GetString("listen")
			if addr == "" {
				addr = "127.0.0.1:8787"
			}

			au.Emit("system", "server.start", map[string]any{"addr": addr})
			return server.Run(addr, st, au)
		},
	}
	cmd.Flags().String("listen", "", "listen address (default 127.0.0.1:8787)")
	rootCmd.AddCommand(cmd)
}

func fmtAny(v any) string { return fmt.Sprintf("%v", v) }
