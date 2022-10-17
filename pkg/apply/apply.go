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
	"strings"
	"time"

	jpl "github.com/mia-platform/jpl/deploy"
	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	vabBuild "github.com/mia-platform/vab/pkg/build"
	"github.com/mia-platform/vab/pkg/logger"
	"golang.org/x/exp/slices"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	filesPermissions  fs.FileMode = 0600
	folderPermissions fs.FileMode = 0700
	defaultContext                = "default"
)

var (
	gvrCRDs = schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}
)

// Apply builds the cluster resources and applies them by calling the jpl deploy function
func Apply(logger logger.LogInterface, configPath string, isDryRun bool, groupName string, clusterName string, contextPath string, options *jpl.Options, crdStatusCheckRetries int) error {
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

		// context, err := getContext(configPath, groupName, cluster)
		if err != nil {
			return fmt.Errorf("error searching for context: %s", err)
		}

		k8sContext, err := getContext(configPath, groupName, cluster)
		if err != nil {
			return fmt.Errorf("error searching for context: %s", err)
		}
		options.Context = k8sContext
		clients := jpl.InitRealK8sClients(options)
		crds, resources, err := jpl.NewResourcesFromBuffer(buffer.Bytes(), "default", jpl.RealSupportedResourcesGetter{}, clients)
		if err != nil {
			logger.V(5).Writef("Error generating resources in %s", targetPath)
			return err
		}

		apply := jpl.DecorateDefaultApplyFunction()
		deployConfig := jpl.DeployConfig{}

		// if there are any CRDs, deploy them first
		if len(crds) != 0 {
			if err := jpl.Deploy(clients, "", crds, deployConfig, apply); err != nil {
				logger.V(5).Writef("Error applying CRDs in %s", targetPath)
				return fmt.Errorf("deploy of crds failed with error: %w", err)
			}
			// wait until all the CRDs satisfy the "Established" condition
			if err := checkCRDsStatus(clients, crdStatusCheckRetries); err != nil {
				logger.V(5).Writef("The check of CRDs status failed", targetPath)
				return fmt.Errorf("crds check failed with error: %w", err)
			}
		}

		if err := jpl.Deploy(clients, "", resources, jpl.DeployConfig{}, apply); err != nil {
			logger.V(5).Writef("Error applying resources in %s", targetPath)
			return fmt.Errorf("deploy of resources failed with error: %w", err)
		}
	}
	return nil
}

// checkCRDsStatus loops over the deployed CRDs to check whether the condition
// `Established` evaluates to true. If the condition is not met for any CRD
// before `retries` times, the function returns an error
func checkCRDsStatus(clients *jpl.K8sClients, retries int) error {
	var establishedCount int
	for ; retries > 0; retries-- {
		establishedCount = 0
		crdList, err := jpl.ListResources(gvrCRDs, clients)
		if err != nil && !apierrors.IsNotFound(err) {
			fmt.Printf("fails to check CRDs: %s", err)
			return err
		}
		for _, crd := range crdList.Items {
			crdStatus := crd.Object["status"]
			if crdStatus == nil {
				continue
			}
			crdConditions := crdStatus.(map[string]interface{})["conditions"]
			if crdConditions == nil {
				continue
			}
			for _, condition := range crdConditions.([]interface{}) {
				conditionType := condition.(map[string]interface{})["type"]
				conditionStatus := condition.(map[string]interface{})["status"]
				if conditionType == "Established" && conditionStatus == "True" {
					establishedCount++
				}
			}
		}
		if len(crdList.Items) == establishedCount {
			fmt.Printf("Established %d CRDs\n", establishedCount)
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("reached limit of max retries for CRDs status check")

}

// getContext retrieves the context for the cluster/group from the config file.
func getContext(configPath string, groupName string, clusterName string) (string, error) {
	config, err := utils.ReadConfig(configPath)
	if err != nil {
		return defaultContext, err
	}

	groupIdx := slices.IndexFunc(config.Spec.Groups, func(g v1alpha1.Group) bool { return g.Name == groupName })
	if groupIdx == -1 {
		return defaultContext, errors.New("Group " + groupName + " not found in configuration")
	}

	clusterIdx := slices.IndexFunc(config.Spec.Groups[groupIdx].Clusters, func(c v1alpha1.Cluster) bool { return c.Name == clusterName })
	if clusterIdx == -1 {
		return defaultContext, errors.New("Cluster " + clusterName + " not found in configuration")
	}

	return config.Spec.Groups[groupIdx].Clusters[clusterIdx].Context, nil
}
