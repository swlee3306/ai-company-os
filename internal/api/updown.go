package api

import (
	"github.com/spf13/cobra"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/driver"
	"github.com/swlee3306/ai-company-os/internal/store"
)

func init() {
	upCmd := &cobra.Command{
		Use:   "up",
		Short: "Start the local runtime (k3d cluster)",
		RunE: func(cmd *cobra.Command, args []string) error {
			st := store.NewFileStore(defaultDataDir())
			au := audit.NewFileAudit(st)
			name, _ := cmd.Flags().GetString("cluster")
			if name == "" {
				name = "company-os"
			}
			au.Emit("cli", "driver.k3d.up", map[string]any{"cluster": name})
			return driver.K3DUp(name)
		},
	}
	upCmd.Flags().String("cluster", "", "k3d cluster name (default company-os)")

	downCmd := &cobra.Command{
		Use:   "down",
		Short: "Stop the local runtime (k3d cluster)",
		RunE: func(cmd *cobra.Command, args []string) error {
			st := store.NewFileStore(defaultDataDir())
			au := audit.NewFileAudit(st)
			name, _ := cmd.Flags().GetString("cluster")
			if name == "" {
				name = "company-os"
			}
			au.Emit("cli", "driver.k3d.down", map[string]any{"cluster": name})
			return driver.K3DDown(name)
		},
	}
	downCmd.Flags().String("cluster", "", "k3d cluster name (default company-os)")

	rootCmd.AddCommand(upCmd)
	rootCmd.AddCommand(downCmd)
}
