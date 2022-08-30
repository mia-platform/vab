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
	"fmt"

	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apply"
	"github.com/mia-platform/vab/pkg/logger"
	"github.com/spf13/cobra"
)

// NewApplyCommand returns a new cobra.Command for building and applying the
// clusters configuration
func NewApplyCommand(logger logger.LogInterface) *cobra.Command {
	applyCmd := &cobra.Command{
		Use:   "apply GROUP [CLUSTER] CONTEXT",
		Short: "Build and apply the local configuration.",
		Long:  `Builds and applies the local configuration to the specified cluster or group, or to all of them.`,
		Args:  cobra.RangeArgs(minArgs, maxArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.V(0).Write("Applying the configuration...")
			cmd.SilenceUsage = true
			group := args[0]
			cluster := ""
			context := args[len(args)-1]
			fmt.Println(group, cluster, context)
			if len(args) == maxArgs {
				cluster = args[1]
			}

			return apply.Apply(logger, flags.Config, flags.Output, flags.DryRun, group, cluster, context)
		},
	}
	applyCmd.Flags().StringVarP(&flags.Output, "output", "o", utils.DefaultOutputDir, "specify a different path for the applied files")
	applyCmd.Flags().BoolVar(&flags.DryRun, "dry-run", false, "sync the project files without downloading new packages")

	return applyCmd
}
