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

package init

import (
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/mia-platform/vab/internal/testutils"
	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/logger"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/kustomize/api/konfig"
)

const (
	testName = "foo"
)

// Return current path if the name arg is the empty string
func TestCurrentPath(t *testing.T) {
	testDirPath := t.TempDir()
	dstPath, err := ensureProjectPath(testDirPath, "")

	if assert.NoError(t, err) {
		assert.Equal(t, dstPath, testDirPath)
	}
}

// Return path of a new project directory named "foo"
func TestNewProjectPath(t *testing.T) {
	testDirPath := t.TempDir()
	expectedPath := path.Join(testDirPath, testName)
	dstPath, err := ensureProjectPath(testDirPath, testName)

	if assert.NoError(t, err) {
		assert.Equal(t, dstPath, expectedPath)
	}
}

// Return ErrNotExists if the path parameter is invalid (empty name)
func TestCurrentInvalidPath(t *testing.T) {
	_, err := ensureProjectPath(testutils.InvalidFolderPath, "")

	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrNotExist)
	}
}

// Return ErrNotExists if the path parameter is invalid (non-empty name)
func TestNewInvalidPath(t *testing.T) {
	_, err := ensureProjectPath(testutils.InvalidFolderPath, testName)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrNotExist)
	}
}

// Return ErrPermission if access to the path is denied (non-empty name)
func TestNewPathErrPermission(t *testing.T) {
	testDirPath := t.TempDir()
	if err := os.Chmod(testDirPath, 0); err != nil {
		t.Fatal(err)
	}
	_, err := ensureProjectPath(testDirPath, testName)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrPermission)
	}
}

// If a directory with the specified name already exists, return the path to that directory
func TestNewExistingPath(t *testing.T) {
	testDirPath := t.TempDir()
	expectedPath := path.Join(testDirPath, testName)
	if err := os.Mkdir(expectedPath, fs.ModePerm); !assert.NoError(t, err) {
		return
	}

	dstPath, err := ensureProjectPath(testDirPath, testName)
	if assert.NoError(t, err) {
		assert.Equal(t, dstPath, expectedPath)
	}
}

// Test that the directory for the new cluster is created correctly with its empty kustomization.yaml
func TestNewClusterOverride(t *testing.T) {
	testDirPath := t.TempDir()
	if err := os.Mkdir(path.Join(testDirPath, utils.ClustersDirName), fs.ModePerm); !assert.NoError(t, err) {
		return
	}

	if err := createClusterOverride(testDirPath, testName); !assert.NoError(t, err) {
		return
	}

	kustomizationPath := path.Join(testDirPath, utils.ClustersDirName, testName, konfig.DefaultKustomizationFileName())
	_, err := os.Stat(kustomizationPath)
	assert.NoError(t, err)
}

// Test that the clusters directory is created if missing
func TestMissingClustersDirectory(t *testing.T) {
	testDirPath := t.TempDir()
	err := createClusterOverride(testDirPath, testName)
	assert.NoError(t, err)
}

// Return ErrPermission if if access to the path is denied
func TestNewClusterErrPermission(t *testing.T) {
	testDirPath := t.TempDir()
	if err := os.Mkdir(path.Join(testDirPath, utils.ClustersDirName), 0); assert.NoError(t, err) {
		return
	}

	err := createClusterOverride(testDirPath, testName)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrPermission)
	}
}

// Return ErrPermission if access to the path is denied
func TestMissingClustersDirErrPermission(t *testing.T) {
	testDirPath := t.TempDir()
	if err := os.Chmod(testDirPath, 0); assert.NoError(t, err) {
		return
	}

	err := createClusterOverride(testDirPath, testName)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrPermission)
	}
}

func TestInitProject(t *testing.T) {
	testDirPath := t.TempDir()
	testProjectName := "foo"
	expectedProjectPath := path.Join(testDirPath, testProjectName)
	logger := logger.DisabledLogger{}
	err := NewProject(logger, testDirPath, testProjectName)

	if !assert.NoError(t, err) {
		return
	}

	expectedConfigPath := path.Join(expectedProjectPath, utils.DefaultConfigFilename)
	_, err = os.Stat(expectedConfigPath)
	if !assert.NoError(t, err) {
		return
	}

	config, err := utils.ReadConfig(expectedConfigPath)

	if assert.NoError(t, err) {
		assert.Equal(t, config.Name, testProjectName, "Unexpected project name")
		_, err = os.Stat(path.Join(expectedProjectPath, utils.ClustersDirName, "all-groups", konfig.DefaultKustomizationFileName()))
		assert.NoError(t, err)
	}
}
