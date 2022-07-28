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
	"fmt"
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

// urlForModule return the git url to use for downloading the files for module
func urlForModule(module v1alpha1.Module) string {
	return defaultRepositoryURL
}

// authForModule return an AuthMethod for the module
func authForModule(module v1alpha1.Module) transport.AuthMethod {
	return nil
}

// tagReferenceForModule return a valid tag reference for moduleName and version
func tagReferenceForModule(moduleName string, version string) plumbing.ReferenceName {
	splitName := strings.Split(moduleName, "/") // TODO da mettere in utils
	return plumbing.NewTagReferenceName("module-" + splitName[0] + "-" + version)
}

// urlForAddon return the git url to use for downloading the files for addon
func urlForAddon(addon v1alpha1.AddOn) string {
	return defaultRepositoryURL
}

// authForAddon return an AuthMethod for the addon
func authForAddon(addon v1alpha1.AddOn) transport.AuthMethod {
	return nil
}

// tagReferenceForAddon return a valid tag reference for addonName and version
func tagReferenceForAddon(addonName string, version string) plumbing.ReferenceName {
	return plumbing.NewTagReferenceName("addon-" + addonName + "-" + version)
}

// cloneOptionsForModule return a the options for cloning moduleName with the module configuration
func cloneOptionsForModule(moduleName string, module v1alpha1.Module) *git.CloneOptions {
	return &git.CloneOptions{
		URL:           urlForModule(module),
		Auth:          authForModule(module),
		ReferenceName: tagReferenceForModule(moduleName, module.Version),
		Depth:         1,
		SingleBranch:  true,
		Tags:          git.NoTags,
	}
}

// cloneOptionsForAddon return a the options for cloning addonName with the addon configuration
func cloneOptionsForAddon(addonName string, addon v1alpha1.AddOn) *git.CloneOptions {
	return &git.CloneOptions{
		URL:           urlForAddon(addon),
		Auth:          authForAddon(addon),
		ReferenceName: tagReferenceForAddon(addonName, addon.Version),
		Depth:         1,
		SingleBranch:  true,
		Tags:          git.NoTags,
	}
}

type cloner interface {
	Clone(addonName string, addon v1alpha1.AddOn, cloneOptions *git.CloneOptions) (billy.Filesystem, error)
}

type clone struct {
}

func (c clone) Clone(addonName string, addon v1alpha1.AddOn, cloneOptions *git.CloneOptions) (billy.Filesystem, error) {
	workTree := memfs.New()
	memStorage := memory.NewStorage()
	_, err := git.Clone(memStorage, workTree, cloneOptions)
	if err != nil {
		return nil, fmt.Errorf("error cloning repo: %w", err)
	}
	return workTree, nil
}

func cloneAddon(addonName string, addon v1alpha1.AddOn, cloneOptions *git.CloneOptions, cloner cloner) (billy.Filesystem, error) {
	workTree := memfs.New()
	return workTree, nil
}

// TODO Aggiungere qui la funzione func(addonName, addon, targetPath )
// TODO prendere solo il basename di modules mentre addons non ha slash

// TODO Interfaccia addon/module

// TODO lanciare make generate se modifichi in pkg apis
