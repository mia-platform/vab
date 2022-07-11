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
	"io/fs"
	"os"
	"path"

	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

const (
	clustersDirName        = "clusters"
	userPermissionBitCheck = 7
)

var emptyKustomization = &kustomizeTypes.Kustomization{
	TypeMeta: kustomizeTypes.TypeMeta{
		Kind:       kustomizeTypes.KustomizationKind,
		APIVersion: kustomizeTypes.KustomizationVersion,
	},
}

// GetProjectPath returns the project's relative path, creating the project
// directory if name is non-empty
func GetProjectPath(currentPath string, name string) (string, error) {
	if name == "" {
		info, err := os.Stat(currentPath)
		if info != nil && info.Mode().Perm()&(1<<(uint(userPermissionBitCheck))) == 0 {
			err = fs.ErrPermission
		}
		return currentPath, err
	}
	dstPath := path.Join(currentPath, name)
	if err := os.Mkdir(dstPath, fs.ModePerm); err != nil && !errors.Is(err, fs.ErrExist) {
		return currentPath, err
	}
	return dstPath, nil
}

// CreateClusterOverride creates the directory for clusterName's override at
// the specified configPath
func CreateClusterOverride(configPath string, clusterName string) error {
	clusterDir := path.Join(configPath, clustersDirName, clusterName)
	if err := os.MkdirAll(clusterDir, os.ModePerm); err != nil {
		return err
	}
	if err := WriteKustomization(*emptyKustomization, clusterDir); err != nil {
		return err
	}
	return nil
}
