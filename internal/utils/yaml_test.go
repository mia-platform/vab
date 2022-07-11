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

const (
	testConfigName = "empty-test"
)

// Test marshalling of config struct
func TestWriteEmptyConfig(t *testing.T) {
	testDirPath := t.TempDir()
	t.Log(testDirPath)

	emptyConfig := v1alpha1.EmptyConfig(testConfigName)

	if err := WriteConfig(emptyConfig, testDirPath); err != nil {
		t.Fatal(err)
	}

	testFileContent, _ := os.ReadFile(path.Join(testDirPath, defaultConfigFileName))
	expectedFileContent, _ := os.ReadFile(path.Join("..", "test_data", "empty_config.yaml"))

	if !bytes.Equal(testFileContent, expectedFileContent) {
		t.Fatal("Unexpected file content.")
	}
}

// Test generation of configuration with custom file name
func TestCustomConfigName(t *testing.T) {
	testDirPath := t.TempDir()
	fileName := "custom_config.yaml"
	filePath := path.Join(testDirPath, fileName)

	emptyConfig := v1alpha1.EmptyConfig(testConfigName)

	if err := WriteConfig(emptyConfig, filePath); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filePath); err != nil {
		t.Fatal(err)
	}
}

// Test that the correct error is returned if the path to the config file is invalid
func TestPathNotExists(t *testing.T) {
	testWrongPath := "/wrong/path/to/config.yaml"

	emptyConfig := v1alpha1.EmptyConfig(testConfigName)
	err := WriteConfig(emptyConfig, testWrongPath)

	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrNotExist)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrNotExist, err)
	}
}

// Test that the correct error is returned if vab does not have permissions to access config path
func TestPathPermError(t *testing.T) {
	testDirPath := t.TempDir()
	if err := os.Chmod(testDirPath, 0); err != nil {
		t.Fatal(err)
	}

	emptyConfig := v1alpha1.EmptyConfig(testConfigName)
	err := WriteConfig(emptyConfig, testDirPath)

	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrPermission)
	}
	if !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrPermission, err)
	}
}

// Test that no error is returned if the config file already exists
func TestEmptyExistingConfig(t *testing.T) {
	testDirPath := t.TempDir()

	if writeErr := os.WriteFile(path.Join(testDirPath, "config.yaml"), []byte{}, defaultFilePermissions); writeErr != nil {
		t.Fatal(writeErr)
	}

	emptyConfig := v1alpha1.EmptyConfig(testConfigName)
	err := WriteConfig(emptyConfig, testDirPath)

	if err != nil {
		t.Fatal(err)
	}
}

// Test marshalling of Kustomization struct
func TestWriteEmptyKustomization(t *testing.T) {
	testDirPath := t.TempDir()
	t.Log(testDirPath)

	if err := WriteKustomization(*emptyKustomization, testDirPath); err != nil {
		t.Fatal(err)
	}

	testFileContent, _ := os.ReadFile(path.Join(testDirPath, kustomizationFileName))
	expectedFileContent, _ := os.ReadFile(path.Join("..", "test_data", "empty_kustomization.yaml"))

	if !bytes.Equal(testFileContent, expectedFileContent) {
		t.Fatal("Unexpected file content.")
	}
}

// Test that the correct error is returned if the file is not named kustomization.yaml
func TestWrongKustomizationFileName(t *testing.T) {
	testDirPath := t.TempDir()
	file, fileErr := os.Create(path.Join(testDirPath, "notkustomization.yaml"))
	if fileErr != nil {
		t.Fatalf("Error while creating test file: %s", fileErr)
	}

	err := WriteKustomization(*emptyKustomization, path.Join(testDirPath, file.Name()))

	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", errKustomizationTarget)
	}
	if !errors.Is(err, errKustomizationTarget) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", errKustomizationTarget, err)
	}
}

// Test that the correct error is returned if the path to the Kustomization file is invalid
func TestKustomizationPathNotExists(t *testing.T) {
	testWrongPath := "/wrong/path/to/kustomization.yaml"

	err := WriteKustomization(*emptyKustomization, testWrongPath)

	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrNotExist)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrNotExist, err)
	}
}

// Test that the correct error is returned if vab does not have permissions to access config path
func TestKustomizationPathPermError(t *testing.T) {
	testDirPath := t.TempDir()
	if err := os.Chmod(testDirPath, 0); err != nil {
		t.Fatal(err)
	}

	err := WriteKustomization(*emptyKustomization, testDirPath)

	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrPermission)
	}
	if !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrPermission, err)
	}
}

// Test that no error is returned if the Kustomization file already exists
func TestEmptyExistingKustomization(t *testing.T) {
	testDirPath := t.TempDir()
	if writeErr := os.WriteFile(path.Join(testDirPath, "kustomization.yaml"), []byte{}, defaultFilePermissions); writeErr != nil {
		t.Fatal(writeErr)
	}

	err := WriteKustomization(*emptyKustomization, testDirPath)

	if err != nil {
		t.Fatal(err)
	}
}
