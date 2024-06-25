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
	"io/fs"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/stretchr/testify/assert"
)

// NewTestFilesGetter return a FilesGetter with a fixed worktree and will not make calls to remote repositories
func NewTestFilesGetter(t *testing.T) (*FilesGetter, billy.Filesystem) {
	t.Helper()

	fg := NewFilesGetter()
	f := memfs.New()
	fg.clonePackage = func(_ v1alpha1.Package) (billy.Filesystem, error) {
		t.Helper()

		populateWorktree(t, f)
		return f, nil
	}

	return fg, f
}

func populateWorktree(t *testing.T, fsys billy.Filesystem) {
	t.Helper()
	assert.NoError(t, fsys.MkdirAll("modules/category/test-module1/test-flavor1", fs.ModePerm))
	assert.NoError(t, fsys.MkdirAll("modules/category/test-module1/test-flavor2", fs.ModePerm))
	assert.NoError(t, fsys.MkdirAll("modules/category/test-module2/test-flavor1", fs.ModePerm))
	assert.NoError(t, fsys.MkdirAll("modules/category/test-module3/test-flavor1", fs.ModePerm))
	assert.NoError(t, fsys.MkdirAll("addons/category/test-addon1/subdir", fs.ModePerm))
	assert.NoError(t, fsys.MkdirAll("addons/category/test-addon2/", fs.ModePerm))
	assert.NoError(t, fsys.MkdirAll("otherdir", fs.ModePerm))

	_, err := fsys.Create("README.md")
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
