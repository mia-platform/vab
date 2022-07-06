package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Build and apply the local configuration.",
	Long: `Builds and applies the local configuration to the specified
					cluster, group, or to all of them.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Applying the configuration...")
	},
}
