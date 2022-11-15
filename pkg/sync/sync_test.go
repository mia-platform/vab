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

package sync

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mia-platform/vab/internal/git"
	kustomizehelper "github.com/mia-platform/vab/internal/kustomize"
	"github.com/mia-platform/vab/internal/testutils"
	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/mia-platform/vab/pkg/logger"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/kustomize/api/konfig"
)

// ClonePackages returns a list of mocked file pointers w/o errors
func TestClonePackage(t *testing.T) {
	logger := logger.DisabledLogger{}
	testModule := v1alpha1.NewModule(t, "category/test-module1/test-flavour1", "1.0.0", false)
	outputFiles, err := ClonePackages(logger, testModule, testutils.FakeFilesGetter{Testing: t})
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, outputFiles, "The returned array of mocked file pointers is empty")
}

// MoveToDisk correctly moves the files from the worktree to disk
func TestMoveToDisk(t *testing.T) {
	logger := logger.DisabledLogger{}
	fakeWorktree := testutils.PrepareFakeWorktree(t)
	input := []*git.File{
		git.NewFile("./modules/category/test-module1/test-flavour1/file1.yaml", "./modules/category/test-module1", *fakeWorktree),
		git.NewFile("./modules/category/test-module1/test-flavour1/file2.yaml", "./modules/category/test-module1", *fakeWorktree),
		git.NewFile("./modules/category/test-module1/test-flavour1/file1.yaml", "./modules/category/test-module1", *fakeWorktree),
	}
	testDirPath := t.TempDir()
	err := MoveToDisk(logger, input, "test-module1/test-flavour1", testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.FileExists(t, filepath.Join(testDirPath, "test-flavour1/file1.yaml"), "Mock file 1 does not exist on disk")
	assert.FileExists(t, filepath.Join(testDirPath, "test-flavour1/file2.yaml"), "Mock file 2 does not exist on disk")
	assert.FileExists(t, filepath.Join(testDirPath, "test-flavour1/file1.yaml"), "Mock file 3 does not exist on disk")
}

// UpdateModules syncs new modules without errors
func TestUpdateModules(t *testing.T) {
	logger := logger.DisabledLogger{}
	modules := make(map[string]v1alpha1.Package)
	modules["module1"] = v1alpha1.NewModule(
		t,
		"category/test-module1/test-flavour1",
		"1.0.0",
		false,
	)
	modules["module2"] = v1alpha1.NewModule(
		t,
		"category/test-module2/test-flavour1",
		"1.0.0",
		false,
	)
	modules["module3"] = v1alpha1.NewModule(
		t,
		"category/test-module3/test-flavour1",
		"1.0.0",
		true,
	)

	testDirPath := t.TempDir()
	err := clonePackagesLocally(logger, modules, testDirPath, testutils.FakeFilesGetter{Testing: t})
	if !assert.NoError(t, err) {
		return
	}
}

// UpdateAddOns syncs new modules without errors
func TestUpdateAddOns(t *testing.T) {
	logger := logger.DisabledLogger{}
	addons := make(map[string]v1alpha1.Package)
	addons["addon1"] = v1alpha1.NewAddon(
		t,
		"category/test-addon1",
		"1.0.0",
		false,
	)
	addons["addon2"] = v1alpha1.NewAddon(
		t,
		"category/test-addon2",
		"1.0.0",
		true,
	)
	testDirPath := t.TempDir()
	err := clonePackagesLocally(logger, addons, testDirPath, testutils.FakeFilesGetter{Testing: t})
	if !assert.NoError(t, err) {
		return
	}
}

// UpdateBases correctly updates the resources list in the all-groups kustomization
func TestUpdateAllGroups(t *testing.T) {
	testDirPath := t.TempDir()
	targetPath := filepath.Join(testDirPath, utils.AllGroupsDirPath)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		return
	}
	modules := make(map[string]v1alpha1.Package)
	modules["module3"] = v1alpha1.NewModule(
		t,
		"category/test-module3/test-flavour1",
		"1.0.0",
		false,
	)
	modules["module2"] = v1alpha1.NewModule(
		t,
		"category/test-module2/test-flavour1",
		"1.0.0",
		false,
	)
	modules["module1"] = v1alpha1.NewModule(
		t,
		"category/test-module1/test-flavour1",
		"1.0.0",
		false,
	)

	addons := make(map[string]v1alpha1.Package)
	addons["addon1"] = v1alpha1.NewAddon(
		t,
		"category/test-addon1",
		"1.0.0",
		false,
	)
	addons["addon2"] = v1alpha1.NewAddon(
		t,
		"category/test-addon2",
		"1.0.0",
		false,
	)
	config := v1alpha1.ClustersConfiguration{}
	config.Spec.Modules = modules
	config.Spec.AddOns = addons
	err := syncAllGroups(&config, testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	expectedKustomization, err := os.ReadFile(testutils.GetTestFile("sync", "outputs", "all_groups.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, expectedKustomization, filepath.Join(targetPath, utils.BasesDir, konfig.DefaultKustomizationFileName()))
}

// CheckClusterPath creates and returns the correct path to a missing cluster folder
// and creates the missing kustomization file
func TestCreateClusterPath(t *testing.T) {
	testDirPath := t.TempDir()
	clusterPath, err := checkClusterPath("test-cluster", testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	expectedPath := filepath.Join(testDirPath, utils.ClustersDirName, "test-cluster")
	assert.Equal(t, expectedPath, clusterPath, "Wrong path to cluster")
	assert.DirExists(t, clusterPath, "The cluster directory does not exist")
	kustomizationPath := filepath.Join(clusterPath, konfig.DefaultKustomizationFileName())
	assert.FileExists(t, kustomizationPath, "The kustomization file does not exist")
	expectedKustomization, err := os.ReadFile(testutils.GetTestFile("sync", "misc", "simple_cluster_kustomization.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, expectedKustomization, kustomizationPath)
}

// CheckClusterPath returns the correct path to an existing cluster folder
// and creates the missing kustomization file
func TestExistingClusterPath(t *testing.T) {
	testDirPath := t.TempDir()
	expectedPath := filepath.Join(testDirPath, utils.ClustersDirName, "test-cluster")
	if err := os.MkdirAll(expectedPath, os.ModePerm); err != nil {
		return
	}
	clusterPath, err := checkClusterPath("test-cluster", testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expectedPath, clusterPath, "Wrong path to cluster")
	kustomizationPath := filepath.Join(clusterPath, konfig.DefaultKustomizationFileName())
	assert.FileExists(t, kustomizationPath, "The kustomization file does not exist")
	expectedKustomization, err := os.ReadFile(testutils.GetTestFile("sync", "misc", "simple_cluster_kustomization.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, expectedKustomization, kustomizationPath)
}

// CheckClusterPath returns the correct path to an existing cluster folder
// and does not alter the existing (and well-formed) kustomization file
func TestExistingClusterPathWithKustomization(t *testing.T) {
	testDirPath := t.TempDir()
	expectedPath := filepath.Join(testDirPath, utils.ClustersDirName, "test-cluster")
	if err := os.MkdirAll(expectedPath, os.ModePerm); err != nil {
		return
	}
	kustomizationPath := filepath.Join(expectedPath, konfig.DefaultKustomizationFileName())
	expectedKustomization, err := os.ReadFile(testutils.GetTestFile("sync", "misc", "cluster_kustomization.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	err = os.WriteFile(kustomizationPath, expectedKustomization, os.ModePerm)
	if !assert.NoError(t, err) {
		return
	}
	clusterPath, err := checkClusterPath("test-cluster", testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expectedPath, clusterPath, "Wrong path to cluster")
	testutils.CompareFile(t, expectedKustomization, kustomizationPath)
}

// CheckClusterPath returns the correct path to an existing cluster folder
// and prepends "bases" to the kustomization resources when missing
func TestExistingClusterPathMissingBases(t *testing.T) {
	testDirPath := t.TempDir()
	expectedPath := filepath.Join(testDirPath, utils.ClustersDirName, "test-cluster")
	if err := os.MkdirAll(expectedPath, os.ModePerm); err != nil {
		return
	}
	// create a kustomization file without "bases" among the resources
	kustomizationPath := filepath.Join(expectedPath, konfig.DefaultKustomizationFileName())
	kustomization, err := os.ReadFile(testutils.GetTestFile("sync", "misc", "missing_bases.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	err = os.WriteFile(kustomizationPath, kustomization, os.ModePerm)
	if !assert.NoError(t, err) {
		return
	}
	clusterPath, err := checkClusterPath("test-cluster", testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expectedPath, clusterPath, "Wrong path to cluster")
	// read the resulting kustomization file and check if the "bases" directory was added to the resources
	resultingKustomization, err := kustomizehelper.ReadKustomization(expectedPath)
	if !assert.NoError(t, err) {
		return
	}
	expectedResources := []string{"bases", "deployments", "services"}
	assert.Equal(t, expectedResources, resultingKustomization.Resources, "Unexpected kustomization resources")
}

// UpdateClusters correctly syncs the clusters' directories according to the config file
func TestUpdateClusters(t *testing.T) {
	testGroups := []v1alpha1.Group{
		{
			Name: "group-1",
			Clusters: []v1alpha1.Cluster{
				{
					Name: "cluster-1",
				},
				{
					Name: "cluster-2",
				},
			},
		},
		{
			Name: "group-2",
			Clusters: []v1alpha1.Cluster{
				{
					Name: "cluster-3",
				},
				{
					Name: "cluster-4",
				},
			},
		},
	}

	logger := logger.DisabledLogger{}
	config := v1alpha1.ClustersConfiguration{}
	config.Spec.Groups = testGroups
	testDirPath := t.TempDir()
	err := syncClusters(logger, &config, testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	expectedKustomization, err := os.ReadFile(testutils.GetTestFile("sync", "outputs", "default_import.yaml"))
	if !assert.NoError(t, err) {
		return
	}

	testutils.CompareFile(t, expectedKustomization, filepath.Join(testDirPath, "clusters/group-1/cluster-1/bases", konfig.DefaultKustomizationFileName()))
	testutils.CompareFile(t, expectedKustomization, filepath.Join(testDirPath, "clusters/group-1/cluster-2/bases", konfig.DefaultKustomizationFileName()))
	testutils.CompareFile(t, expectedKustomization, filepath.Join(testDirPath, "clusters/group-2/cluster-3/bases", konfig.DefaultKustomizationFileName()))
	testutils.CompareFile(t, expectedKustomization, filepath.Join(testDirPath, "clusters/group-2/cluster-4/bases", konfig.DefaultKustomizationFileName()))
}

// UpdateClusterModules returns the correct map of modules (w/o overrides)
func TestUpdateClusterModulesNoOverrides(t *testing.T) {
	logger := logger.DisabledLogger{}
	defaultModules := make(map[string]v1alpha1.Package)
	defaultModules["module3"] = v1alpha1.NewModule(
		t,
		"test-module3/test-flavour1",
		"1.0.0",
		false,
	)
	defaultModules["module2"] = v1alpha1.NewModule(
		t,
		"test-module2/test-flavour1",
		"1.0.0",
		false,
	)
	defaultModules["module1"] = v1alpha1.NewModule(
		t,
		"test-module1/test-flavour1",
		"1.0.0",
		false,
	)
	overrides := make(map[string]v1alpha1.Package)
	expectedOutput := make(map[string]v1alpha1.Package)
	expectedOutput["module3"] = v1alpha1.NewModule(
		t,
		"test-module3/test-flavour1",
		"1.0.0",
		false,
	)
	expectedOutput["module2"] = v1alpha1.NewModule(
		t,
		"test-module2/test-flavour1",
		"1.0.0",
		false,
	)
	expectedOutput["module1"] = v1alpha1.NewModule(
		t,
		"test-module1/test-flavour1",
		"1.0.0",
		false,
	)
	output := mergePackages(logger, defaultModules, overrides)
	assert.Equal(t, expectedOutput, output)
}

// UpdateClusterModules returns the correct map of modules (w/ overrides)
func TestUpdateClusterModules(t *testing.T) {
	logger := logger.DisabledLogger{}
	defaultModules := make(map[string]v1alpha1.Package)
	defaultModules["module3"] = v1alpha1.NewModule(
		t,
		"category/test-module3/test-flavour1",
		"1.0.0",
		false,
	)
	defaultModules["module2"] = v1alpha1.NewModule(
		t,
		"category/test-module2/test-flavour1",
		"1.0.0",
		false,
	)
	defaultModules["module1"] = v1alpha1.NewModule(
		t,
		"category/test-module1/test-flavour1",
		"1.0.0",
		false,
	)
	overrides := make(map[string]v1alpha1.Package)
	overrides["module3"] = v1alpha1.NewModule(
		t,
		"category/test-module3/test-flavour1",
		"1.0.1",
		false,
	)
	overrides["module2"] = v1alpha1.NewModule(
		t,
		"category/test-module2/test-flavour1",
		"",
		true,
	)
	overrides["module1"] = v1alpha1.NewModule(
		t,
		"category/test-module1/test-flavour1",
		"1.0.0",
		false,
	)
	output := mergePackages(logger, defaultModules, overrides)
	expectedOutput := make(map[string]v1alpha1.Package)
	expectedOutput["module1"] = v1alpha1.NewModule(
		t,
		"category/test-module1/test-flavour1",
		"1.0.0",
		false,
	)
	expectedOutput["module3"] = v1alpha1.NewModule(
		t,
		"category/test-module3/test-flavour1",
		"1.0.1",
		false,
	)
	assert.Equal(t, expectedOutput, output, "Unexpected map of modules")
}

// UpdateClusterAddOns returns the correct map of add-ons (w/o overrides)
func TestUpdateClusterAddOnsNoOverrides(t *testing.T) {
	defaultAddOns := make(map[string]v1alpha1.Package)
	defaultAddOns["addon1"] = v1alpha1.NewAddon(
		t,
		"test-addon1",
		"1.0.0",
		false,
	)
	defaultAddOns["addon2"] = v1alpha1.NewAddon(
		t,
		"test-addon2",
		"1.0.0",
		false,
	)
	overrides := make(map[string]v1alpha1.Package)
	expectedOutput := make(map[string]v1alpha1.Package)
	expectedOutput["addon1"] = v1alpha1.NewAddon(
		t,
		"test-addon1",
		"1.0.0",
		false,
	)
	expectedOutput["addon2"] = v1alpha1.NewAddon(
		t,
		"test-addon2",
		"1.0.0",
		false,
	)
	logger := logger.DisabledLogger{}
	output := mergePackages(logger, defaultAddOns, overrides)
	assert.Equal(t, expectedOutput, output)
}

// UpdateClusterAddOns returns the correct map of add-ons (w/ overrides)
func TestUpdateClusterAddOns(t *testing.T) {
	defaultAddOns := make(map[string]v1alpha1.Package)
	defaultAddOns["addon1"] = v1alpha1.NewAddon(
		t,
		"test-addon1",
		"1.0.0",
		false,
	)
	defaultAddOns["addon2"] = v1alpha1.NewAddon(
		t,
		"test-addon2",
		"1.0.0",
		false,
	)
	overrides := make(map[string]v1alpha1.Package)
	overrides["addon1"] = v1alpha1.NewAddon(
		t,
		"test-addon1",
		"1.0.1",
		false,
	)
	overrides["addon2"] = v1alpha1.NewAddon(
		t,
		"test-addon2",
		"",
		true,
	)
	logger := logger.DisabledLogger{}
	output := mergePackages(logger, defaultAddOns, overrides)
	expectedOutput := make(map[string]v1alpha1.Package)
	expectedOutput["addon1"] = v1alpha1.NewAddon(
		t,
		"test-addon1",
		"1.0.1",
		false,
	)
	assert.Equal(t, expectedOutput, output, "Unexpected map of add-ons")
}

// Sync correctly updates the project according to the configuration file (w/o overrides)
func TestSyncNoOverrides(t *testing.T) {
	logger := logger.DisabledLogger{}
	testDirPath := t.TempDir()
	configPath := testutils.GetTestFile("sync", "inputs", "basic.yaml")
	err := Sync(logger, testutils.FakeFilesGetter{Testing: t}, configPath, testDirPath, true)
	if !assert.NoError(t, err) {
		return
	}
	expectedKustomization, err := os.ReadFile(testutils.GetTestFile("sync", "outputs", "default_import.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, expectedKustomization, filepath.Join(testDirPath, "clusters/group-1/cluster-1/bases", konfig.DefaultKustomizationFileName()))
	testutils.CompareFile(t, expectedKustomization, filepath.Join(testDirPath, "clusters/group-1/cluster-2/bases", konfig.DefaultKustomizationFileName()))
	testutils.CompareFile(t, expectedKustomization, filepath.Join(testDirPath, "clusters/group-2/cluster-3/bases", konfig.DefaultKustomizationFileName()))
	testutils.CompareFile(t, expectedKustomization, filepath.Join(testDirPath, "clusters/group-2/cluster-4/bases", konfig.DefaultKustomizationFileName()))
}

// Sync correctly updates the project according to the configuration file (w/ overrides)
func TestSync(t *testing.T) {
	logger := logger.DisabledLogger{}
	testDirPath := t.TempDir()
	configPath := testutils.GetTestFile("sync", "inputs", "advanced.yaml")
	config, err := utils.ReadConfig(configPath)
	if !assert.NoError(t, err) {
		return
	}
	defaultSpec := *config.Spec.DeepCopy()
	err = Sync(logger, testutils.FakeFilesGetter{Testing: t}, configPath, testDirPath, true)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, defaultSpec, config.Spec, "The config spec has changed: EXPECTED %+v\n ACTUAL %+v\n", defaultSpec, config.Spec)
	allGroups, err := os.ReadFile(testutils.GetTestFile("sync", "outputs", "all_groups.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, allGroups, filepath.Join(testDirPath, utils.AllGroupsDirPath, utils.BasesDir, konfig.DefaultKustomizationFileName()))
	cluster1, err := os.ReadFile(testutils.GetTestFile("sync", "outputs", "advanced_g1c1.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, cluster1, filepath.Join(testDirPath, "clusters/group-1/cluster-1/bases", konfig.DefaultKustomizationFileName()))
	cluster2, err := os.ReadFile(testutils.GetTestFile("sync", "outputs", "advanced_g1c2.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, cluster2, filepath.Join(testDirPath, "clusters/group-1/cluster-2/bases", konfig.DefaultKustomizationFileName()))
	cluster3, err := os.ReadFile(testutils.GetTestFile("sync", "outputs", "advanced_g2c3.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, cluster3, filepath.Join(testDirPath, "clusters/group-2/cluster-3/bases", konfig.DefaultKustomizationFileName()))
	cluster4, err := os.ReadFile(testutils.GetTestFile("sync", "outputs", "default_import.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, cluster4, filepath.Join(testDirPath, "clusters/group-2/cluster-4/bases", konfig.DefaultKustomizationFileName()))
}
