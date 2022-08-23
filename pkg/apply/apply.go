// Copyright 2022 Mia-Platform

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//    http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apply

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	vabBuild "github.com/mia-platform/vab/pkg/build"
	"github.com/mia-platform/vab/pkg/logger"
	"golang.org/x/exp/slices"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/kubectl/pkg/cmd/apply"
	"k8s.io/kubectl/pkg/cmd/util"
)

func Apply(logger logger.LogInterface, configPath string, outputDir string, groupName string, clusterName string, contextPath string) error {
	cleanedContextPath := path.Clean(contextPath)
	contextInfo, err := os.Stat(cleanedContextPath)
	if err != nil {
		return err
	}
	if !contextInfo.IsDir() {
		return fmt.Errorf("the target path %s is not a directory", cleanedContextPath)
	}

	targetPaths, err := utils.BuildPaths(configPath, groupName, clusterName)
	if err != nil {
		return err
	}

	for _, clusterPath := range targetPaths {
		buffer := new(bytes.Buffer)
		pathArray := strings.Split(clusterPath, "/")
		cluster := pathArray[len(pathArray)-1]

		targetPath := path.Join(cleanedContextPath, clusterPath)
		if err := vabBuild.RunKustomizeBuild(targetPath, buffer); err != nil {
			logger.V(5).Writef("Error building kustomize in %s", targetPath)
			return err
		}

		os.MkdirAll(outputDir, 0700)
		filePath := filepath.Join(outputDir, cluster)
		fmt.Println("creating resources at ", filePath)
		err = ioutil.WriteFile(filePath, buffer.Bytes(), 0644)
		if err != nil {
			return err
		}

		context, err := getContext(configPath, groupName, cluster)
		if err != nil {
			return err
		}

		err = runKubectlApply(logger, cluster, context)
		if err != nil {
			return err
		}
	}
	return nil

}

func getContext(configPath string, groupName string, clusterName string) (string, error) {
	config, err := utils.ReadConfig(configPath)
	if err != nil {
		return "", err
	}

	groupIdx := slices.IndexFunc(config.Spec.Groups, func(g v1alpha1.Group) bool { return g.Name == groupName })
	if groupIdx == -1 {
		return "", errors.New("Group " + groupName + " not found in configuration")
	}

	clusterIdx := slices.IndexFunc(config.Spec.Groups[groupIdx].Clusters, func(c v1alpha1.Cluster) bool { return c.Name == clusterName })
	if clusterIdx == -1 {
		return "", errors.New("Cluster " + clusterName + " not found in configuration")
	}

	return config.Spec.Groups[groupIdx].Clusters[clusterIdx].Context, nil

}

func runKubectlApply(logger logger.LogInterface, fileName string, context string) error {
	//default configflags
	configFlags := genericclioptions.NewConfigFlags(false)
	//the kubeconfig context used is equal to the fileName
	configFlags.Context = &context

	factory := util.NewFactory(configFlags)
	streams := genericclioptions.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
	args := []string{
		"-f",
		fileName,
	}
	cmd := apply.NewCmdApply("kubectl", factory, streams)
	cmd.SetArgs(args)

	//dry-run for testing purposes
	test, _ := regexp.Match("test*", []byte(fileName))

	if !test {
		err := cmd.Execute()
		if err != nil {
			return err
		}
	} else {
		logger.V(5).Writef("Skipping apply on ", fileName, "...")
	}

	return nil
}
