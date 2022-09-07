//go:build !e2e
// +build !e2e

// Copyright 2022 Mia-Platform

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package build

import (
	"bytes"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"testing"

	"github.com/mia-platform/vab/internal/testutils"
	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/logger"
	"github.com/stretchr/testify/assert"
)

// Test that the function returns the correct kustomized configuration
func TestRunKustomizeBuild(t *testing.T) {
	targetPath := testutils.GetTestFile("build", testutils.KustomizeTestDirName)

	buffer := new(bytes.Buffer)
	if err := RunKustomizeBuild(targetPath, buffer); assert.NoError(t, err) {
		assert.Equal(t, buffer.String(), expectedKustomizeResult)
	}
}

// Returns an error if the path is invalid
func TestInvalidKustomizeBuildPath(t *testing.T) {
	buffer := new(bytes.Buffer)
	err := RunKustomizeBuild(testutils.InvalidFolderPath, buffer)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrNotExist)
	}
}

func TestBuildFunctionForASingleCluster(t *testing.T) {
	log := logger.DisabledLogger{}
	buffer := new(bytes.Buffer)
	configPath := testutils.GetTestFile("build", testBuildFolder, testConfigFileName)
	err := Build(log, configPath, testutils.TestGroupName2, testutils.TestClusterName1, testutils.GetTestFile("build", testBuildFolder), buffer)
	if !assert.NoError(t, err) {
		return
	}

	writtenLogs := buffer.String()
	writtenLines := strings.Split(writtenLogs, "\n")

	expectedStarterMarker := fmt.Sprintf(startMarkerFormat, path.Join(utils.ClustersDirName, testutils.TestGroupName2, testutils.TestClusterName1))
	assert.NotEqual(t, writtenLines, 15, "Unexpected line length")
	assert.Contains(t, writtenLines, expectedStarterMarker, "Start marker for cluster not found")
	assert.Contains(t, writtenLines, endMarkerString, "End marker not found")
}

func TestBuildFunctionForAGroup(t *testing.T) {
	log := logger.DisabledLogger{}
	buffer := new(bytes.Buffer)
	configPath := testutils.GetTestFile("build", testBuildFolder, testConfigFileName)
	err := Build(log, configPath, testutils.TestGroupName2, "", testutils.GetTestFile("build", testBuildFolder), buffer)

	if !assert.NoError(t, err) {
		return
	}

	writtenLogs := buffer.String()
	writtenLines := strings.Split(writtenLogs, "\n")

	expectedStarterMarker := fmt.Sprintf(startMarkerFormat, path.Join(utils.ClustersDirName, testutils.TestGroupName2, testutils.TestClusterName1))
	expectedSeparationMarker := fmt.Sprintf(startMarkerFormat, path.Join(utils.ClustersDirName, testutils.TestGroupName2, testutils.TestClusterName2))
	assert.NotEqual(t, writtenLines, 29, "Unexpected line length")
	assert.Contains(t, writtenLines, expectedStarterMarker, "Start marker for cluster not found")
	assert.Contains(t, writtenLines, expectedSeparationMarker, "Sepration for the two clusters not found")
	assert.Contains(t, writtenLines, endMarkerString, "End marker not found")
}

func TestWrongContextPath(t *testing.T) {
	log := logger.DisabledLogger{}
	buffer := new(bytes.Buffer)
	configPath := testutils.GetTestFile("build", testBuildFolder, testConfigFileName)
	err := Build(log, configPath, "", "", testutils.InvalidFolderPath, buffer)

	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrNotExist)
	}

	err = Build(log, configPath, "", "", configPath, buffer)
	if assert.Error(t, err) {
		assert.Equal(t, err.Error(), fmt.Sprintf("the target path %s is not a directory", configPath))
	}
}

func TestBuildInvalidConfigPath(t *testing.T) {
	log := logger.DisabledLogger{}
	buffer := new(bytes.Buffer)
	contextPath := testutils.GetTestFile("build", testBuildFolder)
	err := Build(log, testutils.InvalidFileName, "", "", contextPath, buffer)

	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrNotExist)
	}
}

func TestBuildInvalidKustomization(t *testing.T) {
	log := logger.DisabledLogger{}
	buffer := new(bytes.Buffer)
	configPath := testutils.GetTestFile("build", testBuildFolder, testConfigFileName)
	err := Build(log, configPath, testutils.TestGroupName1, testutils.TestClusterName1, testutils.GetTestFile("build", testBuildFolder), buffer)
	assert.Error(t, err)
}

const (
	testBuildFolder         = "build-test"
	testConfigFileName      = "testconfig.yaml"
	startMarkerFormat       = "### BUILD RESULTS FOR: %s ###"
	endMarkerString         = "---"
	expectedKustomizeResult = `apiVersion: apps/v1
kind: Service
metadata:
  name: test
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 80
  selector:
    app: test
  type: ClusterIp
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
spec:
  replicas: 1
  selector:
    app: test
  template:
    metadata:
      labels:
        app: test
    spec:
      containers:
      - image: nginx
        name: test
        resources:
          limits:
            cpu: 10m
            memory: 10Mi
`
)
