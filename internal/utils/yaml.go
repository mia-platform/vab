package utils

import (
	"bytes"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"gopkg.in/yaml.v3"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

const (
	defaultConfigFileName  = "config.yaml"
	defaultFilePermissions = 0644
	kustomizationFileName  = "kustomization.yaml"
	yamlDefaultIndent      = 2
)

var errKustomizationTarget = errors.New("The target file must be a kustomization.yaml")

func writeYamlFile(file interface{}, dstPath string) error {
	var b bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&b)
	yamlEncoder.SetIndent(yamlDefaultIndent)

	if err := yamlEncoder.Encode(&file); err != nil {
		return err
	}

	if writeErr := os.WriteFile(dstPath, b.Bytes(), defaultFilePermissions); writeErr != nil {
		return writeErr
	}

	return nil
}

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

	return writeYamlFile(config, dstPath)
}

func WriteKustomization(kustomization kustomizeTypes.Kustomization, dirOrFilePath string) error {
	dirOrFile, err := os.Stat(dirOrFilePath)

	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	var dstPath string
	var dstPathCond bool
	switch dstPathCond {
	case dstPathCond == (err == nil && dirOrFile.IsDir()):
		dstPath = path.Join(dirOrFilePath, kustomizationFileName)
	case dstPathCond == (filepath.Base(dirOrFilePath) != kustomizationFileName):
		return errKustomizationTarget
	default:
		dstPath = dirOrFilePath
	}

	return writeYamlFile(kustomization, dstPath)
}
