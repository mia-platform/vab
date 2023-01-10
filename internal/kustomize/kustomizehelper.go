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

package kustomizehelper

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/api/konfig"
	kustomize "sigs.k8s.io/kustomize/api/types"
)

// SyncKustomization read the kustomization file ath path and sync its resources and components properties with
// the content of modules and addOns
func SyncAllClusterKustomization(basePath string, modules, addOns map[string]v1alpha1.Package) error {
	allGroupsPath := filepath.Join(basePath, utils.AllGroupsDirPath, utils.BasesDir)
	kustomization, err := ReadKustomization(allGroupsPath)
	if err != nil {
		return fmt.Errorf("failed to read kustomization at path %s: %w", allGroupsPath, err)
	}

	resources, err := sortedPackagesPathList(modules, allGroupsPath, filepath.Join(basePath, utils.VendorsModulesPath))
	if err != nil {
		return fmt.Errorf("failed to create a sorted list for modules: %w", err)
	}

	components, err := sortedPackagesPathList(addOns, allGroupsPath, filepath.Join(basePath, utils.VendorsAddOnsPath))
	if err != nil {
		return fmt.Errorf("failed to create a sorted list for addons: %w", err)
	}

	kustomization.Resources = resources
	kustomization.Components = components

	if err := utils.WriteKustomization(*kustomization, allGroupsPath); err != nil {
		return fmt.Errorf("failed to write kustomization at path %s: %w", allGroupsPath, err)
	}
	return nil
}

// SyncKustomization read the kustomization file ath path and sync its resources and components properties with
// the content of modules and addOns specifically for clusters
func SyncClusterKustomization(basePath, clusterPath string, modules, addOns map[string]v1alpha1.Package) error {
	kustomization, err := ReadKustomization(clusterPath)
	if err != nil {
		return fmt.Errorf("failed to read kustomization at path %s: %w", clusterPath, err)
	}

	var resources, components []string
	allGroupsDirPath := filepath.Join(basePath, utils.AllGroupsDirPath)
	if len(modules) == 0 && len(addOns) == 0 {
		allGroupsPath, err := filepath.Rel(clusterPath, allGroupsDirPath)
		if err != nil {
			return fmt.Errorf("failed to create link to all groups bases: %w", err)
		}
		resources = []string{allGroupsPath}
	} else {
		resources, err = sortedPackagesPathList(modules, clusterPath, filepath.Join(basePath, utils.VendorsModulesPath))
		if err != nil {
			return fmt.Errorf("failed to create a sorted list for modules: %w", err)
		}

		components, err = sortedPackagesPathList(addOns, clusterPath, filepath.Join(basePath, utils.VendorsAddOnsPath))
		if err != nil {
			return fmt.Errorf("failed to create a sorted list for addons: %w", err)
		}

		customResourcePath := filepath.Join(allGroupsDirPath, utils.CustomResourcesDir)
		customResourceseRelativePath, err := filepath.Rel(clusterPath, customResourcePath)
		if err != nil {
			return fmt.Errorf("failed to create link to all custom resources: %w", err)
		}
		components = append(components, customResourceseRelativePath)
	}

	kustomization.Resources = resources
	kustomization.Components = components

	if err := utils.WriteKustomization(*kustomization, clusterPath); err != nil {
		return fmt.Errorf("failed to write kustomization at path %s: %w", clusterPath, err)
	}
	return nil
}

// sortedPackagesPathList return a sorted array of path for the given packages map with paths relative to basePath
func sortedPackagesPathList(packages map[string]v1alpha1.Package, basePath, targetPath string) ([]string, error) {
	sortedList := make([]string, 0)

	for _, pkg := range packages {
		if !pkg.Disable {
			var pkgPath string
			versionedPath := pkg.GetName() + "-" + pkg.Version
			if pkg.IsModule() {
				pkgPath = filepath.Join(versionedPath, pkg.GetFlavorName())
			} else {
				pkgPath = versionedPath
			}
			modulePath, err := filepath.Rel(basePath, filepath.Join(targetPath, pkgPath))
			if err != nil {
				return nil, err
			}
			sortedList = append(sortedList, modulePath)
		}
	}

	sort.SliceStable(sortedList, func(i, j int) bool {
		return sortedList[i] < sortedList[j]
	})

	return sortedList, nil
}

// ReadKustomization reads a kustomization file given its path
func ReadKustomization(targetPath string) (*kustomize.Kustomization, error) {
	// create the path to the kustomization file if it does not exist
	// useful when creating clusters' sub-directories
	if err := utils.ValidatePath(targetPath); err != nil {
		return nil, err
	}
	// create the kustomization file if it does not exist
	kustomizationPath, err := getKustomizationFilePath(targetPath)
	if err != nil {
		return nil, fmt.Errorf("error getting kustomization file path for %s: %w", targetPath, err)
	}
	kustomization, err := os.ReadFile(kustomizationPath)
	if err != nil {
		return nil, fmt.Errorf("error reading kustomization file %s: %w", kustomizationPath, err)
	}
	output := &kustomize.Kustomization{}
	err = yaml.Unmarshal(kustomization, output)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling kustomization file %s: %w", targetPath, err)
	}

	return output, nil
}

// getKustomizationFilePath checks if a kustomization file exists and creates it if missing,
// initializing the TypeMeta fields
func getKustomizationFilePath(targetPath string) (string, error) {
	for _, validFileName := range konfig.RecognizedKustomizationFileNames() {
		kustomizationPath := filepath.Join(targetPath, validFileName)
		_, err := os.Stat(kustomizationPath)
		switch {
		case err == nil:
			return kustomizationPath, nil // if there is a match, return the valid path to the kustomization file
		case errors.Is(err, os.ErrNotExist):
			continue
		default:
			return "", fmt.Errorf("error while checking kustomization path %s: %w", kustomizationPath, err)
		}
	}
	// If the execution gets here, it means that no kustomization file with a valid name
	// was found. A new kustomization file is created (with initialized TypeMeta)
	kustomizationPath := filepath.Join(targetPath, konfig.DefaultKustomizationFileName())
	newKustomization := utils.EmptyKustomization()
	newKustomization.TypeMeta = kustomize.TypeMeta{
		Kind:       kustomize.KustomizationKind,
		APIVersion: kustomize.KustomizationVersion,
	}
	if err := utils.WriteKustomization(newKustomization, kustomizationPath); err != nil {
		return "", fmt.Errorf("error writing kustomization file %s: %w", targetPath, err)
	}
	return kustomizationPath, nil
}

// PackagesMapForPaths return the package map with key the disk path for kustomize
func PackagesMapForPaths(packages map[string]v1alpha1.Package) map[string]v1alpha1.Package {
	pathsMap := make(map[string]v1alpha1.Package, len(packages))

	for _, pkg := range packages {
		newKey := pkg.GetName() + "-" + pkg.Version
		pathsMap[newKey] = pkg
	}

	return pathsMap
}
