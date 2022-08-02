package kustomizehelper

import (
	"testing"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/stretchr/testify/assert"
	kustomize "sigs.k8s.io/kustomize/api/types"
)

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
	modules["m2"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  2,
	}
	modules["m0"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  10,
		Disable: true,
	}

	expectedList := []string{"m1", "m2", "m3", "m4"}
	list := sortedModulesList(&modules)

	assert.Equal(t, expectedList, list, "Unexpected modules list.")
}

func TestSyncModules(t *testing.T) {
	emptyKustomization := kustomize.Kustomization{}
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
	modules["m2"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  2,
	}
	modules["m0"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  10,
		Disable: true,
	}

	finalKustomization := SyncModules(&modules, emptyKustomization)
	expectedResources := []string{"m1", "m2", "m3", "m4"}

	assert.Equal(t, expectedResources, finalKustomization.Resources, "Unexpected resources in Kustomization.")
	assert.NotEqual(t, emptyKustomization, expectedResources, "The original Kustomization struct should remain unchanged.")

}
