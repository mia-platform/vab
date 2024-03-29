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

package v1alpha1

import "testing"

const (
	// Valid value for the kind property of the configuration
	Kind = "ClustersConfiguration"
	// Valid value for the apiVersion property of the configuration
	Version = "vab.mia-platform.eu/v1alpha1"
)

// EmptyConfig generates an empty ClustersConfiguration with provided name
func EmptyConfig(name string) *ClustersConfiguration {
	return &ClustersConfiguration{
		TypeMeta: TypeMeta{
			Kind:       Kind,
			APIVersion: Version,
		},
		Name: name,
		Spec: ConfigSpec{
			Modules: make(map[string]Package),
			AddOns:  make(map[string]Package),
			Groups:  make([]Group, 0),
		},
	}
}

func NewModule(t *testing.T, name string, version string, disable bool) Package {
	t.Helper()
	return Package{
		name:     moduleName(name),
		Version:  version,
		Disable:  disable,
		isModule: true,
		flavor:   moduleFlavorName(name),
	}
}

func NewAddon(t *testing.T, name string, version string, disable bool) Package {
	t.Helper()
	return Package{
		name:     name,
		Version:  version,
		Disable:  disable,
		isModule: false,
	}
}
