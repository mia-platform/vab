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

package create

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"

	cmdutil "github.com/mia-platform/vab/pkg/cmd/util"
)

const (
	shortCmd = "Initialize a vab Project"
	longCmd  = `Initialize a vab project with a preliminary directory structure, together
	with the skeleton of the configuration file.

	The project directory will contain the clusters directory (including the all-groups
	folder with a minimal kustomize configuration), and the configuration file.`

	pathArgHelpText     = "Specify the path where to create the project"
	tooManyArgsHelpText = "Too many arguments"
)

// Flags contains all the flags for the `create` command. They will be converted to Options
// that contains all runtime options for the command.
type Flags struct{}

// Options have the data required to perform the create operation
type Options struct {
	path   string
	logger logr.Logger
}

// NewCommand return the command for creating a new configuration file and basic folder structures
func NewCommand() *cobra.Command {
	flags := &Flags{}
	cmd := &cobra.Command{
		Use:                   "create PATH",
		Aliases:               []string{"init"},
		Short:                 heredoc.Doc(shortCmd),
		Long:                  heredoc.Doc(longCmd),
		DisableFlagsInUseLine: true,

		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: validArgs,

		Run: func(cmd *cobra.Command, args []string) {
			options, err := flags.ToOptions(args)
			cobra.CheckErr(err)
			cobra.CheckErr(options.Run(cmd.Context()))
		},
	}

	return cmd
}

// ToOptions transform the command flags in command runtime arguments
func (f *Flags) ToOptions(args []string) (*Options, error) {
	contextPath, err := cmdutil.ValidateContextPath(args[0])
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}

	return &Options{
		path: contextPath,
	}, nil
}

// Run execute the create command
func (o *Options) Run(ctx context.Context) error {
	o.logger = logr.FromContextOrDiscard(ctx)

	name := filepath.Base(o.path)
	return cmdutil.InitializeConfiguration(name, o.path)
}

func validArgs(_ *cobra.Command, args []string, _ string) ([]string, cobra.ShellCompDirective) {
	var comps []string
	var directive cobra.ShellCompDirective
	switch len(args) {
	case 0:
		comps = cobra.AppendActiveHelp(comps, pathArgHelpText)
		directive = cobra.ShellCompDirectiveDefault
	default:
		comps = cobra.AppendActiveHelp(comps, tooManyArgsHelpText)
		directive = cobra.ShellCompDirectiveNoFileComp
	}
	return comps, directive
}
