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

package util

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/kustomize/api/konfig"
	kustomize "sigs.k8s.io/kustomize/api/types"
)

func TestInitializeConfiguration(t *testing.T) {
	t.Parallel()

	testDirPath := t.TempDir()
	err := InitializeConfiguration("test", testDirPath)
	if !assert.NoError(t, err) {
		return
	}

	testStructure(t, testDirPath, filepath.Join("testdata", "sync", "init"))
}

func TestSyncDirectories(t *testing.T) {
	t.Parallel()

	testdata := filepath.Join("testdata", "sync")
	tests := map[string]struct {
		config             v1alpha1.ConfigSpec
		path               string
		expectedResultPath string
	}{
		"sync empty project": {
			config: v1alpha1.ConfigSpec{
				Modules: map[string]v1alpha1.Package{
					"test/module/base":    v1alpha1.NewModule(t, "test/module/base", "v1.28.0", false),
					"test/module2/flavor": v1alpha1.NewModule(t, "test/module2/flavor", "v1.28.0", false),
				},
				AddOns: map[string]v1alpha1.Package{
					"test/addon":  v1alpha1.NewAddon(t, "test/addon", "v1.0.0", false),
					"test/addon2": v1alpha1.NewAddon(t, "test/addon2", "v1.5.0", false),
					"test/addon3": v1alpha1.NewAddon(t, "test/addon2", "v1.5.0", true),
				},
				Groups: []v1alpha1.Group{
					{
						Name: "group1",
						Clusters: []v1alpha1.Cluster{
							{
								Name: "cluster",
							},
						},
					},
					{
						Name: "group2",
						Clusters: []v1alpha1.Cluster{
							{
								Name: "cluster",
								Modules: map[string]v1alpha1.Package{
									"test/module2/flavor":  v1alpha1.NewModule(t, "test/module2/flavor", "", true),
									"test/module2/flavor2": v1alpha1.NewModule(t, "test/module2/flavor2", "v1.28.0", false),
								},
								AddOns: map[string]v1alpha1.Package{
									"test/addon": v1alpha1.NewAddon(t, "test/addon", "", true),
								},
							},
						},
					},
				},
			},
			path:               t.TempDir(),
			expectedResultPath: filepath.Join(testdata, "empty"),
		},
		// "sync project with old config": {
		// 	config:             v1alpha1.ConfigSpec{},
		// 	path:               t.TempDir(),
		// 	expectedResultPath: filepath.Join(testdata, "old"),
		// },
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := SyncDirectories(test.config, test.path)
			require.NoError(t, err)
			testStructure(t, test.path, test.expectedResultPath)
		})
	}
}

func testStructure(t *testing.T, pathToTest, expectationPath string) {
	t.Helper()

	_ = filepath.WalkDir(pathToTest, func(path string, d fs.DirEntry, err error) error {
		require.NoError(t, err)
		cleanPath := strings.TrimPrefix(path, pathToTest)

		testPath := filepath.Join(expectationPath, cleanPath)
		if d.IsDir() {
			require.DirExists(t, testPath)
		} else {
			data, err := os.ReadFile(path)
			require.NoError(t, err)
			expectedData, err := os.ReadFile(testPath)
			require.NoError(t, err)
			assert.Equal(t, string(expectedData), string(data), testPath)
		}

		return nil
	})
}

func TestReadConfig(t *testing.T) {
	t.Parallel()

	testdata := "testdata"
	tempDir := t.TempDir()
	tests := map[string]struct {
		configPath     string
		expectedConfig *v1alpha1.ClustersConfiguration
		expectedError  string
	}{
		"empty config": {
			configPath:     filepath.Join(testdata, "empty.yaml"),
			expectedConfig: v1alpha1.EmptyConfig("empty-test"),
		},
		"read config": {
			configPath: filepath.Join(testdata, "config.yaml"),
			expectedConfig: &v1alpha1.ClustersConfiguration{
				TypeMeta: v1alpha1.TypeMeta{
					Kind:       v1alpha1.Kind,
					APIVersion: v1alpha1.Version,
				},
				Name: "test",
				Spec: v1alpha1.ConfigSpec{
					Modules: make(map[string]v1alpha1.Package),
					AddOns:  make(map[string]v1alpha1.Package),
					Groups: []v1alpha1.Group{
						{
							Name: "test-group",
							Clusters: []v1alpha1.Cluster{
								{
									Name:    "test-cluster",
									Modules: make(map[string]v1alpha1.Package),
									AddOns:  make(map[string]v1alpha1.Package),
								},
							},
						},
					},
				},
			},
		},
		"invalid path": {
			configPath:    filepath.Join(tempDir, "missing.yaml"),
			expectedError: fmt.Sprintf("open %s", filepath.Join(tempDir, "missing.yaml")),
		},
		"invalid yaml": {
			configPath:    filepath.Join(testdata, "invalid.yaml"),
			expectedError: "could not find expected ':'",
		},
		"empty path would use default path": {
			configPath:    "",
			expectedError: fmt.Sprintf("open %s", defaultConfigFileName),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			config, err := ReadConfig(test.configPath)
			if len(test.expectedError) > 0 {
				assert.ErrorContains(t, err, test.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.expectedConfig, config)
		})
	}
}

func TestWriteFile(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		path          string
		data          interface{}
		expectedError string
	}{
		"empty configuration": {
			path: filepath.Join(t.TempDir(), "emptyfile.yaml"),
			data: v1alpha1.EmptyConfig("test"),
		},
		"data with top comment": {
			path: filepath.Join(t.TempDir(), konfig.DefaultKustomizationFileName()),
			data: &kustomize.Kustomization{
				TypeMeta: kustomize.TypeMeta{
					Kind:       kustomize.ComponentKind,
					APIVersion: kustomize.ComponentVersion,
				},
			},
		},
		"path don't exists return error": {
			path:          filepath.Join("invalid", "path"),
			data:          v1alpha1.EmptyConfig("test"),
			expectedError: "no such file or directory",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := writeYamlFile(test.path, test.data)
			if len(test.expectedError) > 0 {
				assert.ErrorContains(t, err, test.expectedError)
				assert.NoFileExists(t, test.path)
				return
			}
			assert.NoError(t, err)
			assert.FileExists(t, test.path)
		})
	}
}
