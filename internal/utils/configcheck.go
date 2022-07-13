package utils

import (
	"errors"
	"path"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"golang.org/x/exp/slices"
)

const (
	maxArgs = 2
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
