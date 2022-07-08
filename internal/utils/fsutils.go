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
		Kind:       "Kustomization",
		APIVersion: "kustomize.config.k8s.io/v1beta1",
	},
}

func GetProjectRelativePath(currentPath string, name string) (string, error) {
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

func CreateClusterOverride(configPath string, clusterName string) error {
	clustersDir := path.Join(configPath, clustersDirName)
	if _, statErr := os.Stat(clustersDir); statErr != nil {
		if errors.Is(statErr, fs.ErrNotExist) {
			if mkdirErr := os.Mkdir(path.Join(configPath, clustersDirName), os.ModePerm); mkdirErr != nil {
				return mkdirErr
			}
		} else {
			return statErr
		}
	}
	if err := os.Mkdir(path.Join(clustersDir, clusterName), os.ModePerm); err != nil {
		return err
	}
	if err := WriteKustomization(*emptyKustomization, path.Join(clustersDir, clusterName)); err != nil {
		return err
	}
	return nil
}
