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

	"github.com/mia-platform/vab/internal/git"
	kustomizehelper "github.com/mia-platform/vab/internal/kustomize"
	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/mia-platform/vab/pkg/logger"
)

// Sync synchronizes modules and add-ons to the latest configuration
func Sync(logger logger.LogInterface, configPath string, filesGetter git.FilesGetter) error { // add basePath as parameter

	// ReadConfig -> get modules and addons
	config, err := utils.ReadConfig(configPath)

	if err != nil {
		return fmt.Errorf("sync error: %w", err)
	}

	// loop on modules and addons
	// GetFilesForPackage + WritePkgToDir
	defaultModules := config.Spec.Modules
	defaultAddons := config.Spec.AddOns

	if err := SyncModules(logger, defaultModules, filesGetter); err != nil {
		return fmt.Errorf("error syncing default modules %+v: %w", defaultModules, err)
	}

	if err := SyncAddons(logger, defaultAddons, filesGetter); err != nil {
		return fmt.Errorf("error syncing default add-ons %+v: %w", defaultAddons, err)
	}

	// loop on all clusters
	// SyncPackages + SyncResources + WriteKustomization
	// TODO: optimize loops with concurrency
	for _, group := range config.Spec.Groups {
		for _, cluster := range group.Clusters {
			// TODO: sync modules and add-ons patches for each cluster
			// if err := SyncPackages(logger, cluster.Modules); err != nil {
			// 	return fmt.Errorf("error syncing modules for cluster %s, %+v: %w", cluster.Name, cluster.Modules, err)
			// }
			// if err := SyncPackages(logger, cluster.AddOns); err != nil {
			// 	return fmt.Errorf("error syncing add-ons for cluster %s, %+v: %w", cluster.Name, cluster.AddOns, err)
			// }

			kustomizationPath := path.Join(utils.ClustersDirName, group.Name, cluster.Name, utils.KustomizationFileName)
			// TODO: check if file exists
			kustomization, err := kustomizehelper.ReadKustomization(kustomizationPath)
			if err != nil {
				return fmt.Errorf("error reading kustomization file for %s/%s: %w", group.Name, cluster.Name, err)
			}
			syncedKustomization := kustomizehelper.SyncKustomizeResources(&defaultModules, &defaultAddons, *kustomization)
			utils.WriteKustomization(syncedKustomization, kustomizationPath)
		}
	}

	return nil
}

// SyncModules synchronizes modules to the latest configuration
// TODO: merge duplicate functions
func SyncModules(logger logger.LogInterface, modules map[string]v1alpha1.Module, filesGetter git.FilesGetter) error {
	for name, v := range modules {
		if v.IsDisabled() {
			continue
		}
		files, err := ClonePackages(logger, name, v, filesGetter)
		if err != nil {
			return err
		}
		err = MoveToDisk(logger, files, name)
		if err != nil {
			return err
		}
	}
	return nil
}

// SyncModules synchronizes add-ons to the latest configuration
// TODO: merge duplicate functions
func SyncAddons(logger logger.LogInterface, addons map[string]v1alpha1.AddOn, filesGetter git.FilesGetter) error {
	for name, v := range addons {
		if v.IsDisabled() {
			continue
		}
		files, err := ClonePackages(logger, name, v, filesGetter)
		if err != nil {
			return err
		}
		err = MoveToDisk(logger, files, name)
		if err != nil {
			return err
		}
	}
	return nil
}

// ClonePackages clones and writes package repos to disk
func ClonePackages(logger logger.LogInterface, packageName string, pkg v1alpha1.Package, filesGetter git.FilesGetter) ([]*git.File, error) {

	files, err := filesGetter(logger, packageName, pkg)

	if err != nil {
		return nil, fmt.Errorf("error getting files for module %s: %w", packageName, err)
	}

	return files, nil

}

// MoveToDisk moves the cloned packages from memory to disk
func MoveToDisk(logger logger.LogInterface, files []*git.File, packageName string) error {

	targetPath := path.Join(utils.VendorsModulesPath, packageName)
	logger.V(10).Writef("Path for module %s: %s", packageName, targetPath)

	if err := WritePkgToDir(files, targetPath); err != nil {
		return fmt.Errorf("error while writing package %s on disk: %w", packageName, err)
	}

	return nil

}
