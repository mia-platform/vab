package utils

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"

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

// TODO: consider using a logger w/ levels instead of fmt.Print

// ValidateConfig reads the config file and checks for errors/inconsistencies
func ValidateConfig(configPath string, writer io.Writer) int {
	if writer == nil {
		writer = os.Stdout
	}
	fmt.Fprint(writer, "Reading the configuration...\n")

	code := 0

	config, readErr := ReadConfig(configPath)
	if readErr != nil {
		fmt.Fprintf(writer, "[error] error while parsing the configuration file: %v\n", readErr)
		code = 1
	}
	if config != nil {
		checkTypeMeta(&config.TypeMeta, writer, &code)
		checkModules(&config.Spec.Modules, "", writer, &code)
		checkAddOns(&config.Spec.AddOns, "", writer, &code)
		checkGroups(&config.Spec.Groups, writer, &code)
	}

	if code == 0 {
		fmt.Fprint(writer, "The configuration is valid!\n")
	} else {
		fmt.Fprint(writer, "The configuration is invalid.\n")
	}

	return code
}

// checkTypeMeta checks the file's Kind and APIVersion
func checkTypeMeta(config *v1alpha1.TypeMeta, writer io.Writer, code *int) {
	if config.Kind != v1alpha1.Kind {
		fmt.Fprintf(writer, "[error] wrong kind: %s - expected: %s\n", config.Kind, v1alpha1.Kind)
		*code = 1
	}
	if config.APIVersion != v1alpha1.Version {
		fmt.Fprintf(writer, "[error] wrong version: %s - expected: %s\n", config.APIVersion, v1alpha1.Version)
		*code = 1
	}
}

// checkModules checks the modules listed in the config file
func checkModules(modules *map[string]v1alpha1.Module, scope string, writer io.Writer, code *int) {
	if scope == "" {
		scope = defaultScope
	}
	if len(*modules) == 0 {
		fmt.Fprintf(writer, "[warn][%s] no module found: check the config file if this behavior is unexpected\n", scope)
	} else {
		for m := range *modules {
			if (*modules)[m].Disable {
				fmt.Fprintf(writer, "[warn][%s] disabling module %s\n", scope, m)
			} else {
				if (*modules)[m].Version == "" {
					fmt.Fprintf(writer, "[error][%s] missing version of module %s\n", scope, m)
					*code = 1
				}
				if (*modules)[m].Weight == 0 {
					fmt.Fprintf(writer, "[warn][%s] missing weight of module %s: setting default (0)\n", scope, m)
				}
			}
		}
	}
}

// checkAddOns checks the add-ons listed in the config file
func checkAddOns(addons *map[string]v1alpha1.AddOn, scope string, writer io.Writer, code *int) {
	if scope == "" {
		scope = defaultScope
	}
	if len(*addons) == 0 {
		fmt.Fprintf(writer, "[warn][%s] no add-on found: check the config file if this behavior is unexpected\n", scope)
	} else {
		for m := range *addons {
			if (*addons)[m].Version == "" {
				fmt.Fprintf(writer, "[error][%s] missing version of add-on %s\n", scope, m)
				*code = 1
			}
			if (*addons)[m].Disable {
				fmt.Fprintf(writer, "[warn][%s] disabling add-on %s\n", scope, m)
			}
		}
	}
}

// checkGroups checks the cluster groups listed in the config file
func checkGroups(groups *[]v1alpha1.Group, writer io.Writer, code *int) {
	if len(*groups) == 0 {
		fmt.Fprint(writer, "[warn] no group found: check the config file if this behavior is unexpected\n")
	} else {
		for _, g := range *groups {
			groupName := g.Name
			if groupName == "" {
				fmt.Fprint(writer, "[error] please specify a valid name for each group\n")
				groupName = "undefined"
				*code = 1
			}
			group := g
			checkClusters(&group, groupName, writer, code)
		}
	}
}

// checkClusters checks the clusters of a group
func checkClusters(group *v1alpha1.Group, groupName string, writer io.Writer, code *int) {
	if len(group.Clusters) == 0 {
		fmt.Fprintf(writer, "[warn][%s] no cluster found in group: check the config file if this behavior is unexpected\n", groupName)
	} else {
		for _, cluster := range group.Clusters {
			clusterName := cluster.Name
			if clusterName == "" {
				fmt.Fprintf(writer, "[error][%s] missing cluster name in group: please specify a valid name for each cluster\n", groupName)
				*code = 1
				clusterName = "undefined"
			}
			if cluster.Context == "" {
				fmt.Fprintf(writer, "[error][%s/%s] missing cluster context: please specify a valid context for each cluster\n", groupName, clusterName)
				*code = 1
			}
			scope := groupName + "/" + clusterName
			checkModules(&cluster.Modules, scope, writer, code)
			checkAddOns(&cluster.AddOns, scope, writer, code)
		}
	}
}
