// Copyright 2022 Mia-Platform

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package git

import (
	"io/fs"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestRemoteURL(t *testing.T) {
	expectedURL := "https://github.com/mia-platform/distribution"

	module := v1alpha1.Module{Version: "1.2.3", Weight: 1, Disable: false}
	assert.Equal(t, remoteURL(module), expectedURL)
	assert.Equal(t, remoteURL(v1alpha1.Module{}), expectedURL)

	addon := v1alpha1.AddOn{Version: "1.2.3", Disable: false}
	assert.Equal(t, remoteURL(addon), expectedURL)
	assert.Equal(t, remoteURL(v1alpha1.AddOn{}), expectedURL)
}

func TestGetAuths(t *testing.T) {
	assert.Nil(t, remoteAuth(v1alpha1.AddOn{}))
	assert.Nil(t, remoteAuth(v1alpha1.Module{}))
}

func TestTagReferences(t *testing.T) {
	addonName := "addon-name"
	addonVersion := "1.0.0"
	expectedReference := "refs/tags/addon-" + addonName + "-" + addonVersion

	addon := v1alpha1.AddOn{
		Version: addonVersion,
		Disable: false,
	}
	tag := tagReferenceForPackage(addonName, addon)
	assert.Equal(t, tag, plumbing.ReferenceName(expectedReference))
	assert.True(t, tag.IsTag(), "The addon reference %s is not a tag reference", tag)

	moduleName := "module-name/flavor"
	moduleVersion := "1.0.0"
	expectedReference = "refs/tags/module-module-name-" + addonVersion
	module := v1alpha1.Module{
		Version: moduleVersion,
		Weight:  10,
		Disable: false,
	}

	tag = tagReferenceForPackage(moduleName, module)
	assert.Equal(t, tag, plumbing.ReferenceName(expectedReference))
	assert.True(t, tag.IsTag(), "The module reference %s is not a tag reference", tag)
}

func TestCloneOptions(t *testing.T) {
	addon := v1alpha1.AddOn{Version: "1.0.0", Disable: false}
	addonName := "addon-name"
	options := cloneOptionsForPackage(addonName, addon)

	assert.Equal(t, options.URL, remoteURL(addon))
	assert.Nil(t, options.Auth)
	assert.Equal(t, options.ReferenceName, tagReferenceForPackage(addonName, addon))

	module := v1alpha1.Module{Version: "1.0.0", Weight: 10, Disable: false}
	moduleName := "module-name/flavor-name"
	options = cloneOptionsForPackage(moduleName, module)

	assert.Equal(t, options.URL, remoteURL(module))
	assert.Nil(t, options.Auth)
	assert.Equal(t, options.ReferenceName, tagReferenceForPackage(moduleName, module))
}

func prepareFakeWorktree(t *testing.T) billy.Filesystem {
	t.Helper()

	var err error
	// Create a new repository
	worktree := memfs.New()

	err = worktree.MkdirAll("modules/test-module1/test-flavour1", fs.FileMode(0755))
	assert.NoError(t, err)
	err = worktree.MkdirAll("modules/test-module1/test-flavour2", fs.FileMode(0755))
	assert.NoError(t, err)
	err = worktree.MkdirAll("modules/test-module2/test-flavour1", fs.FileMode(0755))
	assert.NoError(t, err)
	err = worktree.MkdirAll("add-ons/test-addon1/subdir", fs.FileMode(0755))
	assert.NoError(t, err)
	err = worktree.MkdirAll("add-ons/test-addon2/", fs.FileMode(0755))
	assert.NoError(t, err)
	err = worktree.MkdirAll("otherdir", fs.FileMode(0755))
	assert.NoError(t, err)
	_, err = worktree.Create("README.md")
	assert.NoError(t, err)
	_, err = worktree.Create("modules/test-module1/test-flavour1/file1.yaml")
	assert.NoError(t, err)
	_, err = worktree.Create("modules/test-module1/test-flavour1/file2.yaml")
	assert.NoError(t, err)
	_, err = worktree.Create("modules/test-module1/test-flavour2/file1.yaml")
	assert.NoError(t, err)
	_, err = worktree.Create("add-ons/test-addon1/file1.yaml")
	assert.NoError(t, err)
	_, err = worktree.Create("add-ons/test-addon1/subdir/file1.yaml")
	assert.NoError(t, err)
	_, err = worktree.Create("add-ons/test-addon2/file1.yaml")
	assert.NoError(t, err)

	return worktree
}

func TestFilterFilesForPackage(t *testing.T) {
	fakeWorktree := prepareFakeWorktree(t)
	assert.NotNil(t, fakeWorktree)
}
