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
	"sort"
	"strings"

	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/api/konfig"
	kustomize "sigs.k8s.io/kustomize/api/types"
)

// SyncKustomizeResources updates the clusters' kustomization resources to the latest sync
func SyncKustomizeResources(modules *map[string]v1alpha1.Module, addons *map[string]v1alpha1.AddOn, k kustomize.Kustomization) *kustomize.Kustomization {
	resourcesList := getSortedModulesList(modules)
	addonsList := getAddOnsList(addons)
	resourcesList = append(resourcesList, addonsList...)

	// If the file already includes a non-empty list of resources, this function
	// collects all the custom modules that were added manually by the user
	// (i.e. all those modules that are not present in the vendors folder, thus
	// without "vendors" in their path). Then, the custom modules are appended
	// to the updated modules list that will substitute the already existing one.
	if k.Resources != nil {
		for _, r := range k.Resources {
			if !strings.Contains(r, "vendors/") {
				resourcesList = append(resourcesList, r)
			}
		}
	}

	k.Resources = resourcesList

	return &k
}

// getSortedModulesList returns the list of module names sorted by weight.
// In case of equal weights, the modules are ordered lexicographically.
func getSortedModulesList(modules *map[string]v1alpha1.Module) []string {
	modulesList := make([]string, 0, len(*modules))

	for m := range *modules {
		if !(*modules)[m].Disable {
			// modulePath := path.Join(utils.VendorsModulesPath, m)
			modulesList = append(modulesList, m)
		}
	}

	sort.SliceStable(modulesList, func(i, j int) bool {
		// If the weights are equal, order the elements lexicographically
		if (*modules)[modulesList[i]].Weight == (*modules)[modulesList[j]].Weight {
			return modulesList[i] < modulesList[j]
		}
		// Otherwise, sort by weight (increasing order)
		return (*modules)[modulesList[i]].Weight < (*modules)[modulesList[j]].Weight
	})

	return *fixResourcesPath(modulesList, true)
}

// getAddOnsList returns the list of addons names in lexicographic order
func getAddOnsList(addons *map[string]v1alpha1.AddOn) []string {
	addonsList := make([]string, 0, len(*addons))

	for ao := range *addons {
		if !(*addons)[ao].Disable {
			addonsList = append(addonsList, ao)
		}
	}

	sort.SliceStable(addonsList, func(i, j int) bool {
		return addonsList[i] < addonsList[j]
	})

	return *fixResourcesPath(addonsList, false)
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
func fixResourcesPath(resourcesList []string, isModulesList bool) *[]string {
	fixedResourcesList := make([]string, 0, len(resourcesList))
	for _, res := range resourcesList {
		if isModulesList {
			fixedResourcesList = append(fixedResourcesList, path.Join("..", "..", "..", utils.VendorsModulesPath, res))
		} else {
			fixedResourcesList = append(fixedResourcesList, path.Join("..", "..", "..", utils.VendorsAddOnsPath, res))
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
