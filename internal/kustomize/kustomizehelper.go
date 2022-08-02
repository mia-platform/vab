package kustomizehelper

import (
	"sort"
	"strings"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	kustomize "sigs.k8s.io/kustomize/api/types"
)

// SyncModules updates the modules in kustomization resources to the latest sync
func SyncModules(modules *map[string]v1alpha1.Module, k kustomize.Kustomization) kustomize.Kustomization {

	modulesList := sortedModulesList(modules)

	// If the file already includes a non-empty list of resources, this function
	// collects all the custom modules that were added manually by the user
	// (i.e. all those modules that are not present in the vendors folder, thus
	// without "vendors" in their path). Then, the custom modules are appended
	// to the updated modules list that will substitute the already existing one.
	if k.Resources != nil {
		for _, r := range k.Resources {
			if !strings.Contains(r, "vendors") {
				modulesList = append(modulesList, r)
			}
		}
	}

	k.Resources = modulesList

	return k
}

// sortedModulesList returns the list of module names sorted by weight
func sortedModulesList(modules *map[string]v1alpha1.Module) []string {

	modulesList := make([]string, 0, len(*modules))

	for m := range *modules {
		if !(*modules)[m].Disable {
			modulesList = append(modulesList, m)
		}
	}

	sort.SliceStable(modulesList, func(i, j int) bool {
		return (*modules)[modulesList[i]].Weight < (*modules)[modulesList[j]].Weight
	})

	return modulesList

}
