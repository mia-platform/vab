package utils

import (
	"errors"
	"io/fs"
	"os"
	"path"
)

func GetProjectRelativePath(currentPath string, name string) (string, error) {
	if name == "" {
		info, err := os.Stat(currentPath)
		if info != nil && info.Mode().Perm()&(1<<(uint(7))) == 0 {
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
