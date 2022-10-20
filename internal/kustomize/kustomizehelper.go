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

package kustomizehelper

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/api/konfig"
	kustomize "sigs.k8s.io/kustomize/api/types"
)

// SyncKustomizeResources updates the clusters' kustomization resources to the latest sync
func SyncKustomizeResources(modules *map[string]v1alpha1.Package, addons *map[string]v1alpha1.Package, k kustomize.Kustomization, targetPath string) *kustomize.Kustomization {
	modulesList := getSortedPackagesList(modules, targetPath)
	addonsList := getSortedPackagesList(addons, targetPath)

	// If the file already includes a non-empty list of resources, this function
	// collects all the custom modules that were added manually by the user
	// (i.e. all those modules that are not present in the vendors folder, thus
	// without "vendors" in their path). Then, the custom modules are appended
	// to the updated modules list that will substitute the already existing one.
	if k.Resources != nil {
		// since we are overriding the vendors we need to drop the references to the
		// vendors contained in all-groups/bases
		// if the only resource in the kustomization is the whole all-groups directory,
		// change it to point to the custom-resources directory only
		if len(k.Resources) == 1 && k.Resources[0] == "../../../all-groups" {
			k.Resources[0] = "../../../all-groups/custom-resources"
		}
		for _, r := range k.Resources {
			if !strings.Contains(r, "vendors/") {
				modulesList = append(modulesList, r)
			}
		}
		for _, r := range k.Components {
			if !strings.Contains(r, "vendors/") {
				addonsList = append(addonsList, r)
			}
		}
	}

	k.Resources = modulesList
	k.Components = addonsList

	return &k
}

// getSortedPackagesList returns the list of modules names in lexicographic order.
func getSortedPackagesList(packages *map[string]v1alpha1.Package, targetPath string) []string {
	sordtedList := make([]string, 0, len(*packages))

	for name, pkg := range *packages {
		if !pkg.Disable {
			var pkgPath string
			if pkg.IsModule() {
				pkgPath = filepath.Join(name, pkg.GetFlavorName())
			} else {
				pkgPath = name
			}
			sordtedList = append(sordtedList, pkgPath)
		}
	}

	sort.SliceStable(sordtedList, func(i, j int) bool {
		return sordtedList[i] < sordtedList[j]
	})

	return *fixResourcesPath(sordtedList, targetPath, true)
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

// fixResourcesPath returns the list of resources with the actual path
func fixResourcesPath(resourcesList []string, targetPath string, isModulesList bool) *[]string {
	fixedResourcesList := make([]string, 0, len(resourcesList))
	for _, res := range resourcesList {
		if isModulesList {
			fixedResourcesList = append(fixedResourcesList, getVendorPackageRelativePath(targetPath, path.Join(utils.VendorsModulesPath, res)))
		} else {
			fixedResourcesList = append(fixedResourcesList, getVendorPackageRelativePath(targetPath, path.Join(utils.VendorsAddOnsPath, res)))
		}
	}
	return &fixedResourcesList
}

// getKustomizationFilePath checks if a kustomization file exists and creates it if missing,
// initializing the TypeMeta fields
func getKustomizationFilePath(targetPath string) (string, error) {
	for _, validFileName := range konfig.RecognizedKustomizationFileNames() {
		kustomizationPath := path.Join(targetPath, validFileName)
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
	kustomizationPath := path.Join(targetPath, konfig.DefaultKustomizationFileName())
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

// getVendorPackageRelativePath returns the relative path to the vendor package
// for the Kustomization file
func getVendorPackageRelativePath(targetPath string, pkgPath string) string {
	var vendorPackageRelativePath string
	if strings.Contains(targetPath, utils.AllGroupsDirPath) {
		vendorPackageRelativePath = path.Join("..", "..", "..", pkgPath)
	} else {
		vendorPackageRelativePath = path.Join("..", "..", "..", "..", pkgPath)
	}
	return vendorPackageRelativePath
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
