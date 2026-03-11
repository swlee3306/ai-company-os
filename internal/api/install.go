package api

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/swlee3306/ai-company-os/internal/audit"
	"github.com/swlee3306/ai-company-os/internal/installer"
	"github.com/swlee3306/ai-company-os/internal/store"
)

func init() {
	installCmd := &cobra.Command{
		Use:   "install",
		Short: "Install optional dependencies for Company OS",
	}

	k3dCmd := &cobra.Command{
		Use:   "k3d",
		Short: "Install k3d",
		RunE: func(cmd *cobra.Command, args []string) error {
			st := store.NewFileStore(defaultDataDir())
			au := audit.NewFileAudit(st)
			au.Emit("cli", "install.k3d", map[string]any{"os": runtime.GOOS})

			dry, _ := cmd.Flags().GetBool("dry-run")
			if dry {
				fmt.Println(installer.K3DPlan(runtime.GOOS))
				return nil
			}
			return installer.InstallK3D(runtime.GOOS)
		},
	}
	k3dCmd.Flags().Bool("dry-run", false, "print install plan without executing")

	installCmd.AddCommand(k3dCmd)
	rootCmd.AddCommand(installCmd)
}
