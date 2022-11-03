//go:build !e2e
// +build !e2e

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
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/mia-platform/vab/internal/testutils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/kustomize/api/konfig"
)

const (
	testConfigName  = "empty-test"
	emptyConfigFile = "empty_config.yaml"
)

// Test marshalling of config struct
func TestWriteEmptyConfig(t *testing.T) {
	testDirPath := t.TempDir()

	emptyConfig := v1alpha1.EmptyConfig(testConfigName)

	if err := WriteConfig(*emptyConfig, testDirPath); assert.NoError(t, err) {
		testFileContent, _ := os.ReadFile(filepath.Join(testDirPath, DefaultConfigFilename))
		expectedFileContent, _ := os.ReadFile(testutils.GetTestFile("utils", emptyConfigFile))
		assert.Equal(t, testFileContent, expectedFileContent, "Unexpected file content.")
	}
}

// Test generation of configuration with custom file name
func TestCustomConfigName(t *testing.T) {
	testDirPath := t.TempDir()
	fileName := "custom_config.yaml"
	filePath := filepath.Join(testDirPath, fileName)

	emptyConfig := v1alpha1.EmptyConfig(testConfigName)

	if err := WriteConfig(*emptyConfig, filePath); assert.NoError(t, err) {
		_, err = os.Stat(filePath)
		assert.NoError(t, err)
	}
}

// Test that the correct error is returned if the path to the config file is invalid
func TestPathNotExists(t *testing.T) {
	testWrongPath := "/wrong/path/to/config.yaml"

	emptyConfig := v1alpha1.EmptyConfig(testConfigName)
	err := WriteConfig(*emptyConfig, testWrongPath)

	if assert.Error(t, err, "Expected: %s", fs.ErrNotExist) {
		assert.ErrorIs(t, err, fs.ErrNotExist)
	}
}

// Test that the correct error is returned if vab does not have permissions to access config path
func TestPathPermError(t *testing.T) {
	testDirPath := t.TempDir()
	if err := os.Chmod(testDirPath, 0); !assert.NoError(t, err) {
		return
	}

	emptyConfig := v1alpha1.EmptyConfig(testConfigName)
	err := WriteConfig(*emptyConfig, testDirPath)

	if assert.Error(t, err, "Expected: %s", fs.ErrPermission) {
		assert.ErrorIs(t, err, fs.ErrPermission)
	}
}

// Test that no error is returned if the config file already exists
func TestEmptyExistingConfig(t *testing.T) {
	testDirPath := t.TempDir()
	writeErr := os.WriteFile(filepath.Join(testDirPath, "config.yaml"), []byte{}, defaultFilePermissions)

	if assert.NoError(t, writeErr) {
		emptyConfig := v1alpha1.EmptyConfig(testConfigName)
		err := WriteConfig(*emptyConfig, testDirPath)
		assert.NoError(t, err)
	}
}

// Test marshalling of Kustomization struct
func TestWriteEmptyKustomization(t *testing.T) {
	testDirPath := t.TempDir()

	if err := WriteKustomization(EmptyKustomization(), testDirPath); assert.NoError(t, err, err) {
		testFileContent, _ := os.ReadFile(filepath.Join(testDirPath, konfig.DefaultKustomizationFileName()))
		expectedFileContent, _ := os.ReadFile(testutils.GetTestFile("utils", "empty_kustomization.yaml"))
		assert.Equal(t, testFileContent, expectedFileContent, "Unexpected file content.")
	}
}

// Test that the correct error is returned if the file is not named kustomization.yaml
func TestWrongKustomizationFileName(t *testing.T) {
	testDirPath := t.TempDir()
	file, fileErr := os.Create(filepath.Join(testDirPath, "notkustomization.yaml"))
	if !assert.NoError(t, fileErr, "Error while creating test file") {
		return
	}

	expectedError := NewWrongFileNameError(konfig.DefaultKustomizationFileName(), filepath.Base(file.Name()))
	err := WriteKustomization(EmptyKustomization(), file.Name())

	if assert.Error(t, err, "Expected: %s", expectedError) {
		assert.ErrorAs(t, err, &WrongFileNameError{})
	}
}

// Test that the correct error is returned if the path to the Kustomization file is invalid
func TestKustomizationPathNotExists(t *testing.T) {
	testWrongPath := "/wrong/path/to/kustomization.yaml"

	err := WriteKustomization(EmptyKustomization(), testWrongPath)

	if assert.Error(t, err, "Expected: %s", fs.ErrNotExist) {
		assert.ErrorIs(t, err, fs.ErrNotExist)
	}
}

// Test that the correct error is returned if vab does not have permissions to access config path
func TestKustomizationPathPermError(t *testing.T) {
	testDirPath := t.TempDir()
	if err := os.Chmod(testDirPath, 0); !assert.NoError(t, err) {
		return
	}

	err := WriteKustomization(EmptyKustomization(), testDirPath)
	if assert.Error(t, err, "Expected: %s", fs.ErrPermission) {
		assert.ErrorIs(t, err, fs.ErrPermission)
	}
}

// Test that no error is returned if the Kustomization file already exists
func TestEmptyExistingKustomization(t *testing.T) {
	testDirPath := t.TempDir()
	writeErr := os.WriteFile(filepath.Join(testDirPath, "kustomization.yaml"), []byte{}, defaultFilePermissions)
	if assert.NoError(t, writeErr) {
		return
	}

	err := WriteKustomization(EmptyKustomization(), testDirPath)
	assert.NoError(t, err)
}

// ReadConfig reads the configuration correctly
func TestReadEmptyConfig(t *testing.T) {
	config, err := ReadConfig(testutils.GetTestFile("utils", emptyConfigFile))
	if assert.NoError(t, err) {
		expectedConfig := v1alpha1.EmptyConfig(testConfigName)
		assert.Equal(t, config, expectedConfig, "Unexpected configuration.")
	}
}

// ReadConfig returns ErrNotExist if the path is invalid
func TestReadConfigInvalidPath(t *testing.T) {
	_, err := ReadConfig(testutils.InvalidFolderPath)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrNotExist)
	}
}

// ReadConfig returns ErrPermission if the path is not accessible
func TestReadConfigErrPermission(t *testing.T) {
	testDirPath := t.TempDir()
	dstPath := filepath.Join(testDirPath, "foo")
	if err := os.Mkdir(dstPath, 0); !assert.NoError(t, err) {
		return
	}

	_, err := ReadConfig(dstPath)
	if assert.Error(t, err) {
		assert.ErrorIs(t, err, fs.ErrPermission)
	}
}

// ReadConfig returns an error if the YAML is not invalid
func TestReadConfigUnmarshalErr(t *testing.T) {
	invalidConfigPath := testutils.GetTestFile("utils", "invalid_yaml.yaml")
	_, err := ReadConfig(invalidConfigPath)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "yaml")
	}
}
