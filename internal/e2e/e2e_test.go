// Copyright Mia srl
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build conformance

package e2e_test

// import (
// 	"context"
// 	"os"
// 	"path/filepath"

// 	jpl "github.com/mia-platform/jpl/deploy"
// 	"github.com/mia-platform/vab/internal/git"
// 	"github.com/mia-platform/vab/pkg/apply"
// 	"github.com/mia-platform/vab/pkg/logger"
// 	"github.com/mia-platform/vab/pkg/sync"
// 	. "github.com/onsi/ginkgo" //revive:disable-line:dot-imports
// 	. "github.com/onsi/gomega" //revive:disable-line:dot-imports
// 	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/runtime/schema"
// 	"k8s.io/client-go/dynamic"
// 	"k8s.io/client-go/rest"
// 	"k8s.io/client-go/tools/clientcmd"
// )

// const (
// 	crdDefaultRetries        = 2
// 	testProjectName          = "test-e2e"
// 	sipleModuleKustomization = `kind: Kustomization
// apiVersion: kustomize.config.k8s.io/v1beta1
// resources:
//   - example.yaml`
// 	moduleWithCRDsKustomization = `kind: Kustomization
// apiVersion: kustomize.config.k8s.io/v1beta1
// resources:
//   - example.yaml
//   - project.crd.yaml
//   - foobar.crd.yaml`
// 	addonKustomization = `kind: Component
// apiVersion: kustomize.config.k8s.io/v1alpha1
// patches:
//   - path: example.yaml`
// 	kustomizationPatch1 = `kind: Kustomization
// apiVersion: kustomize.config.k8s.io/v1beta1
// resources:
//   - bases
// patches:
//   - path: module.patch.yaml`
// 	kustomizationPatch2 = `kind: Kustomization
// apiVersion: kustomize.config.k8s.io/v1beta1
// resources:
//   - bases
// patches:
//   - path: addon.patch.yaml`
// )

// var log logger.LogInterface
// var jplClientsCluster1 dynamic.Interface
// var jplClientsCluster2 dynamic.Interface
// var options *jpl.Options
// var testDirPath string
// var configPath string
// var projectPath string
// var clustersDirPath string
// var allGroupsDirPath string
// var modulePath1 string
// var modulePath2 string
// var moduleOverridePath1 string
// var moduleOverridePath2 string
// var addOnPath string
// var addOnOverridePath string
// var depsGvr schema.GroupVersionResource

// var _ = BeforeSuite(func() {
// 	By("setting up the test environment...", func() {
// 		// initialize configs and clients for the test clusters
// 		homeDir, err := os.UserHomeDir()
// 		Expect(err).ToNot(HaveOccurred())

// 		kubeConfigPath := filepath.Join(homeDir, ".kube/config")

// 		cluster1Cfg, err := buildConfigFromFlags("kind-vab-cluster-1", kubeConfigPath)
// 		Expect(err).ToNot(HaveOccurred())
// 		Expect(cluster1Cfg).ToNot(BeNil())

// 		cluster2Cfg, err := buildConfigFromFlags("kind-vab-cluster-2", kubeConfigPath)
// 		Expect(err).ToNot(HaveOccurred())
// 		Expect(cluster2Cfg).ToNot(BeNil())

// 		jplClientsCluster1 = dynamic.NewForConfigOrDie(cluster1Cfg)
// 		jplClientsCluster2 = dynamic.NewForConfigOrDie(cluster2Cfg)

// 		options = jpl.NewOptions()
// 		options.Context = "kind-vab-cluster-1"

// 		// initialize global paths and vars
// 		testDirPath = os.TempDir()
// 		// testDirPath = "."
// 		projectPath = filepath.Join(testDirPath, testProjectName)
// 		configPath = filepath.Join(projectPath, "config.yaml")
// 		clustersDirPath = filepath.Join(projectPath, "clusters")
// 		allGroupsDirPath = filepath.Join(clustersDirPath, "all-groups")
// 		modulePath1 = filepath.Join(projectPath, "vendors", "modules", "module1-0.1.0", "flavor1")
// 		moduleOverridePath1 = filepath.Join(projectPath, "vendors", "modules", "module1-0.1.1", "flavor1")
// 		modulePath2 = filepath.Join(projectPath, "vendors", "modules", "module2-0.1.0", "flavor1")
// 		moduleOverridePath2 = filepath.Join(projectPath, "vendors", "modules", "module2-0.1.1", "flavor1")
// 		addOnPath = filepath.Join(projectPath, "vendors", "addons", "addon1-0.1.0")
// 		addOnOverridePath = filepath.Join(projectPath, "vendors", "addons", "addon1-0.1.1")
// 		depsGvr = schema.GroupVersionResource{
// 			Group:    "apps",
// 			Version:  "v1",
// 			Resource: "deployments",
// 		}

// 		// initialize project
// 		log = logger.DisabledLogger{}
// 		err = initProj.NewProject(log, testDirPath, testProjectName)
// 		Expect(err).NotTo(HaveOccurred())
// 	})
// }, 60)

// var _ = AfterSuite(func() {
// 	By("tearing down the test environment...", func() {
// 		os.RemoveAll(testDirPath)
// 	})
// }, 60)

// var _ = Describe("setup vab project", func() {
// 	Context("1 module (w/ override)", func() {
// 		It("syncs the project without errors", func() {
// 			config := `kind: ClustersConfiguration
// apiVersion: vab.mia-platform.eu/v1alpha1
// name: test-project
// spec:
//   modules:
//     module1/flavor1:
//       version: 0.1.0
//   addOns: {}
//   groups:
//   - name: group1
//     clusters:
//     - name: cluster1
//       context: kind-vab-cluster-1
//       modules:
//         module1/flavor1:
//           version: 0.1.1`
// 			err := os.WriteFile(configPath, []byte(config), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())

// 			err = sync.Sync(log, git.RealFilesGetter{}, configPath, projectPath, true)
// 			Expect(err).NotTo(HaveOccurred())
// 		})
// 		It("returns an error due to CRD status check", func() {
// 			sampleFile1 := `apiVersion: apps/v1
// kind: Deployment
// metadata:
//   name: module1-flavor1
//   namespace: default
// spec:
//   replicas: 1
//   selector:
//     matchLabels:
//       app: module1-flavor1
//   template:
//     metadata:
//       labels:
//         app: module1-flavor1
//         version: 0.1.0
//     spec:
//       containers:
//       - image: k8s.gcr.io/echoserver:1.4
//         name: echoserver
//         ports:
//         - containerPort: 8080`
// 			crd1 := `apiVersion: apiextensions.k8s.io/v1
// kind: CustomResourceDefinition
// metadata:
//   name: projects.example.vab.com
// spec:
//   group: example.vab.com
//   versions:
//     - name: v1
//       served: true
//       storage: true
//       schema:
//         openAPIV3Schema:
//           required: [spec]
//           type: object
//           properties:
//             spec:
//               required: [replicas]
//               type: object
//               properties:
//                 replicas:
//                   type: integer
//                   minimum: 1
//   scope: Namespaced
//   names:
//     plural: projects
//     singular: project
//     kind: Project
//     shortNames:
//     - pj`
// 			brokenCrd := `apiVersion: apiextensions.k8s.io/v1
// kind: CustomResourceDefinition
// metadata:
//   name: foobars.example.vab.com
// spec:
//   group: example.vab.com
//   versions:
//     - name: v1
//       served: true
//       storage: true
//       schema:
//         openAPIV3Schema:
//           required: [spec]
//           type: object
//           properties:
//             spec:
//               required: [replicas]
//               type: object
//               properties:
//                 replicas:
//                   type: integer
//                   minimum: 1
//   scope: Namespaced
//   names:
//     plural: foobars
//     singular: project
//     kind: FooBar
//     shortNames:
//     - pj`

// 			err := os.MkdirAll(modulePath1, os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(modulePath1, "example.yaml"), []byte(sampleFile1), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(modulePath1, "project.crd.yaml"), []byte(crd1), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(modulePath1, "foobar.crd.yaml"), []byte(brokenCrd), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(modulePath1, "kustomization.yaml"), []byte(moduleWithCRDsKustomization), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())

// 			sampleFile2 := `apiVersion: apps/v1
// kind: Deployment
// metadata:
//   name: module1-flavor1
//   namespace: default
// spec:
//   replicas: 1
//   selector:
//     matchLabels:
//       app: module1-flavor1
//   template:
//     metadata:
//       labels:
//         app: module1-flavor1
//         version: 0.1.1
//     spec:
//       containers:
//       - image: k8s.gcr.io/echoserver:1.4
//         name: echoserver
//         ports:
//         - containerPort: 8080`
// 			err = os.MkdirAll(moduleOverridePath1, os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(moduleOverridePath1, "example.yaml"), []byte(sampleFile2), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(moduleOverridePath1, "project.crd.yaml"), []byte(crd1), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(moduleOverridePath1, "foobar.crd.yaml"), []byte(brokenCrd), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(moduleOverridePath1, "kustomization.yaml"), []byte(moduleWithCRDsKustomization), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())

// 			err = apply.Apply(log, configPath, "group1", "cluster1", projectPath, options, crdDefaultRetries)
// 			Expect(err).To(HaveOccurred())
// 			Expect(err.Error()).To(BeIdenticalTo("crds check failed with error: reached limit of max retries for CRDs status check"))
// 		})
// 		It("applies the configuration to the kind cluster", func() {
// 			crd2 := `apiVersion: apiextensions.k8s.io/v1
// kind: CustomResourceDefinition
// metadata:
//   name: foobars.example.vab.com
// spec:
//   group: example.vab.com
//   versions:
//     - name: v1
//       served: true
//       storage: true
//       schema:
//         openAPIV3Schema:
//           required: [spec]
//           type: object
//           properties:
//             spec:
//               required: [replicas]
//               type: object
//               properties:
//                 replicas:
//                   type: integer
//                   minimum: 1
//   scope: Namespaced
//   names:
//     plural: foobars
//     singular: foobar
//     kind: FooBar
//     shortNames:
//     - fb`

// 			err := os.WriteFile(filepath.Join(modulePath1, "foobar.crd.yaml"), []byte(crd2), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())

// 			err = os.WriteFile(filepath.Join(moduleOverridePath1, "foobar.crd.yaml"), []byte(crd2), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())

// 			err = apply.Apply(log, configPath, "group1", "cluster1", projectPath, options, crdDefaultRetries)
// 			Expect(err).NotTo(HaveOccurred())

// 			dep, err := jplClientsCluster1.Resource(depsGvr).Namespace("default").Get(context.Background(), "module1-flavor1", v1.GetOptions{})
// 			Expect(dep).NotTo(BeNil())
// 			Expect(err).NotTo(HaveOccurred())
// 			modVer := dep.Object["spec"].(map[string]interface{})["template"].(map[string]interface{})["metadata"].(map[string]interface{})["labels"].(map[string]interface{})["version"]
// 			Expect(modVer).To(BeIdenticalTo("0.1.1"))
// 		})
// 	})
// 	Context("1 module (w/ override and patch)", func() {
// 		It("updates the resources on the kind cluster", func() {
// 			patch := `apiVersion: apps/v1
// kind: Deployment
// metadata:
//   name: module1-flavor1
// spec:
//   replicas: 2`
// 			pathToCluster := filepath.Join(clustersDirPath, "group1", "cluster1")
// 			err := os.WriteFile(filepath.Join(pathToCluster, "module.patch.yaml"), []byte(patch), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(pathToCluster, "kustomization.yaml"), []byte(kustomizationPatch1), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())

// 			err = apply.Apply(log, configPath, "group1", "cluster1", projectPath, options, crdDefaultRetries)
// 			Expect(err).NotTo(HaveOccurred())

// 			dep, err := jplClientsCluster1.Resource(depsGvr).Namespace("default").Get(context.Background(), "module1-flavor1", v1.GetOptions{})
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(dep).NotTo(BeNil())
// 			Expect(dep.Object["spec"].(map[string]interface{})["replicas"]).Should(BeNumerically("==", 2))
// 		})
// 	})
// 	Context("1 module (w/ override and patch), 1 add-on (w/o overrides)", func() {
// 		It("syncs the project without errors", func() {
// 			config := `kind: ClustersConfiguration
// apiVersion: vab.mia-platform.eu/v1alpha1
// name: test-project
// spec:
//   modules:
//     module1/flavor1:
//       version: 0.1.0
//   addOns:
//     addon1:
//       version: 0.1.0
//   groups:
//   - name: group1
//     clusters:
//     - name: cluster1
//       context: kind-vab-cluster-1
//       modules:
//         module1/flavor1:
//           version: 0.1.1`
// 			err := os.WriteFile(configPath, []byte(config), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())

// 			err = sync.Sync(log, git.RealFilesGetter{}, configPath, projectPath, true)
// 			Expect(err).NotTo(HaveOccurred())
// 		})
// 		It("updates the resources on the kind cluster", func() {
// 			sampleFile := `apiVersion: apps/v1
// kind: Deployment
// metadata:
//   name: module1-flavor1
// spec:
//   selector:
//     matchLabels:
//       app: module1-flavor1
//   template:
//     spec:
//       containers:
//       - image: k8s.gcr.io/echoserver:1.4
//         name: sidecar
//         ports:
//         - containerPort: 8080`
// 			err := os.MkdirAll(addOnPath, os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(addOnPath, "example.yaml"), []byte(sampleFile), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(addOnPath, "kustomization.yaml"), []byte(addonKustomization), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())

// 			err = apply.Apply(log, configPath, "group1", "cluster1", projectPath, options, crdDefaultRetries)
// 			Expect(err).NotTo(HaveOccurred())

// 			depMod, err := jplClientsCluster1.Resource(depsGvr).Namespace("default").Get(context.Background(), "module1-flavor1", v1.GetOptions{})
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(depMod).NotTo(BeNil())
// 			// module patched
// 			replicas := depMod.Object["spec"].(map[string]interface{})["replicas"]
// 			Expect(replicas).Should(BeNumerically("==", 2))
// 			// add-on deployed
// 			containersCount := len(depMod.Object["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{}))
// 			Expect(containersCount).Should(BeNumerically("==", 2))
// 		})
// 	})
// 	Context("1 module (w/ override and patch), 1 and add-on (w/ override)", func() {
// 		It("syncs the project without errors", func() {
// 			config := `kind: ClustersConfiguration
// apiVersion: vab.mia-platform.eu/v1alpha1
// name: test-project
// spec:
//   modules:
//     module1/flavor1:
//       version: 0.1.0
//   addOns: {}
//   groups:
//   - name: group1
//     clusters:
//     - name: cluster1
//       context: kind-vab-cluster-1
//       modules:
//         module1/flavor1:
//           version: 0.1.1
//       addOns:
//         addon1:
//           version: 0.1.1`
// 			err := os.WriteFile(configPath, []byte(config), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())

// 			err = sync.Sync(log, git.RealFilesGetter{}, configPath, projectPath, true)
// 			Expect(err).NotTo(HaveOccurred())
// 		})
// 		It("applies the configuration to the kind cluster", func() {
// 			sampleFile := `apiVersion: apps/v1
// kind: Deployment
// metadata:
//   name: module1-flavor1
// spec:
//   selector:
//     matchLabels:
//       app: module1-flavor1
//   template:
//     spec:
//       containers:
//       - image: k8s.gcr.io/echoserver:1.4
//         name: sidecar-v2
//         ports:
//         - containerPort: 8080`
// 			err := os.MkdirAll(addOnOverridePath, os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(addOnOverridePath, "example.yaml"), []byte(sampleFile), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(addOnOverridePath, "kustomization.yaml"), []byte(addonKustomization), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())

// 			err = apply.Apply(log, configPath, "group1", "cluster1", projectPath, options, crdDefaultRetries)
// 			Expect(err).NotTo(HaveOccurred())

// 			depMod, err := jplClientsCluster1.Resource(depsGvr).Namespace("default").Get(context.Background(), "module1-flavor1", v1.GetOptions{})
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(depMod).NotTo(BeNil())
// 			// module patched
// 			replicas := depMod.Object["spec"].(map[string]interface{})["replicas"]
// 			Expect(replicas).Should(BeNumerically("==", 2))
// 			// add-on deployed
// 			containersCount := len(depMod.Object["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{}))
// 			Expect(containersCount).Should(BeNumerically("==", 2))
// 			// add-on overridden
// 			containerName := depMod.Object["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[0].(map[string]interface{})["name"]
// 			Expect(containerName).To(BeIdenticalTo("sidecar-v2"))
// 		})
// 	})
// 	Context("1 module, 1 add-on (w/ overrides and patches)", func() {
// 		It("syncs the project without errors", func() {
// 			patch := `apiVersion: apps/v1
// kind: Deployment
// metadata:
//   name: module1-flavor1
// spec:
//   replicas: 3
//   template:
//     spec:
//       containers:
//       - name: sidecar-v2
//         ports:
//         - containerPort: 9000`
// 			pathToCluster := filepath.Join(clustersDirPath, "group1", "cluster1")
// 			err := os.WriteFile(filepath.Join(pathToCluster, "addon.patch.yaml"), []byte(patch), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(pathToCluster, "kustomization.yaml"), []byte(kustomizationPatch2), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())

// 			err = sync.Sync(log, git.RealFilesGetter{}, configPath, projectPath, true)
// 			Expect(err).NotTo(HaveOccurred())
// 		})
// 		It("updates the resources on the kind cluster", func() {
// 			err := apply.Apply(log, configPath, "group1", "cluster1", projectPath, options, crdDefaultRetries)
// 			Expect(err).NotTo(HaveOccurred())

// 			depMod, err := jplClientsCluster1.Resource(depsGvr).Namespace("default").Get(context.Background(), "module1-flavor1", v1.GetOptions{})
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(depMod).NotTo(BeNil())
// 			// module patched
// 			replicas := depMod.Object["spec"].(map[string]interface{})["replicas"]
// 			Expect(replicas).Should(BeNumerically("==", 3))
// 			// add-on deployed
// 			containersCount := len(depMod.Object["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{}))
// 			Expect(containersCount).Should(BeNumerically("==", 2))
// 			// add-on patched
// 			newSidecarPort := depMod.Object["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[0].(map[string]interface{})["ports"].([]interface{})[0].(map[string]interface{})["containerPort"]
// 			Expect(newSidecarPort).Should(BeNumerically("==", 9000))
// 		})
// 	})
// 	Context("broken CRD", func() {

// 	})
// 	Context("2 clusters, same group", func() {
// 		It("syncs the project without errors", func() {
// 			// clean up cluster 1
// 			err := jplClientsCluster1.Resource(depsGvr).Namespace("default").Delete(context.Background(), "module1-flavor1", v1.DeleteOptions{})
// 			Expect(err).NotTo(HaveOccurred())
// 			_, err = jplClientsCluster1.Resource(depsGvr).Namespace("default").Get(context.Background(), "module1-flavor1", v1.GetOptions{})
// 			Expect(err).To(HaveOccurred())

// 			config := `kind: ClustersConfiguration
// apiVersion: vab.mia-platform.eu/v1alpha1
// name: test-project
// spec:
//   modules:
//     module1/flavor1:
//       version: 0.1.0
//     module2/flavor1:
//       version: 0.1.0
//   addOns:
//     addon1:
//       version: 0.1.0
//   groups:
//   - name: group1
//     clusters:
//     - name: cluster1
//       context: kind-vab-cluster-1
//       modules:
//         module1/flavor1:
//           version: 0.1.1
//       addOns:
//         addon1:
//           version: 0.1.1
//     - name: cluster2
//       context: kind-vab-cluster-2
//       modules:
//         module2/flavor1:
//           version: 0.1.1
//       addOns:
//         addon1:
//           disable: true`
// 			err = os.WriteFile(configPath, []byte(config), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())

// 			err = sync.Sync(log, git.RealFilesGetter{}, configPath, projectPath, true)
// 			Expect(err).NotTo(HaveOccurred())
// 		})
// 		It("applies the configuration to the correct cluster and context", func() {
// 			sampleFile1 := `apiVersion: apps/v1
// kind: Deployment
// metadata:
//   name: module2-flavor1
//   namespace: default
// spec:
//   replicas: 1
//   selector:
//     matchLabels:
//       app: module2-flavor1
//   template:
//     metadata:
//       labels:
//         app: module2-flavor1
//         version: 0.1.0
//     spec:
//       containers:
//       - image: k8s.gcr.io/echoserver:1.4
//         name: echoserver
//         ports:
//         - containerPort: 8080`
// 			err := os.MkdirAll(modulePath2, os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(modulePath2, "example.yaml"), []byte(sampleFile1), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(modulePath2, "kustomization.yaml"), []byte(sipleModuleKustomization), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())

// 			sampleFile2 := `apiVersion: apps/v1
// kind: Deployment
// metadata:
//   name: module2-flavor1
//   namespace: default
// spec:
//   replicas: 1
//   selector:
//     matchLabels:
//       app: module2-flavor1
//   template:
//     metadata:
//       labels:
//         app: module2-flavor1
//         version: 0.1.1
//     spec:
//       containers:
//       - image: k8s.gcr.io/echoserver:1.4
//         name: echoserver
//         ports:
//         - containerPort: 8080`
// 			err = os.MkdirAll(moduleOverridePath2, os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(moduleOverridePath2, "example.yaml"), []byte(sampleFile2), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())
// 			err = os.WriteFile(filepath.Join(moduleOverridePath2, "kustomization.yaml"), []byte(sipleModuleKustomization), os.ModePerm)
// 			Expect(err).NotTo(HaveOccurred())

// 			err = apply.Apply(log, configPath, "group1", "", projectPath, options, crdDefaultRetries)
// 			Expect(err).NotTo(HaveOccurred())

// 			// cluster 1: module1-flavor1 deployed and patched, addon1 deployed (replicas == 3)
// 			depMod, err := jplClientsCluster1.Resource(depsGvr).Namespace("default").Get(context.Background(), "module1-flavor1", v1.GetOptions{})
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(depMod).NotTo(BeNil())
// 			Expect(depMod.Object["spec"].(map[string]interface{})["replicas"]).Should(BeNumerically("==", 3))
// 			// cluster 1: addon1 patched
// 			newSidecarPort := depMod.Object["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{})[0].(map[string]interface{})["ports"].([]interface{})[0].(map[string]interface{})["containerPort"]
// 			Expect(newSidecarPort).Should(BeNumerically("==", 9000))

// 			// cluster 2: module2-flavor1 deployed and overridden
// 			depMod, err = jplClientsCluster2.Resource(depsGvr).Namespace("default").Get(context.Background(), "module2-flavor1", v1.GetOptions{})
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(depMod).NotTo(BeNil())
// 			modVer := depMod.Object["spec"].(map[string]interface{})["template"].(map[string]interface{})["metadata"].(map[string]interface{})["labels"].(map[string]interface{})["version"]
// 			Expect(modVer).To(BeIdenticalTo("0.1.1"))
// 			// cluster 2: no module patch, addon-1 disabled (1 replica, no sidecar container)
// 			depMod, err = jplClientsCluster2.Resource(depsGvr).Namespace("default").Get(context.Background(), "module1-flavor1", v1.GetOptions{})
// 			Expect(err).NotTo(HaveOccurred())
// 			Expect(depMod).NotTo(BeNil())
// 			Expect(depMod.Object["spec"].(map[string]interface{})["replicas"]).Should(BeNumerically("==", 1))
// 			containersCount := len(depMod.Object["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})["containers"].([]interface{}))
// 			Expect(containersCount).Should(BeNumerically("==", 1))
// 		})
// 	})
// })

// // buildConfigFromFlags supports the switch between multiple kubecontext
// func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
// 	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
// 		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
// 		&clientcmd.ConfigOverrides{
// 			CurrentContext: context,
// 		}).ClientConfig()
// }
