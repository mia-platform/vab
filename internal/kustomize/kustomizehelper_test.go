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

package kustomizehelper

import (
	"testing"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/stretchr/testify/assert"
	kustomize "sigs.k8s.io/kustomize/api/types"
)

// getSortedModules returns the list of modules sorted correctly
func TestSortedModulesList(t *testing.T) {
	modules := make(map[string]v1alpha1.Module)
	modules["m4"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  4,
	}
	modules["m1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  1,
	}
	modules["m3"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  3,
	}
	modules["m2b"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  2,
	}
	modules["m2a"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  2,
	}
	modules["m0"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  10,
		Disable: true,
	}

	expectedList := []string{"m1", "m2a", "m2b", "m3", "m4"}
	list := getSortedModulesList(&modules)

	assert.Equal(t, expectedList, list, "Unexpected modules list.")
}

// SyncResources appends the correct resources in the kustomization.yaml
// when the existing resources list is empty
func TestSyncEmptyKustomization(t *testing.T) {
	emptyKustomization := kustomize.Kustomization{}
	modules := make(map[string]v1alpha1.Module)
	modules["m1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  1,
	}
	modules["m3"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  3,
	}
	modules["m2"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  2,
	}
	modules["m0"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  10,
		Disable: true,
	}
	addons := make(map[string]v1alpha1.AddOn)
	addons["ao1"] = v1alpha1.AddOn{
		Version: "1.0.0",
	}
	addons["ao2"] = v1alpha1.AddOn{
		Version: "1.0.0",
	}
	addons["ao0"] = v1alpha1.AddOn{
		Version: "1.0.0",
		Disable: true,
	}

	finalKustomization := SyncResources(&modules, &addons, emptyKustomization)
	expectedResources := []string{"m1", "m2", "m3", "ao1", "ao2"}

	assert.Equal(t, expectedResources, finalKustomization.Resources, "Unexpected resources in Kustomization.")
	assert.NotEqual(t, emptyKustomization, expectedResources, "The original Kustomization struct should remain unchanged.")
}

// SyncResources appends the correct resources in the kustomization.yaml
// when the existing resources list is not empty
func TestSyncExistingKustomization(t *testing.T) {
	kustomization := kustomize.Kustomization{}
	kustomization.Resources = []string{
		"vendors/modules/mod1-1.0.0",
		"vendors/modules/mod2-1.0.0",
		"vendors/modules/mod3-1.0.0",
		"local/mod-1.0.0",
		"vendors/addon/ao1-1.0.0",
		"vendors/addon/ao2-1.0.0",
		"vendors/addon/ao3-1.0.0",
		"local/ao-1.0.0",
	}
	modules := make(map[string]v1alpha1.Module)
	// change mod1 version
	modules["vendors/modules/mod1-2.0.0"] = v1alpha1.Module{
		Version: "2.0.0",
		Weight:  1,
	}
	// disable mod2
	modules["vendors/modules/mod2-1.0.0"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  2,
		Disable: true,
	}
	// unchanged module
	modules["vendors/modules/mod3-1.0.0"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  3,
	}
	addons := make(map[string]v1alpha1.AddOn)
	// change ao1 version
	addons["vendors/addon/ao1-2.0.0"] = v1alpha1.AddOn{
		Version: "2.0.0",
	}
	// disable ao2
	addons["vendors/addon/ao2-1.0.0"] = v1alpha1.AddOn{
		Version: "1.0.0",
		Disable: true,
	}
	// unchanged add-on
	addons["vendors/addon/ao3-1.0.0"] = v1alpha1.AddOn{
		Version: "1.0.0",
	}

	finalKustomization := SyncResources(&modules, &addons, kustomization)
	expectedResources := []string{
		"vendors/modules/mod1-2.0.0",
		"vendors/modules/mod3-1.0.0",
		"vendors/addon/ao1-2.0.0",
		"vendors/addon/ao3-1.0.0",
		"local/mod-1.0.0",
		"local/ao-1.0.0",
	}

	assert.Equal(t, expectedResources, finalKustomization.Resources, "Unexpected resources in Kustomization.")
}
