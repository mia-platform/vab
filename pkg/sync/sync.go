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
	"fmt"
	"os"
	"path/filepath"

	"github.com/mia-platform/vab/internal/git"
	kustomizehelper "github.com/mia-platform/vab/internal/kustomize"
	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/mia-platform/vab/pkg/logger"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// Sync synchronizes modules and add-ons to the latest configuration
func Sync(logger logger.LogInterface, filesGetter git.FilesGetter, configPath string, basePath string, dryRun bool) error {
	// ReadConfig -> get default modules and addons
	config, err := utils.ReadConfig(configPath)
	if err != nil {
		return fmt.Errorf("sync error: %w", err)
	}

	if err := syncAllGroups(config, basePath); err != nil {
		return err
	}

	if err := syncClusters(logger, config, basePath); err != nil {
		return err
	}

	if err := downloadPackages(logger, config, basePath, filesGetter, dryRun); err != nil {
		return err
	}
	return nil
}

func syncAllGroups(config *v1alpha1.ClustersConfiguration, basePath string) error {
	if err := kustomizehelper.SyncAllClusterKustomization(basePath, config.Spec.Modules, config.Spec.AddOns); err != nil {
		return fmt.Errorf("error updating all-groups kustomize file: %w", err)
	}
	return nil
}

func syncClusters(logger logger.LogInterface, config *v1alpha1.ClustersConfiguration, basePath string) error {
	groups := config.Spec.Groups
	modules := config.Spec.Modules
	addons := config.Spec.AddOns

	for _, group := range groups {
		for _, cluster := range group.Clusters {
			fullClusterName := filepath.Join(group.Name, cluster.Name)
			clusterPath, err := checkClusterPath(fullClusterName, basePath)
			if err != nil {
				return fmt.Errorf("error retrieving path for cluster %s: %w", fullClusterName, err)
			}

			var clusterModules, clusterAddOns map[string]v1alpha1.Package
			if len(cluster.Modules) == 0 && len(cluster.AddOns) == 0 {
				clusterModules = make(map[string]v1alpha1.Package, 0)
				clusterAddOns = make(map[string]v1alpha1.Package, 0)
			} else {
				clusterModules = mergePackages(logger, modules, cluster.Modules)
				clusterAddOns = mergePackages(logger, addons, cluster.AddOns)
			}

			clusterBasePath := filepath.Join(clusterPath, utils.BasesDir)
			if err := kustomizehelper.SyncClusterKustomization(basePath, clusterBasePath, clusterModules, clusterAddOns); err != nil {
				return fmt.Errorf("error updating %s cluster kustomize file: %w", fullClusterName, err)
			}
		}
	}
	return nil
}

func downloadPackages(logger logger.LogInterface, config *v1alpha1.ClustersConfiguration, path string, filesGetter git.FilesGetter, dryRun bool) error {
	if dryRun {
		return nil
	}

	if err := os.RemoveAll(filepath.Join(path, utils.VendorsModulesPath)); err != nil {
		return fmt.Errorf("failed to remove vendors folder for modules: %w", err)
	}
	if err := os.RemoveAll(filepath.Join(path, utils.VendorsAddOnsPath)); err != nil {
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

	return clonePackagesLocally(logger, mergedPackages, path, filesGetter)
}

// clonePackagesLocally download packages using filesGetter
func clonePackagesLocally(logger logger.LogInterface, packages map[string]v1alpha1.Package, path string, filesGetter git.FilesGetter) error {
	for _, pkg := range packages {
		files, err := ClonePackages(logger, pkg, filesGetter)
		if err != nil {
			return fmt.Errorf("error cloning packages for %s %s: %w", pkg.PackageType(), pkg.GetName(), err)
		}

		var vendorsPath string
		if pkg.IsModule() {
			vendorsPath = utils.VendorsModulesPath
		} else {
			vendorsPath = utils.VendorsAddOnsPath
		}

		pkgPath := filepath.Join(path, vendorsPath, pkg.GetName()+"-"+pkg.Version)
		logger.V(10).Writef("disk path for package %s: %s", pkg.GetName(), pkgPath)
		if err := MoveToDisk(logger, files, pkg.GetName(), pkgPath); err != nil {
			return fmt.Errorf("error moving packages to disk for %s %s: %w", pkg.PackageType(), pkg.GetName(), err)
		}
	}
	return nil
}

// ClonePackages clones and writes package repos to disk
func ClonePackages(logger logger.LogInterface, pkg v1alpha1.Package, filesGetter git.FilesGetter) ([]*git.File, error) {
	files, err := git.GetFilesForPackage(logger, filesGetter, pkg)

	if err != nil {
		return nil, fmt.Errorf("error getting files for module %s: %w", pkg.GetName(), err)
	}

	return files, nil
}

// MoveToDisk moves the cloned packages from memory to disk
func MoveToDisk(logger logger.LogInterface, files []*git.File, packageName string, targetPath string) error {
	logger.V(10).Writef("Path for module %s: %s", packageName, targetPath)

	if err := WritePkgToDir(files, targetPath); err != nil {
		return fmt.Errorf("error while writing package %s on disk: %w", packageName, err)
	}

	return nil
}

// checkClusterPath returns the path to the cluster folder, or creates it if it does not exist;
// it also initializes the cluster kustomization file for the user
// clusterName must be <group-name>/<cluster-name>
func checkClusterPath(clusterName string, basePath string) (string, error) {
	clusterPath := filepath.Join(basePath, utils.ClustersDirName, clusterName)
	if err := utils.ValidatePath(clusterPath); err != nil {
		return "", fmt.Errorf("error validating cluster path %s: %w", clusterPath, err)
	}
	// initialize cluster kustomization if not present, importing the "bases" directory by default
	clusterKustomization, err := kustomizehelper.ReadKustomization(clusterPath)
	if err != nil {
		return "", fmt.Errorf("error getting kustomization for cluster %s: %w", clusterName, err)
	}
	if !slices.Contains(clusterKustomization.Resources, utils.BasesDir) {
		clusterKustomization.Resources = append([]string{utils.BasesDir}, clusterKustomization.Resources...)
		if err := utils.WriteKustomization(*clusterKustomization, clusterPath); err != nil {
			return "", fmt.Errorf("error writing kustomization file for cluster %s: %w", clusterName, err)
		}
	}
	return clusterPath, nil
}

// mergePackages return a map of merged packages excluding disabled ones, if second has no elements return nil
func mergePackages(logger logger.LogInterface, first, second map[string]v1alpha1.Package) map[string]v1alpha1.Package {
	mergedMap := make(map[string]v1alpha1.Package)
	maps.Copy(mergedMap, first)
	for name, pkg := range second {
		// if the current package is disabled remove it from the map
		if pkg.Disable {
			logger.V(10).Writef("Disable %s %s for cluster", pkg.PackageType(), pkg.GetName())
			delete(mergedMap, name)
		} else {
			logger.V(10).Writef("Override %s %s for cluser", pkg.PackageType(), pkg.GetName())
			mergedMap[name] = pkg
		}
	}

	// return the list of packages with the on disk path as key
	return mergedMap
}
