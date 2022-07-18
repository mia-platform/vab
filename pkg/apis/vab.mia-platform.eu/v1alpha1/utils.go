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

package v1alpha1

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
			Modules: make(map[string]Module),
			AddOns:  make(map[string]AddOn),
			Groups:  make([]Group, 0),
		},
	}
}
