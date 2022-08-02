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
	"path"
	"strings"

	"github.com/go-git/go-billy/v5"
)

type File struct {
	path       string
	baseFolder string
	fs         billy.Filesystem
	billy.File
}

func NewFile(path string, baseFolder string, fs billy.Filesystem) *File {
	return &File{
		path:       path,
		fs:         fs,
		baseFolder: baseFolder,
	}
}

func (f *File) Open() (billy.File, error) {
	return f.fs.Open(f.path)
}

func (f *File) FilePath() string {
	return path.Clean(path.Join(".", strings.TrimPrefix(f.path, f.baseFolder)))
}

func (f *File) String() string {
	return f.path
}
