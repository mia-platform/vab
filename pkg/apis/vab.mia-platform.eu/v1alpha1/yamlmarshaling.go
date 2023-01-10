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

import (
	"crypto/sha1" //#nosec G505
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// moduleName return the actual name of a module from the key used in the configuration file
func moduleName(moduleKey string) string {
	splitStrings := strings.Split(moduleKey, "/")
	return strings.Join(splitStrings[:len(splitStrings)-1], "/")
}

// moduleFlavorName return the flavor name of a module from the key used in the configuration file
func moduleFlavorName(moduleKey string) string {
	splitStrings := strings.Split(moduleKey, "/")
	return splitStrings[len(splitStrings)-1]
}

// mapKeyForName return an abstract key for a package name consisting in the first 7 character of its sha1 shasum
func mapKeyForName(name string) string {
	// using of sha1 is not a problem because we don't use this string for secure operations
	// but only for generating unique ids for maps
	sha1 := sha1.Sum([]byte(name)) //#nosec G401
	return fmt.Sprintf("%x", sha1)[0:7]
}

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

// UnmarshalYAML conform to Unmarshaler interface for customize the modules and addons
// maps enriching the Package structs with additionals information
func (configSpec *ConfigSpec) UnmarshalYAML(value *yaml.Node) error {
	var temporaryConfig shadowConfigSpec

	if err := value.Decode(&temporaryConfig); err != nil {
		return err
	}

	configSpec.Groups = temporaryConfig.Groups

	newModules := map[string]Package{}
	for key, module := range temporaryConfig.Modules {
		moduleName := moduleName(key)
		module.name = moduleName
		module.isModule = true
		module.flavor = moduleFlavorName(key)
		newModules[mapKeyForName(moduleName)] = module
	}
	configSpec.Modules = newModules

	newAddons := map[string]Package{}
	for key, addon := range temporaryConfig.AddOns {
		addon.name = key
		addon.isModule = false
		newAddons[mapKeyForName(key)] = addon
	}
	configSpec.AddOns = newAddons

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

// UnmarshalYAML conform to Unmarshaler interface for customize the modules and addons
// maps enriching the Package structs with additionals information
func (cluster *Cluster) UnmarshalYAML(value *yaml.Node) error {
	var temporaryCluster shadowCluster

	if err := value.Decode(&temporaryCluster); err != nil {
		return err
	}

	cluster.Name = temporaryCluster.Name
	cluster.Context = temporaryCluster.Context

	newModules := map[string]Package{}
	for key, module := range temporaryCluster.Modules {
		moduleName := moduleName(key)
		module.name = moduleName
		module.isModule = true
		module.flavor = moduleFlavorName(key)
		newModules[mapKeyForName(moduleName)] = module
	}
	cluster.Modules = newModules

	newAddons := map[string]Package{}
	for key, addon := range temporaryCluster.AddOns {
		addon.name = key
		addon.isModule = false
		newAddons[mapKeyForName(key)] = addon
	}
	cluster.AddOns = newAddons

	return nil
}
