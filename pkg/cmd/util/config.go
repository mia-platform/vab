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
	"os"
	"path/filepath"

	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"sigs.k8s.io/kustomize/api/types"
)

func InitializeConfiguration(name, path string) error {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}

	config := v1alpha1.EmptyConfig(name)
	if err := utils.WriteConfig(*config, path); err != nil {
		return fmt.Errorf("failed to write new config: %w", err)
	}

	return InitGroups(config, path)
}

func InitGroups(config *v1alpha1.ClustersConfiguration, path string) error {
	if err := InitGroup(filepath.Join(path, utils.AllGroupsDirPath)); err != nil {
		return err
	}

	for _, group := range config.Spec.Groups {
		for _, cluster := range group.Clusters {
			if err := InitGroup(filepath.Join(path, group.Name, cluster.Name)); err != nil {
				return err
			}
		}
	}

	return nil
}

func InitGroup(path string) error {
	basesDir := filepath.Join(path, utils.BasesDir)
	customResourcesDir := filepath.Join(path, utils.CustomResourcesDir)

	for _, path := range []string{path, basesDir, customResourcesDir} {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return fmt.Errorf("error initializing group: %w", err)
		}
		var kustomization types.Kustomization
		switch path {
		case basesDir:
			kustomization = utils.EmptyKustomization()
			kustomization.Resources = []string{}
		case customResourcesDir:
			kustomization = utils.EmptyComponent()
			kustomization.Resources = []string{}
		default:
			kustomization = utils.EmptyKustomization()
			kustomization.Resources = append(kustomization.Resources, utils.BasesDir)
			kustomization.Components = append(kustomization.Components, utils.CustomResourcesDir)
		}

		if err := utils.WriteKustomization(kustomization, path, false); err != nil {
			return err
		}
	}

	return nil
}
