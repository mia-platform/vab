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
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/MakeNowJust/heredoc/v2"
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

	dryRunDefaultValue = false
	dryRunFlagName     = "dry-run"
	dryRunUsage        = "if true no files will be downloaded"
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

		Run: func(_ *cobra.Command, args []string) {
			options, err := flags.ToOptions(cf, args)
			cobra.CheckErr(err)
			cobra.CheckErr(options.Run())
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

	var err error
	var cleanedContextPath string
	if cleanedContextPath, err = filepath.Abs(args[0]); err != nil {
		return nil, err
	}

	var contextInfo fs.FileInfo
	if contextInfo, err = os.Stat(cleanedContextPath); err != nil {
		return nil, fmt.Errorf("error locating files: %w", err)
	}

	if !contextInfo.IsDir() {
		return nil, fmt.Errorf("the target path %q is not a directory", cleanedContextPath)
	}

	return &Options{
		contextPath: cleanedContextPath,
		configPath:  configPath,
		dryRun:      f.dryRun,
	}, nil
}

// Run execute the create command
func (o *Options) Run() error {
	config, err := util.ReadConfig(o.configPath)
	if err != nil {
		return fmt.Errorf("cannot read config file: %w", err)
	}

	if err := util.SyncDirectories(config.Spec, o.contextPath); err != nil {
		return err
	}

	return o.downloadPackages(config, git.RealFilesGetter{})
}

func (o *Options) downloadPackages(config *v1alpha1.ClustersConfiguration, filesGetter git.FilesGetter) error {
	if o.dryRun {
		return nil
	}

	if err := os.RemoveAll(filepath.Join(o.contextPath, util.VendoredModulePath(""))); err != nil {
		return fmt.Errorf("failed to remove vendors folder for modules: %w", err)
	}
	if err := os.RemoveAll(filepath.Join(o.contextPath, util.VendoredAddOnPath(""))); err != nil {
		return fmt.Errorf("failed to remove vendors folder for add-ons: %w", err)
	}

	mergedPackages := make(map[string]v1alpha1.Package)
	for _, pkg := range config.Spec.Modules {
		if !pkg.Disable {
			mergedPackages[pkg.GetName()+pkg.GetFlavorName()+"_"+pkg.Version] = pkg
		}
	}
	for _, pkg := range config.Spec.AddOns {
		if !pkg.Disable {
			mergedPackages[pkg.GetName()+"_"+pkg.Version] = pkg
		}
	}

	for _, group := range config.Spec.Groups {
		for _, cluster := range group.Clusters {
			for _, pkg := range cluster.Modules {
				if !pkg.Disable {
					mergedPackages[pkg.GetName()+pkg.GetFlavorName()+"_"+pkg.Version] = pkg
				}
			}
			for _, pkg := range cluster.AddOns {
				if !pkg.Disable {
					mergedPackages[pkg.GetName()+"_"+pkg.Version] = pkg
				}
			}
		}
	}

	return clonePackagesLocally(mergedPackages, o.contextPath, filesGetter)
}

// clonePackagesLocally download packages using filesGetter
func clonePackagesLocally(packages map[string]v1alpha1.Package, path string, filesGetter git.FilesGetter) error {
	for _, pkg := range packages {
		files, err := ClonePackages(pkg, filesGetter)
		if err != nil {
			return fmt.Errorf("error cloning packages for %s %s: %w", pkg.PackageType(), pkg.GetName(), err)
		}

		pkgName := pkg.GetName() + "-" + pkg.Version
		var pkgPath string
		if pkg.IsModule() {
			pkgPath = util.VendoredModulePath(pkgName)
		} else {
			pkgPath = util.VendoredAddOnPath(pkgName)
		}

		if err := MoveToDisk(files, pkg.GetName(), filepath.Join(path, pkgPath)); err != nil {
			return fmt.Errorf("error moving packages to disk for %s %s: %w", pkg.PackageType(), pkg.GetName(), err)
		}
	}
	return nil
}

// ClonePackages clones and writes package repos to disk
func ClonePackages(pkg v1alpha1.Package, filesGetter git.FilesGetter) ([]*git.File, error) {
	files, err := git.GetFilesForPackage(filesGetter, pkg)

	if err != nil {
		return nil, fmt.Errorf("error getting files for module %s: %w", pkg.GetName(), err)
	}

	return files, nil
}

// MoveToDisk moves the cloned packages from memory to disk
func MoveToDisk(files []*git.File, packageName string, targetPath string) error {
	if err := WritePkgToDir(files, targetPath); err != nil {
		return fmt.Errorf("error while writing package %s on disk: %w", packageName, err)
	}

	return nil
}

// WritePkgToDir writes the files in memory to the target path on disk
func WritePkgToDir(files []*git.File, targetPath string) error {
	for _, gitFile := range files {
		err := os.MkdirAll(filepath.Dir(filepath.Join(targetPath, gitFile.FilePath())), os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory: %s : %w", filepath.Dir(gitFile.FilePath()), err)
		}

		err = gitFile.Open()
		if err != nil {
			return fmt.Errorf("error opening file: %s : %w", gitFile.String(), err)
		}
		outFile, err := os.Create(filepath.Join(targetPath, gitFile.FilePath()))
		if err != nil {
			return fmt.Errorf("error opening file: %s : %w", filepath.Join(targetPath, gitFile.FilePath()), err)
		}

		r := bufio.NewReader(gitFile)
		w := bufio.NewWriter(outFile)

		_, err = r.WriteTo(w)
		if err != nil {
			return fmt.Errorf("error writing: %s : %w", outFile.Name(), err)
		}

		err = gitFile.Close()
		if err != nil {
			return fmt.Errorf("error closing: %s : %w", gitFile.String(), err)
		}

		err = outFile.Close()
		if err != nil {
			return fmt.Errorf("error closing: %s : %w", outFile.Name(), err)
		}
	}
	return nil
}
