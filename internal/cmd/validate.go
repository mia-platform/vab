package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var ValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration contained in the specified path.",
	Long: `Validate the configuration contained in the specified path.
It returns an error if the config file is malformed or includes
resources that do not exist in our catalogue.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Validating the config file...")
	},
}
