package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewBuildCommand() *cobra.Command {
	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Run kustomize build for the specified cluster or group.",
		Long: `Run kustomize build for the specified cluster or group. It returns the full configuration locally without applying it to
the cluster, allowing the user to check if all the resources are generated correctly for the target cluster.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Building the configuration...")
			return nil
		},
	}
	return buildCmd
}
