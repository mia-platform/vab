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

package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mia-platform/vab/internal/testutils"
	"github.com/stretchr/testify/assert"
)

const (
	testGroupsFile = "test_groups.yaml"
)

// Test that the correct path is returned given valid group and cluster
func TestGetClusterPath(t *testing.T) {
	configPath := testutils.GetTestFile("utils", testGroupsFile)
	buildPath, err := BuildPaths(configPath, testutils.TestGroupName1, testutils.TestClusterName1)
	if assert.NoError(t, err) {
		expectedPath := filepath.Join(ClustersDirName, testutils.TestGroupName1, testutils.TestClusterName1)
		assert.Equal(t, buildPath[0], expectedPath)
	}
}

// Test that the correct paths are returned given valid group
func TestGetGroupPath(t *testing.T) {
	configPath := testutils.GetTestFile("utils", testGroupsFile)
	buildPaths, err := BuildPaths(configPath, testutils.TestGroupName1, "")
	assert.Nil(t, err, err)
	if assert.NoError(t, err) {
		clusterPath1 := filepath.Join(ClustersDirName, testutils.TestGroupName1, testutils.TestClusterName1)
		clusterPath2 := filepath.Join(ClustersDirName, testutils.TestGroupName1, testutils.TestClusterName2)
		expectedPaths := []string{clusterPath1, clusterPath2}
		assert.Equal(t, buildPaths, expectedPaths, "Unexpected paths. Expected: %v, actual: %v", expectedPaths, buildPaths)
	}
}

// Returns an error if the specified group doesn't exist
func TestBuildPathsWrongGroup(t *testing.T) {
	configPath := testutils.GetTestFile("utils", testGroupsFile)
	_, err := BuildPaths(configPath, testutils.InvalidGroupName, "")
	if assert.Error(t, err, "Expected: Group "+testutils.InvalidGroupName+" not found in configuration") {
		assert.Contains(t, err.Error(), "not found in configuration", "Unexpected error: %s", err)
	}
}

// Returns an error if the specified cluster doesn't exist
func TestBuildPathsWrongCluster(t *testing.T) {
	configPath := testutils.GetTestFile("utils", testGroupsFile)
	_, err := BuildPaths(configPath, testutils.InvalidGroupName, testutils.InvalidClusterName)
	if assert.Error(t, err, "Expected: Cluster "+testutils.InvalidGroupName+" not found in configuration") {
		assert.Contains(t, err.Error(), "not found in configuration", "Unexpected error: %s", err)
	}
}

// ValidatePath creates the correct path if missing
func TestValidateMissingPath(t *testing.T) {
	testDirPath := t.TempDir()
	expectedPath := filepath.Join(testDirPath, "dir", "another_dir")
	err := ValidatePath(expectedPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.DirExists(t, expectedPath, "The path does not exist")
}

// ValidatePath returns no error if the path exists
func TestValidateExistingPath(t *testing.T) {
	testDirPath := t.TempDir()
	expectedPath := filepath.Join(testDirPath, "dir", "another_dir")
	if err := os.MkdirAll(expectedPath, os.ModePerm); err != nil {
		return
	}
	assert.DirExists(t, expectedPath, "The path does not exist") // ensure the dir exists before calling ValidatePath
	err := ValidatePath(expectedPath)
	if !assert.NoError(t, err) {
		return
	}
}
