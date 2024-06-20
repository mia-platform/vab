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

package sync

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/go-logr/logr"
	"github.com/mia-platform/vab/internal/git"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/mia-platform/vab/pkg/cmd/util"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	shortCmd = "Fetch module and addon for current configurations"
	longCmd  = `Fetches new and updated vendor versions, and updates the clusters configuration
	locally to the latest changes of the configuration file.

	After the execution, the vendors folder will include the new and updated
	modules/add-ons (if not already present), and the directory structure
	inside the clusters folder will be updated according to the current configuration.`
	cmdUsage = "sync CONTEXT"

	dryRunDefaultValue = true
	dryRunFlagName     = "download-packages"
	dryRunUsage        = "if false packages files will not be downloaded"
)

// Flags contains all the flags for the `sync` command. They will be converted to Options
// that contains all runtime options for the command.
type Flags struct {
	dryRun bool
}

// AddFlags set the connection between Flags property to command line flags
func (f *Flags) AddFlags(flags *pflag.FlagSet) {
	flags.BoolVar(&f.dryRun, dryRunFlagName, dryRunDefaultValue, heredoc.Doc(dryRunUsage))
}

// Options have the data required to perform the sync operation
type Options struct {
	contextPath string
	configPath  string
	dryRun      bool
	filesGetter *git.FilesGetter
	logger      logr.Logger
}

// NewCommand return the command for creating a new configuration file and basic folder structures
func NewCommand(cf *util.ConfigFlags) *cobra.Command {
	flags := &Flags{}
	cmd := &cobra.Command{
		Use:     cmdUsage,
		Aliases: []string{"init"},
		Short:   heredoc.Doc(shortCmd),
		Long:    heredoc.Doc(longCmd),

		Args: cobra.ExactArgs(1),

		Run: func(cmd *cobra.Command, args []string) {
			options, err := flags.ToOptions(cf, args)
			cobra.CheckErr(err)
			cobra.CheckErr(options.Run(cmd.Context()))
		},
	}

	return cmd
}

// ToOptions transform the command flags in command runtime arguments
func (f *Flags) ToOptions(cf *util.ConfigFlags, args []string) (*Options, error) {
	configPath := ""
	if cf.ConfigPath != nil && len(*cf.ConfigPath) > 0 {
		configPath = filepath.Clean(*cf.ConfigPath)
	}

	contextPath, err := util.ValidateContextPath(args[0])
	if err != nil {
		return nil, err
	}

	return &Options{
		contextPath: contextPath,
		configPath:  configPath,
		dryRun:      f.dryRun,
		filesGetter: git.NewFilesGetter(),
	}, nil
}

// Run execute the create command
func (o *Options) Run(ctx context.Context) error {
	o.logger = logr.FromContextOrDiscard(ctx)

	config, err := util.ReadConfig(o.configPath)
	if err != nil {
		return fmt.Errorf("reading config file: %w", err)
	}

	if err := util.SyncDirectories(config.Spec, o.contextPath); err != nil {
		return err
	}

	return o.downloadPackages(config)
}

func (o *Options) downloadPackages(config *v1alpha1.ClustersConfiguration) error {
	if err := os.RemoveAll(filepath.Join(o.contextPath, util.VendoredModulePath(""))); err != nil {
		return fmt.Errorf("removing vendors folder for modules: %w", err)
	}
	if err := os.RemoveAll(filepath.Join(o.contextPath, util.VendoredAddOnPath(""))); err != nil {
		return fmt.Errorf("removing vendors folder for add-ons: %w", err)
	}

	if o.dryRun {
		return nil
	}

	mergedPackages := make(map[string]v1alpha1.Package)

	addPackages := func(packages map[string]v1alpha1.Package) {
		for _, pkg := range packages {
			if !pkg.Disable {
				mergedPackages[pkg.GetName()+pkg.GetFlavorName()+"_"+pkg.Version] = pkg
			}
		}
	}

	addPackages(config.Spec.Modules)
	addPackages(config.Spec.AddOns)

	for _, group := range config.Spec.Groups {
		for _, cluster := range group.Clusters {
			addPackages(cluster.Modules)
			addPackages(cluster.AddOns)
		}
	}

	return o.clonePackagesLocally(mergedPackages, o.contextPath, o.filesGetter)
}

// clonePackagesLocally download packages using filesGetter
func (o *Options) clonePackagesLocally(packages map[string]v1alpha1.Package, path string, filesGetter *git.FilesGetter) error {
	for _, pkg := range packages {
		files, err := filesGetter.GetFilesForPackage(pkg)
		if err != nil {
			return fmt.Errorf("cloning packages for %s %s: %w", pkg.PackageType(), pkg.GetName(), err)
		}

		pkgName := pkg.GetName() + "-" + pkg.Version
		var pkgPath string
		if pkg.IsModule() {
			pkgPath = util.VendoredModulePath(pkgName)
		} else {
			pkgPath = util.VendoredAddOnPath(pkgName)
		}

		if err := o.writePackageToDisk(files, filepath.Join(path, pkgPath)); err != nil {
			return fmt.Errorf("writing %s %s on disk: %w", pkg.PackageType(), pkg.GetName(), err)
		}
	}
	return nil
}

// writePackageToDisk writes the files in memory to the target path on disk
func (o *Options) writePackageToDisk(files []*git.File, targetPath string) error {
	for _, gitFile := range files {
		return gitFile.WriteContent(targetPath)
	}

	return nil
}
