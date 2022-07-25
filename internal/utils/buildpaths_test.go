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

package utils

import (
	"path"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mia-platform/vab/internal/testutils"
)

const (
	testGroupsFile = "test_groups.yaml"
)

// Test that the correct path is returned given valid group and cluster
func TestGetClusterPath(t *testing.T) {
	configPath := testutils.GetTestFile("utils", testGroupsFile)
	buildPath, err := BuildPaths(configPath, testutils.TestGroupName1, testutils.TestClusterName1)
	if err != nil {
		t.Fatal(err)
	}

	expectedPath := path.Join(ClustersDirName, testutils.TestGroupName1, testutils.TestClusterName1)
	if buildPath[0] != expectedPath {
		t.Fatalf("Unexpected path. Expected: %s, actual: %s", expectedPath, buildPath[0])
	}
}

// Test that the correct paths are returned given valid group
func TestGetGroupPath(t *testing.T) {
	configPath := testutils.GetTestFile("utils", testGroupsFile)
	buildPaths, err := BuildPaths(configPath, testutils.TestGroupName1, "")
	if err != nil {
		t.Fatal(err)
	}

	clusterPath1 := path.Join(ClustersDirName, testutils.TestGroupName1, testutils.TestClusterName1)
	clusterPath2 := path.Join(ClustersDirName, testutils.TestGroupName1, testutils.TestClusterName2)
	expectedPaths := []string{clusterPath1, clusterPath2}
	if !cmp.Equal(buildPaths, expectedPaths) {
		t.Fatalf("Unexpected paths. Expected: %v, actual: %v", expectedPaths, buildPaths)
	}
}

// Returns an error if the specified group doesn't exist
func TestBuildPathsWrongGroup(t *testing.T) {
	configPath := testutils.GetTestFile("utils", testGroupsFile)
	_, err := BuildPaths(configPath, testutils.InvalidGroupName, "")
	if err == nil {
		t.Fatal("No error was returned. Expected: Group " + testutils.InvalidGroupName + " not found in configuration")
	}

	if !strings.Contains(err.Error(), "not found in configuration") {
		t.Fatalf("Unexpected error: %s", err)
	}
}

// Returns an error if the specified cluster doesn't exist
func TestBuildPathsWrongCluster(t *testing.T) {
	configPath := testutils.GetTestFile("utils", testGroupsFile)
	_, err := BuildPaths(configPath, testutils.InvalidGroupName, testutils.InvalidClusterName)
	if err == nil {
		t.Fatal("No error was returned. Expected: Cluster " + testutils.InvalidGroupName + " not found in configuration")
	}

	if !strings.Contains(err.Error(), "not found in configuration") {
		t.Fatalf("Unexpected error: %s", err)
	}
}
