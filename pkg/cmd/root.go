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
	"context"
	"fmt"
	"log"
	"runtime"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/go-logr/logr"
	"github.com/go-logr/stdr"
	"github.com/mia-platform/vab/pkg/cmd/apply"
	"github.com/mia-platform/vab/pkg/cmd/build"
	"github.com/mia-platform/vab/pkg/cmd/create"
	"github.com/mia-platform/vab/pkg/cmd/sync"
	"github.com/mia-platform/vab/pkg/cmd/util"
	"github.com/mia-platform/vab/pkg/cmd/validate"
	"github.com/spf13/cobra"
)

var (
	// Version is dynamically set by the ci or overridden by the Makefile.
	Version = "DEV"
	// BuildDate is dynamically set at build time by the cli or overridden in the Makefile.
	BuildDate = "" // YYYY-MM-DD
)

const (
	vabCmdShort = "vab is used for managing and installing the Magellano k8s distro on your cluster(s)"
	vabCmdLong  = `vab is used for managing and installing the Magellano k8s distro on your cluster(s)

	It will manage folders for separating kustomize patches for your clusters, downloads
	the modules and add-ons, and then apply the resulting manifests to your clusters.

	More information about the Magellano k8s distribution can be found here:
		https://github.com/mia-platform/distribution`
)

// NewVabCommand creates the `vab` command and its nested children.
func NewVabCommand() *cobra.Command {
	configFlags := util.NewConfigFlags()
	cmd := &cobra.Command{
		Use: "vab",

		Short: heredoc.Doc(vabCmdShort),
		Long:  heredoc.Doc(vabCmdLong),

		SilenceErrors: true,
		Version:       versionString(),

		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		PersistentPreRun: func(*cobra.Command, []string) {
			stdr.SetVerbosity(*configFlags.Verbose)
		},
	}

	cmd.SetContext(logr.NewContext(context.Background(), stdr.New(log.Default())))
	configFlags.AddFlags(cmd.PersistentFlags())

	cmd.AddCommand(
		create.NewCommand(),
		apply.NewCommand(configFlags),
		build.NewCommand(configFlags),
		validate.NewCommand(configFlags),
		sync.NewCommand(configFlags),
	)
	return cmd
}

// versionString format a complete version string to output to the user
func versionString() string {
	version := Version

	if BuildDate != "" {
		version = fmt.Sprintf("%s (%s)", version, BuildDate)
	}

	return fmt.Sprintf("%s, Go Version: %s", version, runtime.Version())
}
