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
	"path"
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
	testModule := v1alpha1.Module{
		Version: "1.0.0",
		Weight:  2,
	}
	outputFiles, err := ClonePackages(logger, "test-module1/test-flavour1", testModule, testutils.FakeFilesGetter{Testing: t})
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
		git.NewFile("./modules/test-module1/test-flavour1/file1.yaml", "./modules/test-module1", *fakeWorktree),
		git.NewFile("./modules/test-module1/test-flavour1/file2.yaml", "./modules/test-module1", *fakeWorktree),
		git.NewFile("./modules/test-module1/test-flavour1/file1.yaml", "./modules/test-module1", *fakeWorktree),
	}
	testDirPath := t.TempDir()
	err := MoveToDisk(logger, input, "test-module1/test-flavour1", testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.FileExists(t, path.Join(testDirPath, "test-flavour1/file1.yaml"), "Mock file 1 does not exist on disk")
	assert.FileExists(t, path.Join(testDirPath, "test-flavour1/file2.yaml"), "Mock file 2 does not exist on disk")
	assert.FileExists(t, path.Join(testDirPath, "test-flavour1/file1.yaml"), "Mock file 3 does not exist on disk")
}

// UpdateModules syncs new modules without errors
func TestUpdateModules(t *testing.T) {
	logger := logger.DisabledLogger{}
	modules := make(map[string]v1alpha1.Module)
	modules["test-module1/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  4,
	}
	modules["test-module2/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  1,
	}
	modules["test-module3/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  3,
		Disable: true,
	}
	testDirPath := t.TempDir()
	err := UpdateModules(logger, modules, testDirPath, testutils.FakeFilesGetter{Testing: t})
	if !assert.NoError(t, err) {
		return
	}
}

// UpdateAddOns syncs new modules without errors
func TestUpdateAddOns(t *testing.T) {
	logger := logger.DisabledLogger{}
	addons := make(map[string]v1alpha1.AddOn)
	addons["test-addon1"] = v1alpha1.AddOn{
		Version: "1.0.0",
	}
	addons["test-addon2"] = v1alpha1.AddOn{
		Version: "1.0.0",
		Disable: true,
	}
	testDirPath := t.TempDir()
	err := UpdateAddOns(logger, addons, testDirPath, testutils.FakeFilesGetter{Testing: t})
	if !assert.NoError(t, err) {
		return
	}
}

// UpdateBases correctly updates the resources list in the all-groups kustomization
func TestUpdateBasesAllGroups(t *testing.T) {
	testDirPath := t.TempDir()
	logger := logger.DisabledLogger{}
	targetPath := path.Join(testDirPath, utils.AllGroupsDirPath)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		return
	}
	modules := make(map[string]v1alpha1.Module)
	modules["test-module3-1.0.0/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  4,
	}
	modules["test-module2-1.0.0/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  1,
	}
	modules["test-module1-1.0.0/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  3,
	}
	addons := make(map[string]v1alpha1.AddOn)
	addons["test-addon1-1.0.0"] = v1alpha1.AddOn{
		Version: "1.0.0",
	}
	addons["test-addon2-1.0.0"] = v1alpha1.AddOn{
		Version: "1.0.0",
	}
	config := v1alpha1.ClustersConfiguration{}
	config.Spec.Modules = modules
	config.Spec.AddOns = addons
	err := UpdateBases(logger, testutils.FakeFilesGetter{Testing: t}, testDirPath, targetPath, modules, addons, &config, true)
	if !assert.NoError(t, err) {
		return
	}
	expectedKustomization, err := os.ReadFile(testutils.GetTestFile("sync", "outputs", "all_groups.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, expectedKustomization, path.Join(targetPath, utils.BasesDir, konfig.DefaultKustomizationFileName()))
}

// UpdateBases correctly initializes the resources list in a cluster's kustomization
func TestUpdateBasesCluster(t *testing.T) {
	testDirPath := t.TempDir()
	logger := logger.DisabledLogger{}
	targetPath := path.Join(testDirPath, "groups/group-1/cluster-1")
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		return
	}
	config := v1alpha1.ClustersConfiguration{}
	err := UpdateBases(logger, testutils.FakeFilesGetter{Testing: t}, testDirPath, targetPath, nil, nil, &config, true)
	if !assert.NoError(t, err) {
		return
	}
	expectedKustomization, err := os.ReadFile(testutils.GetTestFile("sync", "outputs", "default_import.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, expectedKustomization, path.Join(targetPath, utils.BasesDir, konfig.DefaultKustomizationFileName()))
}

// CheckClusterPath creates and returns the correct path to a missing cluster folder
// and creates the missing kustomization file
func TestCreateClusterPath(t *testing.T) {
	testDirPath := t.TempDir()
	clusterPath, err := CheckClusterPath("test-cluster", testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	expectedPath := path.Join(testDirPath, utils.ClustersDirName, "test-cluster")
	assert.Equal(t, expectedPath, clusterPath, "Wrong path to cluster")
	assert.DirExists(t, clusterPath, "The cluster directory does not exist")
	kustomizationPath := path.Join(clusterPath, konfig.DefaultKustomizationFileName())
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
	expectedPath := path.Join(testDirPath, utils.ClustersDirName, "test-cluster")
	if err := os.MkdirAll(expectedPath, os.ModePerm); err != nil {
		return
	}
	clusterPath, err := CheckClusterPath("test-cluster", testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expectedPath, clusterPath, "Wrong path to cluster")
	kustomizationPath := path.Join(clusterPath, konfig.DefaultKustomizationFileName())
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
	expectedPath := path.Join(testDirPath, utils.ClustersDirName, "test-cluster")
	if err := os.MkdirAll(expectedPath, os.ModePerm); err != nil {
		return
	}
	kustomizationPath := path.Join(expectedPath, konfig.DefaultKustomizationFileName())
	expectedKustomization, err := os.ReadFile(testutils.GetTestFile("sync", "misc", "cluster_kustomization.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	err = os.WriteFile(kustomizationPath, expectedKustomization, os.ModePerm)
	if !assert.NoError(t, err) {
		return
	}
	clusterPath, err := CheckClusterPath("test-cluster", testDirPath)
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
	expectedPath := path.Join(testDirPath, utils.ClustersDirName, "test-cluster")
	if err := os.MkdirAll(expectedPath, os.ModePerm); err != nil {
		return
	}
	// create a kustomization file without "bases" among the resources
	kustomizationPath := path.Join(expectedPath, konfig.DefaultKustomizationFileName())
	kustomization, err := os.ReadFile(testutils.GetTestFile("sync", "misc", "missing_bases.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	err = os.WriteFile(kustomizationPath, kustomization, os.ModePerm)
	if !assert.NoError(t, err) {
		return
	}
	clusterPath, err := CheckClusterPath("test-cluster", testDirPath)
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
	logger := logger.DisabledLogger{}
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
	config := v1alpha1.ClustersConfiguration{}
	config.Spec.Groups = testGroups
	testDirPath := t.TempDir()
	err := UpdateClusters(logger, testutils.FakeFilesGetter{Testing: t}, &config, testDirPath, true)
	if !assert.NoError(t, err) {
		return
	}
	expectedKustomization, err := os.ReadFile(testutils.GetTestFile("sync", "outputs", "default_import.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, expectedKustomization, path.Join(testDirPath, "clusters/group-1/cluster-1/bases", konfig.DefaultKustomizationFileName()))
	testutils.CompareFile(t, expectedKustomization, path.Join(testDirPath, "clusters/group-1/cluster-2/bases", konfig.DefaultKustomizationFileName()))
	testutils.CompareFile(t, expectedKustomization, path.Join(testDirPath, "clusters/group-2/cluster-3/bases", konfig.DefaultKustomizationFileName()))
	testutils.CompareFile(t, expectedKustomization, path.Join(testDirPath, "clusters/group-2/cluster-4/bases", konfig.DefaultKustomizationFileName()))
}

// UpdateClusterModules returns the correct map of modules (w/o overrides)
func TestUpdateClusterModulesNoOverrides(t *testing.T) {
	defaultModules := make(map[string]v1alpha1.Module)
	defaultModules["test-module3/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  4,
	}
	defaultModules["test-module2/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  1,
	}
	defaultModules["test-module1/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  3,
	}
	overrides := make(map[string]v1alpha1.Module)
	output := UpdateClusterModules(overrides, defaultModules)
	assert.Equal(t, 0, len(output), "The output should be nil")
}

// UpdateClusterModules returns the correct map of modules (w/ overrides)
func TestUpdateClusterModules(t *testing.T) {
	defaultModules := make(map[string]v1alpha1.Module)
	defaultModules["test-module3/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  4,
	}
	defaultModules["test-module2/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  1,
	}
	defaultModules["test-module1/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  3,
	}
	overrides := make(map[string]v1alpha1.Module)
	overrides["test-module3/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.1",
		Weight:  4,
	}
	overrides["test-module2/test-flavour1"] = v1alpha1.Module{
		Disable: true,
	}
	overrides["test-module1/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  1,
	}
	output := UpdateClusterModules(overrides, defaultModules)
	expectedOutput := make(map[string]v1alpha1.Module)
	expectedOutput["test-module1-1.0.0/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  1,
	}
	expectedOutput["test-module3-1.0.1/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.1",
		Weight:  4,
	}
	assert.Equal(t, expectedOutput, output, "Unexpected map of modules")
}

// UpdateClusterAddOns returns the correct map of add-ons (w/o overrides)
func TestUpdateClusterAddOnsNoOverrides(t *testing.T) {
	defaultAddOns := make(map[string]v1alpha1.AddOn)
	defaultAddOns["test-addon1"] = v1alpha1.AddOn{
		Version: "1.0.0",
	}
	defaultAddOns["test-addon2"] = v1alpha1.AddOn{
		Version: "1.0.0",
	}
	overrides := make(map[string]v1alpha1.AddOn)
	output := UpdateClusterAddOns(overrides, defaultAddOns)
	assert.Equal(t, 0, len(output), "The output should be nil")
}

// UpdateClusterAddOns returns the correct map of add-ons (w/ overrides)
func TestUpdateClusterAddOns(t *testing.T) {
	defaultAddOns := make(map[string]v1alpha1.AddOn)
	defaultAddOns["test-addon1"] = v1alpha1.AddOn{
		Version: "1.0.0",
	}
	defaultAddOns["test-addon2"] = v1alpha1.AddOn{
		Version: "1.0.0",
	}
	overrides := make(map[string]v1alpha1.AddOn)
	overrides["test-addon1"] = v1alpha1.AddOn{
		Version: "1.0.1",
	}
	overrides["test-addon2"] = v1alpha1.AddOn{
		Disable: true,
	}
	output := UpdateClusterAddOns(overrides, defaultAddOns)
	expectedOutput := make(map[string]v1alpha1.AddOn)
	expectedOutput["test-addon1-1.0.1"] = v1alpha1.AddOn{
		Version: "1.0.1",
	}
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
	testutils.CompareFile(t, expectedKustomization, path.Join(testDirPath, "clusters/group-1/cluster-1/bases", konfig.DefaultKustomizationFileName()))
	testutils.CompareFile(t, expectedKustomization, path.Join(testDirPath, "clusters/group-1/cluster-2/bases", konfig.DefaultKustomizationFileName()))
	testutils.CompareFile(t, expectedKustomization, path.Join(testDirPath, "clusters/group-2/cluster-3/bases", konfig.DefaultKustomizationFileName()))
	testutils.CompareFile(t, expectedKustomization, path.Join(testDirPath, "clusters/group-2/cluster-4/bases", konfig.DefaultKustomizationFileName()))
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
	testutils.CompareFile(t, allGroups, path.Join(testDirPath, utils.AllGroupsDirPath, utils.BasesDir, konfig.DefaultKustomizationFileName()))
	cluster1, err := os.ReadFile(testutils.GetTestFile("sync", "outputs", "advanced_g1c1.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, cluster1, path.Join(testDirPath, "clusters/group-1/cluster-1/bases", konfig.DefaultKustomizationFileName()))
	cluster2, err := os.ReadFile(testutils.GetTestFile("sync", "outputs", "advanced_g1c2.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, cluster2, path.Join(testDirPath, "clusters/group-1/cluster-2/bases", konfig.DefaultKustomizationFileName()))
	cluster3, err := os.ReadFile(testutils.GetTestFile("sync", "outputs", "advanced_g2c3.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, cluster3, path.Join(testDirPath, "clusters/group-2/cluster-3/bases", konfig.DefaultKustomizationFileName()))
	cluster4, err := os.ReadFile(testutils.GetTestFile("sync", "outputs", "default_import.yaml"))
	if !assert.NoError(t, err) {
		return
	}
	testutils.CompareFile(t, cluster4, path.Join(testDirPath, "clusters/group-2/cluster-4/bases", konfig.DefaultKustomizationFileName()))
}
