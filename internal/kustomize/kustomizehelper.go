package kustomizehelper

import (
	"sort"
	"strings"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	kustomize "sigs.k8s.io/kustomize/api/types"
)

// SyncKustomizations updates the kustomization files to the latest sync
func SyncModules(modules *map[string]v1alpha1.Module, k kustomize.Kustomization) kustomize.Kustomization {

	modulesList := sortedModulesList(modules)

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
