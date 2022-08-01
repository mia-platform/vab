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
	"io/fs"
	"strings"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
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

type Package interface {
	v1alpha1.Module | v1alpha1.AddOn
}

// remoteUrl return the git url to use for downloading the files for a package (module or addon)
func remoteURL[P Package](pkg P) string {
	return defaultRepositoryURL
}

// remoteAuth return an AuthMethod for package (module or addon)
func remoteAuth[P Package](pkg P) transport.AuthMethod {
	return nil
}

// tagReferenceForPackage return a valid tag reference for the package name and version
func tagReferenceForPackage[P Package](pkgName string, pkg P) plumbing.ReferenceName {
	var tag string
	switch pkg := (interface{})(pkg).(type) {
	case v1alpha1.Module:
		tag = "module-" + strings.Split(pkgName, "/")[0] + "-" + pkg.Version
	case v1alpha1.AddOn:
		tag = "addon-" + pkgName + "-" + pkg.Version
	}

	return plumbing.NewTagReferenceName(tag)
}

// cloneOptionsForPackage return the options for cloning the package with pkgName with pkg configuaration
func cloneOptionsForPackage[P Package](pkgName string, pkg P) *git.CloneOptions {
	return &git.CloneOptions{
		URL:           remoteURL(pkg),
		Auth:          remoteAuth(pkg),
		ReferenceName: tagReferenceForPackage(pkgName, pkg),
		Depth:         1,
		SingleBranch:  true,
		Tags:          git.NoTags,
	}
}

// worktreeForPackage return a worktree from the cloned repository for the package with pkgName
func worktreeForPackage[P Package](pkgName string, pkg P) (billy.Filesystem, error) {
	cloneOptions := cloneOptionsForPackage(pkgName, pkg)
	fs := memfs.New()
	storage := memory.NewStorage()
	if _, err := git.Clone(storage, fs, cloneOptions); err != nil {
		return nil, err
	}

	return fs, nil
}

func filterWorktreeForPackage[P Package](worktree billy.Filesystem, pkgName string, pkg P) ([]fs.FileInfo, error) {
	var packageFolder string
	switch (interface{})(pkg).(type) {
	case v1alpha1.Module:
		packageFolder = "./modules/" + strings.Split(pkgName, "/")[0]
	case v1alpha1.AddOn:
		packageFolder = "./add-ons/" + pkgName
	}

	dirElements, err := worktree.ReadDir(packageFolder)
	if err != err {
		return nil, err
	}

	return dirElements, nil
}
