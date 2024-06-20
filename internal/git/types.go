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
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
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

type FilesGetter struct {
	fs           billy.Filesystem
	storage      *memory.Storage
	clonePackage func(billy.Filesystem, storage.Storer, v1alpha1.Package) (billy.Filesystem, error)
}

// NewFilesGetter create a new FilesGetter instance configured for downloading from remote repository using
// an in memory storage
func NewFilesGetter() *FilesGetter {
	return &FilesGetter{
		fs:      memfs.New(),
		storage: memory.NewStorage(),
		clonePackage: func(fs billy.Filesystem, storage storage.Storer, pkg v1alpha1.Package) (billy.Filesystem, error) {
			cloneOptions := cloneOptionsForPackage(pkg)
			if _, err := git.Clone(storage, fs, cloneOptions); err != nil {
				return nil, fmt.Errorf("error cloning repository %w", err)
			}

			return fs, nil
		},
	}
}

// GetFilesForPackage clones the pkg from the remote repository and return all the files relative for the package
// or an error otherwise
func (r *FilesGetter) GetFilesForPackage(pkg v1alpha1.Package) ([]*File, error) {
	memFs, err := r.clonePackage(r.fs, r.storage, pkg)
	if err != nil {
		return nil, err
	}

	return filterWorktreeForPackage(memFs, pkg)
}
