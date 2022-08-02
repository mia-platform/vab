package kustomizehelper

import (
	"sort"
	"strings"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	kustomize "sigs.k8s.io/kustomize/api/types"
)

// SyncResources updates the clusters' kustomization resources to the latest sync
func SyncResources(modules *map[string]v1alpha1.Module, addons *map[string]v1alpha1.AddOn, k kustomize.Kustomization) kustomize.Kustomization {

	modulesList := getSortedModulesList(modules)
	addonsList := getAddOnsList(addons)
	resourcesList := append(modulesList, addonsList...)

	// If the file already includes a non-empty list of resources, this function
	// collects all the custom modules that were added manually by the user
	// (i.e. all those modules that are not present in the vendors folder, thus
	// without "vendors" in their path). Then, the custom modules are appended
	// to the updated modules list that will substitute the already existing one.
	if k.Resources != nil {
		for _, r := range k.Resources {
			if !strings.Contains(r, "vendors") {
				resourcesList = append(resourcesList, r)
			}
		}
	}

	k.Resources = resourcesList

	return k
}

// getSortedModulesList returns the list of module names sorted by weight.
// In case of equal weights, the modules are ordered lexicographically.
func getSortedModulesList(modules *map[string]v1alpha1.Module) []string {

	modulesList := make([]string, 0, len(*modules))

	for m := range *modules {
		if !(*modules)[m].Disable {
			modulesList = append(modulesList, m)
		}
	}

	sort.SliceStable(modulesList, func(i, j int) bool {
		// If the weights are equal, order the elements lexicographically
		if (*modules)[modulesList[i]].Weight == (*modules)[modulesList[j]].Weight {
			return modulesList[i] < modulesList[j]
		}
		// Otherwise, sort by weight (increasing order)
		return (*modules)[modulesList[i]].Weight < (*modules)[modulesList[j]].Weight
	})

	return modulesList

}

// getAddOnsList returns the list of addons names in lexicographic order
func getAddOnsList(addons *map[string]v1alpha1.AddOn) []string {

	addonsList := make([]string, 0, len(*addons))

	for ao := range *addons {
		if !(*addons)[ao].Disable {
			addonsList = append(addonsList, ao)
		}
	}

	sort.SliceStable(addonsList, func(i, j int) bool {
		return addonsList[i] < addonsList[j]
	})

	return addonsList
}
