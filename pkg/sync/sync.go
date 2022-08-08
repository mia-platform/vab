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
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"

	"github.com/mia-platform/vab/internal/git"
	kustomizehelper "github.com/mia-platform/vab/internal/kustomize"
	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/mia-platform/vab/pkg/logger"
)

// Sync synchronizes modules and add-ons to the latest configuration
func Sync(logger logger.LogInterface, filesGetter git.FilesGetter, configPath string, basePath string) error {

	// ReadConfig -> get default modules and addons
	config, err := utils.ReadConfig(configPath)
	if err != nil {
		return fmt.Errorf("sync error: %w", err)
	}
	defaultModules := config.Spec.Modules
	defaultAddons := config.Spec.AddOns

	// sync default modules and add-ons in all-groups "bases" folder
	if err := SyncModules(logger, defaultModules, basePath, filesGetter); err != nil {
		return fmt.Errorf("error syncing default modules %+v: %w", defaultModules, err)
	}
	if err := SyncAddons(logger, defaultAddons, basePath, filesGetter); err != nil {
		return fmt.Errorf("error syncing default add-ons %+v: %w", defaultAddons, err)
	}
	// update the default bases in the all-groups directory
	if err := UpdateBases(utils.AllGroupsDirPath, defaultModules, defaultAddons); err != nil {
		return fmt.Errorf("error updating kustomize bases for all-groups: %w", err)
	}
	// synchronize clusters to the latest configuration
	if err := SyncClusters(&config.Spec.Groups, basePath); err != nil {
		return fmt.Errorf("error syncing clusters: %w", err)
	}

	return nil
}

// SyncModules synchronizes modules to the latest configuration
// TODO: merge duplicate functions
func SyncModules(logger logger.LogInterface, modules map[string]v1alpha1.Module, basePath string, filesGetter git.FilesGetter) error {
	for name, v := range modules {
		if v.IsDisabled() {
			continue
		}
		files, err := ClonePackages(logger, name, v, filesGetter)
		if err != nil {
			return err
		}
		modulePath := path.Join(basePath, utils.VendorsModulesPath, name)
		err = MoveToDisk(logger, files, name, modulePath)
		if err != nil {
			return err
		}
	}
	return nil
}

// SyncModules synchronizes add-ons to the latest configuration
// TODO: merge duplicate functions
func SyncAddons(logger logger.LogInterface, addons map[string]v1alpha1.AddOn, basePath string, filesGetter git.FilesGetter) error {
	for name, v := range addons {
		if v.IsDisabled() {
			continue
		}
		files, err := ClonePackages(logger, name, v, filesGetter)
		if err != nil {
			return err
		}
		addonPath := path.Join(basePath, utils.VendorsAddonsPath, name)
		err = MoveToDisk(logger, files, name, addonPath)
		if err != nil {
			return err
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
func UpdateBases(targetPath string, modules map[string]v1alpha1.Module, addons map[string]v1alpha1.AddOn) error {
	targetKustomizationPath := path.Join(targetPath, "bases")
	kustomization, err := kustomizehelper.ReadKustomization(targetKustomizationPath)
	if err != nil {
		return fmt.Errorf("error reading kustomization file for %s: %w", targetPath, err)
	}
	// if the path contains "clusters/all-groups", it is the path to the default configurations
	// otherwise, it is the path to a single cluster
	if strings.Contains(targetPath, utils.AllGroupsDirPath) {
		syncedKustomization := kustomizehelper.SyncKustomizeResources(&modules, &addons, *kustomization)
		utils.WriteKustomization(*syncedKustomization, targetKustomizationPath)
	} else {
		// case in which the cluster does not override the default configuration
		if modules == nil && addons == nil {
			// overwrite the kustomization to contain only the path to all-groups
			syncedKustomization := utils.EmptyKustomization()
			syncedKustomization.Resources = append(syncedKustomization.Resources, "../../../all-groups")
			utils.WriteKustomization(syncedKustomization, targetKustomizationPath)
		}
	}
	return nil
}

// GetClusterPath returns the path to the cluster folder, or creates it if it does not exist
// clusterName must be <group-name>/<cluster-name>
func GetClusterPath(clusterName string, basePath string) (string, error) {
	clusterPath := path.Join(basePath, utils.ClustersDirName, clusterName)
	if _, err := os.Stat(clusterPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			os.MkdirAll(clusterPath, fs.ModePerm)
		} else {
			return "", fmt.Errorf("error getting cluster path for %s: %w", clusterName, err)
		}
	}
	return clusterPath, nil
}

// SyncClusters synchronizes the clusters to the latest configuration
// Currently this function does not handle cluster overrides
func SyncClusters(groups *[]v1alpha1.Group, basePath string) error {
	for _, group := range *groups {
		for _, cluster := range group.Clusters {
			fullClusterName := path.Join(group.Name, cluster.Name)
			clusterPath, err := GetClusterPath(fullClusterName, basePath)
			if err != nil {
				return fmt.Errorf("error retrieving path for cluster %s: %w", fullClusterName, err)
			}
			// TODO: handle cluster overrides
			if err := UpdateBases(clusterPath, cluster.Modules, cluster.AddOns); err != nil {
				return fmt.Errorf("error updating kustomize bases for all-groups: %w", err)
			}
		}
	}
	return nil
}
