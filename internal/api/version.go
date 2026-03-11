package api

import (
	"fmt"

	"github.com/spf13/cobra"
)

// These can be set via -ldflags at build time.
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func init() {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("company-os %s\ncommit: %s\nbuilt: %s\n", Version, Commit, BuildTime)
		},
	}
	rootCmd.AddCommand(cmd)
}
