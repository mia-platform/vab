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
	clusterDir := path.Join(configPath, clustersDirName, clusterName)
	if err := os.MkdirAll(clusterDir, os.ModePerm); err != nil {
		return err
	}
	if err := WriteKustomization(*emptyKustomization, clusterDir); err != nil {
		return err
	}
	return nil
}
