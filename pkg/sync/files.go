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

package sync

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mia-platform/vab/internal/git"
)

// WritePkgToDir writes the files in memory to the target path on disk
func WritePkgToDir(files []*git.File, targetPath string) error {
	for _, gitFile := range files {
		err := os.MkdirAll(filepath.Dir(filepath.Join(targetPath, gitFile.FilePath())), os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory: %s : %w", filepath.Dir(gitFile.FilePath()), err)
		}

		err = gitFile.Open()
		if err != nil {
			return fmt.Errorf("error opening file: %s : %w", gitFile.String(), err)
		}
		outFile, err := os.Create(filepath.Join(targetPath, gitFile.FilePath()))
		if err != nil {
			return fmt.Errorf("error opening file: %s : %w", filepath.Join(targetPath, gitFile.FilePath()), err)
		}

		r := bufio.NewReader(gitFile)
		w := bufio.NewWriter(outFile)

		_, err = r.WriteTo(w)
		if err != nil {
			return fmt.Errorf("error writing: %s : %w", outFile.Name(), err)
		}

		err = gitFile.Close()
		if err != nil {
			return fmt.Errorf("error closing: %s : %w", gitFile.String(), err)
		}

		err = outFile.Close()
		if err != nil {
			return fmt.Errorf("error closing: %s : %w", outFile.Name(), err)
		}
	}
	return nil
}
