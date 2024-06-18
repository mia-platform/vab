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
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestRemoteURL(t *testing.T) {
	expectedURL := "https://github.com/mia-platform/distribution"

	// module := v1alpha1.Package{Version: "1.2.3", Disable: false}
	assert.Equal(t, remoteURL(), expectedURL)
	assert.Equal(t, remoteURL(), expectedURL)

	// addon := v1alpha1.Package{Version: "1.2.3", Disable: false}
	assert.Equal(t, remoteURL(), expectedURL)
	assert.Equal(t, remoteURL(), expectedURL)
}

func TestGetAuths(t *testing.T) {
	assert.Nil(t, remoteAuth())
	assert.Nil(t, remoteAuth())
}

func TestTagReferences(t *testing.T) {
	expectedReference := "refs/tags/addon-category-addon-name-1.0.0"

	addon := v1alpha1.NewAddon(t, "category/addon-name", "1.0.0", false)
	tag := tagReferenceForPackage(addon)
	assert.Equal(t, tag, plumbing.ReferenceName(expectedReference))
	assert.True(t, tag.IsTag(), "The addon reference %s is not a tag reference", tag)

	expectedReference = "refs/tags/module-category-module-name-1.0.0"
	module := v1alpha1.NewModule(t, "category/module-name/flavor", "1.0.0", false)

	tag = tagReferenceForPackage(module)
	assert.Equal(t, plumbing.ReferenceName(expectedReference), tag)
	assert.True(t, tag.IsTag(), "The module reference %s is not a tag reference", tag)
}

func TestCloneOptions(t *testing.T) {
	addon := v1alpha1.NewAddon(t, "category/addon-name", "1.0.0", false)
	options := cloneOptionsForPackage(addon)

	assert.Equal(t, options.URL, remoteURL())
	assert.Nil(t, options.Auth)
	assert.Equal(t, options.ReferenceName, tagReferenceForPackage(addon))

	module := v1alpha1.NewModule(t, "category/module-name/flavor-name", "1.0.0", false)
	options = cloneOptionsForPackage(module)

	assert.Equal(t, options.URL, remoteURL())
	assert.Nil(t, options.Auth)
	assert.Equal(t, options.ReferenceName, tagReferenceForPackage(module))
}

func TestFilterFilesForPackage(t *testing.T) {
	fakeWorktree := prepareFakeWorktree(t)

	t.Run("filter module files", func(t *testing.T) {
		module := v1alpha1.NewModule(t, "category/test-module1/test-flavor1", "1.0.0", false)

		expectedArray := []*File{
			NewFile("modules/category/test-module1/test-flavor1/file1.yaml", "modules/category/test-module1", *fakeWorktree),
			NewFile("modules/category/test-module1/test-flavor1/file2.yaml", "modules/category/test-module1", *fakeWorktree),
			NewFile("modules/category/test-module1/test-flavor2/file1.yaml", "modules/category/test-module1", *fakeWorktree),
		}
		files, err := filterWorktreeForPackage(fakeWorktree, module)
		assert.NoError(t, err)
		assert.Equal(t, expectedArray, files)
	})

	t.Run("filter addon files", func(t *testing.T) {
		addon := v1alpha1.NewAddon(t, "category/test-addon1", "1.0.0", false)

		expectedArray := []*File{
			NewFile("addons/category/test-addon1/file1.yaml", "addons/category/test-addon1", *fakeWorktree),
			NewFile("addons/category/test-addon1/subdir/file1.yaml", "addons/category/test-addon1", *fakeWorktree),
		}
		files, err := filterWorktreeForPackage(fakeWorktree, addon)
		assert.NoError(t, err)
		assert.Equal(t, expectedArray, files)
	})
}

func TestFilterError(t *testing.T) {
	fakeWorktree := prepareFakeWorktree(t)

	addon := v1alpha1.NewAddon(t, "category/test-addon4", "1.0.0", false)

	files, err := filterWorktreeForPackage(fakeWorktree, addon)
	assert.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.Nil(t, files)
}

func TestGetFilesForPackage(t *testing.T) {
	module := v1alpha1.NewModule(t, "category/test-module1/test-flavor1", "1.0.0", false)

	files, err := GetFilesForPackage(fakeFilesGetter{Testing: t}, module)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, files, "Nil output file references")
}
