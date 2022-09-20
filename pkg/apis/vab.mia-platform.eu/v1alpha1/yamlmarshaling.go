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

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type shadowConfigSpec struct {

	// Dictionary of Modules
	// These modules will be installed on every cluster
	// unless otherwise specified
	// Modules in the dictionary are referenced by module-name/flavor-name
	// For example: ingress/traefik, cni/cilium, etc.
	Modules map[string]Package `yaml:"modules"`

	// Dictionary of AddOns
	// These add-ons will be installed on every cluster
	// unless otherwise specified
	// AddOns in the dictionary are referenced by their name
	AddOns map[string]Package `yaml:"addOns"`

	// Groups contains the list of cluster groups
	Groups []Group `yaml:"groups"`
}

func (configSpec *ConfigSpec) UnmarshalYAML(value *yaml.Node) error {
	var temporaryConfig shadowConfigSpec

	if err := value.Decode(&temporaryConfig); err != nil {
		return err
	}

	fmt.Println("Hi!")
	configSpec.Groups = temporaryConfig.Groups

	newModules := map[string]Package{}
	for key, module := range temporaryConfig.Modules {
		module.name = key
		module.isModule = true
		newModules[key] = module
	}
	configSpec.Modules = newModules

	newAddos := map[string]Package{}
	for key, addon := range temporaryConfig.AddOns {
		addon.name = key
		addon.isModule = false
		newAddos[key] = addon
	}
	configSpec.AddOns = newAddos

	return nil
}

type shadowCluster struct {

	// The cluster name
	// It is required to reference the cluster directory
	Name string `yaml:"name,omitempty"`

	// Name of the context used by the cluster
	Context string `yaml:"context,omitempty"`

	// Dictionary of Modules
	// This field can be used to add a new module
	// or patch/disable a default module
	// Modules in the dictionary are referenced by "module-name/flavor-name"
	// For example: ingress/traefik, cni/cilium, etc.
	Modules map[string]Package `yaml:"modules,omitempty"`

	// Dictionary of AddOns
	// This field can be used to add a new add-on
	// or patch/disable a default add-on
	// AddOns in the dictionary are referenced by their name
	AddOns map[string]Package `yaml:"addOns,omitempty"`
}

func (cluster *Cluster) UnmarshalYAML(value *yaml.Node) error {
	var temporaryCluster shadowCluster

	if err := value.Decode(&temporaryCluster); err != nil {
		return err
	}

	cluster.Name = temporaryCluster.Name
	cluster.Context = temporaryCluster.Context

	newModules := map[string]Package{}
	for key, module := range temporaryCluster.Modules {
		module.name = key
		module.isModule = true
		newModules[key] = module
	}
	cluster.Modules = newModules

	newAddos := map[string]Package{}
	for key, addon := range temporaryCluster.AddOns {
		addon.name = key
		addon.isModule = false
		newAddos[key] = addon
	}
	cluster.AddOns = newAddos

	return nil
}
