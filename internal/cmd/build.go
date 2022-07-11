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

	"github.com/spf13/cobra"
)

// NewBuildCommand returns a new cobra.Command for building the clusters
// configuration with Kustomize
func NewBuildCommand() *cobra.Command {
	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Run kustomize build for the specified cluster or group.",
		Long: `Run kustomize build for the specified cluster or group. It returns the full configuration locally without applying it to
the cluster, allowing the user to check if all the resources are generated correctly for the target cluster.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Building the configuration...")
			return nil
		},
	}
	return buildCmd
}
