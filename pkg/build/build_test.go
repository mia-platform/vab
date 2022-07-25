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
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strings"
	"testing"

	"github.com/mia-platform/vab/internal/testutils"
	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/logger"
	"golang.org/x/exp/slices"
)

// Test that the function returns the correct kustomized configuration
func TestRunKustomizeBuild(t *testing.T) {
	targetPath := testutils.GetTestFile("build", testutils.KustomizeTestDirName)

	buffer := new(bytes.Buffer)
	if err := runKustomizeBuild(targetPath, buffer); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(buffer.Bytes(), []byte(expectedKustomizeResult)) {
		t.Fatalf("Unexpected Kustomize result:\n%s", buffer.String())
	}
}

// Returns an error if the path is invalid
func TestInvalidKustomizeBuildPath(t *testing.T) {
	buffer := new(bytes.Buffer)
	err := runKustomizeBuild(testutils.InvalidFolderPath, buffer)
	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrNotExist)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrNotExist, err)
	}
}

func TestBuildFunctionForASingleCluster(t *testing.T) {
	log := logger.DisabledLogger{}
	buffer := new(bytes.Buffer)
	configPath := testutils.GetTestFile("build", testBuildFolder, testConfigFileName)
	err := Build(log, configPath, testutils.TestGroupName2, testutils.TestClusterName1, testutils.GetTestFile("build", testBuildFolder), buffer)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	writtenLogs := buffer.String()
	writtenLines := strings.Split(writtenLogs, "\n")

	expectedStarterMarker := fmt.Sprintf(startMarkerFormat, path.Join(utils.ClustersDirName, testutils.TestGroupName2, testutils.TestClusterName1))
	if len(writtenLines) < 15 {
		t.Fatalf("Unexpected line length %d", len(writtenLines))
	}

	if !slices.Contains(writtenLines, expectedStarterMarker) {
		t.Log(writtenLines[0])
		t.Fatal("Start marker for cluster not found")
	}

	if !slices.Contains(writtenLines, endMarkerString) {
		t.Log(writtenLines[len(writtenLines)-1])
		t.Fatal("End marker not found")
	}
}

func TestBuildFunctionForAGroup(t *testing.T) {
	log := logger.DisabledLogger{}
	buffer := new(bytes.Buffer)
	configPath := testutils.GetTestFile("build", testBuildFolder, testConfigFileName)
	err := Build(log, configPath, testutils.TestGroupName2, "", testutils.GetTestFile("build", testBuildFolder), buffer)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	writtenLogs := buffer.String()
	writtenLines := strings.Split(writtenLogs, "\n")

	expectedStarterMarker := fmt.Sprintf(startMarkerFormat, path.Join(utils.ClustersDirName, testutils.TestGroupName2, testutils.TestClusterName1))
	if len(writtenLines) < 29 {
		t.Log(writtenLogs)
		t.Fatalf("Unexpected line length %d", len(writtenLines))
	}

	if !slices.Contains(writtenLines, expectedStarterMarker) {
		t.Log(writtenLines[0])
		t.Fatal("Start marker for cluster not found")
	}

	expectedSeparationMarker := fmt.Sprintf(startMarkerFormat, path.Join(utils.ClustersDirName, testutils.TestGroupName2, testutils.TestClusterName2))
	if !slices.Contains(writtenLines, expectedSeparationMarker) {
		t.Log(expectedSeparationMarker)
		t.Fatal("Sepration for the two clusters not found")
	}

	if !slices.Contains(writtenLines, endMarkerString) {
		t.Log(writtenLines[len(writtenLines)-1])
		t.Fatal("End marker not found")
	}
}

func TestWrongContextPath(t *testing.T) {
	log := logger.DisabledLogger{}
	buffer := new(bytes.Buffer)
	configPath := testutils.GetTestFile("build", testBuildFolder, testConfigFileName)
	err := Build(log, configPath, "", "", testutils.InvalidFolderPath, buffer)

	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrNotExist)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrNotExist, err)
	}

	err = Build(log, configPath, "", "", configPath, buffer)
	if err == nil {
		t.Fatalf("No error was returned")
	}

	if err.Error() != fmt.Sprintf("the target path %s is not a directory", configPath) {
		t.Fatalf("Unexpected Error: %s", err)
	}
}

func TestBuildInvalidConfigPath(t *testing.T) {
	log := logger.DisabledLogger{}
	buffer := new(bytes.Buffer)
	contextPath := testutils.GetTestFile("build", testBuildFolder)
	err := Build(log, testutils.InvalidFileName, "", "", contextPath, buffer)

	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrNotExist)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrNotExist, err)
	}
}

func TestBuildInvalidKustomization(t *testing.T) {
	log := logger.DisabledLogger{}
	buffer := new(bytes.Buffer)
	configPath := testutils.GetTestFile("build", testBuildFolder, testConfigFileName)
	err := Build(log, configPath, testutils.TestGroupName1, testutils.TestClusterName1, testutils.GetTestFile("build", testBuildFolder), buffer)

	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrNotExist)
	}
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
