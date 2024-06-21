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
	"bufio"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
)

// File rappresent a file downloaded in the in memory store
type File struct {
	path         string
	internalPath string
	fs           billy.Filesystem
}

// WriteContent copy the file content to targetPath mantaining the folder structure
func (f *File) WriteContent(targetPath string) error {
	onDiskPath := filepath.Join(targetPath, f.path)

	if err := os.MkdirAll(filepath.Dir(onDiskPath), os.ModePerm); err != nil {
		return err
	}

	file, err := f.fs.Open(f.internalPath)
	if err != nil {
		_ = os.Remove(onDiskPath)
		return err
	}
	defer file.Close()

	outFile, err := os.Create(onDiskPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	r := bufio.NewReader(file)
	w := bufio.NewWriter(outFile)

	_, err = r.WriteTo(w)
	return err
}
