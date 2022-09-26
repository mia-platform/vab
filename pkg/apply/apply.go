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

package apply

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	jpl "github.com/mia-platform/jpl/deploy"
	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	vabBuild "github.com/mia-platform/vab/pkg/build"
	"github.com/mia-platform/vab/pkg/logger"
	"golang.org/x/exp/slices"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/util"
	k8syaml "sigs.k8s.io/yaml"
)

const (
	filesPermissions  fs.FileMode = 0600
	folderPermissions fs.FileMode = 0700
)

func ApplyWithJpl(logger logger.LogInterface, configPath string, outputDir string, isDryRun bool, groupName string, clusterName string, contextPath string, options *jpl.Options) error {
	cleanedContextPath := path.Clean(contextPath)
	contextInfo, err := os.Stat(cleanedContextPath)
	if err != nil {
		return err
	}
	if !contextInfo.IsDir() {
		return fmt.Errorf("the target path %s is not a directory", cleanedContextPath)
	}
	targetPaths, err := utils.BuildPaths(configPath, groupName, clusterName)
	if err != nil {
		return err
	}
	for _, clusterPath := range targetPaths {
		buffer := new(bytes.Buffer)
		// pathArray := strings.Split(clusterPath, "/")
		// cluster := pathArray[len(pathArray)-1]

		targetPath := path.Join(cleanedContextPath, clusterPath)
		if err := vabBuild.RunKustomizeBuild(targetPath, buffer); err != nil {
			logger.V(5).Writef("Error building kustomize in %s", targetPath)
			return err
		}

		// context, err := getContext(configPath, groupName, cluster)
		if err != nil {
			return fmt.Errorf("error searching for context: %s", err)
		}

		resources, err := NewResourcesFromBuffer(buffer.Bytes(), "default")
		if err != nil {
			logger.V(5).Writef("Error generating resources in %s", targetPath)
			return err
		}

		apply := jpl.DecorateDefaultApplyFunction()

		if err := jpl.Deploy(jpl.InitRealK8sClients(options), "default", resources, jpl.DeployConfig{}, apply); err != nil {
			logger.V(5).Writef("Error applying resources in %s", targetPath)
			return err
		}
	}
	return nil
}

// Apply builds the selected Kustomize resources and apply them in the context specified. The resources built are saved in files at the specified path.
func Apply(logger logger.LogInterface, configPath string, outputDir string, isDryRun bool, groupName string, clusterName string, contextPath string) error {
	cleanedContextPath := path.Clean(contextPath)
	contextInfo, err := os.Stat(cleanedContextPath)
	if err != nil {
		return err
	}
	if !contextInfo.IsDir() {
		return fmt.Errorf("the target path %s is not a directory", cleanedContextPath)
	}

	targetPaths, err := utils.BuildPaths(configPath, groupName, clusterName)
	if err != nil {
		return err
	}

	for _, clusterPath := range targetPaths {
		buffer := new(bytes.Buffer)
		pathArray := strings.Split(clusterPath, "/")
		cluster := pathArray[len(pathArray)-1]

		targetPath := path.Join(cleanedContextPath, clusterPath)
		if err := vabBuild.RunKustomizeBuild(targetPath, buffer); err != nil {
			logger.V(5).Writef("Error building kustomize in %s", targetPath)
			return err
		}

		crdFilename := cluster + "-crds"
		resourcesFilename := cluster + "-res"
		crdsFilePath := filepath.Join(outputDir, crdFilename)
		resourcesFilepath := filepath.Join(outputDir, resourcesFilename)

		err := createResourcesFiles(outputDir, crdsFilePath, resourcesFilepath, *buffer)
		if err != nil {
			return fmt.Errorf("error creating resource files: %s", err)
		}

		context, err := getContext(configPath, groupName, cluster)
		if err != nil {
			return fmt.Errorf("error searching for context: %s", err)
		}

		if _, err := os.Stat(crdsFilePath); err == nil {
			err = runKubectlApply(logger, crdsFilePath, context, isDryRun)
			if err != nil {
				return fmt.Errorf("error applying crds at %s: %s", crdsFilePath, err)
			}
		}

		if _, err := os.Stat(resourcesFilepath); err == nil {
			err = runKubectlApply(logger, resourcesFilepath, context, isDryRun)
			if err != nil {
				return fmt.Errorf("error applying resources at %s: %s", resourcesFilepath, err)
			}
		}
	}
	return nil
}

// getContext retrieves the context for the cluster/group from the config file.
func getContext(configPath string, groupName string, clusterName string) (string, error) {
	config, err := utils.ReadConfig(configPath)
	if err != nil {
		return "default", err
	}

	groupIdx := slices.IndexFunc(config.Spec.Groups, func(g v1alpha1.Group) bool { return g.Name == groupName })
	if groupIdx == -1 {
		return "default", errors.New("Group " + groupName + " not found in configuration")
	}

	clusterIdx := slices.IndexFunc(config.Spec.Groups[groupIdx].Clusters, func(c v1alpha1.Cluster) bool { return c.Name == clusterName })
	if clusterIdx == -1 {
		return "default", errors.New("Cluster " + clusterName + " not found in configuration")
	}

	return config.Spec.Groups[groupIdx].Clusters[clusterIdx].Context, nil
}

// runKubectlApply instantiates and executes the kubectl Apply command, with the correct parameters.
func runKubectlApply(logger logger.LogInterface, fileName string, context string, isDryRun bool) error {
	// default configflags
	configFlags := genericclioptions.NewConfigFlags(false)
	// the kubeconfig context used is equal to the fileName
	configFlags.Context = &context

	factory := util.NewFactory(configFlags)
	streams := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	args := []string{
		"-f",
		fileName,
		"--wait",
		"--server-side",
	}
	cmd := apply.NewCmdApply("kubectl", factory, streams)
	cmd.SetArgs(args)

	if !isDryRun {
		fmt.Println("Apply")
		err := cmd.Execute()
		if err != nil {
			return err
		}
	} else {
		logger.V(5).Writef("Skipping apply on ", fileName, "...")
	}

	return nil
}

// NewResourcesFromBuffer (temp) reimplements the NewResources fn from jpl library
func NewResourcesFromBuffer(stream []byte, namespace string) ([]jpl.Resource, error) {
	var resources []jpl.Resource
	re := regexp.MustCompile(`\n---\n`)
	for _, resourceYAML := range re.Split(string(stream), -1) {

		if len(resourceYAML) == 0 {
			continue
		}

		u := unstructured.Unstructured{Object: map[string]interface{}{}}
		if err := k8syaml.Unmarshal([]byte(resourceYAML), &u.Object); err != nil {
			return nil, err
		}
		gvk := u.GroupVersionKind()
		u.SetNamespace(namespace)

		resources = append(resources,
			jpl.Resource{
				Filepath:         "default",
				GroupVersionKind: &gvk,
				Object:           u,
			})
	}
	fmt.Printf("PRIMA: %+v\n", resources)
	resources = jpl.SortResourcesByKind(resources, nil)
	fmt.Printf("DOPO: %+v\n", resources)
	return resources, nil
}
