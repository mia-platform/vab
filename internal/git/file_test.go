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
	"path/filepath"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteContent(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		targetPath    string
		filePath      string
		internalPath  string
		expectedError string
	}{
		"save file": {
			targetPath:   t.TempDir(),
			filePath:     "test-flavor2/file1.yaml",
			internalPath: "modules/category/test-module1/test-flavor2/file1.yaml",
		},
		"error creating file locally": {
			targetPath: func() string {
				dir := t.TempDir()
				require.NoError(t, os.Chmod(dir, 0666))
				return dir
			}(),
			filePath:      "README.md",
			internalPath:  "README.md",
			expectedError: os.ErrPermission.Error(),
		},
		"error finding file in memory": {
			targetPath:    t.TempDir(),
			filePath:      "test-flavor2/file1.yaml",
			internalPath:  "missing/test-module1/test-flavor2/file1.yaml",
			expectedError: os.ErrExist.Error(),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			fs := memfs.New()
			populateWorktree(t, fs)

			f := &File{
				path:         test.filePath,
				internalPath: test.internalPath,
				fs:           fs,
			}

			err := f.WriteContent(test.targetPath)
			switch len(test.expectedError) {
			case 0:
				assert.NoError(t, err)
				assert.FileExists(t, filepath.Join(test.targetPath, test.filePath))
			default:
				assert.ErrorContains(t, err, "")
				assert.NoFileExists(t, filepath.Join(test.targetPath, test.filePath))
			}
		})
	}
}
