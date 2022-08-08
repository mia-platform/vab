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
	"os"
	"path"
	"testing"

	"github.com/mia-platform/vab/internal/git"
	"github.com/mia-platform/vab/internal/testutils"
	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/mia-platform/vab/pkg/logger"
	"github.com/stretchr/testify/assert"
)

const (
	testPackageName       = "test-module1/test-flavour1"
	expectedKustomization = `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
  - ../../../all-groups
`
)

func TestClonePackage(t *testing.T) {
	logger := logger.DisabledLogger{}
	testModule := v1alpha1.Module{
		Version: "1.0.0",
		Weight:  2,
	}
	outputFiles, err := ClonePackages(logger, testPackageName, testModule, testutils.FakeFilesGetter{Testing: t})
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, outputFiles, "The returned array of mocked file pointers is empty")
}

func TestMoveToDisk(t *testing.T) {
	logger := logger.DisabledLogger{}
	fakeWorktree := testutils.PrepareFakeWorktree(t)
	input := []*git.File{
		git.NewFile("./modules/test-module1/test-flavour1/file1.yaml", "./modules/test-module1", *fakeWorktree),
		git.NewFile("./modules/test-module1/test-flavour1/file2.yaml", "./modules/test-module1", *fakeWorktree),
		git.NewFile("./modules/test-module1/test-flavour2/file1.yaml", "./modules/test-module1", *fakeWorktree),
	}
	testDirPath := t.TempDir()
	err := MoveToDisk(logger, input, testPackageName, testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.FileExists(t, path.Join(testDirPath, "test-flavour1/file1.yaml"), "Mock file 1 does not exist on disk")
	assert.FileExists(t, path.Join(testDirPath, "test-flavour1/file2.yaml"), "Mock file 2 does not exist on disk")
	assert.FileExists(t, path.Join(testDirPath, "test-flavour2/file1.yaml"), "Mock file 3 does not exist on disk")
}

func TestSyncModules(t *testing.T) {
	logger := logger.DisabledLogger{}
	modules := make(map[string]v1alpha1.Module)
	modules["test-module1/test-flavour3"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  4,
	}
	modules["test-module2/test-flavour2"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  1,
	}
	modules["test-module2/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  3,
		Disable: true,
	}
	testDirPath := t.TempDir()
	err := SyncModules(logger, modules, testDirPath, testutils.FakeFilesGetter{Testing: t})
	if !assert.NoError(t, err) {
		return
	}
}

func TestSyncAddons(t *testing.T) {
	logger := logger.DisabledLogger{}
	addons := make(map[string]v1alpha1.AddOn)
	addons["test-addon1"] = v1alpha1.AddOn{
		Version: "1.0.0",
	}
	addons["test-addon2"] = v1alpha1.AddOn{
		Version: "1.0.0",
		Disable: true,
	}
	testDirPath := t.TempDir()
	err := SyncAddons(logger, addons, testDirPath, testutils.FakeFilesGetter{Testing: t})
	if !assert.NoError(t, err) {
		return
	}
}

func TestUpdateBasesAllGroups(t *testing.T) {
	testDirPath := t.TempDir()
	targetPath := path.Join(testDirPath, utils.AllGroupsDirPath)
	os.MkdirAll(targetPath, os.ModePerm)
	modules := make(map[string]v1alpha1.Module)
	modules["test-module1/test-flavour3"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  4,
	}
	modules["test-module2/test-flavour2"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  1,
	}
	modules["test-module2/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  3,
		Disable: true,
	}
	addons := make(map[string]v1alpha1.AddOn)
	addons["test-addon1"] = v1alpha1.AddOn{
		Version: "1.0.0",
	}
	addons["test-addon2"] = v1alpha1.AddOn{
		Version: "1.0.0",
		Disable: true,
	}
	err := UpdateBases(targetPath, modules, addons)
	if !assert.NoError(t, err) {
		return
	}
}

func TestUpdateBasesCluster(t *testing.T) {
	testDirPath := t.TempDir()
	targetPath := path.Join(testDirPath, "groups/group-1/cluster-1")
	os.MkdirAll(targetPath, os.ModePerm)
	err := UpdateBases(targetPath, nil, nil)
	if !assert.NoError(t, err) {
		return
	}
	file, _ := os.ReadFile(path.Join(targetPath, "bases", utils.KustomizationFileName))
	t.Log(string(file))
	compareFile(t, []byte(expectedKustomization), path.Join(targetPath, "bases", utils.KustomizationFileName))
}
