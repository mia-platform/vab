package utils

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path"
	"testing"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
)

// Test marshalling of config struct
func TestWriteEmptyConfig(t *testing.T) {
	testDirPath := t.TempDir()

	emptyConfig := &v1alpha1.ClustersConfiguration{

		TypeMeta: v1alpha1.TypeMeta{
			Kind:       "ClustersConfiguration",
			APIVersion: "vab.mia-platform.eu/v1alpha1",
		},
		Name: "empty-test",
		Spec: v1alpha1.ConfigSpec{
			Modules: make(map[string]v1alpha1.Module),
			AddOns:  make(map[string]v1alpha1.AddOn),
			Groups:  make([]v1alpha1.Group, 0),
		},
	}

	WriteConfig(*emptyConfig, testDirPath)

	testFileContent, _ := os.ReadFile(path.Join(testDirPath, defaultConfigFileName))
	expectedFileContent, _ := os.ReadFile(path.Join("..", "test_data", "empty.yaml"))

	if !bytes.Equal(testFileContent, expectedFileContent) {
		t.Fatal("unexpected file content")
	}
}

// Test generation of configuration with custom file name
func TestCustomConfigName(t *testing.T) {
	testDirPath := t.TempDir()
	fileName := "custom_config.yaml"
	filePath := path.Join(testDirPath, fileName)

	emptyConfig := &v1alpha1.ClustersConfiguration{}

	WriteConfig(*emptyConfig, filePath)

	if _, err := os.Stat(filePath); err != nil {
		t.Fatal(err)
	}
}

// Test that the correct error is returned if the path to the config file is invalid
func TestPathNotExists(t *testing.T) {
	testWrongPath := "/wrong/path/to/config.yaml"

	emptyConfig := &v1alpha1.ClustersConfiguration{}
	err := WriteConfig(*emptyConfig, testWrongPath)

	if err == nil {
		t.Fatalf("should return the following error: %s", fs.ErrNotExist)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("wrong error: %s", err)
	}
}

// Test that the correct error is returned if the path is only partially valid
func TestPartiallyValidPath(t *testing.T) {
	testDirPath := t.TempDir()
	testWrongPath := path.Join(testDirPath, "missing-foo", "config.yaml")

	emptyConfig := &v1alpha1.ClustersConfiguration{}
	err := WriteConfig(*emptyConfig, testWrongPath)

	if err == nil {
		t.Fatalf("should return the following error: %s", fs.ErrNotExist)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("wrong error: %s", err)
	}
}

// Test that the correct error is returned if vab does not have permissions to access config path
func TestPathPermError(t *testing.T) {
	testDirPath := t.TempDir()
	os.Chmod(testDirPath, 0)

	emptyConfig := &v1alpha1.ClustersConfiguration{}
	err := WriteConfig(*emptyConfig, testDirPath)

	if err == nil {
		t.Fatalf("should return the following error: %s", fs.ErrPermission)
	}
	if !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("wrong error: %s", err)
	}
}

// Test that no error is returned if the config file already exists
func TestEmptyExistingFile(t *testing.T) {
	testDirPath := t.TempDir()
	if writeErr := os.WriteFile(path.Join(testDirPath, "config.yaml"), []byte{}, defaultFilePermissions); writeErr != nil {
		t.Fatal(writeErr)
	}

	emptyConfig := &v1alpha1.ClustersConfiguration{}
	err := WriteConfig(*emptyConfig, testDirPath)

	if err != nil {
		t.Fatal(err)
	}
}
