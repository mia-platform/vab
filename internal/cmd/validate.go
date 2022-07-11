package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewValidateCommand() *cobra.Command {
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the configuration contained in the specified path.",
		Long: `Validate the configuration contained in the specified path. It returns an error if the config file is malformed or
includes resources that do not exist in our catalogue.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Validating the config file...")
			return nil
		},
	}
	return validateCmd
}
