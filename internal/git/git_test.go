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
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"gotest.tools/assert"
)

func TestUrlStringFromModule(t *testing.T) {
	expectedURL := "https://github.com/mia-platform/distribution"
	module := v1alpha1.Module{Version: "1.2.3", Weight: 1, Disable: false}

	url := urlForModule(module)
	if url != expectedURL {
		t.Fatalf("Unexpected url: %s", url)
	}

	module = v1alpha1.Module{}
	url = urlForModule(module)
	if url != expectedURL {
		t.Fatalf("Unexpected url for empty module: %s", url)
	}
}

func TestUrlStringFromAddOn(t *testing.T) {
	expectedURL := "https://github.com/mia-platform/distribution"
	addon := v1alpha1.AddOn{Version: "1.2.3", Disable: false}

	url := urlForAddon(addon)
	if url != expectedURL {
		t.Fatalf("Unexpected url: %s", url)
	}

	addon = v1alpha1.AddOn{}
	url = urlForAddon(addon)
	if url != expectedURL {
		t.Fatalf("Unexpected url for empty add-on: %s", url)
	}
}

func TestGetAuths(t *testing.T) {
	addonAuth := authForAddon(v1alpha1.AddOn{})
	if addonAuth != nil {
		t.Fatalf("Unexpected auth configuration %s", addonAuth)
	}

	moduleAuth := authForModule(v1alpha1.Module{})
	if moduleAuth != nil {
		t.Fatalf("Unexpected auth configuration %s", addonAuth)
	}
}

func TestTagReferences(t *testing.T) {
	addonName := "addon-name/with-slash"
	addonVersion := "1.0.0"
	expectedReference := "refs/tags/addon-" + addonName + "-" + addonVersion
	tag := tagReferenceForAddon(addonName, addonVersion)
	if tag != plumbing.ReferenceName(expectedReference) {
		t.Fatalf("Unexpected addon tag reference %s, expected %s", tag, expectedReference)
	}
	if !tag.IsTag() {
		t.Fatalf("The addon reference %s is not a tag reference", tag)
	}

	moduleName := "module-name/flavor"
	moduleVersion := "1.0.0"
	expectedReference = "refs/tags/module-module-name-" + addonVersion
	tag = tagReferenceForModule(moduleName, moduleVersion)
	if tag != plumbing.ReferenceName(expectedReference) {
		t.Fatalf("Unexpected module tag reference %s, expected %s", tag, expectedReference)
	}
	if !tag.IsTag() {
		t.Fatalf("The module reference %s is not a tag reference", tag)
	}
}

func TestCloneOptions(t *testing.T) {
	addon := v1alpha1.AddOn{Version: "1.0.0", Disable: false}
	addonName := "addon-name"
	options := cloneOptionsForAddon(addonName, addon)

	if options.URL != urlForAddon(addon) {
		t.Fatalf("Unexpected URL for addon %s: %s", addonName, options.URL)
	}
	if options.Auth != nil {
		t.Fatalf("Unexpected Auth for addon %s: %s", addonName, options.Auth)
	}
	if options.ReferenceName != tagReferenceForAddon(addonName, addon.Version) {
		t.Fatalf("Unexpected reference name for addon %s: %s", addonName, options.ReferenceName)
	}
	if !options.ReferenceName.IsTag() {
		t.Fatalf("Reference created for addon %s is not a branch: %s", addonName, options.ReferenceName)
	}

	module := v1alpha1.Module{Version: "1.0.0", Weight: 10, Disable: false}
	moduleName := "module-name/flavor-name"
	options = cloneOptionsForModule(moduleName, module)

	if options.URL != urlForModule(module) {
		t.Fatalf("Unexpected URL for module %s: %s", moduleName, options.URL)
	}
	if options.Auth != nil {
		t.Fatalf("Unexpected Auth for module %s: %s", moduleName, options.Auth)
	}
	if options.ReferenceName != tagReferenceForModule(moduleName, addon.Version) {
		t.Fatalf("Unexpected reference name for module %s: %s", moduleName, options.ReferenceName)
	}
	if !options.ReferenceName.IsTag() {
		t.Fatalf("Reference created for module %s is not a branch: %s", moduleName, options.ReferenceName)
	}
}

func TestProvaClone(t *testing.T) {
	mem := memfs.New()
	memStorage := memory.NewStorage()
	_, err := git.Clone(memStorage, mem, &git.CloneOptions{
		URL:           "https://github.com/go-git/go-git.git",
		ReferenceName: plumbing.NewTagReferenceName("v5.1.0"),
		Depth:         1,
		Tags:          git.NoTags,
		SingleBranch:  true,
	})
	assert.NilError(t, err)
	fileInfos, err := mem.ReadDir("/")
	assert.NilError(t, err, "readdir")
	for _, v := range fileInfos {
		t.Logf("%s\n", v.Name())
	}
}

type fakeCloner struct {
	t *testing.T
}

func createDirs(t *testing.T, workTree billy.Filesystem, dirs []string) {
	t.Helper()
	for _, dir := range dirs {
		err := workTree.MkdirAll(dir, fs.FileMode(0755))
		assert.NilError(t, err, dir)
	}
}

func createFiles(t *testing.T, workTree billy.Filesystem, files []string) {
	t.Helper()
	for _, file := range files {
		_, err := workTree.Create(file)
		assert.NilError(t, err, file)
	}
}

func (c fakeCloner) Clone(addonName string, addon v1alpha1.AddOn, cloneOptions *git.CloneOptions) (billy.Filesystem, error) {
	workTree := memfs.New()

	createDirs(c.t, workTree, []string{
		"/addon1/subdir1",
		"/addon1/subdir2",
		"/addon2/subdir1",
	})

	createFiles(c.t, workTree, []string{
		"/addon1/file1",
		"/addon1/subdir1/file2",
		"/addon2/subdir1/file3",
	})
	return workTree, nil
}

func TestCloneAddon(t *testing.T) {
	addon := v1alpha1.AddOn{
		Version: "1.0.0",
	}
	fc := fakeCloner{t: t}
	outFs, err := cloneAddon("addon1", addon, nil, fc)
	assert.NilError(t, err)
	expectedFs := memfs.New()
	createDirs(t, expectedFs, []string{
		"/addon1/subdir1",
		"/addon1/subdir2",
	})
	createFiles(t, expectedFs, []string{
		"/addon1/file1",
		"/addon1/subdir1/file2",
	})

	outFiles, err := outFs.ReadDir("/")
	assert.NilError(t, err)
	expectedFiles, err := expectedFs.ReadDir("/")
	assert.NilError(t, err)
	assert.DeepEqual(t, outFiles, expectedFiles)
	// fs.WalkDir(outFs, "/", func(path string, d fs.DirEntry, err error) error {})
}
