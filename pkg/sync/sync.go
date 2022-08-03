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
func Sync(logger logger.LogInterface, configPath string) error { // add basePath as parameter

	// ReadConfig -> get modules and addons
	config, err := utils.ReadConfig(configPath)

	if err != nil {
		return fmt.Errorf("sync error: %w", err)
	}

	// loop on modules and addons
	// GetFilesForPackage + WritePkgToDir
	defaultModules := config.Spec.Modules
	defaultAddons := config.Spec.AddOns

	if err := SyncModules(logger, defaultModules); err != nil {
		return fmt.Errorf("error syncing default modules %+v: %w", defaultModules, err)
	}

	if err := SyncAddons(logger, defaultAddons); err != nil {
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

func SyncModules(logger logger.LogInterface, modules map[string]v1alpha1.Module) error {
	for name, v := range modules {
		if v.IsDisabled() {
			continue
		}
		err := SyncPackages(logger, name, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func SyncAddons(logger logger.LogInterface, addons map[string]v1alpha1.AddOn) error {
	for name, v := range addons {
		if v.IsDisabled() {
			continue
		}
		err := SyncPackages(logger, name, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// SyncPackages clones and writes package repos to disk
func SyncPackages(logger logger.LogInterface, packageName string, pkg v1alpha1.Package) error {

	files, err := git.GetFilesForPackage(logger, packageName, pkg)

	if err != nil {
		return fmt.Errorf("error getting files for module %s: %w", packageName, err)
	}

	targetPath := path.Join(utils.VendorsModulesPath, packageName)
	logger.V(10).Writef("Path for module %s: %s", packageName, targetPath)

	if err := WritePkgToDir(files, targetPath); err != nil {
		return fmt.Errorf("error while writing package %s on disk: %w", packageName, err)
	}

	return nil

}
