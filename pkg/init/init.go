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

package init

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/mia-platform/vab/pkg/logger"
)

func NewProject(logger logger.LogInterface, currentPath string, optionalName string) error {
	logger.V(10).Write("Ensuring that the target path exists...")
	configPath, err := ensureProjectPath(currentPath, optionalName)
	if err != nil {
		logger.V(10).Write("Error while ensuring the project path")
		return err
	}

	name := filepath.Base(configPath)
	logger.V(5).Writef("Selected project name: %s", name)

	logger.V(10).Write("Writing empty configuration...")
	if err := utils.WriteConfig(*v1alpha1.EmptyConfig(name), configPath); err != nil {
		logger.V(10).Write("Error while writing the configuration file")
		return err
	}

	logger.V(10).Write("Initializing the all-groups directory...")
	if err := initAllGroups(configPath); err != nil {
		logger.V(10).Write("Error while writing the kustomize file")
		return err
	}

	return nil
}

// ensureProjectPath will return a cleaned and complete path based on currentPath and optional name
// ensuring that the appropriate folders are present on file system
func ensureProjectPath(basePath string, name string) (string, error) {
	projectPath := filepath.Clean(basePath)
	if name != "" {
		projectPath = filepath.Join(projectPath, name)
	}

	if err := os.Mkdir(projectPath, fs.ModePerm); err != nil && !errors.Is(err, fs.ErrExist) {
		return "", err
	}

	return projectPath, nil
}

// initAllGroups initializes the all-groups directory
func initAllGroups(configPath string) error {
	allGroupsDir := filepath.Join(configPath, utils.AllGroupsDirPath)
	if err := os.MkdirAll(allGroupsDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating path %s: %w", allGroupsDir, err)
	}
	basesDir := filepath.Join(allGroupsDir, utils.BasesDir)
	customResourcesDir := filepath.Join(allGroupsDir, utils.CustomResourcesDir)
	if err := os.Mkdir(basesDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %w", basesDir, err)
	}
	if err := os.Mkdir(customResourcesDir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory %s: %w", customResourcesDir, err)
	}
	if err := utils.WriteKustomization(utils.EmptyKustomization(), basesDir); err != nil {
		return fmt.Errorf("error writing kustomization file in %s: %w", basesDir, err)
	}
	if err := utils.WriteKustomization(utils.EmptyComponent(), customResourcesDir); err != nil {
		return fmt.Errorf("error writing kustomization file in %s: %w", customResourcesDir, err)
	}
	allGroupsKustomization := utils.EmptyKustomization()
	allGroupsKustomization.Resources = append(allGroupsKustomization.Resources, utils.BasesDir)
	allGroupsKustomization.Components = append(allGroupsKustomization.Components, utils.CustomResourcesDir)
	if err := utils.WriteKustomization(allGroupsKustomization, allGroupsDir); err != nil {
		return fmt.Errorf("error writing kustomization file in %s: %w", allGroupsDir, err)
	}
	return nil
}
