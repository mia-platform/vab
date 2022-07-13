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

	"github.com/mia-platform/vab/internal/utils"
	"github.com/spf13/cobra"
)

const (
	minArgs           = 1
	maxArgs           = 2
	defaultConfigPath = "./config.yaml"
)

// NewBuildCommand returns a new cobra.Command for building the clusters
// configuration with Kustomize
func NewBuildCommand() *cobra.Command {
	buildCmd := &cobra.Command{
		Use:   "build GROUP [CLUSTER] [flags]",
		Short: "Run kustomize build for the specified cluster or group.",
		Long: `Run kustomize build for the specified cluster or group. It returns the full configuration locally without applying it to
the cluster, allowing the user to check if all the resources are generated correctly for the target cluster.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < minArgs {
				return fmt.Errorf("at least the cluster group is required")
			}
			if len(args) > maxArgs {
				return fmt.Errorf("too many args")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			fmt.Println("Building the configuration...")

			targetPath, err := utils.GetBuildPath(args, flags.Config)
			if err != nil {
				return err
			}

			for _, clusterPath := range targetPath {
				if err := utils.RunKustomizeBuild(clusterPath, nil); err != nil {
					return err
				}
			}

			return nil
		},
	}

	buildCmd.Flags().StringVarP(&flags.Config, "config", "c", defaultConfigPath, "specify a different path for the configuration file")
	return buildCmd
}
