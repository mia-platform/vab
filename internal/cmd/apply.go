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
	"github.com/spf13/cobra"

	"github.com/mia-platform/vab/internal/logger"
)

// NewApplyCommand returns a new cobra.Command for building and applying the
// clusters configuration
func NewApplyCommand(logger logger.LogInterface) *cobra.Command {
	applyCmd := &cobra.Command{
		Use:   "apply",
		Short: "Build and apply the local configuration.",
		Long:  `Builds and applies the local configuration to the specified cluster or group, or to all of them.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.V(0).Info("Applying the configuration...")
			return nil
		},
	}
	return applyCmd
}
