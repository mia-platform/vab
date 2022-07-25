// Copyright 2022 Mia-Platform

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"
	"path"

	"github.com/mia-platform/vab/pkg/logger"
	"github.com/mia-platform/vab/pkg/validate"
	"github.com/spf13/cobra"
)

// NewValidateCommand returns a new cobra.Command for validating the
// configuration file
func NewValidateCommand(logger logger.LogInterface) *cobra.Command {
	validateCmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate the configuration contained in the specified path.",
		Long: `Validate the configuration contained in the specified path. It returns an error if the config file is malformed or
includes resources that do not exist in our catalogue.`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			configPath := path.Clean(flags.Config)
			logger.V(0).Writef("Validating configuration at %s...", configPath)
			return validate.ConfigurationFile(logger, configPath, os.Stdout)
		},
	}

	return validateCmd
}
