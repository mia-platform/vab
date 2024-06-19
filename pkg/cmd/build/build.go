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

package build

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/go-logr/logr"
	"github.com/mia-platform/vab/pkg/cmd/util"
	"github.com/spf13/cobra"
)

const (
	shortCmd = "Show local files that will be applied"
	longCmd  = `Run kustomize build for the specified cluster or group searching
	in the given context.
	It returns the full configuration locally without applying it to the cluster,
	allowing the user to check if all the resources are generated correctly for
	the target cluster.

	The configurations will be searched inside the path passed as context`
	cmdUsage = "build GROUP [CLUSTER] CONTEXT"

	minArgs = 2
	maxArgs = 3
)

// Flags contains all the flags for the `build` command. They will be converted to Options
// that contains all runtime options for the command
type Flags struct{}

// Options have the data required to perform the apply operation
type Options struct {
	group       string
	cluster     string
	contextPath string
	configPath  string
	writer      io.Writer
	logger      logr.Logger
}

// NewCommand return the command for showing the manifests for every group and cluster requested
// that will be applied to the remote server
func NewCommand(cf *util.ConfigFlags) *cobra.Command {
	flags := &Flags{}

	cmd := &cobra.Command{
		Use:   cmdUsage,
		Short: heredoc.Doc(shortCmd),
		Long:  heredoc.Doc(longCmd),

		Args: cobra.RangeArgs(minArgs, maxArgs),
		Run: func(cmd *cobra.Command, args []string) {
			options, err := flags.ToOptions(cf, args, cmd.OutOrStdout())
			cobra.CheckErr(err)
			cobra.CheckErr(options.Run(cmd.Context()))
		},
	}

	return cmd
}

// ToOptions transform the command flags in command runtime arguments
func (f *Flags) ToOptions(cf *util.ConfigFlags, args []string, writer io.Writer) (*Options, error) {
	if len(args) < minArgs {
		return nil, fmt.Errorf("at least %d arguments are needed", minArgs)
	}

	group := args[0]
	cluster := ""
	contextPath := args[len(args)-1]
	if len(args) >= maxArgs {
		cluster = args[1]
	}

	var err error
	var cleanedContextPath string
	if cleanedContextPath, err = filepath.Abs(contextPath); err != nil {
		return nil, err
	}

	var contextInfo fs.FileInfo
	if contextInfo, err = os.Stat(cleanedContextPath); err != nil {
		return nil, fmt.Errorf("error locating files: %w", err)
	}

	if !contextInfo.IsDir() {
		return nil, fmt.Errorf("target path %q is not a directory", cleanedContextPath)
	}

	configPath := ""
	if cf.ConfigPath != nil && len(*cf.ConfigPath) > 0 {
		configPath = filepath.Clean(*cf.ConfigPath)
	}

	return &Options{
		group:       group,
		cluster:     cluster,
		contextPath: cleanedContextPath,
		configPath:  configPath,
		writer:      writer,
	}, nil
}

// Run execute the build command
func (o *Options) Run(ctx context.Context) error {
	o.logger = logr.FromContextOrDiscard(ctx)

	group, err := util.GroupFromConfig(o.group, o.configPath)
	if err != nil {
		return err
	}

	found := false
	str := new(strings.Builder)
	for _, cluster := range group.Clusters {
		clusterName := cluster.Name
		if o.cluster != "" && clusterName != o.cluster {
			continue
		}

		found = true
		path := filepath.Join(o.contextPath, util.ClusterPath(o.group, clusterName))

		clusterID := fmt.Sprintf("%s/%s", o.group, clusterName)
		str.WriteString("---\n")
		str.WriteString(fmt.Sprintf("### BUILD RESULTS FOR: %q ###\n", clusterID))
		if err := util.WriteKustomizationData(path, str); err != nil {
			return fmt.Errorf("building resources for %q: %w", clusterID, err)
		}
	}

	switch {
	case !found && len(o.cluster) == 0:
		return fmt.Errorf("group %q doesn't have any cluster", o.group)
	case !found && len(o.cluster) != 0:
		return fmt.Errorf("group %q doesn't have cluster %q", o.group, o.cluster)
	}

	fmt.Fprint(o.writer, str.String())
	return nil
}
