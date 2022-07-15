// Copyright 2022 Mia-Platform

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/mia-platform/vab/internal/utils"
	"github.com/spf13/cobra"
)

// NewValidateCommand returns a new cobra.Command for validating the
// configuration file
func NewValidateCommand() *cobra.Command {
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the configuration contained in the specified path.",
		Long: `Validate the configuration contained in the specified path. It returns an error if the config file is malformed or
includes resources that do not exist in our catalogue.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return nil
			}
			return fmt.Errorf("no argument %s expected here", args[0])
		},
		Run: func(cmd *cobra.Command, args []string) { // TODO: consider integrating a logger
			os.Exit(utils.ValidateConfig(flags.Config, nil))
		},
	}

	validateCmd.Flags().StringVarP(&flags.Config, "config", "c", defaultConfigPath, "specify a different path for the configuration file")
	return validateCmd
}
