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

package validate

import (
	"fmt"
	"io"

	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/mia-platform/vab/pkg/logger"
)

const (
	defaultScope = "default"
)

func ConfigurationFile(logger logger.LogInterface, configurationPath string, writer io.Writer) error {
	code := 0

	config, readErr := utils.ReadConfig(configurationPath)
	if readErr != nil {
		return fmt.Errorf("error while parsing the configuration file: %v", readErr)
	}

	feedbackString := checkTypeMeta(logger, &config.TypeMeta, &code)
	logger.V(5).Writef("Checking TypeMeta for config ended with %d", code)
	feedbackString += checkModules(logger, &config.Spec.Modules, "", &code)
	logger.V(5).Writef("Checking configuration modules ended with %d", code)
	feedbackString += checkAddOns(logger, &config.Spec.AddOns, "", &code)
	logger.V(5).Writef("Checking configuration add-ons ended with %d", code)
	feedbackString += checkGroups(logger, &config.Spec.Groups, &code)
	logger.V(5).Writef("Checking configuration groups add-ons ended with %d", code)

	fmt.Fprint(writer, feedbackString)
	if code > 0 {
		return fmt.Errorf("the configuration is invalid")
	}

	fmt.Fprint(writer, "The configuration is valid!\n")
	return nil
}

// checkTypeMeta checks the file's Kind and APIVersion
func checkTypeMeta(logger logger.LogInterface, config *v1alpha1.TypeMeta, code *int) string {
	outString := ""
	if config.Kind != v1alpha1.Kind {
		outString += fmt.Sprintf("[error] wrong kind: %s - expected: %s\n", config.Kind, v1alpha1.Kind)
		*code = 1
	}

	if config.APIVersion != v1alpha1.Version {
		outString += fmt.Sprintf("[error] wrong version: %s - expected: %s\n", config.APIVersion, v1alpha1.Version)
		*code = 1
	}

	return outString
}

// checkModules checks the modules listed in the config file
func checkModules(logger logger.LogInterface, modules *map[string]v1alpha1.Module, scope string, code *int) string {
	if scope == "" {
		scope = defaultScope
	}

	outString := ""
	if len(*modules) == 0 {
		outString += fmt.Sprintf("[warn][%s] no module found: check the config file if this behavior is unexpected\n", scope)
		return outString
	}

	for m := range *modules {
		// TODO: add check for modules' uniqueness (only one flavor per module is allowed)
		if (*modules)[m].Disable {
			outString += fmt.Sprintf("[info][%s] disabling module %s\n", scope, m)
			continue
		}

		if (*modules)[m].Version == "" {
			outString += fmt.Sprintf("[error][%s] missing version of module %s\n", scope, m)
			*code = 1
		}
		if (*modules)[m].Weight == 0 {
			outString += fmt.Sprintf("[warn][%s] missing weight of module %s: setting default (0)\n", scope, m)
		}
	}

	return outString
}

// checkAddOns checks the add-ons listed in the config file
func checkAddOns(logger logger.LogInterface, addons *map[string]v1alpha1.AddOn, scope string, code *int) string {
	if scope == "" {
		scope = defaultScope
	}

	outString := ""
	if len(*addons) == 0 {
		outString += fmt.Sprintf("[warn][%s] no add-on found: check the config file if this behavior is unexpected\n", scope)
		return outString
	}

	for m := range *addons {
		if (*addons)[m].Disable {
			outString += fmt.Sprintf("[info][%s] disabling add-on %s\n", scope, m)
		} else if (*addons)[m].Version == "" {
			outString += fmt.Sprintf("[error][%s] missing version of add-on %s\n", scope, m)
			*code = 1
		}
	}
	return outString
}

// checkGroups checks the cluster groups listed in the config file
func checkGroups(logger logger.LogInterface, groups *[]v1alpha1.Group, code *int) string {
	outString := ""

	if len(*groups) == 0 {
		outString += "[warn] no group found: check the config file if this behavior is unexpected\n"
		return outString
	}

	for _, g := range *groups {
		groupName := g.Name
		if groupName == "" {
			outString += "[error] please specify a valid name for each group\n"
			groupName = "undefined"
			*code = 1
		}

		group := g
		outString += checkClusters(logger, &group, groupName, code)
		logger.V(5).Writef("Checking group %s clusters ended with %d\n", groupName, *code)
	}

	return outString
}

// checkClusters checks the clusters of a group
func checkClusters(logger logger.LogInterface, group *v1alpha1.Group, groupName string, code *int) string {
	outString := ""
	if len(group.Clusters) == 0 {
		outString += fmt.Sprintf("[warn][%s] no cluster found in group: check the config file if this behavior is unexpected\n", groupName)
		return outString
	}

	for _, cluster := range group.Clusters {
		clusterName := cluster.Name
		if clusterName == "" {
			outString += fmt.Sprintf("[error][%s] missing cluster name in group: please specify a valid name for each cluster\n", groupName)
			*code = 1
			clusterName = "undefined"
		}

		if cluster.Context == "" {
			outString += fmt.Sprintf("[error][%s/%s] missing cluster context: please specify a valid context for each cluster\n", groupName, clusterName)
			*code = 1
		}

		scope := groupName + "/" + clusterName
		outString += checkModules(logger, &cluster.Modules, scope, code)
		logger.V(5).Writef("Checking cluster %s modules ended with %d", scope, *code)
		outString += checkAddOns(logger, &cluster.AddOns, scope, code)
		logger.V(5).Writef("Checking cluster %s add-on ended with %d", scope, *code)
	}

	return outString
}
