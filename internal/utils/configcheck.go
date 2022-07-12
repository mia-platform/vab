package utils

import (
	"errors"
	"path"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"golang.org/x/exp/slices"
)

// GetBuildPath returns the target path for the kustomize build command
func GetBuildPath(args []string, configPath string) (string, error) {

	if configPath == "" {
		configPath = defaultConfigFileName
	}
	config, err := ReadConfig(configPath)
	if err != nil {
		return "", err
	}

	var group, cluster string

	if len(args) < 1 {
		return "", errors.New("at least the cluster group is required")
	}
	if len(args) > 2 {
		return "", errors.New("too many args")
	} else {
		group = args[0]
		groupIdx := slices.IndexFunc(config.Spec.Groups, func(g v1alpha1.Group) bool { return g.Name == group })
		if groupIdx == -1 {
			return "", errors.New("Group " + group + " not found in configuration")
		}
		if len(args) == 2 {
			cluster = args[1]
			clusterIdx := slices.IndexFunc(config.Spec.Groups[groupIdx].Clusters, func(c v1alpha1.Cluster) bool { return c.Name == cluster })
			if clusterIdx == -1 {
				return "", errors.New("Cluster " + cluster + " not found in configuration")
			}
		}
	}

	targetPath := path.Join(clustersDirName, group, cluster)

	return targetPath, nil
}
