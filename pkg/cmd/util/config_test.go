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
	"os"
	"path/filepath"
	"testing"

	"github.com/mia-platform/vab/internal/utils"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/kustomize/api/konfig"
)

func TestInitAllGroups(t *testing.T) {
	t.Parallel()

	testDirPath := t.TempDir()
	err := InitGroup(testDirPath)
	if !assert.NoError(t, err) {
		return
	}

	basesDir := filepath.Join(testDirPath, utils.BasesDir)
	customResourcesDir := filepath.Join(testDirPath, utils.CustomResourcesDir)
	assert.DirExists(t, basesDir, "The bases directory does not exist")
	assert.DirExists(t, customResourcesDir, "The custom-resources directory does not exist")

	allGroupsKustomizationPath := filepath.Join(testDirPath, konfig.DefaultKustomizationFileName())
	assert.FileExists(t, allGroupsKustomizationPath, "Missing kustomization file in the all-groups directory")
	expectedKustomization, err := os.ReadFile(filepath.Join("testdata", "all_groups_kustomization.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	actualKustomization, err := os.ReadFile(allGroupsKustomizationPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expectedKustomization, actualKustomization, "Unexpected file content")

	basesDirKustomizationPath := filepath.Join(basesDir, konfig.DefaultKustomizationFileName())
	assert.FileExists(t, basesDirKustomizationPath, "Missing kustomization file in the bases directory")

	customResourcesKustomizationPath := filepath.Join(customResourcesDir, konfig.DefaultKustomizationFileName())
	assert.FileExists(t, customResourcesKustomizationPath, "Missing kustomization file in the custom-resources directory")
}
