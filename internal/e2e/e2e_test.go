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
	"fmt"
	"os"
	"path"

	"github.com/mia-platform/vab/internal/cmd"
	"github.com/mia-platform/vab/pkg/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"
)

const (
	testProjectName     = "test-e2e"
	sampleKustomization = `kind: Kustomization
apiVersion: kustomize.config.k8s.io/v1beta1
resources:
  - example.yaml`
)

var log logger.LogInterface
var testDirPath string
var rootCmd *cobra.Command
var configPath string
var projectPath string
var clustersDirPath string
var allGroupsDirPath string
var sampleModulePath string
var err error

var _ = BeforeSuite(func() {

	log = logger.DisabledLogger{}
	testDirPath = GinkgoT().TempDir()
	rootCmd = cmd.NewRootCommand()
	projectPath = path.Join(testDirPath, testProjectName)
	configPath = path.Join(testDirPath, "config.yaml")
	clustersDirPath = path.Join(projectPath, "clusters")
	allGroupsDirPath = path.Join(clustersDirPath, "all-groups")
	sampleModulePath = path.Join(projectPath, "vendors", "modules", "ingress", "traefik-base")

})

var _ = Describe("setup vab project", func() {
	Context("initialize new project", func() {
		It("creates a preliminar directory structure", func() {
			rootCmd.SetArgs([]string{
				"init",
				fmt.Sprintf("--path=%s", testDirPath),
				fmt.Sprintf("--name=%s", testProjectName),
			})
			err = rootCmd.Execute()
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
	Context("validate a sample configuration", func() {
		It("returns that the configuration is valid", func() {
			config := `kind: ClustersConfiguration
apiVersion: vab.mia-platform.eu/v1alpha1
name: test-project
spec:
  modules:
    ingress/traefik-base:
      version: 0.1.0
      weight: 1
  addOns: {}
  groups:
  - name: test-group1
    clusters:
    - name: test-g1c1
      context: kind-kind`
			err = os.WriteFile(configPath, []byte(config), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			rootCmd.SetArgs([]string{
				"validate",
				fmt.Sprintf("--config=%s", configPath),
			})
			err = rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
	})
	Context("sync and build", func() {
		It("updates the directories according to the config", func() {
			rootCmd.SetArgs([]string{
				"sync",
				fmt.Sprintf("--path=%s", projectPath),
				"--dry-run",
			})
			err = rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
		It("builds the configuration", func() {
			sampleFile := `apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "1"
  labels:
    app: ingress-traefik
  name: ingress-traefik
  namespace: default
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: ingress-traefik
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      creationTimestamp: null
      labels:
        app: ingress-traefik
    spec:
      containers:
      - image: k8s.gcr.io/echoserver:1.4
        imagePullPolicy: IfNotPresent
        name: echoserver
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      terminationGracePeriodSeconds: 30`

			err = os.MkdirAll(sampleModulePath, os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(path.Join(sampleModulePath, "example.yaml"), []byte(sampleFile), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			err = os.WriteFile(path.Join(sampleModulePath, "kustomization.yaml"), []byte(sampleKustomization), os.ModePerm)
			Expect(err).NotTo(HaveOccurred())
			rootCmd.SetArgs([]string{
				"build",
				"test-group1",
				projectPath,
				fmt.Sprintf("--path=%s", projectPath),
			})
			err = rootCmd.Execute()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
