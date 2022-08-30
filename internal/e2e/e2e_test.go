//go:build e2e
// +build e2e

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

package e2e_test

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/mia-platform/vab/internal/cmd"
	"github.com/mia-platform/vab/pkg/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

const (
	testProjectName     = "test-e2e"
	sampleKustomization = `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
  - example.yaml`
	patchKustomization = `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
  - bases
patches:
  - path: patch.yaml`
)

var log logger.LogInterface
var cfg *rest.Config
var dynamicClient dynamic.Interface
var testEnv *envtest.Environment
var testDirPath string
var rootCmd *cobra.Command
var configPath string
var projectPath string
var clustersDirPath string
var allGroupsDirPath string
var sampleModulePath string
var sampleAddOnPath string

var _ = BeforeSuite(func() {
	By("setting up the test environment...", func() {
		logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))
		// initialize test environment
		useCluster := true
		testEnv = &envtest.Environment{
			UseExistingCluster:       &useCluster,
			AttachControlPlaneOutput: true,
		}
		var err error
		fmt.Println("Starting test environment...")
		cfg, err = testEnv.Start()
		Expect(err).ToNot(HaveOccurred())
		Expect(cfg).ToNot(BeNil())

		dynamicClient = dynamic.NewForConfigOrDie(cfg)

		// initialize vab logger and root command
		log = logger.DisabledLogger{}
		rootCmd = cmd.NewRootCommand()

		// initialize paths
		testDirPath = os.TempDir()
		// testDirPath = "."
		defer os.RemoveAll(testDirPath)
		projectPath = path.Join(testDirPath, testProjectName)
		configPath = path.Join(projectPath, "config.yaml")
		clustersDirPath = path.Join(projectPath, "clusters")
		allGroupsDirPath = path.Join(clustersDirPath, "all-groups")
		sampleModulePath = path.Join(projectPath, "vendors", "modules", "example", "sample-module1")
		sampleAddOnPath = path.Join(projectPath, "vendors", "add-ons", "sample-addon1")
	})
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment...")
	if testEnv != nil {
		err := testEnv.Stop()
		Expect(err).NotTo(HaveOccurred())
	}
}, 60)

var _ = Describe("setup vab project", func() {
	Context("initialize new project", func() {
		It("creates a preliminar directory structure", func() {
			rootCmd.SetArgs([]string{
				"init",
				fmt.Sprintf("--path=%s", testDirPath),
				fmt.Sprintf("--name=%s", testProjectName),
			})
			err := rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			// check that the path /tmpdir/test-e2e/clusters/all-groups/bases exists
			info, err := os.Stat(path.Join(allGroupsDirPath, "bases"))
			Expect(err).NotTo(HaveOccurred())
			Expect(info.IsDir()).To(BeTrue())
			// check that the path /tmpdir/test-e2e/clusters/all-groups/custom-resources exists
			info, err = os.Stat((path.Join(allGroupsDirPath, "custom-resources")))
			Expect(err).NotTo(HaveOccurred())
			Expect(info.IsDir()).To(BeTrue())
		})
	})
	Context("config with module (w/o overrides)", func() {
		It("validates the config file without errors", func() {
			config := `kind: ClustersConfiguration
apiVersion: vab.mia-platform.eu/v1alpha1
name: test-project
spec:
  modules:
    example/sample-module1:
      version: 0.1.0
      weight: 1
  addOns: {}
  groups:
  - name: group1
    clusters:
    - name: cluster1
      context: kind-kind`
			err := os.WriteFile(configPath, []byte(config), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			rootCmd.SetArgs([]string{
				"validate",
				fmt.Sprintf("--config=%s", configPath),
			})
			err = rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("syncs the project without errors", func() {
			rootCmd.SetArgs([]string{
				"sync",
				fmt.Sprintf("--path=%s", projectPath),
				"--dry-run",
			})
			err := rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("builds the configuration without errors", func() {
			sampleFile := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample-module1
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: sample-module1
  template:
    metadata:
      labels:
        app: sample-module1
    spec:
      containers:
      - image: k8s.gcr.io/echoserver:1.4
        name: echoserver
        ports:
        - containerPort: 8080`
			err := os.MkdirAll(sampleModulePath, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(path.Join(sampleModulePath, "example.yaml"), []byte(sampleFile), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(path.Join(sampleModulePath, "kustomization.yaml"), []byte(sampleKustomization), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			rootCmd.SetArgs([]string{
				"build",
				"group1",
				projectPath,
				fmt.Sprintf("--path=%s", projectPath),
			})
			err = rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("config with module (w/ override)", func() {
		It("validates the config file without errors", func() {
			config := `kind: ClustersConfiguration
apiVersion: vab.mia-platform.eu/v1alpha1
name: test-project
spec:
  modules:
    example/sample-module1:
      version: 0.1.0
      weight: 1
  addOns: {}
  groups:
  - name: group1
    clusters:
    - name: cluster1
      context: kind-kind
      modules:
        example/sample-module1:
          version: 0.1.1
          weight: 1`
			err := os.WriteFile(configPath, []byte(config), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			rootCmd.SetArgs([]string{
				"validate",
				fmt.Sprintf("--config=%s", configPath),
			})
			err = rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("syncs the project without errors", func() {
			rootCmd.SetArgs([]string{
				"sync",
				fmt.Sprintf("--path=%s", projectPath),
				"--dry-run",
			})
			err := rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("builds the configuration without errors", func() {
			rootCmd.SetArgs([]string{
				"build",
				"group1",
				"cluster1",
				projectPath,
			})
			err := rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("applies the configuration to the kind cluster", func() {
			rootCmd.SetArgs([]string{
				"apply",
				"group1",
				"cluster1",
				projectPath,
			})
			err := rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())

			depsGvr := schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			}
			dep, err := dynamicClient.Resource(depsGvr).Namespace("default").Get(context.Background(), "sample-module1", v1.GetOptions{})
			Expect(dep).NotTo(BeNil())
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("config with module (w/ override and patch)", func() {
		It("syncs the project without errors", func() {
			patch := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample-module1
spec:
  replicas: 2`
			pathToCluster := path.Join(clustersDirPath, "group1", "cluster1")
			err := os.WriteFile(path.Join(pathToCluster, "patch.yaml"), []byte(patch), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(path.Join(pathToCluster, "kustomization.yaml"), []byte(patchKustomization), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			rootCmd.SetArgs([]string{
				"sync",
				fmt.Sprintf("--path=%s", projectPath),
				"--dry-run",
			})
			err = rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
			k, err := os.ReadFile(path.Join(pathToCluster, "kustomization.yaml"))
			Expect(err).NotTo(HaveOccurred())
			Expect(k).To(BeEquivalentTo([]byte(patchKustomization)))
		})
		It("builds the configuration without errors", func() {
			rootCmd.SetArgs([]string{
				"build",
				"group1",
				"cluster1",
				projectPath,
			})
			err := rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("updates the resources on the kind cluster", func() {
			rootCmd.SetArgs([]string{
				"apply",
				"group1",
				"cluster1",
				projectPath,
			})
			err := rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())

			depsGvr := schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			}
			dep, err := dynamicClient.Resource(depsGvr).Namespace("default").Get(context.Background(), "sample-module1", v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(dep).NotTo(BeNil())
			Expect(dep.Object["spec"].(map[string]interface{})["replicas"]).Should(BeNumerically("==", 2))
		})
	})
	Context("config with module (w/ override and patch) and add-on (w/o overrides)", func() {
		It("validates the config file without errors", func() {
			config := `kind: ClustersConfiguration
apiVersion: vab.mia-platform.eu/v1alpha1
name: test-project
spec:
  modules:
    example/sample-module1:
      version: 0.1.0
      weight: 1
  addOns:
    sample-addon1:
      version: 0.1.0
  groups:
  - name: group1
    clusters:
    - name: cluster1
      context: kind-kind
      modules:
        example/sample-module1:
          version: 0.1.1
          weight: 1`
			err := os.WriteFile(configPath, []byte(config), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			rootCmd.SetArgs([]string{
				"validate",
				fmt.Sprintf("--config=%s", configPath),
			})
			err = rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("syncs the project without errors", func() {
			rootCmd.SetArgs([]string{
				"sync",
				fmt.Sprintf("--path=%s", projectPath),
				"--dry-run",
			})
			err := rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("builds the configuration without errors", func() {
			sampleFile := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: sample-addon1
  namespace: default
spec:
  replicas: 1
  selector:
    matchLabels:
      app: sample-addon1
  template:
    metadata:
      labels:
        app: sample-addon1
    spec:
      containers:
      - image: k8s.gcr.io/echoserver:1.4
        name: echoserver
        ports:
        - containerPort: 8080`
			err := os.MkdirAll(sampleAddOnPath, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(path.Join(sampleAddOnPath, "example.yaml"), []byte(sampleFile), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(path.Join(sampleAddOnPath, "kustomization.yaml"), []byte(sampleKustomization), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			rootCmd.SetArgs([]string{
				"build",
				"group1",
				"cluster1",
				projectPath,
			})
			err = rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("updates the resources on the kind cluster", func() {
			rootCmd.SetArgs([]string{
				"apply",
				"group1",
				"cluster1",
				projectPath,
			})
			err := rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())

			depsGvr := schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			}
			depMod, err := dynamicClient.Resource(depsGvr).Namespace("default").Get(context.Background(), "sample-module1", v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(depMod).NotTo(BeNil())
			Expect(depMod.Object["spec"].(map[string]interface{})["replicas"]).Should(BeNumerically("==", 2))
			depAddOn, err := dynamicClient.Resource(depsGvr).Namespace("default").Get(context.Background(), "sample-addon1", v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(depAddOn).NotTo(BeNil())
		})
	})
	Context("config with module (w/ override and patch) and add-on (w/ override)", func() {
		It("validates the config file without errors", func() {
			config := `kind: ClustersConfiguration
apiVersion: vab.mia-platform.eu/v1alpha1
name: test-project
spec:
  modules:
    example/sample-module1:
      version: 0.1.0
      weight: 1
  addOns: {}
  groups:
  - name: group1
    clusters:
    - name: cluster1
      context: kind-kind
      modules:
        example/sample-module1:
          version: 0.1.1
          weight: 1
      addOns:
        sample-addon1:
          version: 0.1.1`
			err := os.WriteFile(configPath, []byte(config), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			rootCmd.SetArgs([]string{
				"validate",
				fmt.Sprintf("--config=%s", configPath),
			})
			err = rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("syncs the project without errors", func() {
			rootCmd.SetArgs([]string{
				"sync",
				fmt.Sprintf("--path=%s", projectPath),
				"--dry-run",
			})
			err := rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("builds the configuration without errors", func() {
			rootCmd.SetArgs([]string{
				"build",
				"group1",
				"cluster1",
				projectPath,
			})
			err := rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("applies the configuration to the kind cluster", func() {
			rootCmd.SetArgs([]string{
				"apply",
				"group1",
				"cluster1",
				projectPath,
			})
			err := rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())

			depsGvr := schema.GroupVersionResource{
				Group:    "apps",
				Version:  "v1",
				Resource: "deployments",
			}
			depMod, err := dynamicClient.Resource(depsGvr).Namespace("default").Get(context.Background(), "sample-module1", v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(depMod).NotTo(BeNil())
			Expect(depMod.Object["spec"].(map[string]interface{})["replicas"]).Should(BeNumerically("==", 2))
			depAddOn, err := dynamicClient.Resource(depsGvr).Namespace("default").Get(context.Background(), "sample-addon1", v1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
			Expect(depAddOn).NotTo(BeNil())
		})
	})
})
