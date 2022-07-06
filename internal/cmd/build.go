package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var BuildCmd = &cobra.Command{
	Use:   "build",
	Short: "Run kustomize build for the specified cluster or group.",
	Long: `Run kustomize build for the specified cluster or group.
It returns the full configuration locally without applying it to the cluster,
allowing the user to check if all the resources are generated correctly
for the target cluster.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Building the configuration...")
	},
}
