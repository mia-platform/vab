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

package apply

import (
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/mia-platform/vab/internal/testutils"
	"github.com/mia-platform/vab/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestBuildFunctionForASingleCluster(t *testing.T) {
	log := logger.DisabledLogger{}
	configPath := testutils.GetTestFile("build", testBuildFolder, testConfigFileName)
	outputDir := "./output"
	err := Apply(log, configPath, outputDir, testutils.TestGroupName2, testutils.TestClusterName1, testutils.GetTestFile("build", testBuildFolder))

	assert.FileExists(t, "./output/test-cluster")

	os.RemoveAll("./output")
	if !assert.NoError(t, err) {
		return
	}
}

func TestBuildFunctionForAGroup(t *testing.T) {
	log := logger.DisabledLogger{}
	configPath := testutils.GetTestFile("build", testBuildFolder, testConfigFileName)
	outputDir := "./output"
	err := Apply(log, configPath, outputDir, testutils.TestGroupName2, "", testutils.GetTestFile("build", testBuildFolder))

	assert.FileExists(t, "./output/test-cluster")
	assert.FileExists(t, "./output/test-cluster2")

	os.RemoveAll("./output")
	if !assert.NoError(t, err) {
		return
	}
}

func TestWrongContextPath(t *testing.T) {
	log := logger.DisabledLogger{}
	configPath := testutils.GetTestFile("build", testBuildFolder, testConfigFileName)
	outputDir := "./output"
	err := Apply(log, configPath, outputDir, "", "", testutils.InvalidFolderPath)

	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrNotExist)
	}

	err = Apply(log, configPath, outputDir, "", "", configPath)
	if assert.Error(t, err) {
		assert.Equal(t, err.Error(), fmt.Sprintf("the target path %s is not a directory", configPath))
	}
}

func TestBuildInvalidConfigPath(t *testing.T) {
	log := logger.DisabledLogger{}
	contextPath := testutils.GetTestFile("build", testBuildFolder)
	outputDir := "./output"
	err := Apply(log, testutils.InvalidFileName, outputDir, "", "", contextPath)

	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrNotExist)
	}
}

func TestBuildInvalidKustomization(t *testing.T) {
	log := logger.DisabledLogger{}
	configPath := testutils.GetTestFile("build", testBuildFolder, testConfigFileName)
	outputDir := "./output"
	err := Apply(log, configPath, outputDir, testutils.TestGroupName1, testutils.TestClusterName1, testutils.GetTestFile("build", testBuildFolder))
	assert.Error(t, err)
}

func TestGetContextError(t *testing.T) {
	configPath := testutils.GetTestFile("build", testBuildFolder, testConfigFileName)
	_, err := getContext(configPath, "notExistent", "test-cluster")
	assert.Error(t, err)

	_, err = getContext(configPath, "test-group", "notExistent")
	assert.Error(t, err)

	configPathError := "notExistent"
	_, err = getContext(configPathError, "test-group", "test-cluster")
	assert.Error(t, err)
}

const (
	testBuildFolder    = "build-test"
	testConfigFileName = "testconfig.yaml"
)
