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
	defaultGitUrl = "https://github.com/mia-platform/distribution"
)

func urlForModule(module v1alpha1.Module) string {
	return defaultGitUrl
}

func urlForAddon(addon v1alpha1.AddOn) string {
	return defaultGitUrl
}
