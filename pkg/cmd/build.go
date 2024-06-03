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

	"github.com/mia-platform/vab/pkg/build"
	"github.com/mia-platform/vab/pkg/logger"
	"github.com/spf13/cobra"
)

const (
	minArgs = 1
	maxArgs = 3
)

// NewBuildCommand returns a new cobra.Command for building the clusters
// configuration with Kustomize
func NewBuildCommand(logger logger.LogInterface) *cobra.Command {
	buildCmd := &cobra.Command{
		Use:   "build GROUP [CLUSTER] CONTEXT",
		Short: "Run kustomize build for the specified cluster or group searching in the given context",
		Long: `Run kustomize build for the specified cluster or group searching in the given context. It returns the full configuration locally without applying it to
the cluster, allowing the user to check if all the resources are generated correctly for the target cluster.
The configurations will be searched inside the path passed as context`,
		Args: cobra.RangeArgs(minArgs, maxArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			group := args[0]
			cluster := ""
			context := args[len(args)-1]
			if len(args) == maxArgs {
				cluster = args[1]
			}
			logger.V(10).Writef("Start build command with group name \"%s\" and cluster name \"%s\"", group, cluster)
			return build.Build(logger, flags.Config, group, cluster, context, os.Stdout)
		},
	}

	return buildCmd
}
