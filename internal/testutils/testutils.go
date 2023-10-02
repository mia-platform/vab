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

package testutils

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/stretchr/testify/assert"
)

const (
	// Invalid resources names
	InvalidFolderPath  = "/invalid/path"
	InvalidFileName    = "invalid.yaml"
	InvalidGroupName   = "invalid-group"
	InvalidClusterName = "invalid-cluster"

	// Valid resources names
	TestGroupName1       = "test-group"
	TestGroupName2       = "test-group2"
	TestClusterName1     = "test-cluster"
	TestClusterName2     = "test-cluster2"
	KustomizeTestDirName = "kustomize-test"
)

func GetTestFile(module string, args ...string) string {
	combinedElements := append([]string{
		"..",
		"..",
		"tests",
		module,
	},
		args...,
	)
	return filepath.Join(combinedElements...)
}

func populateWorktree(t *testing.T, fsys billy.Filesystem) {
	t.Helper()
	err := fsys.MkdirAll("modules/category/test-module1/test-flavor1", fs.ModePerm)
	assert.NoError(t, err)
	err = fsys.MkdirAll("modules/category/test-module1/test-flavor2", fs.ModePerm)
	assert.NoError(t, err)
	err = fsys.MkdirAll("modules/category/test-module2/test-flavor1", fs.ModePerm)
	assert.NoError(t, err)
	err = fsys.MkdirAll("modules/category/test-module3/test-flavor1", fs.ModePerm)
	assert.NoError(t, err)
	err = fsys.MkdirAll("addons/category/test-addon1/subdir", fs.ModePerm)
	assert.NoError(t, err)
	err = fsys.MkdirAll("addons/category/test-addon2/", fs.ModePerm)
	assert.NoError(t, err)
	err = fsys.MkdirAll("otherdir", fs.ModePerm)
	assert.NoError(t, err)
	_, err = fsys.Create("README.md")
	assert.NoError(t, err)
	f, err := fsys.Create("modules/category/test-module1/test-flavor1/file1.yaml")
	assert.NoError(t, err)
	_, err = f.Write([]byte("file1-1-1 content\n"))
	assert.NoError(t, err)
	err = f.Close()
	assert.NoError(t, err)
	f, err = fsys.Create("modules/category/test-module1/test-flavor1/file2.yaml")
	assert.NoError(t, err)
	_, err = f.Write([]byte("file1-1-2 content\n"))
	assert.NoError(t, err)
	err = f.Close()
	assert.NoError(t, err)
	f, err = fsys.Create("modules/category/test-module1/test-flavor2/file1.yaml")
	assert.NoError(t, err)
	_, err = f.Write([]byte("file1-2-1 content\n"))
	assert.NoError(t, err)
	err = f.Close()
	assert.NoError(t, err)
	_, err = fsys.Create("modules/category/test-module2/test-flavor1/file1.yaml")
	assert.NoError(t, err)
	_, err = fsys.Create("addons/category/test-addon1/file1.yaml")
	assert.NoError(t, err)
	_, err = fsys.Create("addons/category/test-addon1/subdir/file1.yaml")
	assert.NoError(t, err)
	_, err = fsys.Create("addons/category/test-addon2/file1.yaml")
	assert.NoError(t, err)
}

func PrepareWorktree(t *testing.T, fsType string) *billy.Filesystem {
	t.Helper()
	var worktree billy.Filesystem
	switch fsType {
	case "osfs":
		worktree = osfs.New(t.TempDir())
	case "memfs":
		worktree = memfs.New()
	default:
		assert.FailNow(t, "fstype not recognized")
	}
	populateWorktree(t, worktree)
	if !assert.NotNil(t, worktree) {
		t.FailNow()
	}
	return &worktree
}

func PrepareFakeWorktree(t *testing.T) *billy.Filesystem {
	t.Helper()
	return PrepareWorktree(t, "memfs")
}

type FakeFilesGetter struct {
	Testing *testing.T
}

func (filesGetter FakeFilesGetter) WorkTreeForPackage(_ v1alpha1.Package) (*billy.Filesystem, error) {
	return PrepareFakeWorktree(filesGetter.Testing), nil
}

func CompareFile(t *testing.T, fileContent []byte, filePath string) {
	t.Helper()
	f, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, fileContent, f)
}
