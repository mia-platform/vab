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
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
	"sigs.k8s.io/kustomize/api/konfig"
	kustomize "sigs.k8s.io/kustomize/api/types"
)

const (
	defaultFilePermissions = 0644
	yamlDefaultIndent      = 2
)

type WrongFileNameError struct {
	expectedFileName string
	actualFileName   string
}

func NewWrongFileNameError(expected string, actual string) error {
	return WrongFileNameError{
		expectedFileName: expected,
		actualFileName:   actual,
	}
}

func (e WrongFileNameError) Error() string {
	return "expected file name " + e.expectedFileName + " but found " + e.actualFileName
}

// ReadConfig reads a configuration file into a ClustersConfiguration struct
func ReadConfig(configPath string) (*v1alpha1.ClustersConfiguration, error) {
	configFile, readErr := os.ReadFile(configPath)
	if readErr != nil {
		return nil, readErr
	}

	output := &v1alpha1.ClustersConfiguration{}
	yamlErr := yaml.Unmarshal(configFile, output)
	if yamlErr != nil {
		return nil, yamlErr
	}

	return output, nil
}

// writeYamlFile marshals the interface passed as argument, and writes it to a
// YAML file
func writeYamlFile(file interface{}, dstPath string) error {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(yamlDefaultIndent)

	if err := yamlEncoder.Encode(&file); err != nil {
		return err
	}

	return os.WriteFile(dstPath, b.Bytes(), defaultFilePermissions)
}

// WriteConfig creates and writes an empty vab configuration file
func WriteConfig(config v1alpha1.ClustersConfiguration, dirOrFilePath string) error {
	dirOrFile, err := os.Stat(dirOrFilePath)

	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	var dstPath string
	if err == nil && dirOrFile.IsDir() {
		dstPath = filepath.Join(dirOrFilePath, DefaultConfigFilename)
	} else {
		dstPath = dirOrFilePath
	}

	return writeYamlFile(config, dstPath)
}

// WriteKustomization creates and writes an empty kustomization file
func WriteKustomization(kustomization kustomize.Kustomization, dirOrFilePath string, override bool) error {
	dirOrFile, err := os.Stat(dirOrFilePath)

	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	var dstPath string
	switch {
	case err == nil && dirOrFile.IsDir():
		dstPath = filepath.Join(dirOrFilePath, konfig.DefaultKustomizationFileName())
	case !slices.Contains(konfig.RecognizedKustomizationFileNames(), filepath.Base(dirOrFilePath)):
		return NewWrongFileNameError(konfig.DefaultKustomizationFileName(), filepath.Base(dirOrFilePath))
	default:
		dstPath = dirOrFilePath
	}

	if !override {
		if _, err := os.Stat(dstPath); !errors.Is(err, fs.ErrNotExist) {
			return err
		}
	}

	return writeYamlFile(kustomization, dstPath)
}

// EmptyKustomization return a valid empty kustomization with valid kind and apiVersion fields
func EmptyKustomization() kustomize.Kustomization {
	// mini hack for generating a valid kustomization structure as kustomize intend
	empty := kustomize.Kustomization{}
	empty.FixKustomizationPostUnmarshalling()
	return empty
}

// EmptyKustomization return a valid empty kustomization with valid kind and apiVersion fields
func EmptyComponent() kustomize.Kustomization {
	// mini hack for generating a valid kustomization structure as kustomize intend
	empty := kustomize.Kustomization{}
	empty.Kind = kustomize.ComponentKind
	empty.FixKustomizationPostUnmarshalling()
	return empty
}
