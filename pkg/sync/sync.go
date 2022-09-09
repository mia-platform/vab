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

package sync

import (
	"fmt"
	"path"
	"strings"

	"github.com/mia-platform/vab/internal/git"
	kustomizehelper "github.com/mia-platform/vab/internal/kustomize"
	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/mia-platform/vab/pkg/logger"
	"golang.org/x/exp/slices"
	"sigs.k8s.io/kustomize/api/types"
)

// Sync synchronizes modules and add-ons to the latest configuration
func Sync(logger logger.LogInterface, filesGetter git.FilesGetter, configPath string, basePath string, dryRun bool) error {
	// ReadConfig -> get default modules and addons
	config, err := utils.ReadConfig(configPath)
	if err != nil {
		return fmt.Errorf("sync error: %w", err)
	}
	defaultModules := kustomizehelper.CompleteModuleNames(config.Spec.Modules)
	defaultAddOns := kustomizehelper.CompleteAddOnNames(config.Spec.AddOns)

	// update the default bases in the all-groups directory
	if err := UpdateBases(logger, filesGetter, basePath, path.Join(basePath, utils.AllGroupsDirPath), defaultModules, defaultAddOns, config, dryRun); err != nil {
		return fmt.Errorf("error updating kustomize bases for all-groups: %w", err)
	}
	// synchronize clusters to the latest configuration
	if err := UpdateClusters(logger, filesGetter, config, basePath, dryRun); err != nil {
		return fmt.Errorf("error syncing clusters: %w", err)
	}

	return nil
}

// UpdateModules synchronizes modules to the latest configuration
// TODO: merge duplicate functions
func UpdateModules(logger logger.LogInterface, modules map[string]v1alpha1.Module, basePath string, filesGetter git.FilesGetter) error {
	for name, v := range modules {
		if v.IsDisabled() {
			continue
		}
		files, err := ClonePackages(logger, name, v, filesGetter)
		if err != nil {
			return fmt.Errorf("error cloning packages for module %s: %w", name, err)
		}
		modulePath := path.Join(basePath, utils.VendorsModulesPath, name)
		if err := MoveToDisk(logger, files, name, modulePath); err != nil {
			return fmt.Errorf("error moving packages to disk for module %s: %w", name, err)
		}
	}
	return nil
}

// SyncModules synchronizes add-ons to the latest configuration
// TODO: merge duplicate functions
func UpdateAddOns(logger logger.LogInterface, addons map[string]v1alpha1.AddOn, basePath string, filesGetter git.FilesGetter) error {
	for name, v := range addons {
		if v.IsDisabled() {
			continue
		}
		files, err := ClonePackages(logger, name, v, filesGetter)
		if err != nil {
			return fmt.Errorf("error cloning packages for add-on %s: %w", name, err)
		}
		addonPath := path.Join(basePath, utils.VendorsAddOnsPath, name)
		if err := MoveToDisk(logger, files, name, addonPath); err != nil {
			return fmt.Errorf("error moving packages to disk for add-on %s: %w", name, err)
		}
	}
	return nil
}

// ClonePackages clones and writes package repos to disk
func ClonePackages(logger logger.LogInterface, packageName string, pkg v1alpha1.Package, filesGetter git.FilesGetter) ([]*git.File, error) {
	files, err := git.GetFilesForPackage(logger, filesGetter, packageName, pkg)

	if err != nil {
		return nil, fmt.Errorf("error getting files for module %s: %w", packageName, err)
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

// UpdateBases updates the kustomize bases in the target path
func UpdateBases(logger logger.LogInterface, filesGetter git.FilesGetter, basePath string, targetPath string, modules map[string]v1alpha1.Module, addons map[string]v1alpha1.AddOn, config *v1alpha1.ClustersConfiguration, dryRun bool) error {
	targetKustomizationPath := path.Join(targetPath, utils.BasesDir)
	kustomization, err := kustomizehelper.ReadKustomization(targetKustomizationPath)
	if err != nil {
		return fmt.Errorf("error reading kustomization file for %s: %w", targetPath, err)
	}
	var syncedKustomization types.Kustomization
	// if modules and add-ons are nil and the path does not contains "clusters/all-groups",
	// the target is a cluster that does not override the default configuration
	if modules == nil && addons == nil && !strings.Contains(targetPath, utils.AllGroupsDirPath) {
		// overwrite the kustomization to contain only the path to all-groups
		syncedKustomization = utils.EmptyKustomization()
		syncedKustomization.Resources = append(syncedKustomization.Resources, "../../../all-groups")
		// in any other case we need to explicitly list the resources,
		// whether it is the all-groups configuration or a single cluster override
	} else {
		// NB: one between the lists of modules and add-ons overrides may still be nil.
		// If that's the case, assign the default list of relative packages to
		// overwrite the corresponding kustomization
		if modules == nil {
			modules = kustomizehelper.CompleteModuleNames(config.Spec.DeepCopy().Modules)
		}
		if addons == nil {
			addons = kustomizehelper.CompleteAddOnNames(config.Spec.DeepCopy().AddOns)
		}
		syncedKustomization = *kustomizehelper.SyncKustomizeResources(&modules, &addons, *kustomization, targetPath)
	}
	// if dryRun is true, skip modules and addons update (ClonePackages + MoveToDisk)
	if dryRun {
		logger.V(0).Writef("[warn] Dry-run mode enabled, skipping package cloning for %s. The following packages may be missing in the vendors directory.\nSkipped modules: %+v\nSkipped add-ons: %+v",
			targetPath, modules, addons)
	} else {
		if err := UpdateModules(logger, modules, basePath, filesGetter); err != nil {
			return fmt.Errorf("error syncing default modules %+v: %w", modules, err)
		}
		if err := UpdateAddOns(logger, addons, basePath, filesGetter); err != nil {
			return fmt.Errorf("error syncing default add-ons %+v: %w", addons, err)
		}
	}
	if err := utils.WriteKustomization(syncedKustomization, targetKustomizationPath); err != nil {
		return fmt.Errorf("error writing kustomization in path %s: %w", targetKustomizationPath, err)
	}
	return nil
}

// CheckClusterPath returns the path to the cluster folder, or creates it if it does not exist;
// it also initializes the cluster kustomization file for the user
// clusterName must be <group-name>/<cluster-name>
func CheckClusterPath(clusterName string, basePath string) (string, error) {
	clusterPath := path.Join(basePath, utils.ClustersDirName, clusterName)
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

// UpdateClusters synchronizes the clusters to the latest configuration
func UpdateClusters(logger logger.LogInterface, filesGetter git.FilesGetter, config *v1alpha1.ClustersConfiguration, basePath string, dryRun bool) error {
	groups := config.Spec.Groups
	for _, group := range groups {
		for _, cluster := range group.Clusters {
			fullClusterName := path.Join(group.Name, cluster.Name)
			clusterPath, err := CheckClusterPath(fullClusterName, basePath)
			if err != nil {
				return fmt.Errorf("error retrieving path for cluster %s: %w", fullClusterName, err)
			}
			syncedModules := UpdateClusterModules(cluster.Modules, config.Spec.DeepCopy().Modules)
			syncedAddOns := UpdateClusterAddOns(cluster.AddOns, config.Spec.DeepCopy().AddOns)
			if err := UpdateBases(logger, filesGetter, basePath, clusterPath, syncedModules, syncedAddOns, config, dryRun); err != nil {
				return fmt.Errorf("error updating kustomize bases for cluster %s: %w", fullClusterName, err)
			}
		}
	}
	return nil
}

// UpdateClusterModules returns the complete map of modules of the given cluster
// TODO: refactor (duplicate of UpdateClusterAddOns)
func UpdateClusterModules(modulesOverrides map[string]v1alpha1.Module, defaultModules map[string]v1alpha1.Module) map[string]v1alpha1.Module {
	// if the cluster does not provide any override, return nil to apply the default configuration
	if len(modulesOverrides) == 0 {
		return nil
	}
	// loop over the cluster modules
	for name, clusterModule := range modulesOverrides {
		// if the cluster module exists in the default modules dictionary
		// and the disable flag is set, delete the it from the list
		if _, exists := defaultModules[name]; exists && clusterModule.Disable {
			delete(defaultModules, name)
		} else {
			// this directive reassigns the clusterModule to its corresponding element
			// in defaultModules if it exists and is enabled, or adds the module to the
			// list if it is not present
			defaultModules[name] = clusterModule
		}
	}
	// return the updated copy of defaultModules
	return kustomizehelper.CompleteModuleNames(defaultModules)
}

// UpdateClusterAddOns returns the complete map of add-ons of the given cluster
// TODO: refactor (duplicate of UpdateClusterModules)
func UpdateClusterAddOns(addonsOverrides map[string]v1alpha1.AddOn, defaultAddOns map[string]v1alpha1.AddOn) map[string]v1alpha1.AddOn {
	// if the cluster does not provide any override, return nil to apply the default configuration
	if len(addonsOverrides) == 0 {
		return nil
	}
	// loop over the cluster add-ons
	for name, clusterAddOn := range addonsOverrides {
		// if the cluster add-on exists in the default add-ons dictionary
		// and the disable flag is set, delete the it from the list
		if _, exists := defaultAddOns[name]; exists && clusterAddOn.Disable {
			delete(defaultAddOns, name)
		} else {
			// this directive reassigns the clusterAddOn to its corresponding element
			// in defaultAddOns if it exists and is enabled, or adds the add-on to the
			// list if it is not present
			defaultAddOns[name] = clusterAddOn
		}
	}
	// return the updated copy of defaultAddOns
	return kustomizehelper.CompleteAddOnNames(defaultAddOns)
}
