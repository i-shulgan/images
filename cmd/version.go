package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	gitRevision = "HEAD"
	buildStamp  = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	// Does not need the static files directory.
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Git Revision: %s\n", gitRevision)
		fmt.Printf("UTC Build Time: %s\n", buildStamp)
		os.Exit(0)
	},
}
