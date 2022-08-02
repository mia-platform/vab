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
	"fmt"
	"os"
	"path"

	"github.com/mia-platform/vab/internal/git"
)

// Lista di File, targetPath ( dentro a module/add-on ), version = crea tutte le cartelle
// in caso di errore cancella tutta la dir

func Readwrite(files []*git.File, targetPath string) error {
	for _, gitFile := range files {
		outFile, err := os.Open(path.Join(targetPath, gitFile.FilePath()))
		if err != nil {
			return fmt.Errorf("error opering file: %s : %w", gitFile.String(), err)
		}

		// TODO - Read and write file portions
		var buffer []byte
		_, err = gitFile.Read(buffer)
		if err != nil {
			return fmt.Errorf("error reading file: %s : %w", gitFile.String(), err)
		}

		_, err = outFile.Write(buffer)
		if err != nil {
			return fmt.Errorf("error writing file: %s : %w", gitFile.String(), err)
		}

		outFile.Close()
		if err != nil {
			return fmt.Errorf("error closing file: %s : %w", gitFile.String(), err)
		}
	}
	return nil
}
