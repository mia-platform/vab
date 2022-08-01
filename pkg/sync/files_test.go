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
	tempWorktree := testutils.PrepareTempWorktree(t)

	input := []*git.File{
		git.NewFile("modules/test-module1/test-flavour1/file1.yaml", "./modules/test-module1", tempWorktree),
		git.NewFile("modules/test-module1/test-flavour1/file2.yaml", "./modules/test-module1", tempWorktree),
		git.NewFile("modules/test-module1/test-flavour2/file1.yaml", "./modules/test-module1", tempWorktree),
		git.NewFile("modules/test-module2/test-flavour1/file1.yaml", "./modules/test-module2", tempWorktree),
	}

	targetPath := path.Join(tempWorktree.Root(), "test-module1/test-flavour1")

	err := Readwrite(input, targetPath)
	assert.NoError(t, err)

	compareFile(t, []byte("file1-1 content\n"), "test-flavour1/file1.yaml")
	compareFile(t, []byte("file1-2 content\n"), "test-flavour1/file2.yaml")
	compareFile(t, []byte("file2-1 content\n"), "test-flavour2/file1.yaml")

	dirList, err := tempWorktree.ReadDir(path.Join(tempWorktree.Root(), "test-module1/test-flavour1"))
	assert.NoError(t, err)
	assert.Equal(t, 3, len(dirList))
}
