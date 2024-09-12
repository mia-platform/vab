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

package git

import (
	"os"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestCloneOptions(t *testing.T) {
	t.Parallel()

	defaultGitURL := "https://github.com/mia-platform/distribution"
	tests := map[string]struct {
		pkgDefinition     v1alpha1.Package
		expectedURL       string
		expectedAuth      transport.AuthMethod
		expectedReference plumbing.ReferenceName
	}{
		"module": {
			pkgDefinition:     v1alpha1.NewModule(t, "category/module-name/flavor-name", "1.0.0", false),
			expectedURL:       defaultGitURL,
			expectedAuth:      nil,
			expectedReference: plumbing.NewTagReferenceName("module-category-module-name-1.0.0"),
		},
		"addon": {
			pkgDefinition:     v1alpha1.NewAddon(t, "category/addon-name", "1.0.0", false),
			expectedURL:       defaultGitURL,
			expectedAuth:      nil,
			expectedReference: plumbing.NewTagReferenceName("addon-category-addon-name-1.0.0"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			options := cloneOptionsForPackage(test.pkgDefinition)
			assert.Equal(t, test.expectedURL, options.URL)
			assert.Equal(t, test.expectedAuth, options.Auth)
			assert.Equal(t, test.expectedReference, options.ReferenceName)
		})
	}
}

func TestGetFiles(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		pkgDefinition v1alpha1.Package
		expectedFiles []*File
		expectedError string
	}{
		"filter module files": {
			pkgDefinition: v1alpha1.NewModule(t, "category/test-module1/test-flavor1", "1.0.0", false),
			expectedFiles: []*File{
				{
					path:         "test-flavor1/file.json",
					internalPath: "modules/category/test-module1/test-flavor1/file.json",
				},
				{
					path:         "test-flavor1/file1.yaml",
					internalPath: "modules/category/test-module1/test-flavor1/file1.yaml",
				},
				{
					path:         "test-flavor1/file2.yaml",
					internalPath: "modules/category/test-module1/test-flavor1/file2.yaml",
				},
				{
					path:         "test-flavor2/file1.yaml",
					internalPath: "modules/category/test-module1/test-flavor2/file1.yaml",
				},
			},
		},
		"filter addon files": {
			pkgDefinition: v1alpha1.NewAddon(t, "category/test-addon1", "1.0.0", false),
			expectedFiles: []*File{
				{
					path:         "file1.yaml",
					internalPath: "addons/category/test-addon1/file1.yaml",
				},
				{
					path:         "subdir/file1.yaml",
					internalPath: "addons/category/test-addon1/subdir/file1.yaml",
				},
			},
		},
		"missing package definition in downloaded files": {
			pkgDefinition: v1alpha1.NewAddon(t, "category/test-addon4", "1.0.0", false),
			expectedError: os.ErrNotExist.Error(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			fg, fs := NewTestFilesGetter(t)

			for _, file := range test.expectedFiles {
				file.fs = fs
			}

			files, err := fg.GetFilesForPackage(test.pkgDefinition)
			switch len(test.expectedError) {
			case 0:
				assert.NoError(t, err)
				assert.Equal(t, test.expectedFiles, files)
			default:
				assert.ErrorContains(t, err, test.expectedError)
				assert.Nil(t, files)
			}
		})
	}
}
