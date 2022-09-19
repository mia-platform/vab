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
	"strings"
)

// TypeMeta partially copies apimachinery/pkg/apis/meta/v1.TypeMeta
type TypeMeta struct {
	Kind       string `yaml:"kind,omitempty"`
	APIVersion string `yaml:"apiVersion,omitempty"`
}

// ClustersConfiguration contains the schema for vab's configuration
type ClustersConfiguration struct {
	TypeMeta `yaml:",inline"`

	// The configuration name
	Name string `yaml:"name,omitempty"`

	// ConfigSpec contains the configuration of the clusters
	// It includes the modules and add-ons installed by default
	// as well as the list of cluster groups
	Spec ConfigSpec `yaml:"spec"`
}

// ConfigSpec contains the configuration of the clusters
type ConfigSpec struct {

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

// Group contains the configuration of a cluster group
type Group struct {

	// The group name
	Name string `yaml:"name,omitempty"`

	// Clusters contains the list of the clusters in the group
	// This field is required to reference the clusters correctly
	// in the directory structure
	Clusters []Cluster `yaml:"clusters,omitempty"`
}

// Cluster contains the configuration of a cluster
// and customizations of its modules/add-ons
type Cluster struct {

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

// Module contains the module's version and priority
type Package struct {

	// name is a private property containing the name of the module
	name string

	// Version of the module to be installed
	Version string `yaml:"version"`

	// Flag that disables the add-on if set to true
	Disable bool `yaml:"disable"`

	// isModule is a private property for setting if a package is a module or an addon
	isModule bool
}

// IsModule return the value of the private property with the same name
func (pkg Package) IsModule() bool {
	return pkg.isModule
}

// GetName return the canonical name of the package
func (pkg Package) GetName() string {
	if pkg.isModule {
		splittedStrings := strings.Split(pkg.name, "/")
		return strings.Join(splittedStrings[:len(splittedStrings)-1], "/")
	}
	return pkg.name
}

// GetFlavorName return the flavor name of the package if is a module or an empty string otherwise
func (pkg Package) GetFlavorName() string {
	if pkg.isModule {
		splittedStrings := strings.Split(pkg.name, "/")
		return splittedStrings[len(splittedStrings)-1]
	}
	return ""
}

// PackageType return the type of the package in string form
func (pkg Package) PackageType() string {
	if pkg.isModule {
		return "module"
	}
	return "addon"
}
