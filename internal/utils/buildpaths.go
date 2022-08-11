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
	"fmt"
	"os"
	"path"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"golang.org/x/exp/slices"
)

// BuildPaths returns a list of clusters paths based on configPath configuration
// selected by groupName and optionally by clusterName
func BuildPaths(configPath string, groupName string, clusterName string) ([]string, error) {
	config, err := ReadConfig(configPath)
	if err != nil {
		return []string{}, err
	}

	var targetPaths []string
	// Check if the specified group exists in the group array.
	// The IndexFunc call below compares the group name passed as argument
	// against the names of the groups in the group array. If there isn't a
	// match, IndexFunc returns -1
	groupIdx := slices.IndexFunc(config.Spec.Groups, func(g v1alpha1.Group) bool { return g.Name == groupName })
	if groupIdx == -1 {
		return nil, errors.New("Group " + groupName + " not found in configuration")
	}
	group := config.Spec.Groups[groupIdx]
	// The second arg, if present, contains the name of the cluster to build.
	// The usage of IndexFunc for clusters is similar to that mentioned above.
	if clusterName != "" {
		clusterIdx := slices.IndexFunc(config.Spec.Groups[groupIdx].Clusters, func(c v1alpha1.Cluster) bool { return c.Name == clusterName })
		if clusterIdx == -1 {
			return nil, errors.New("Cluster " + clusterName + " not found in configuration")
		}
		targetPaths = append(targetPaths, path.Join(ClustersDirName, groupName, clusterName))
	} else {
		// If no cluster is specified and the group exists, return all the paths to
		// the clusters in the group.
		for _, cluster := range group.Clusters {
			targetPaths = append(targetPaths, path.Join(ClustersDirName, groupName, cluster.Name))
		}
	}

	return targetPaths, nil
}

// ValidatePath checks if the path exists and creates it eventually
func ValidatePath(targetPath string) error {
	if _, err := os.Stat(targetPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
				return fmt.Errorf("error creating directories for path %s: %w", targetPath, err)
			}
		} else {
			return fmt.Errorf("error accessing path %s: %w", targetPath, err)
		}
	}
	return nil
}
