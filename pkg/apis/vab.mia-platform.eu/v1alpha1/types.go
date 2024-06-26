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

// TypeMeta partially copies apimachinery/pkg/apis/meta/v1.TypeMeta
type TypeMeta struct {
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
}

// ClustersConfiguration contains the schema for vab's configuration
type ClustersConfiguration struct {
	TypeMeta `json:",inline" yaml:",inline"`

	// The configuration name
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// ConfigSpec contains the configuration of the clusters
	// It includes the modules and add-ons installed by default
	// as well as the list of cluster groups
	Spec ConfigSpec `json:"spec" yaml:"spec"`
}

// ConfigSpec contains the configuration of the clusters
type ConfigSpec struct {

	// Dictionary of Modules
	// These modules will be installed on every cluster
	// unless otherwise specified
	// Modules in the dictionary are referenced by module-name/flavor-name
	// For example: ingress/traefik, cni/cilium, etc.
	Modules map[string]Package `json:"modules" yaml:"modules"`

	// Dictionary of AddOns
	// These add-ons will be installed on every cluster
	// unless otherwise specified
	// AddOns in the dictionary are referenced by their name
	AddOns map[string]Package `json:"addOns" yaml:"addOns"`

	// Groups contains the list of cluster groups
	Groups []Group `json:"groups" yaml:"groups"`
}

// Group contains the configuration of a cluster group
type Group struct {

	// The group name
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Clusters contains the list of the clusters in the group
	// This field is required to reference the clusters correctly
	// in the directory structure
	Clusters []Cluster `json:"clusters,omitempty" yaml:"clusters,omitempty"`
}

// Cluster contains the configuration of a cluster
// and customizations of its modules/add-ons
type Cluster struct {

	// The cluster name
	// It is required to reference the cluster directory
	Name string `json:"name,omitempty" yaml:"name,omitempty"`

	// Name of the context used by the cluster
	Context string `json:"context,omitempty" yaml:"context,omitempty"`

	// Dictionary of Modules
	// This field can be used to add a new module
	// or patch/disable a default module
	// Modules in the dictionary are referenced by "module-name/flavor-name"
	// For example: ingress/traefik, cni/cilium, etc.
	Modules map[string]Package `json:"modules,omitempty" yaml:"modules,omitempty"`

	// Dictionary of AddOns
	// This field can be used to add a new add-on
	// or patch/disable a default add-on
	// AddOns in the dictionary are referenced by their name
	AddOns map[string]Package `json:"addOns,omitempty" yaml:"addOns,omitempty"`
}

// Module contains the module's version and priority
type Package struct {

	// name is a private property containing the name of the module
	name string

	// Version of the module to be installed
	Version string `json:"version" yaml:"version"`

	// Flag that disables the add-on if set to true
	Disable bool `json:"disable" yaml:"disable"`

	// isModule is a private property for setting if a package is a module or an addon
	isModule bool

	// flavor is a private property that contains the flavor name if is a module, or is an empty string otherwise
	flavor string
}

// IsModule return the value of the private property with the same name
func (pkg Package) IsModule() bool {
	return pkg.isModule
}

// GetName return the canonical name of the package
func (pkg Package) GetName() string {
	return pkg.name
}

// GetFlavorName return the flavor name of the package if is a module or an empty string otherwise
func (pkg Package) GetFlavorName() string {
	return pkg.flavor
}

// PackageType return the type of the package in string form
func (pkg Package) PackageType() string {
	if pkg.isModule {
		return "module"
	}
	return "addon"
}
