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
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// WriteKustomizationData read kustomize configuration file at path and output the kustomize build result to writer
func WriteKustomizationData(path string, writer io.Writer) error {
	kOpts := krusty.MakeDefaultOptions()
	kOpts.Reorder = krusty.ReorderOptionLegacy
	k := krusty.MakeKustomizer(kOpts)
	resourceMap, err := k.Run(filesys.MakeFsOnDisk(), path)
	if err != nil {
		return err
	}

	yamlData, err := resourceMap.AsYaml()
	if err != nil {
		return err
	}

	_, err = writer.Write(yamlData)
	return err
}

// groupFromConfig return a Group struct if a group with groupName is found inside the configuration at path.
// Will return an error if the file cannot be read or groupName is not found
func GroupFromConfig(groupName string, path string) (v1alpha1.Group, error) {
	var group v1alpha1.Group
	config, err := ReadConfig(path)
	if err != nil {
		return group, err
	}

	found := false
	for _, configGroup := range config.Spec.Groups {
		if configGroup.Name == groupName {
			found = true
			group = configGroup
			break
		}
	}

	if !found {
		return group, fmt.Errorf("no %q group in config at path %q", groupName, path)
	}

	return group, nil
}

// ValidateContextPath will validate contextPath that is a valid existing path, and that is a directory
// it will also return the path in absolute form
func ValidateContextPath(contextPath string) (string, error) {
	var cleanedContextPath string
	var err error
	if cleanedContextPath, err = filepath.Abs(contextPath); err != nil {
		return "", err
	}

	var contextInfo fs.FileInfo
	if contextInfo, err = os.Stat(cleanedContextPath); err != nil {
		return "", err
	}

	if !contextInfo.IsDir() {
		return "", fmt.Errorf(" %q is not a directory", cleanedContextPath)
	}
	return cleanedContextPath, nil
}

// ClusterID return a cluster identifier for group and cluster name
func ClusterID(group, cluster string) string {
	return fmt.Sprintf("%s/%s", group, cluster)
}

// ClusterPath return the canonical path for a cluster given the group and name for the cluster
func ClusterPath(group, cluster string) string {
	return filepath.Join(clustersDirName, group, cluster)
}

// VendoredModulePath return a vendored path for module with packageName
func VendoredModulePath(packageName string) string {
	return filepath.Join(modulesDirPath, packageName)
}

// VendoredAddOnPath return a vendored path for addon with packageName
func VendoredAddOnPath(packageName string) string {
	return filepath.Join(addOnsDirPath, packageName)
}
