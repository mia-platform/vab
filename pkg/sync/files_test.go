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

package sync

import (
	"os"
	"path"
	"testing"

	"github.com/mia-platform/vab/internal/git"
	"github.com/mia-platform/vab/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func compareFile(t *testing.T, fileContent []byte, filePath string) {
	t.Helper()
	f, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, fileContent, f)
}

func TestReadWrite(t *testing.T) {
	fakeWorktree := testutils.PrepareFakeWorktree(t)

	input := []*git.File{
		git.NewFile(path.Join(testModuleBasePath, testFileName1), testModuleBasePath, fakeWorktree),
		git.NewFile(path.Join(testModuleBasePath, testFileName2), testModuleBasePath, fakeWorktree),
		git.NewFile(path.Join(testModuleBasePath, testFileName3), testModuleBasePath, fakeWorktree),
	}

	tempdir := t.TempDir()

	err := WritePkgToDir(input, tempdir)
	assert.NoError(t, err)

	compareFile(t, []byte("file1-1-1 content\n"), path.Join(tempdir, testFileName1))
	compareFile(t, []byte("file1-1-2 content\n"), path.Join(tempdir, testFileName2))
	compareFile(t, []byte("file1-2-1 content\n"), path.Join(tempdir, testFileName3))

	dirList, err := os.ReadDir(path.Join(tempdir, "test-flavour1/"))
	assert.NoError(t, err)
	assert.Equal(t, 2, len(dirList))

	dirList, err = os.ReadDir(path.Join(tempdir, "test-flavour2/"))
	assert.NoError(t, err)
	assert.Equal(t, 1, len(dirList))
}
