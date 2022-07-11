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
	dstPath, err := GetProjectRelativePath(testDirPath, "")

	if err != nil {
		t.Fatal(err)
	}
	if dstPath != testDirPath {
		t.Fatalf("Unexpected relative path. Expected: %s, actual: %s", testDirPath, dstPath)
	}
}

// Return relative path of a new project directory named "foo"
func TestNewProjectPath(t *testing.T) {
	testDirPath := t.TempDir()
	expectedPath := path.Join(testDirPath, testName)
	dstPath, err := GetProjectRelativePath(testDirPath, testName)
	if err != nil {
		t.Fatal(err)
	}
	if dstPath != expectedPath {
		t.Fatalf("Unexpected relative path. Expected: %s, actual: %s", expectedPath, dstPath)
	}
}

// Return ErrNotExists if the path parameter is invalid (empty name)
func TestCurrentInvalidPath(t *testing.T) {
	_, err := GetProjectRelativePath(invalidPath, "")
	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrNotExist)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrNotExist, err)
	}
}

// Return ErrNotExists if the path parameter is invalid (non-empty name)
func TestNewInvalidPath(t *testing.T) {
	_, err := GetProjectRelativePath(invalidPath, testName)
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
	_, err := GetProjectRelativePath(testDirPath, "")
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
	_, err := GetProjectRelativePath(testDirPath, testName)
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
	dstPath, err := GetProjectRelativePath(testDirPath, testName)
	if err != nil {
		t.Fatal(err)
	}
	if dstPath != expectedPath {
		t.Fatalf("Unexpected relative path. Expected: %s, actual: %s", expectedPath, dstPath)
	}
}

// Test that the directory for the new cluster is created correctly with its empty kustomization.yaml
func TestNewClusterOverride(t *testing.T) {
	testDirPath := t.TempDir()
	if err := os.Mkdir(path.Join(testDirPath, clustersDirName), fs.ModePerm); err != nil {
		t.Fatal(err)
	}
	if err := CreateClusterOverride(testDirPath, testName); err != nil {
		t.Fatal(err)
	}
	kustomizationPath := path.Join(testDirPath, clustersDirName, testName, kustomizationFileName)
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
	if err := os.Mkdir(path.Join(testDirPath, clustersDirName), 0); err != nil {
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
