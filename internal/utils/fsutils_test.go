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
	"errors"
	"io/fs"
	"os"
	"path"
	"testing"
)

const (
	invalidPath = "/invalid/path"
	testName    = "foo"
)

// Return current path if the name arg is the empty string
func TestCurrentPath(t *testing.T) {
	testDirPath := t.TempDir()
	dstPath, err := GetProjectPath(testDirPath, "")

	if err != nil {
		t.Fatal(err)
	}
	if dstPath != testDirPath {
		t.Fatalf("Unexpected path. Expected: %s, actual: %s", testDirPath, dstPath)
	}
}

// Return path of a new project directory named "foo"
func TestNewProjectPath(t *testing.T) {
	testDirPath := t.TempDir()
	expectedPath := path.Join(testDirPath, testName)
	dstPath, err := GetProjectPath(testDirPath, testName)
	if err != nil {
		t.Fatal(err)
	}
	if dstPath != expectedPath {
		t.Fatalf("Unexpected path. Expected: %s, actual: %s", expectedPath, dstPath)
	}
}

// Return ErrNotExists if the path parameter is invalid (empty name)
func TestCurrentInvalidPath(t *testing.T) {
	_, err := GetProjectPath(invalidPath, "")
	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrNotExist)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrNotExist, err)
	}
}

// Return ErrNotExists if the path parameter is invalid (non-empty name)
func TestNewInvalidPath(t *testing.T) {
	_, err := GetProjectPath(invalidPath, testName)
	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrNotExist)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrNotExist, err)
	}
}

// Return ErrPermission if access to the path is denied (empty name)
func TestCurrentPathErrPermission(t *testing.T) {
	testDirPath := t.TempDir()
	if err := os.Chmod(testDirPath, 0); err != nil {
		t.Fatal(err)
	}
	_, err := GetProjectPath(testDirPath, "")
	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrPermission)
	}
	if !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrPermission, err)
	}
}

// Return ErrPermission if access to the path is denied (non-empty name)
func TestNewPathErrPermission(t *testing.T) {
	testDirPath := t.TempDir()
	if err := os.Chmod(testDirPath, 0); err != nil {
		t.Fatal(err)
	}
	_, err := GetProjectPath(testDirPath, testName)
	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrPermission)
	}
	if !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrPermission, err)
	}
}

// If a directory with the specified name already exists, return the path to that directory
func TestNewExistingPath(t *testing.T) {
	testDirPath := t.TempDir()
	expectedPath := path.Join(testDirPath, testName)
	if err := os.Mkdir(expectedPath, fs.ModePerm); err != nil {
		t.Fatal(err)
	}
	dstPath, err := GetProjectPath(testDirPath, testName)
	if err != nil {
		t.Fatal(err)
	}
	if dstPath != expectedPath {
		t.Fatalf("Unexpected path. Expected: %s, actual: %s", expectedPath, dstPath)
	}
}

// Test that the directory for the new cluster is created correctly with its empty kustomization.yaml
func TestNewClusterOverride(t *testing.T) {
	testDirPath := t.TempDir()
	if err := os.Mkdir(path.Join(testDirPath, ClustersDirName), fs.ModePerm); err != nil {
		t.Fatal(err)
	}
	if err := CreateClusterOverride(testDirPath, testName); err != nil {
		t.Fatal(err)
	}
	kustomizationPath := path.Join(testDirPath, ClustersDirName, testName, kustomizationFileName)
	if _, err := os.Stat(kustomizationPath); err != nil {
		t.Fatal(err)
	}
}

// Test that the clusters directory is created if missing
func TestMissingClustersDirectory(t *testing.T) {
	testDirPath := t.TempDir()
	if err := CreateClusterOverride(testDirPath, testName); err != nil {
		t.Fatal(err)
	}
}

// Return ErrPermission if if access to the path is denied
func TestNewClusterErrPermission(t *testing.T) {
	testDirPath := t.TempDir()
	if err := os.Mkdir(path.Join(testDirPath, ClustersDirName), 0); err != nil {
		t.Fatal(err)
	}
	err := CreateClusterOverride(testDirPath, testName)
	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrPermission)
	}
	if !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrPermission, err)
	}
}

// Return ErrPermission if if access to the path is denied
func TestMissingClustersDirErrPermission(t *testing.T) {
	testDirPath := t.TempDir()
	if err := os.Chmod(testDirPath, 0); err != nil {
		t.Fatal(err)
	}
	err := CreateClusterOverride(testDirPath, testName)
	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrPermission)
	}
	if !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrPermission, err)
	}
}
