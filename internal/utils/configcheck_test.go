package utils

import (
	"path"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const (
	testGroup     = "test-group"
	testCluster1  = "test-cluster"
	testCluster2  = "another-cluster"
	wrongResource = "wrong-group"
)

// Test that the correct path is returned given valid group and cluster
func TestGetClusterPath(t *testing.T) {
	args := []string{testGroup, testCluster1}
	configPath := path.Join("..", "test_data", "test_groups.yaml")
	buildPath, err := GetBuildPath(args, configPath)
	if err != nil {
		t.Fatal(err)
	}
	expectedPath := path.Join(clustersDirName, testGroup, testCluster1)
	if buildPath[0] != expectedPath {
		t.Fatalf("Unexpected path. Expected: %s, actual: %s", expectedPath, buildPath[0])
	}
}

// Test that the correct paths are returned given valid group
func TestGetGroupPath(t *testing.T) {
	args := []string{testGroup}
	configPath := path.Join("..", "test_data", "test_groups.yaml")
	buildPaths, err := GetBuildPath(args, configPath)
	if err != nil {
		t.Fatal(err)
	}
	clusterPath1 := path.Join(clustersDirName, testGroup, testCluster1)
	clusterPath2 := path.Join(clustersDirName, testGroup, testCluster2)
	expectedPaths := []string{clusterPath1, clusterPath2}
	if !cmp.Equal(buildPaths, expectedPaths) {
		t.Fatalf("Unexpected paths. Expected: %v, actual: %v", expectedPaths, buildPaths)
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
