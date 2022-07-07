package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"gopkg.in/yaml.v3"
)

const (
	yamlDefaultIndent      = 2
	defaultConfigFileName  = "config.yaml"
	defaultFilePermissions = 0644
)

func WriteConfig(config v1alpha1.ClustersConfiguration, dirOrFilePath string) error {
	dirOrFile, err := os.Stat(dirOrFilePath)

	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	var dstPath string
	if err == nil && dirOrFile.IsDir() {
		dstPath = path.Join(dirOrFilePath, defaultConfigFileName)
	} else {
		dstPath = dirOrFilePath
	}

	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(yamlDefaultIndent)
	yamlEncoder.Encode(&config)

	if writeErr := os.WriteFile(dstPath, b.Bytes(), defaultFilePermissions); writeErr != nil {
		fmt.Println(writeErr)
		return writeErr
	}

	return nil
}
