// Copyright Mia srl
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"

	initProj "github.com/mia-platform/vab/pkg/init"
	"github.com/mia-platform/vab/pkg/logger"
	"github.com/spf13/cobra"
)

// NewInitCommand returns a new cobra.Command for initializing the project
func NewInitCommand(logger logger.LogInterface) *cobra.Command {
	initCmd := &cobra.Command{

		Use:   "init",
		Short: "Initialize the vab project",
		Long: `Creates the project folder with a preliminary directory structure, together with the skeleton of the configuration file.
The project directory will contain the clusters directory (including the all-groups folder with a minimal kustomize
configuration), and the configuration file.`,

		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SilenceUsage = true
			logger.V(0).Writef("Initializing...")

			currentPath, err := os.Getwd()
			if err != nil {
				logger.V(10).Write("Error trying to access the current path")
				return err
			}

			return initProj.NewProject(logger, currentPath, flags.Name)
		},
	}

	initCmd.Flags().StringVarP(&flags.Name, "name", "n", "", "project name, defaults to current directory name")
	return initCmd
}
