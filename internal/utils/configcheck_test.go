package utils

import (
	"path"
	"strings"
	"testing"
)

const (
	testGroup     = "test-group"
	testCluster   = "test-cluster"
	wrongResource = "wrong-group"
)

// Test that the correct path is returned given valid group and cluster
func TestGetClusterPath(t *testing.T) {
	args := []string{testGroup, testCluster}
	configPath := path.Join("..", "test_data", "test_groups.yaml")
	configPath, err := GetBuildPath(args, configPath)
	if err != nil {
		t.Fatal(err)
	}
	expectedPath := path.Join(clustersDirName, testGroup, testCluster)
	if configPath != expectedPath {
		t.Fatalf("Unexpected path. Expected: %s, actual: %s", expectedPath, configPath)
	}
}

// Test that the correct path is returned given valid group
func TestGetGroupPath(t *testing.T) {
	args := []string{testGroup}
	configPath := path.Join("..", "test_data", "test_groups.yaml")
	configPath, err := GetBuildPath(args, configPath)
	if err != nil {
		t.Fatal(err)
	}
	expectedPath := path.Join(clustersDirName, testGroup)
	if configPath != expectedPath {
		t.Fatalf("Unexpected path. Expected: %s, actual: %s", expectedPath, configPath)
	}
}

// Returns an error if the specified group doesn't exist
func TestGetBuildPathWrongGroup(t *testing.T) {
	args := []string{wrongResource}
	configPath := path.Join("..", "test_data", "test_groups.yaml")
	_, err := GetBuildPath(args, configPath)
	if err == nil {
		t.Fatal("No error was returned. Expected: Group " + args[0] + " not found in configuration")
	}
	if !strings.Contains(err.Error(), "not found in configuration") {
		t.Fatalf("Unexpected error: %s", err)
	}
}

// Returns an error if the specified cluster doesn't exist
func TestGetBuildPathWrongCluster(t *testing.T) {
	args := []string{wrongResource, wrongResource}
	configPath := path.Join("..", "test_data", "test_groups.yaml")
	_, err := GetBuildPath(args, configPath)
	if err == nil {
		t.Fatal("No error was returned. Expected: Cluster " + args[0] + " not found in configuration")
	}
	if !strings.Contains(err.Error(), "not found in configuration") {
		t.Fatalf("Unexpected error: %s", err)
	}
}

// Returns an error if the args array is empty
func TestGetBuildPathNoArgs(t *testing.T) {
	args := []string{}
	configPath := path.Join("..", "test_data", "test_groups.yaml")
	_, err := GetBuildPath(args, configPath)
	if err == nil {
		t.Fatal("No error was returned. Expected: at least the cluster group is required")
	}
	if !strings.Contains(err.Error(), "required") {
		t.Fatalf("Unexpected error: %s", err)
	}
}

// Returns an error if the args are too many (> 2)
func TestGetBuildPathTooManyArgs(t *testing.T) {
	args := []string{testGroup, testCluster, wrongResource}
	configPath := path.Join("..", "test_data", "test_groups.yaml")
	_, err := GetBuildPath(args, configPath)
	if err == nil {
		t.Fatal("No error was returned. Expected: too many args")
	}
	if !strings.Contains(err.Error(), "too many") {
		t.Fatalf("Unexpected error: %s", err)
	}
}
