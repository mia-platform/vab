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

package sync

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/mia-platform/vab/internal/git"
	"github.com/mia-platform/vab/pkg/cmd/util"
	"github.com/stretchr/testify/assert"
)

func TestCommand(t *testing.T) {
	t.Parallel()

	configFlags := util.NewConfigFlags()

	cmd := NewCommand(configFlags)
	assert.NotNil(t, cmd)
}

func TestToOptions(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()
	testCases := map[string]struct {
		args             []string
		downloadPackages bool
		configPath       string
		expectedOptions  *Options
		expectedError    string
	}{
		"invalid context path": {
			args:          []string{filepath.Join("invalid", "path")},
			expectedError: "no such file or directory",
		},
		"return options": {
			args:             []string{tempDir},
			configPath:       "custom.yaml",
			downloadPackages: true,
			expectedOptions: &Options{
				contextPath:      tempDir,
				downloadPackages: true,
				configPath:       "custom.yaml",
			},
		},
		"no config path": {
			args: []string{tempDir},
			expectedOptions: &Options{
				contextPath:      tempDir,
				downloadPackages: false,
				configPath:       "",
			},
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			flags := Flags{
				downloadPackages: testCase.downloadPackages,
			}
			configFlags := util.NewConfigFlags()
			configFlags.ConfigPath = &testCase.configPath

			options, err := flags.ToOptions(configFlags, testCase.args)
			switch len(testCase.expectedError) {
			case 0:
				assert.NoError(t, err)
				assert.NotNil(t, options.filesGetter)
				options.filesGetter = nil
				assert.Equal(t, testCase.expectedOptions, options)
			default:
				assert.ErrorContains(t, err, testCase.expectedError)
			}
		})
	}
}

func TestRun(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join("testdata", "config.yaml")
	tests := map[string]struct {
		options       *Options
		expectedError string
		expectedPaths []string
	}{
		"clone packages": {
			options: &Options{
				configPath:       configPath,
				contextPath:      t.TempDir(),
				downloadPackages: true,
				filesGetter: func() *git.FilesGetter {
					fg, _ := git.NewTestFilesGetter(t)
					return fg
				}(),
			},
			expectedPaths: append(folderStruct, vendorStruct...),
		},
		"don't clone packages": {
			options: &Options{
				configPath:       configPath,
				contextPath:      t.TempDir(),
				downloadPackages: false,
			},
			expectedPaths: folderStruct,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := test.options.Run(t.Context())
			if len(test.expectedError) > 0 {
				assert.ErrorContains(t, err, test.expectedError)
				return
			}

			assert.NoError(t, err)

			err = fs.WalkDir(os.DirFS(test.options.contextPath), ".", func(path string, _ fs.DirEntry, err error) error {
				assert.Contains(t, test.expectedPaths, path)
				return err
			})

			assert.NoError(t, err)
		})
	}
}

var (
	folderStruct = []string{
		".",
		"clusters",
		"clusters/all-groups",
		"clusters/all-groups/bases",
		"clusters/all-groups/bases/kustomization.yaml",
		"clusters/all-groups/custom-resources",
		"clusters/all-groups/custom-resources/kustomization.yaml",
		"clusters/all-groups/kustomization.yaml",
		"clusters/group",
		"clusters/group/cluster",
		"clusters/group/cluster/bases",
		"clusters/group/cluster/bases/kustomization.yaml",
		"clusters/group/cluster/custom-resources",
		"clusters/group/cluster/custom-resources/kustomization.yaml",
		"clusters/group/cluster/kustomization.yaml",
	}

	vendorStruct = []string{
		"vendors",
		"vendors/addons",
		"vendors/addons/category",
		"vendors/addons/category/test-addon2-v1.0.0",
		"vendors/addons/category/test-addon2-v1.0.0/file1.yaml",
		"vendors/modules",
		"vendors/modules/category",
		"vendors/modules/category/test-module1-v1.0.0",
		"vendors/modules/category/test-module1-v1.0.0/test-flavor1",
		"vendors/modules/category/test-module1-v1.0.0/test-flavor1/file1.yaml",
		"vendors/modules/category/test-module1-v1.0.0/test-flavor1/file2.yaml",
		"vendors/modules/category/test-module1-v1.0.0/test-flavor1/file.json",
		"vendors/modules/category/test-module1-v1.0.0/test-flavor2",
		"vendors/modules/category/test-module1-v1.0.0/test-flavor2/file1.yaml",
	}
)
