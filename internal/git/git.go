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

import "github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"

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

// urlForAddon return the git url to use for downloading the files for addon
func urlForAddon(addon v1alpha1.AddOn) string {
	return defaultRepositoryURL
}
