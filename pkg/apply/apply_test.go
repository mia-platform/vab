//go:build !e2e
// +build !e2e

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
	configPath := testutils.GetTestFile("apply", testBuildFolder, testConfigFileName)
	outputDir := "./outputT1"
	err := Apply(log, configPath, outputDir, true, testutils.TestGroupName2, testutils.TestClusterName1, testutils.GetTestFile("apply", testBuildFolder))

	assert.FileExists(t, "./outputT1/test-cluster-res")

	os.RemoveAll("./outputT1")
	if !assert.NoError(t, err) {
		return
	}
}

func TestBuildFunctionForAGroup(t *testing.T) {
	log := logger.DisabledLogger{}
	configPath := testutils.GetTestFile("apply", testBuildFolder, testConfigFileName)
	outputDir := "./outputT2"
	err := Apply(log, configPath, outputDir, true, testutils.TestGroupName2, "", testutils.GetTestFile("apply", testBuildFolder))

	assert.FileExists(t, "./outputT2/test-cluster-res")
	assert.FileExists(t, "./outputT2/test-cluster2-res")

	os.RemoveAll("./outputT2")
	if !assert.NoError(t, err) {
		return
	}
}

func TestWrongContextPath(t *testing.T) {
	log := logger.DisabledLogger{}
	configPath := testutils.GetTestFile("apply", testBuildFolder, testConfigFileName)
	outputDir := "./outputT3"
	err := Apply(log, configPath, outputDir, true, "", "", testutils.InvalidFolderPath)

	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrNotExist)
	}

	err = Apply(log, configPath, outputDir, true, "", "", configPath)
	if assert.Error(t, err) {
		assert.Equal(t, err.Error(), fmt.Sprintf("the target path %s is not a directory", configPath))
	}

	os.RemoveAll("./outputT3")
}

func TestBuildInvalidConfigPath(t *testing.T) {
	log := logger.DisabledLogger{}
	contextPath := testutils.GetTestFile("apply", testBuildFolder)
	outputDir := "./outputT4"
	err := Apply(log, testutils.InvalidFileName, outputDir, true, "", "", contextPath)

	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrNotExist)
	}
	os.RemoveAll("./outputT4")
}

func TestBuildInvalidKustomization(t *testing.T) {
	log := logger.DisabledLogger{}
	configPath := testutils.GetTestFile("apply", testBuildFolder, testConfigFileName)
	outputDir := "./outputT5"
	err := Apply(log, configPath, outputDir, true, testutils.TestGroupName1, testutils.TestClusterName1, testutils.GetTestFile("apply", testBuildFolder))
	assert.Error(t, err)
	os.RemoveAll("./outputT5")
}

func TestGetContextError(t *testing.T) {
	configPath := testutils.GetTestFile("apply", testBuildFolder, testConfigFileName)
	_, err := getContext(configPath, "notExistent", "test-cluster")
	assert.Error(t, err)

	_, err = getContext(configPath, "test-group2", "notExistent")
	assert.Error(t, err)

	configPathError := "notExistent"
	_, err = getContext(configPathError, "test-group", "test-cluster")
	assert.Error(t, err)
}

const (
	testBuildFolder    = "apply-test"
	testConfigFileName = "testconfig.yaml"
)
