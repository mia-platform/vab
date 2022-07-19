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
	"path"

	"github.com/mia-platform/vab/internal/logger"
	"github.com/mia-platform/vab/internal/utils"
	"github.com/spf13/cobra"
)

const (
	minArgs = 1
	maxArgs = 2
)

// NewBuildCommand returns a new cobra.Command for building the clusters
// configuration with Kustomize
func NewBuildCommand(logger logger.LogInterface) *cobra.Command {
	buildCmd := &cobra.Command{
		Use:   "build GROUP [CLUSTER] [flags]",
		Short: "Run kustomize build for the specified cluster or group.",
		Long: `Run kustomize build for the specified cluster or group. It returns the full configuration locally without applying it to
the cluster, allowing the user to check if all the resources are generated correctly for the target cluster.`,
		Args: cobra.RangeArgs(minArgs, maxArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true

			targetPath, err := utils.GetBuildPath(args, flags.Config)
			if err != nil {
				return err
			}

			for _, clusterPath := range targetPath {
				logger.V(0).Info("### BUILD RESULTS FOR: " + clusterPath + " ###")
				targetPath := path.Join(utils.ClustersDirName, clusterPath)
				if err := utils.RunKustomizeBuild(targetPath, nil); err != nil {
					return err
				}
				logger.V(0).Infof("---")
			}

			return nil
		},
	}

	buildCmd.Flags().StringVarP(&flags.Config, "config", "c", utils.DefaultConfigFilename, "specify a different path for the configuration file")
	return buildCmd
}
