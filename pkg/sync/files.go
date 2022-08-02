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
	"bufio"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/mia-platform/vab/internal/git"
)

func Readwrite(files []*git.File, targetPath string) error {
	for _, gitFile := range files {
		fmt.Printf("filepath: %s\n", gitFile.FilePath())

		err := os.MkdirAll(path.Dir(path.Join(targetPath, gitFile.FilePath())), os.ModePerm)
		if err != nil {
			return fmt.Errorf("error creating directory: %s : %w", path.Dir(gitFile.FilePath()), err)
		}

		inFile, err := gitFile.Open()
		if err != nil {
			return fmt.Errorf("error opering file: %s : %w", inFile.Name(), err)
		}
		outFile, err := os.Create(path.Join(targetPath, gitFile.FilePath()))
		if err != nil {
			return fmt.Errorf("error opering file: %s : %w", path.Join(targetPath, gitFile.FilePath()), err)
		}

		r := bufio.NewReader(inFile)
		w := bufio.NewWriter(outFile)

		buf := make([]byte, 1024)
		for {
			n, err := r.Read(buf)
			if err != nil && err != io.EOF {
				return fmt.Errorf("error reading file: %s : %w", inFile.Name(), err)
			}
			if n == 0 {
				break
			}

			if _, err := w.Write(buf[:n]); err != nil {
				return fmt.Errorf("error writing: %s : %w", outFile.Name(), err)
			}
		}

		if err = w.Flush(); err != nil {
			return fmt.Errorf("error flushing file: %s : %w", outFile.Name(), err)
		}

		err = inFile.Close()
		if err != nil {
			return fmt.Errorf("error closing: %s : %w", inFile.Name(), err)
		}

		err = outFile.Close()
		if err != nil {
			return fmt.Errorf("error closing: %s : %w", outFile.Name(), err)
		}

	}
	return nil
}
