//go:build !e2e
// +build !e2e

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
	"os"
	"path"
	"testing"

	"github.com/mia-platform/vab/internal/testutils"
	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/kustomize/api/konfig"
	kustomize "sigs.k8s.io/kustomize/api/types"
)

// getSortedModules returns the list of modules sorted correctly
func TestSortedModulesList(t *testing.T) {
	modules := make(map[string]v1alpha1.Package)
	modules["m1-1.0.0/f1"] = v1alpha1.NewModule(
		t,
		"m1/f1",
		"1.0.0",
		false,
	)
	modules["m2-1.0.0/f1"] = v1alpha1.NewModule(
		t,
		"m2/f1",
		"1.0.0",
		false,
	)
	modules["m3-1.0.0/f1"] = v1alpha1.NewModule(
		t,
		"m3/f1",
		"1.0.0",
		false,
	)
	modules["m4b-1.0.0/f1"] = v1alpha1.NewModule(
		t,
		"m4b/f1",
		"1.0.0",
		false,
	)
	modules["m4a-1.0.0/f1"] = v1alpha1.NewModule(
		t,
		"m4a/f1",
		"1.0.0",
		false,
	)
	modules["m0-1.0.0/f1"] = v1alpha1.NewModule(
		t,
		"m0/f1",
		"1.0.0",
		true,
	)

	expectedList := []string{"../../../vendors/modules/m1-1.0.0/f1", "../../../vendors/modules/m2-1.0.0/f1", "../../../vendors/modules/m3-1.0.0/f1", "../../../vendors/modules/m4a-1.0.0/f1", "../../../vendors/modules/m4b-1.0.0/f1"}
	list := getSortedPackagesList(&modules, utils.AllGroupsDirPath)

	assert.Equal(t, expectedList, list, "Unexpected modules list.")
}

// SyncResources appends the correct resources in the kustomization.yaml
// when the existing resources list is empty
func TestSyncEmptyKustomization(t *testing.T) {
	emptyKustomization := kustomize.Kustomization{}
	modules := make(map[string]v1alpha1.Package)
	modules["m1-1.0.0/f1"] = v1alpha1.NewModule(
		t,
		"m1/f1",
		"1.0.0",
		false,
	)
	modules["m3-1.0.0/f1"] = v1alpha1.NewModule(
		t,
		"m3/f1",
		"1.0.0",
		false,
	)
	modules["m2-1.0.0/f1"] = v1alpha1.NewModule(
		t,
		"m2/f1",
		"1.0.0",
		false,
	)
	modules["m0-1.0.0/f1"] = v1alpha1.NewModule(
		t,
		"m0/f1",
		"1.0.0",
		true,
	)
	addons := make(map[string]v1alpha1.Package)
	addons["ao1-1.0.0"] = v1alpha1.NewAddon(
		t,
		"ao1",
		"1.0.0",
		false,
	)
	addons["ao2-1.0.0"] = v1alpha1.NewAddon(
		t,
		"ao2",
		"1.0.0",
		false,
	)
	addons["ao3-1.0.0"] = v1alpha1.NewAddon(
		t,
		"ao1",
		"1.0.0",
		true,
	)

	finalKustomization := SyncKustomizeResources(&modules, &addons, emptyKustomization, utils.AllGroupsDirPath)
	expectedResources := []string{"../../../vendors/modules/m1-1.0.0/f1", "../../../vendors/modules/m2-1.0.0/f1", "../../../vendors/modules/m3-1.0.0/f1"}
	expectedComponents := []string{"../../../vendors/add-ons/ao1-1.0.0", "../../../vendors/add-ons/ao2-1.0.0"}

	assert.Equal(t, expectedResources, finalKustomization.Resources, "Unexpected resources in Kustomization.")
	assert.Equal(t, expectedComponents, finalKustomization.Components, "Unexpected resources in Kustomization.")
	assert.NotEqual(t, emptyKustomization.Resources, expectedResources, "The original Kustomization struct should remain unchanged.")
	assert.NotEqual(t, emptyKustomization.Components, expectedComponents, "The original Kustomization struct should remain unchanged.")
}

// SyncResources appends the correct resources in the kustomization.yaml
// when the existing resources list is not empty
func TestSyncExistingKustomization(t *testing.T) {
	kustomization := kustomize.Kustomization{}
	kustomization.Resources = []string{
		"../../../vendors/modules/mod1-1.0.0/f1",
		"../../../vendors/modules/mod2-1.0.0/f1",
		"../../../vendors/modules/mod3-1.0.0/f1",
		"./local/mod-1.0.0/f1",
	}
	kustomization.Components = []string{
		"../../../vendors/add-ons/ao1-1.0.0",
		"../../../vendors/add-ons/ao2-1.0.0",
		"../../../vendors/add-ons/ao3-1.0.0",
		"./local/ao-1.0.0",
	}
	modules := make(map[string]v1alpha1.Package)
	// change mod1 version
	modules["mod1-2.0.0/f1"] = v1alpha1.Package{
		Version: "2.0.0",
	}
	// disable mod2
	modules["mod2-1.0.0/f1"] = v1alpha1.Package{
		Version: "1.0.0",
		Disable: true,
	}
	// unchanged module
	modules["mod3-1.0.0/f1"] = v1alpha1.Package{
		Version: "1.0.0",
	}
	addons := make(map[string]v1alpha1.Package)
	// change ao1 version
	addons["ao1-2.0.0"] = v1alpha1.Package{
		Version: "2.0.0",
	}
	// disable ao2
	addons["ao2-1.0.0"] = v1alpha1.Package{
		Version: "1.0.0",
		Disable: true,
	}
	// unchanged add-on
	addons["ao3-1.0.0"] = v1alpha1.Package{
		Version: "1.0.0",
	}

	finalKustomization := SyncKustomizeResources(&modules, &addons, kustomization, utils.AllGroupsDirPath)
	expectedResources := []string{
		"../../../vendors/modules/mod1-2.0.0/f1",
		"../../../vendors/modules/mod3-1.0.0/f1",
		"./local/mod-1.0.0/f1",
	}
	expectedComponents := []string{
		"../../../vendors/add-ons/ao1-2.0.0",
		"../../../vendors/add-ons/ao3-1.0.0",
		"./local/ao-1.0.0",
	}

	assert.Equal(t, expectedResources, finalKustomization.Resources, "Unexpected resources in Kustomization.")
	assert.Equal(t, expectedComponents, finalKustomization.Components, "Unexpected components in Kustomization.")
}

// fixResourcesPath appends the correct prefix to modules
func TestFixModulesPath(t *testing.T) {
	modulesList := []string{
		"test-module1-1.0.0/test-flavour1",
		"test-module2-1.0.0/test-flavour2",
		"test-module3-1.0.0/test-flavour3",
	}
	fixedList := fixResourcesPath(modulesList, utils.AllGroupsDirPath, true)
	expectedList := []string{
		"../../../vendors/modules/test-module1-1.0.0/test-flavour1",
		"../../../vendors/modules/test-module2-1.0.0/test-flavour2",
		"../../../vendors/modules/test-module3-1.0.0/test-flavour3",
	}
	assert.Equal(t, expectedList, *fixedList, "Unexpected elements in the resulting list of paths")
}

// fixResourcesPath appends the correct prefix to add-ons
func TestFixAddOnsPath(t *testing.T) {
	modulesList := []string{
		"test-addon1-1.0.0",
		"test-addon2-1.0.0",
	}
	fixedList := fixResourcesPath(modulesList, utils.AllGroupsDirPath, false)
	expectedList := []string{
		"../../../vendors/add-ons/test-addon1-1.0.0",
		"../../../vendors/add-ons/test-addon2-1.0.0",
	}
	assert.Equal(t, expectedList, *fixedList, "Unexpected elements in the resulting list of paths")
}

// ReadKustomization creates the empty kustomization file if missing
func TestReadKustomizationCreatePath(t *testing.T) {
	testDirPath := t.TempDir()
	basesPath := path.Join(testDirPath, utils.BasesDir)
	kustomization, err := ReadKustomization(basesPath)
	if !assert.NoError(t, err) {
		return
	}
	expectedKustomizationObject := utils.EmptyKustomization()
	kustomizationFilePath := path.Join(basesPath, konfig.DefaultKustomizationFileName())
	assert.Equal(t, expectedKustomizationObject, *kustomization, "Unexpected kustomization object")
	assert.FileExists(t, kustomizationFilePath, "Missing kustomization file")
	actualKustomization, _ := os.ReadFile(kustomizationFilePath)
	expectedKustomization, _ := os.ReadFile(testutils.GetTestFile("utils", "empty_kustomization.yaml"))
	assert.Equal(t, expectedKustomization, actualKustomization, "Unexpected file content")
}

// getKustomizationFilePath returns the correct path to the valid kustomization
func TestGetExistingKustomizationFilePath(t *testing.T) {
	testDirPath := t.TempDir()
	// existing file name: kustomization.yaml
	expectedPaths := []string{
		path.Join(testDirPath, konfig.RecognizedKustomizationFileNames()[0]),
		path.Join(testDirPath, konfig.RecognizedKustomizationFileNames()[1]),
		path.Join(testDirPath, konfig.RecognizedKustomizationFileNames()[2]),
	}
	_, err := os.Create(expectedPaths[0])
	if err != nil {
		return
	}
	kustomizationPath, err := getKustomizationFilePath(testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expectedPaths[0], kustomizationPath, "Unexpected kustomization path")
	os.Remove(expectedPaths[0])
	// existing file name: kustomization.yml
	_, err = os.Create(expectedPaths[1])
	if err != nil {
		return
	}
	kustomizationPath, err = getKustomizationFilePath(testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expectedPaths[1], kustomizationPath, "Unexpected kustomization path")
	os.Remove(expectedPaths[1])
	// existing file name: Kustomization
	_, err = os.Create(expectedPaths[2])
	if err != nil {
		return
	}
	kustomizationPath, err = getKustomizationFilePath(testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expectedPaths[2], kustomizationPath, "Unexpected kustomization path")
	os.Remove(expectedPaths[2])
}

// getKustomizationFilePath creates the file if missing and returns the correct path
func TestGetMissingKustomizationPath(t *testing.T) {
	testDirPath := t.TempDir()
	expectedPath := path.Join(testDirPath, konfig.DefaultKustomizationFileName())
	kustomizationPath, err := getKustomizationFilePath(testDirPath)
	if assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expectedPath, kustomizationPath, "Unexpected kustomization path")
	assert.FileExists(t, kustomizationPath, "Missing kustomization file")
	actualKustomization, _ := os.ReadFile(path.Join(kustomizationPath, konfig.DefaultKustomizationFileName()))
	expectedKustomization, _ := os.ReadFile(testutils.GetTestFile("utils", "empty_kustomization.yaml"))
	assert.Equal(t, expectedKustomization, actualKustomization, "Unexpected file content")
}

// getModuleCompleteName returns the string in the correct format <module>-<semver>/<flavour>
func TestGetModuleCompleteName(t *testing.T) {
	modules := make(map[string]v1alpha1.Package)
	modules["m1/f1"] = v1alpha1.NewModule(
		t,
		"m1/f1",
		"1.0.0",
		true,
	)
	modules["m2/f1"] = v1alpha1.NewModule(
		t,
		"m2/f1",
		"1.0.0",
		true,
	)
	expectedModules := make(map[string]v1alpha1.Package)
	expectedModules["m1-1.0.0/f1"] = v1alpha1.NewModule(
		t,
		"m1/f1",
		"1.0.0",
		true,
	)
	expectedModules["m2-1.0.0/f1"] = v1alpha1.NewModule(
		t,
		"m2/f1",
		"1.0.0",
		true,
	)
	updatedModules := PackagesMapForPaths(modules)
	assert.Equal(t, expectedModules, updatedModules, "Unexpected module name")
}
