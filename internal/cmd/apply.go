package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewApplyCommand() *cobra.Command {
	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Build and apply the local configuration.",
		Long:  `Builds and applies the local configuration to the specified cluster or group, or to all of them.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Applying the configuration...")
		},
	}
	return applyCmd
}
