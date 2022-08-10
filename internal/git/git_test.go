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
	"os"
	"testing"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/mia-platform/vab/internal/testutils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/mia-platform/vab/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestRemoteURL(t *testing.T) {
	expectedURL := "https://github.com/mia-platform/distribution"

	// module := v1alpha1.Module{Version: "1.2.3", Weight: 1, Disable: false}
	assert.Equal(t, remoteURL(), expectedURL)
	assert.Equal(t, remoteURL(), expectedURL)

	// addon := v1alpha1.AddOn{Version: "1.2.3", Disable: false}
	assert.Equal(t, remoteURL(), expectedURL)
	assert.Equal(t, remoteURL(), expectedURL)
}

func TestGetAuths(t *testing.T) {
	assert.Nil(t, remoteAuth())
	assert.Nil(t, remoteAuth())
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

	assert.Equal(t, options.URL, remoteURL())
	assert.Nil(t, options.Auth)
	assert.Equal(t, options.ReferenceName, tagReferenceForPackage(addonName, addon))

	module := v1alpha1.Module{Version: "1.0.0", Weight: 10, Disable: false}
	moduleName := "module-name/flavor-name"
	options = cloneOptionsForPackage(moduleName, module)

	assert.Equal(t, options.URL, remoteURL())
	assert.Nil(t, options.Auth)
	assert.Equal(t, options.ReferenceName, tagReferenceForPackage(moduleName, module))
}

func TestFilterFilesForPackage(t *testing.T) {
	fakeWorktree := testutils.PrepareFakeWorktree(t)

	logger := logger.DisabledLogger{}
	t.Run("filter module files", func(t *testing.T) {
		moduleName := "test-module1/test-flavour1"
		module := v1alpha1.Module{
			Version: "1.0.0",
		}

		expectedArray := []*File{
			NewFile("modules/test-module1/test-flavour1/file1.yaml", "./modules/test-module1", *fakeWorktree),
			NewFile("modules/test-module1/test-flavour1/file2.yaml", "./modules/test-module1", *fakeWorktree),
			NewFile("modules/test-module1/test-flavour2/file1.yaml", "./modules/test-module1", *fakeWorktree),
		}
		files, err := filterWorktreeForPackage(logger, fakeWorktree, moduleName, module)
		assert.NoError(t, err)
		assert.Equal(t, files, expectedArray)
	})

	t.Run("filter addon files", func(t *testing.T) {
		addonName := "test-addon1"
		addon := v1alpha1.AddOn{
			Version: "1.0.0",
		}

		expectedArray := []*File{
			NewFile("add-ons/test-addon1/file1.yaml", "./add-ons/test-addon1", *fakeWorktree),
			NewFile("add-ons/test-addon1/subdir/file1.yaml", "./add-ons/test-addon1", *fakeWorktree),
		}
		files, err := filterWorktreeForPackage(logger, fakeWorktree, addonName, addon)
		assert.NoError(t, err)
		assert.Equal(t, files, expectedArray)
	})
}

func TestFilterError(t *testing.T) {
	fakeWorktree := testutils.PrepareFakeWorktree(t)

	logger := logger.DisabledLogger{}
	addonName := "test-addon4"
	addon := v1alpha1.AddOn{
		Version: "1.0.0",
	}

	files, err := filterWorktreeForPackage(logger, fakeWorktree, addonName, addon)
	assert.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.Nil(t, files)
}

func TestGetFilesForPackage(t *testing.T) {
	logger := logger.DisabledLogger{}
	moduleName := "test-module1/test-flavour1"
	module := v1alpha1.Module{
		Version: "1.0.0",
	}

	files, err := GetFilesForPackage(logger, testutils.FakeFilesGetter{Testing: t}, moduleName, module)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, files, "Nil output file references")
}
