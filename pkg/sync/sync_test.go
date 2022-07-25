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
	"github.com/mia-platform/vab/internal/testutils"
	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/mia-platform/vab/pkg/logger"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/kustomize/api/konfig"
)

const (
	testPackageName                = "test-module1/test-flavour1"
	clusterName                    = "test-cluster"
	expectedKustomizationAllGroups = `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
  - vendors/modules/test-module2/test-flavour2
  - vendors/modules/test-module1/test-flavour1
  - vendors/modules/test-module3/test-flavour3
  - vendors/add-ons/test-addon1
  - vendors/add-ons/test-addon2
`
	expectedKustomization = `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
  - ../../../all-groups
`
)

// ClonePackages returns a list of mocked file pointers w/o errors
func TestClonePackage(t *testing.T) {
	logger := logger.DisabledLogger{}
	testModule := v1alpha1.Module{
		Version: "1.0.0",
		Weight:  2,
	}
	outputFiles, err := ClonePackages(logger, testPackageName, testModule, testutils.FakeFilesGetter{Testing: t})
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
		git.NewFile("./modules/test-module1/test-flavour2/file1.yaml", "./modules/test-module1", *fakeWorktree),
	}
	testDirPath := t.TempDir()
	err := MoveToDisk(logger, input, testPackageName, testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.FileExists(t, path.Join(testDirPath, "test-flavour1/file1.yaml"), "Mock file 1 does not exist on disk")
	assert.FileExists(t, path.Join(testDirPath, "test-flavour1/file2.yaml"), "Mock file 2 does not exist on disk")
	assert.FileExists(t, path.Join(testDirPath, "test-flavour2/file1.yaml"), "Mock file 3 does not exist on disk")
}

// UpdateModules syncs new modules without errors
func TestUpdateModules(t *testing.T) {
	logger := logger.DisabledLogger{}
	modules := make(map[string]v1alpha1.Module)
	modules["test-module1/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  4,
	}
	modules["test-module2/test-flavour2"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  1,
	}
	modules["test-module3/test-flavour3"] = v1alpha1.Module{
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

// UpdateAddons syncs new modules without errors
func TestUpdateAddons(t *testing.T) {
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
	err := UpdateAddons(logger, addons, testDirPath, testutils.FakeFilesGetter{Testing: t})
	if !assert.NoError(t, err) {
		return
	}
}

// UpdateBases correctly updates the resources list in the all-groups kustomization
func TestUpdateBasesAllGroups(t *testing.T) {
	testDirPath := t.TempDir()
	targetPath := path.Join(testDirPath, utils.AllGroupsDirPath)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		return
	}
	modules := make(map[string]v1alpha1.Module)
	modules["test-module3/test-flavour3"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  4,
	}
	modules["test-module2/test-flavour2"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  1,
	}
	modules["test-module1/test-flavour1"] = v1alpha1.Module{
		Version: "1.0.0",
		Weight:  3,
	}
	addons := make(map[string]v1alpha1.AddOn)
	addons["test-addon1"] = v1alpha1.AddOn{
		Version: "1.0.0",
	}
	addons["test-addon2"] = v1alpha1.AddOn{
		Version: "1.0.0",
	}
	err := UpdateBases(targetPath, modules, addons)
	if !assert.NoError(t, err) {
		return
	}
	compareFile(t, []byte(expectedKustomizationAllGroups), path.Join(targetPath, "bases", konfig.DefaultKustomizationFileName()))
}

// UpdateBases correctly updates the resources list in a cluster's kustomization
func TestUpdateBasesCluster(t *testing.T) {
	testDirPath := t.TempDir()
	targetPath := path.Join(testDirPath, "groups/group-1/cluster-1")
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		return
	}
	err := UpdateBases(targetPath, nil, nil)
	if !assert.NoError(t, err) {
		return
	}
	compareFile(t, []byte(expectedKustomization), path.Join(targetPath, "bases", konfig.DefaultKustomizationFileName()))
}

// GetClusterPath creates and returns the correct path to a missing cluster folder
func TestCreateClusterPath(t *testing.T) {
	testDirPath := t.TempDir()
	clusterPath, err := GetClusterPath(clusterName, testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	expectedPath := path.Join(testDirPath, utils.ClustersDirName, clusterName)
	assert.Equal(t, expectedPath, clusterPath, "wrong path to cluster")
}

// GetClusterPath returns the correct path to an existing cluster folder
func TestExistingClusterPath(t *testing.T) {
	testDirPath := t.TempDir()
	expectedPath := path.Join(testDirPath, utils.ClustersDirName, clusterName)
	if err := os.MkdirAll(expectedPath, os.ModePerm); err != nil {
		return
	}
	clusterPath, err := GetClusterPath(clusterName, testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, expectedPath, clusterPath, "wrong path to cluster")
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
	testDirPath := t.TempDir()
	err := UpdateClusters(&testGroups, testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	compareFile(t, []byte(expectedKustomization), path.Join(testDirPath, "clusters/group-1/cluster-1/bases", konfig.DefaultKustomizationFileName()))
	compareFile(t, []byte(expectedKustomization), path.Join(testDirPath, "clusters/group-1/cluster-2/bases", konfig.DefaultKustomizationFileName()))
	compareFile(t, []byte(expectedKustomization), path.Join(testDirPath, "clusters/group-2/cluster-3/bases", konfig.DefaultKustomizationFileName()))
	compareFile(t, []byte(expectedKustomization), path.Join(testDirPath, "clusters/group-2/cluster-4/bases", konfig.DefaultKustomizationFileName()))
}

// Sync correctly updates the project according to the configuration file
func TestSync(t *testing.T) {
	logger := logger.DisabledLogger{}
	testDirPath := t.TempDir()
	configPath := testutils.GetTestFile("utils", "test_sync.yaml")
	err := Sync(logger, testutils.FakeFilesGetter{Testing: t}, configPath, testDirPath)
	if !assert.NoError(t, err) {
		return
	}
	compareFile(t, []byte(expectedKustomization), path.Join(testDirPath, "clusters/group-1/cluster-1/bases", konfig.DefaultKustomizationFileName()))
	compareFile(t, []byte(expectedKustomization), path.Join(testDirPath, "clusters/group-1/cluster-2/bases", konfig.DefaultKustomizationFileName()))
	compareFile(t, []byte(expectedKustomization), path.Join(testDirPath, "clusters/group-2/cluster-3/bases", konfig.DefaultKustomizationFileName()))
	compareFile(t, []byte(expectedKustomization), path.Join(testDirPath, "clusters/group-2/cluster-4/bases", konfig.DefaultKustomizationFileName()))
}
