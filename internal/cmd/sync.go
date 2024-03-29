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
	"github.com/mia-platform/vab/internal/git"
	"github.com/mia-platform/vab/pkg/logger"
	"github.com/mia-platform/vab/pkg/sync"
	"github.com/spf13/cobra"
)

// NewSyncCommand returns a new cobra.Command for synchronizing the clusters
// configuration locally
func NewSyncCommand(logger logger.LogInterface) *cobra.Command {
	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Fetches new and updated vendor versions, and updates the clusters configuration locally.",
		Long: `Fetches new and updated vendor versions, and updates the clusters configuration locally to the latest changes of the
configuration file. After the execution, the vendors folder will include the new and updated modules/add-ons (if not
already present), and the directory structure inside the clusters folder will be updated according to the current
configuration.`,
		RunE: func(_ *cobra.Command, _ []string) error {
			logger.V(0).Writef("Synchronizing configuration at %s...", flags.Config)
			gitFilesGetter := git.RealFilesGetter{}
			return sync.Sync(logger, gitFilesGetter, flags.Config, ".", flags.DryRun)
		},
	}

	syncCmd.Flags().BoolVar(&flags.DryRun, "dry-run", false, "sync the project files without downloading new packages")
	return syncCmd
}
