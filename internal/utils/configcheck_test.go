package utils

import (
	"bytes"
	"path"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mia-platform/vab/internal/logger"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"golang.org/x/exp/slices"
)

const (
	testGroup        = "test-group"
	testCluster1     = "test-cluster"
	testCluster2     = "another-cluster"
	wrongResource    = "wrong-group"
	testData         = "test_data"
	testGroupsFile   = "test_groups.yaml"
	invalidYamlFile  = "invalid_yaml.yaml"
	testValidateFile = "test_validate.yaml"
	expectedOutput1  = `Reading the configuration...
[warn][default] no module found: check the config file if this behavior is unexpected
[warn][default] no add-on found: check the config file if this behavior is unexpected
[warn] no group found: check the config file if this behavior is unexpected
The configuration is valid!
`
	expectedOutput2 = `[error] wrong kind: WrongKind - expected: ClustersConfiguration
[error] wrong version: wrong.version.io/v1 - expected: vab.mia-platform.eu/v1alpha1
`
	expectedOutput3 = `Reading the configuration...
[error][default] missing version of module module-1/flavor-1
[warn][default] missing weight of module module-1/flavor-1: setting default (0)
[info][default] disabling module module-2/flavor-2
[error][default] missing version of add-on addon-1
[info][default] disabling add-on addon-2
[error] please specify a valid name for each group
[error][undefined] missing cluster name in group: please specify a valid name for each cluster
[error][undefined/undefined] missing cluster context: please specify a valid context for each cluster
[warn][undefined/undefined] no module found: check the config file if this behavior is unexpected
[warn][undefined/undefined] no add-on found: check the config file if this behavior is unexpected
[warn][undefined/cluster-1] no module found: check the config file if this behavior is unexpected
[warn][undefined/cluster-1] no add-on found: check the config file if this behavior is unexpected
[warn][group-1] no cluster found in group: check the config file if this behavior is unexpected
The configuration is invalid.
`
)

// Test that the correct path is returned given valid group and cluster
func TestGetClusterPath(t *testing.T) {
	args := []string{testGroup, testCluster1}
	configPath := path.Join("..", testData, testGroupsFile)
	buildPath, err := GetBuildPath(args, configPath)
	if err != nil {
		t.Fatal(err)
	}
	expectedPath := path.Join(testGroup, testCluster1)
	if buildPath[0] != expectedPath {
		t.Fatalf("Unexpected path. Expected: %s, actual: %s", expectedPath, buildPath[0])
	}
}

// Test that the correct paths are returned given valid group
func TestGetGroupPath(t *testing.T) {
	args := []string{testGroup}
	configPath := path.Join("..", testData, testGroupsFile)
	buildPaths, err := GetBuildPath(args, configPath)
	if err != nil {
		t.Fatal(err)
	}
	clusterPath1 := path.Join(testGroup, testCluster1)
	clusterPath2 := path.Join(testGroup, testCluster2)
	expectedPaths := []string{clusterPath1, clusterPath2}
	if !cmp.Equal(buildPaths, expectedPaths) {
		t.Fatalf("Unexpected paths. Expected: %v, actual: %v", expectedPaths, buildPaths)
	}
}

// Returns an error if the specified group doesn't exist
func TestGetBuildPathWrongGroup(t *testing.T) {
	args := []string{wrongResource}
	configPath := path.Join("..", testData, testGroupsFile)
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
	configPath := path.Join("..", testData, testGroupsFile)
	_, err := GetBuildPath(args, configPath)
	if err == nil {
		t.Fatal("No error was returned. Expected: Cluster " + args[0] + " not found in configuration")
	}
	if !strings.Contains(err.Error(), "not found in configuration") {
		t.Fatalf("Unexpected error: %s", err)
	}
}

// Test parsing error returned from ReadConfig
func TestValidateParseError(t *testing.T) {
	targetPath := path.Join("..", testData, invalidYamlFile)
	buffer := new(bytes.Buffer)
	errBuffer := new(bytes.Buffer)
	logger := logger.NewLogger(logger.LogStreams{OutStream: buffer, ErrStream: errBuffer})

	code := ValidateConfig(logger, targetPath)
	if code != 1 {
		t.Fatalf("Unexpected exit code: %d", code)
	}
	if !strings.Contains(errBuffer.String(), "yaml") {
		t.Fatalf("Unexpected output: %s", buffer.String())
	}
}

// Test validation of valid empty config
func TestValidateEmptySpec(t *testing.T) {
	targetPath := path.Join("..", testData, emptyConfigFile)
	buffer := new(bytes.Buffer)
	errBuffer := new(bytes.Buffer)
	logger := logger.NewLogger(logger.LogStreams{OutStream: buffer, ErrStream: errBuffer})
	code := ValidateConfig(logger, targetPath)
	if code != 0 {
		t.Fatalf("Unexpected exit code: %d", code)
	}
	if !bytes.Equal(buffer.Bytes(), []byte(expectedOutput1)) {
		t.Fatalf("Unexpected output: %s", buffer.String())
	}
}

// Test validation of wrong Kind/APIVersion
func TestCheckTypeMeta(t *testing.T) {
	invalidTypeMeta := v1alpha1.TypeMeta{
		Kind:       "WrongKind",
		APIVersion: "wrong.version.io/v1",
	}
	buffer := new(bytes.Buffer)
	errBuffer := new(bytes.Buffer)
	logger := logger.NewLogger(logger.LogStreams{OutStream: buffer, ErrStream: errBuffer})
	code := 0
	checkTypeMeta(logger, &invalidTypeMeta, &code)

	if code != 1 {
		t.Fatalf("Unexpected exit code: %d", code)
	}
	if !bytes.Equal(buffer.Bytes(), []byte(expectedOutput2)) {
		t.Fatalf("Unexpected output: %s", buffer.String())
	}
}

// Test validate with ad-hoc invalid file for max coverage
func TestValidateOutput(t *testing.T) {
	targetPath := path.Join("..", testData, testValidateFile)
	buffer := new(bytes.Buffer)
	errBuffer := new(bytes.Buffer)
	logger := logger.NewLogger(logger.LogStreams{OutStream: buffer, ErrStream: errBuffer})
	code := ValidateConfig(logger, targetPath)
	if code != 1 {
		t.Fatalf("Unexpected exit code: %d", code)
	}
	loggedLines := strings.Split(buffer.String(), "\n")
	expectedOutput3Array := strings.Split(expectedOutput3, "\n")
	if len(loggedLines) != len(expectedOutput3Array) {
		t.Fatalf("Expected %d log lines, %d founded", len(expectedOutput3Array), len(loggedLines))
	}

	for _, line := range loggedLines {
		if !slices.Contains(expectedOutput3Array, line) {
			t.Fatalf("Unexpected logged line %s", line)
		}
	}
}
