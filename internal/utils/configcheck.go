// Copyright 2022 Mia-Platform

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// 	http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"errors"
	"path"

	"github.com/mia-platform/vab/internal/logger"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"golang.org/x/exp/slices"
)

const (
	maxArgs      = 2
	defaultScope = "default"
)

// GetBuildPath returns a list of target paths for the kustomize build command
func GetBuildPath(args []string, configPath string) ([]string, error) {
	config, err := ReadConfig(configPath)
	if err != nil {
		return []string{}, err
	}

	var groupName, clusterName string
	var targetPaths []string

	groupName = args[0]
	// Check if the specified group exists in the group array.
	// The IndexFunc call below compares the group name passed as argument
	// against the names of the groups in the group array. If there isn't a
	// match, IndexFunc returns -1
	groupIdx := slices.IndexFunc(config.Spec.Groups, func(g v1alpha1.Group) bool { return g.Name == groupName })
	if groupIdx == -1 {
		return []string{}, errors.New("Group " + groupName + " not found in configuration")
	}
	group := config.Spec.Groups[groupIdx]
	// The second arg, if present, contains the name of the cluster to build.
	// The usage of IndexFunc for clusters is similar to that mentioned above.
	if len(args) == maxArgs {
		clusterName = args[1]
		clusterIdx := slices.IndexFunc(config.Spec.Groups[groupIdx].Clusters, func(c v1alpha1.Cluster) bool { return c.Name == clusterName })
		if clusterIdx == -1 {
			return []string{}, errors.New("Cluster " + clusterName + " not found in configuration")
		}
		targetPaths = append(targetPaths, path.Join(groupName, clusterName))
	} else {
		// If no cluster is specified and the group exists, return all the paths to
		// the clusters in the group.
		for _, cluster := range group.Clusters {
			targetPaths = append(targetPaths, path.Join(groupName, cluster.Name))
		}
	}

	return targetPaths, nil
}

// ValidateConfig reads the config file and checks for errors/inconsistencies
func ValidateConfig(logger logger.LogInterface, configPath string) int {
	logger.V(0).Infof("Reading the configuration...")

	code := 0

	config, readErr := ReadConfig(configPath)
	if readErr != nil {
		logger.Warnf("[error] error while parsing the configuration file: %v", readErr)
		return 1
	}

	if config != nil {
		checkTypeMeta(logger, &config.TypeMeta, &code)
		logger.V(5).Infof("Checking TypeMeta for config ended with %d", code)
		checkModules(logger, &config.Spec.Modules, "", &code)
		logger.V(5).Infof("Checking configuration modules ended with %d", code)
		checkAddOns(logger, &config.Spec.AddOns, "", &code)
		logger.V(5).Infof("Checking configuration add-ons ended with %d", code)
		checkGroups(logger, &config.Spec.Groups, &code)
		logger.V(5).Infof("Checking configuration groups add-ons ended with %d", code)
	}

	if code == 0 {
		logger.V(0).Info("The configuration is valid!")
	} else {
		logger.V(0).Info("The configuration is invalid.")
	}

	return code
}

// checkTypeMeta checks the file's Kind and APIVersion
func checkTypeMeta(logger logger.LogInterface, config *v1alpha1.TypeMeta, code *int) {
	if config.Kind != v1alpha1.Kind {
		logger.V(0).Infof("[error] wrong kind: %s - expected: %s", config.Kind, v1alpha1.Kind)
		*code = 1
	}
	if config.APIVersion != v1alpha1.Version {
		logger.V(0).Infof("[error] wrong version: %s - expected: %s", config.APIVersion, v1alpha1.Version)
		*code = 1
	}
}

// checkModules checks the modules listed in the config file
func checkModules(logger logger.LogInterface, modules *map[string]v1alpha1.Module, scope string, code *int) {
	if scope == "" {
		scope = defaultScope
	}
	if len(*modules) == 0 {
		logger.V(0).Infof("[warn][%s] no module found: check the config file if this behavior is unexpected", scope)
	} else {
		for m := range *modules {
			// TODO: add check for modules' uniqueness (only one flavor per module is allowed)
			if (*modules)[m].Disable {
				logger.V(0).Infof("[info][%s] disabling module %s", scope, m)
			} else {
				if (*modules)[m].Version == "" {
					logger.V(0).Infof("[error][%s] missing version of module %s", scope, m)
					*code = 1
				}
				if (*modules)[m].Weight == 0 {
					logger.V(0).Infof("[warn][%s] missing weight of module %s: setting default (0)", scope, m)
				}
			}
		}
	}
}

// checkAddOns checks the add-ons listed in the config file
func checkAddOns(logger logger.LogInterface, addons *map[string]v1alpha1.AddOn, scope string, code *int) {
	if scope == "" {
		scope = defaultScope
	}
	if len(*addons) == 0 {
		logger.V(0).Infof("[warn][%s] no add-on found: check the config file if this behavior is unexpected", scope)
	} else {
		for m := range *addons {
			if (*addons)[m].Disable {
				logger.V(0).Infof("[info][%s] disabling add-on %s", scope, m)
			} else if (*addons)[m].Version == "" {
				logger.V(0).Infof("[error][%s] missing version of add-on %s", scope, m)
				*code = 1
			}
		}
	}
}

// checkGroups checks the cluster groups listed in the config file
func checkGroups(logger logger.LogInterface, groups *[]v1alpha1.Group, code *int) {
	if len(*groups) == 0 {
		logger.V(0).Info("[warn] no group found: check the config file if this behavior is unexpected")
	} else {
		for _, g := range *groups {
			groupName := g.Name
			if groupName == "" {
				logger.V(0).Info("[error] please specify a valid name for each group")
				groupName = "undefined"
				*code = 1
			}
			group := g
			checkClusters(logger, &group, groupName, code)
			logger.V(5).Infof("Checking group %s clusters ended with %d", groupName, *code)
		}
	}
}

// checkClusters checks the clusters of a group
func checkClusters(logger logger.LogInterface, group *v1alpha1.Group, groupName string, code *int) {
	if len(group.Clusters) == 0 {
		logger.V(0).Infof("[warn][%s] no cluster found in group: check the config file if this behavior is unexpected", groupName)
	} else {
		for _, cluster := range group.Clusters {
			clusterName := cluster.Name
			if clusterName == "" {
				logger.V(0).Infof("[error][%s] missing cluster name in group: please specify a valid name for each cluster", groupName)
				*code = 1
				clusterName = "undefined"
			}
			if cluster.Context == "" {
				logger.V(0).Infof("[error][%s/%s] missing cluster context: please specify a valid context for each cluster", groupName, clusterName)
				*code = 1
			}
			scope := groupName + "/" + clusterName
			checkModules(logger, &cluster.Modules, scope, code)
			logger.V(5).Infof("Checking cluster %s modules ended with %d", scope, *code)
			checkAddOns(logger, &cluster.AddOns, scope, code)
			logger.V(5).Infof("Checking cluster %s add-on ended with %d", scope, *code)
		}
	}
}
