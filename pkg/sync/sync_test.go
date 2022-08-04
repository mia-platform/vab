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
	"path"
	"testing"

	"github.com/mia-platform/vab/internal/git"
	"github.com/mia-platform/vab/internal/testutils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/mia-platform/vab/pkg/logger"
	"github.com/stretchr/testify/assert"
)

const (
	testPackageName    = "test-package-1.0.0"
	testModuleBasePath = "modules/test-module1"
	testFileName1      = "test-flavour1/file1.yaml"
	testFileName2      = "test-flavour1/file2.yaml"
	testFileName3      = "test-flavour2/file1.yaml"
)

func TestClonePackage(t *testing.T) {
	logger := logger.DisabledLogger{}
	testModule := v1alpha1.Module{
		Version: "1.0.0",
		Weight:  2,
	}
	outputFiles, err := ClonePackages(logger, testPackageName, testModule, git.MockGetFilesForPackage)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, outputFiles, "The returned array of mocked file pointers is empty")
}

func TestMoveToDisk(t *testing.T) {
	logger := logger.DisabledLogger{}
	fakeWorktree := testutils.PrepareFakeWorktree(t)
	input := []*git.File{
		git.NewFile(path.Join(testModuleBasePath, testFileName1), testModuleBasePath, fakeWorktree),
		git.NewFile(path.Join(testModuleBasePath, testFileName2), testModuleBasePath, fakeWorktree),
		git.NewFile(path.Join(testModuleBasePath, testFileName3), testModuleBasePath, fakeWorktree),
	}
	testDirPath := t.TempDir()
	err := MoveToDisk(logger, input, testPackageName, testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.FileExists(t, path.Join(testDirPath, testFileName1), "Mock file 1 does not exist on disk")
	assert.FileExists(t, path.Join(testDirPath, testFileName2), "Mock file 2 does not exist on disk")
	assert.FileExists(t, path.Join(testDirPath, testFileName3), "Mock file 3 does not exist on disk")
}

func TestSyncAddons(t *testing.T) {

}
