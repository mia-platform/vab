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

package util

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"slices"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"sigs.k8s.io/kustomize/api/konfig"
	kustomize "sigs.k8s.io/kustomize/api/types"
	yaml "sigs.k8s.io/yaml/goyaml.v3"
)

const (
	defaultConfigFileName = "config.yaml"

	clustersDirName = "clusters"
	vendorsDirName  = "vendors"

	basesDirName           = "bases"
	customResourcesDirName = "custom-resources"

	kustomization = kustomize.KustomizationKind
	component     = kustomize.ComponentKind

	filePermission      = 0644
	yamlFileIndentation = 2

	doNotEditComment = "File generated by vab. DO NOT EDIT."
)

var (
	allGroupsDirPath = filepath.Join(clustersDirName, "all-groups")
	modulesDirPath   = filepath.Join(vendorsDirName, "modules")
	addOnsDirPath    = filepath.Join(vendorsDirName, "addons")
)

// ReadConfig reads a configuration file into a ClustersConfiguration struct
func ReadConfig(configPath string) (*v1alpha1.ClustersConfiguration, error) {
	if len(configPath) == 0 {
		configPath = defaultConfigFileName
	}

	configFile, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	output := &v1alpha1.ClustersConfiguration{}
	if err := yaml.Unmarshal(configFile, output); err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	return output, nil
}

// InitializeConfiguration will create an empty configuration file at path and then create all the folder
// structure
func InitializeConfiguration(name, path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	config := v1alpha1.EmptyConfig(name)
	if err := writeYamlFile(filepath.Join(path, defaultConfigFileName), config); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return SyncDirectories(config.Spec, path)
}

// SyncDirectories will create all the folders and kustomization files needed by the config data, it will leave
// alone already present file in the custom-resources folder if they already exists, it will override everything else
func SyncDirectories(config v1alpha1.ConfigSpec, path string) error {
	if err := ensureFolderContent(path, allGroupsDirPath, config.Modules, config.AddOns); err != nil {
		return err
	}

	addons := config.AddOns
	modules := config.Modules
	for _, group := range config.Groups {
		for _, cluster := range group.Clusters {
			var clusterModules, clusterAddOns map[string]v1alpha1.Package
			if len(cluster.Modules) != 0 || len(cluster.AddOns) != 0 {
				clusterModules = mergePackages(modules, cluster.Modules)
				clusterAddOns = mergePackages(addons, cluster.AddOns)
			}

			clusterPath := ClusterPath(group.Name, cluster.Name)
			if err := ensureFolderContent(path, clusterPath, clusterModules, clusterAddOns); err != nil {
				return err
			}
		}
	}

	return nil
}

// ensureFolderContent will create the folder structure if needed and create/override the contents of
// the kustomization file under the bases folder, and ensure the presence of the custom-resource folder
// with its kustomization file if they don't exists
func ensureFolderContent(basePath string, clusterPath string, modules, addOns map[string]v1alpha1.Package) error {
	path := filepath.Join(basePath, clusterPath)
	name := filepath.Base(path)
	basesDir := filepath.Join(path, basesDirName)
	customResourcesDir := filepath.Join(path, customResourcesDirName)

	for _, dir := range []string{basesDir, customResourcesDir} {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("creating folders: %w", err)
		}
	}

	// write root kustomization file
	err := writeKustomizationFile(name,
		path,
		kustomization,
		[]string{basesDirName},
		[]string{customResourcesDirName},
		true,
	)
	if err != nil {
		return err
	}

	sortedModules := sortedPackagesPath(basesDir, filepath.Join(basePath, modulesDirPath), modules)
	sortedAddons := sortedPackagesPath(basesDir, filepath.Join(basePath, addOnsDirPath), addOns)
	switch {
	case len(sortedModules) == 0 && clusterPath != allGroupsDirPath:
		sortedModules = append(sortedModules, relativeModulePath(basesDir, filepath.Join(basePath, allGroupsDirPath)))
	case len(sortedAddons) > 0 && clusterPath != allGroupsDirPath:
		sortedAddons = append(sortedAddons, relativeModulePath(basesDir, filepath.Join(basePath, allGroupsDirPath, customResourcesDirName)))
	}

	// write bases file
	err = writeKustomizationFile(fmt.Sprintf("%s - %s", name, basesDirName),
		basesDir,
		kustomization,
		sortedModules,
		sortedAddons,
		true,
	)
	if err != nil {
		return err
	}

	// write custom-resources file if does not exists
	if _, err := os.Stat(filepath.Join(customResourcesDir, konfig.DefaultKustomizationFileName())); errors.Is(err, os.ErrNotExist) {
		return writeKustomizationFile(fmt.Sprintf("%s - %s", name, customResourcesDirName),
			customResourcesDir,
			component,
			[]string{},
			[]string{},
			false,
		)
	}

	return nil
}

// writeKustomizationFile create a new kustomization file of kind.
// It will also set the resources and components property and add a top head comment if generated is true
func writeKustomizationFile(name, path, kind string, resources, components []string, generated bool) error {
	kustomization := &kustomize.Kustomization{}
	kustomization.Kind = kind
	kustomization.MetaData = &kustomize.ObjectMeta{Name: name} // weird trick to allow empty files
	kustomization.FixKustomization()
	kustomization.Resources = resources
	kustomization.Components = components

	node := new(yaml.Node)
	if err := node.Encode(kustomization); err != nil {
		return fmt.Errorf("writing kustomize file: %w", err)
	}

	if generated {
		node.HeadComment = doNotEditComment
	}

	err := writeYamlFile(filepath.Join(path, konfig.DefaultKustomizationFileName()), node)
	if err != nil {
		return fmt.Errorf("writing kustomize file: %w", err)
	}

	return nil
}

// sortedPackagesPath return an array of packages path relative to basePath orderd alphabetically
func sortedPackagesPath(basePath, packagesPath string, packages map[string]v1alpha1.Package) []string {
	paths := make([]string, 0, len(packages))

	for _, pkg := range packages {
		if pkg.Disable {
			continue
		}

		pkgPath := pkg.GetName() + "-" + pkg.Version
		if pkg.IsModule() {
			pkgPath = filepath.Join(pkgPath, pkg.GetFlavorName())
		}

		paths = append(paths, relativeModulePath(basePath, filepath.Join(packagesPath, pkgPath)))
	}

	slices.SortStableFunc(paths, cmp.Compare)
	return paths
}

// mergePackages return a map of merged packages excluding disabled ones, if second has no elements return nil
func mergePackages(first, second map[string]v1alpha1.Package) map[string]v1alpha1.Package {
	mergedMap := make(map[string]v1alpha1.Package, 0)
	maps.Copy(mergedMap, first)
	for name, pkg := range second {
		// if the current package is disabled remove it from the map
		if pkg.Disable {
			delete(mergedMap, name)
		} else {
			mergedMap[name] = pkg
		}
	}

	// return the list of packages with the on disk path as key
	return mergedMap
}

// writeYamlFile marshals the interface passed as argument, and writes it to a YAML file
func writeYamlFile(path string, data interface{}) error {
	buffer := new(bytes.Buffer)
	encoder := yaml.NewEncoder(buffer)
	encoder.SetIndent(yamlFileIndentation)
	encoder.CompactSeqIndent()

	if err := encoder.Encode(data); err != nil {
		return err
	}
	if err := encoder.Close(); err != nil {
		return err
	}

	return os.WriteFile(path, buffer.Bytes(), filePermission)
}

// relativeModulePath return the relative path of a module to basePath from targetPath
func relativeModulePath(basePath, targetPath string) string {
	modulePath, err := filepath.Rel(basePath, targetPath)
	if err != nil {
		panic(err) // we don't expect an error because the paths are computed by us
	}

	return modulePath
}
