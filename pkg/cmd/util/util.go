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
	"path/filepath"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

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
		return group, fmt.Errorf("reding config file: %w", err)
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

// ClusterPath return the canonical path for a cluster given the group and name for the cluster
func ClusterPath(group, cluster string) string {
	return filepath.Join(clustersDirName, group, cluster)
}

func VendoredModulePath(packageName string) string {
	return filepath.Join(modulesDirPath, packageName)
}

func VendoredAddOnPath(packageName string) string {
	return filepath.Join(addOnsDirPath, packageName)
}
