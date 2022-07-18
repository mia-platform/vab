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
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mia-platform/vab/internal/testutils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
)

const (
	testConfigName  = "empty-test"
	emptyConfigFile = "empty_config.yaml"
)

// Test marshalling of config struct
func TestWriteEmptyConfig(t *testing.T) {
	testDirPath := t.TempDir()
	t.Log(testDirPath)

	emptyConfig := v1alpha1.EmptyConfig(testConfigName)

	if err := WriteConfig(*emptyConfig, testDirPath); err != nil {
		t.Fatal(err)
	}

	testFileContent, _ := os.ReadFile(path.Join(testDirPath, defaultConfigFileName))
	expectedFileContent, _ := os.ReadFile(testutils.GetTestFile("utils", emptyConfigFile))

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

	if err := WriteConfig(*emptyConfig, filePath); err != nil {
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
	err := WriteConfig(*emptyConfig, testWrongPath)

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
	err := WriteConfig(*emptyConfig, testDirPath)

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
	err := WriteConfig(*emptyConfig, testDirPath)

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
	expectedFileContent, _ := os.ReadFile(testutils.GetTestFile("utils", "empty_kustomization.yaml"))

	if !bytes.Equal(testFileContent, expectedFileContent) {
		t.Fatal("Unexpected file content.")
	}
}

// Test that the correct error is returned if the file is not named kustomization.yaml
func TestWrongKustomizationFileName(t *testing.T) {
	testDirPath := t.TempDir()
	file, fileErr := os.Create(testutils.GetTestFile("utils", "notkustomization.yaml"))
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

// ReadConfig reads the configuration correctly
func TestReadEmptyConfig(t *testing.T) {
	config, err := ReadConfig(testutils.GetTestFile("utils", emptyConfigFile))
	if err != nil {
		t.Fatal(err)
	}
	expectedConfig := v1alpha1.EmptyConfig(testConfigName)
	if !cmp.Equal(config, expectedConfig) {
		t.Fatal("Unexpected configuration.")
	}
}

// ReadConfig returns ErrNotExist if the path is invalid
func TestReadConfigInvalidPath(t *testing.T) {
	_, err := ReadConfig(invalidPath)
	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrNotExist)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrNotExist, err)
	}
}

// ReadConfig returns ErrPermission if the path is not accessible
func TestReadConfigErrPermission(t *testing.T) {
	testDirPath := t.TempDir()
	dstPath := path.Join(testDirPath, testName)
	if err := os.Mkdir(dstPath, 0); err != nil {
		t.Fatal(err)
	}
	_, err := ReadConfig(dstPath)
	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrPermission)
	}
	if !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrPermission, err)
	}
}

// ReadConfig returns an error if the YAML is not invalid
func TestReadConfigUnmarshalErr(t *testing.T) {
	invalidConfigPath := testutils.GetTestFile("utils", "invalid_yaml.yaml")
	_, err := ReadConfig(invalidConfigPath)
	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrNotExist)
	}
	if !strings.Contains(err.Error(), "yaml") {
		t.Fatalf("Unexpected error: %s", err)
	}
}
