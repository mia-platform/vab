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

package validate

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/go-logr/logr"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/mia-platform/vab/pkg/cmd/util"
	"github.com/spf13/cobra"
)

const (
	shortCmd = "Validate configuration file"
	longCmd  = `Validate the configuration contained in the specified path.

	It returns an error if the config file is malformed or includes resources
	that do not exist in our catalogue.
`

	defaultScope = "default"
)

// Flags contains all the flags for the `validate` command. They will be converted to Options
// that contains all runtime options for the command.
type Flags struct{}

// Options have the data required to perform the validate operation
type Options struct {
	configPath string
	writer     io.Writer
	logger     logr.Logger
}

// NewCommand return the command for validating the information inserted in the configuration file
func NewCommand(cf *util.ConfigFlags) *cobra.Command {
	flags := &Flags{}
	cmd := &cobra.Command{
		Use:   "validate",
		Short: heredoc.Doc(shortCmd),
		Long:  heredoc.Doc(longCmd),

		Args: cobra.NoArgs,

		Run: func(cmd *cobra.Command, _ []string) {
			options, err := flags.ToOptions(cf, cmd.OutOrStdout())
			cobra.CheckErr(err)
			cobra.CheckErr(options.Run(cmd.Context()))
		},
	}

	return cmd
}

// ToOptions transform the command flags in command runtime arguments
func (f *Flags) ToOptions(cf *util.ConfigFlags, writer io.Writer) (*Options, error) {
	configPath := ""
	if cf.ConfigPath != nil && len(*cf.ConfigPath) > 0 {
		configPath = filepath.Clean(*cf.ConfigPath)
	}

	return &Options{
		configPath: configPath,
		writer:     writer,
	}, nil
}

// Run execute the create command
func (o *Options) Run(ctx context.Context) error {
	o.logger = logr.FromContextOrDiscard(ctx)
	code := 0

	config, err := util.ReadConfig(o.configPath)
	if err != nil {
		return fmt.Errorf("parsing configuration file: %v", err)
	}

	feedbackString := o.checkTypeMeta(&config.TypeMeta, &code)
	o.logger.V(5).Info("checking TypeMeta for config", "code", code)
	feedbackString += o.checkModules(&config.Spec.Modules, "", &code)
	o.logger.V(5).Info("checking configuration modules", "code", code)
	feedbackString += o.checkAddOns(&config.Spec.AddOns, "", &code)
	o.logger.V(5).Info("checking configuration addons", "code", code)
	feedbackString += o.checkGroups(&config.Spec.Groups, &code)
	o.logger.V(5).Info("checking configuration groups", "code", code)

	fmt.Fprint(o.writer, feedbackString)
	if code > 0 {
		return fmt.Errorf("configuration is invalid")
	}

	fmt.Fprintln(o.writer, "The configuration is valid!")
	return nil
}

// checkTypeMeta checks the file's Kind and APIVersion
func (o *Options) checkTypeMeta(config *v1alpha1.TypeMeta, code *int) string {
	outString := ""
	if config.Kind != v1alpha1.Kind {
		outString += fmt.Sprintf("[error] wrong kind: %s - expected: %s\n", config.Kind, v1alpha1.Kind)
		*code = 1
	}

	if config.APIVersion != v1alpha1.Version {
		outString += fmt.Sprintf("[error] wrong version: %s - expected: %s\n", config.APIVersion, v1alpha1.Version)
		*code = 1
	}

	return outString
}

// checkModules checks the modules listed in the config file
func (o *Options) checkModules(packages *map[string]v1alpha1.Package, scope string, code *int) string {
	if scope == "" {
		scope = defaultScope
	}

	outString := ""
	if len(*packages) == 0 {
		outString += fmt.Sprintf("[warn][%s] no module found: check the config file if this behavior is unexpected\n", scope)
		return outString
	}

	for _, pkg := range *packages {
		// TODO: add check for modules' uniqueness (only one flavor per module is allowed)
		if pkg.Disable {
			outString += fmt.Sprintf("[info][%s] disabling %s %s\n", scope, pkg.PackageType(), pkg.GetName())
			continue
		}

		if pkg.Version == "" {
			outString += fmt.Sprintf("[error][%s] missing version of %s %s\n", scope, pkg.PackageType(), pkg.GetName())
			*code = 1
		}
	}

	return outString
}

// checkAddOns checks the addons listed in the config file
func (o *Options) checkAddOns(packages *map[string]v1alpha1.Package, scope string, code *int) string {
	if scope == "" {
		scope = defaultScope
	}

	outString := ""
	if len(*packages) == 0 {
		outString += fmt.Sprintf("[warn][%s] no addon found: check the config file if this behavior is unexpected\n", scope)
		return outString
	}

	for _, pkg := range *packages {
		if pkg.Disable {
			outString += fmt.Sprintf("[info][%s] disabling %s %s\n", scope, pkg.PackageType(), pkg.GetName())
		} else if pkg.Version == "" {
			outString += fmt.Sprintf("[error][%s] missing version of %s %s\n", scope, pkg.PackageType(), pkg.GetName())
			*code = 1
		}
	}
	return outString
}

// checkGroups checks the cluster groups listed in the config file
func (o *Options) checkGroups(groups *[]v1alpha1.Group, code *int) string {
	outString := ""

	if len(*groups) == 0 {
		outString += "[warn] no group found: check the config file if this behavior is unexpected\n"
		return outString
	}

	for _, g := range *groups {
		groupName := g.Name
		if groupName == "" {
			outString += "[error] please specify a valid name for each group\n"
			groupName = "undefined"
			*code = 1
		}

		group := g
		outString += o.checkClusters(&group, groupName, code)
		o.logger.V(5).Info(fmt.Sprintf("checking group %s clusters", groupName), "", *code)
	}

	return outString
}

// checkClusters checks the clusters of a group
func (o *Options) checkClusters(group *v1alpha1.Group, groupName string, code *int) string {
	outString := ""
	if len(group.Clusters) == 0 {
		outString += fmt.Sprintf("[warn][%s] no cluster found in group: check the config file if this behavior is unexpected\n", groupName)
		return outString
	}

	for _, cluster := range group.Clusters {
		clusterName := cluster.Name
		if clusterName == "" {
			outString += fmt.Sprintf("[error][%s] missing cluster name in group: please specify a valid name for each cluster\n", groupName)
			*code = 1
			clusterName = "undefined"
		}

		clusterID := util.ClusterID(groupName, clusterName)
		if cluster.Context == "" {
			outString += fmt.Sprintf("[error][%s] missing cluster context: please specify a valid context for each cluster\n", clusterID)
			*code = 1
		}

		outString += o.checkModules(&cluster.Modules, clusterID, code)
		o.logger.V(5).Info(fmt.Sprintf("checking cluster %s modules", clusterID), "code", *code)
		outString += o.checkAddOns(&cluster.AddOns, clusterID, code)
		o.logger.V(5).Info(fmt.Sprintf("checking cluster %s addon", clusterID), "code", *code)
	}

	return outString
}
