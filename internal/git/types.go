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

package git

import (
	"errors"
	"io"
	"path"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
)

type File struct {
	path       string
	baseFolder string
	fs         billy.Filesystem
	file       billy.File
	io.ReadCloser
}

func NewFile(path string, baseFolder string, fs billy.Filesystem) *File {
	return &File{
		path:       path,
		fs:         fs,
		baseFolder: baseFolder,
	}
}

func (f *File) Open() error {
	file, err := f.fs.Open(f.path)
	f.file = file
	return err
}

func (f *File) Read(p []byte) (n int, err error) {
	if f.file == nil {
		return 0, errors.New("error reading file: file is nil")
	}
	return f.file.Read(p)
}

func (f *File) Close() error {
	if f.file == nil {
		return nil
	}
	return f.file.Close()
}

func (f *File) FilePath() string {
	return path.Clean(path.Join(".", strings.TrimPrefix(f.path, f.baseFolder)))
}

func (f *File) String() string {
	return f.path
}

type FilesGetter interface {
	WorkTreeForPackage(pkgName string, pkg v1alpha1.Package) (*billy.Filesystem, error)
}

type RealFilesGetter struct{}

func (filesGetter RealFilesGetter) WorkTreeForPackage(pkgName string, pkg v1alpha1.Package) (*billy.Filesystem, error) {
	return worktreeForPackage(pkgName, pkg)
}
