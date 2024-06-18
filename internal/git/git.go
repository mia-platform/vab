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
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	billyutil "github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
)

const (
	// defaultGitUrl is the default git remote
	defaultGitURL = "https://github.com"
	// defaultRepositoryUrl is the default repository to use if no other is specified
	defaultRepositoryURL = defaultGitURL + "/mia-platform/distribution"
)

// remoteUrl return the git url to use for downloading the files for a package (module or addon)
func remoteURL() string {
	return defaultRepositoryURL
}

// remoteAuth return an AuthMethod for package (module or addon)
func remoteAuth() transport.AuthMethod {
	return nil
}

// tagReferenceForPackage return a valid tag reference for the package name and version
func tagReferenceForPackage(pkg v1alpha1.Package) plumbing.ReferenceName {
	tag := pkg.PackageType() + "-" + strings.ReplaceAll(pkg.GetName(), "/", "-") + "-" + pkg.Version
	return plumbing.NewTagReferenceName(tag)
}

// cloneOptionsForPackage return the options for cloning the package with pkgName with pkg configuaration
func cloneOptionsForPackage(pkg v1alpha1.Package) *git.CloneOptions {
	return &git.CloneOptions{
		URL:           remoteURL(),
		Auth:          remoteAuth(),
		ReferenceName: tagReferenceForPackage(pkg),
		Depth:         1,
		SingleBranch:  true,
		Tags:          git.NoTags,
	}
}

// worktreeForPackage return a worktree from the cloned repository for the package with pkgName
func worktreeForPackage(pkg v1alpha1.Package) (*billy.Filesystem, error) {
	cloneOptions := cloneOptionsForPackage(pkg)
	fs := memfs.New()
	storage := memory.NewStorage()
	if _, err := git.Clone(storage, fs, cloneOptions); err != nil {
		return nil, fmt.Errorf("error cloning repository %w", err)
	}

	return &fs, nil
}

func filterWorktreeForPackage(worktree *billy.Filesystem, pkg v1alpha1.Package) ([]*File, error) {
	packageFolder := filepath.Join(pkg.PackageType()+"s", pkg.GetName())

	files := []*File{}
	err := billyutil.Walk(*worktree, packageFolder, func(filePath string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error finding file %s, %w", filePath, err)
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, NewFile(filePath, packageFolder, *worktree))
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}

// GetFilesForPackage clones the package in memory
func GetFilesForPackage(filesGetter FilesGetter, pkg v1alpha1.Package) ([]*File, error) {
	memFs, err := filesGetter.WorkTreeForPackage(pkg)
	if err != nil {
		return nil, err
	}

	return filterWorktreeForPackage(memFs, pkg)
}
