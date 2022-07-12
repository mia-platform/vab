package utils

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path"
	"testing"
)

const (
	testDeploymentFileName = "test.deployment.yaml"
	testDeployment         = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
spec:
  replicas: 1`
	testServiceFileName = "test.svc.yaml"
	testService         = `apiVersion: apps/v1
kind: Service
metadata:
  name: test
spec:
  ports:
  - protocol: TCP
    port: 80
    targetPort: 9376`
	testKustomization = `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- test.deployment.yaml
- test.svc.yaml`
	testConfig = `kind: ClustersConfiguration
apiVersion: vab.mia-platform.eu/v1alpha1
name: test
spec:
  modules: {}
  addOns: {}
  groups:
  - name: test-group
    clusters:
    - name: test-cluster
    - name: another-cluster
  - name: another-group`
	expectedResult = `apiVersion: apps/v1
kind: Service
metadata:
  name: test
spec:
  ports:
  - port: 80
    protocol: TCP
    targetPort: 9376
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test
spec:
  replicas: 1
`
)

// Test that the function returns the correct kustomized configuration
func TestRunKustomizeBuild(t *testing.T) {
	testDirPath := t.TempDir()

	if err := writeYamlFile(testConfig, path.Join(testDirPath, defaultConfigFileName)); err != nil {
		t.Fatal(err)
	}

	targetPath := path.Join(testDirPath, clustersDirName, testGroup, testCluster)
	if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(path.Join(targetPath, testDeploymentFileName), []byte(testDeployment), defaultFilePermissions); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path.Join(targetPath, testServiceFileName), []byte(testService), defaultFilePermissions); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path.Join(targetPath, kustomizationFileName), []byte(testKustomization), defaultFilePermissions); err != nil {
		t.Fatal(err)
	}

	buffer := new(bytes.Buffer)
	if err := RunKustomizeBuild(targetPath, buffer); err != nil {
		t.Fatal(err)
	}

	t.Log(buffer.String())
	if !bytes.Equal(buffer.Bytes(), []byte(expectedResult)) {
		t.Fatal("Unexpected Kustomize result.")
	}
}

// Returns an error if the path is invalid
func TestInvalidKustomizeBuildPath(t *testing.T) {
	buffer := new(bytes.Buffer)
	err := RunKustomizeBuild(invalidPath, buffer)
	if err == nil {
		t.Fatalf("No error was returned. Expected: %s", fs.ErrNotExist)
	}
	if !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("Unexpected error. Expected: %s, actual: %s", fs.ErrNotExist, err)
	}
}
