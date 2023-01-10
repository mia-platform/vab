// Copyright Mia srl
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package build

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mia-platform/vab/internal/utils"
	"github.com/mia-platform/vab/pkg/logger"
	"sigs.k8s.io/kustomize/kustomize/v4/commands/build"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// Build kustomization configurations for a given groupName clusters or a single clusterName in groupName
// based on the passed configPath
func Build(logger logger.LogInterface, configPath string, groupName string, clusterName string, contextPath string, writer io.Writer) error {
	cleanedContextPath := filepath.Clean(contextPath)
	contextInfo, err := os.Stat(cleanedContextPath)
	if err != nil {
		return err
	}
	if !contextInfo.IsDir() {
		return fmt.Errorf("the target path %s is not a directory", cleanedContextPath)
	}

	logger.V(10).Writef("Read configuration from %s", configPath)
	targetPaths, err := utils.BuildPaths(configPath, groupName, clusterName)
	if err != nil {
		return err
	}

	logger.V(10).Writef("Found the following paths %s", targetPaths)
	for _, clusterPath := range targetPaths {
		fmt.Fprintf(writer, "### BUILD RESULTS FOR: %s ###\n", clusterPath)
		targetPath := filepath.Join(cleanedContextPath, clusterPath)
		if err := RunKustomizeBuild(targetPath, writer); err != nil {
			logger.V(5).Writef("Error building kustomize in %s", targetPath)
			return err
		}
		fmt.Fprint(writer, "---\n")
	}

	logger.V(10).Writef("Built all configurations in %s for group \"%s\", cluster\"%s\"", targetPaths, groupName, clusterName)
	return nil
}

// runKustomizeBuild runs the kustomize build command in targetPath
func RunKustomizeBuild(targetPath string, writer io.Writer) error {
	kustomizeCmd := build.NewCmdBuild(
		filesys.MakeFsOnDisk(),
		&build.Help{},
		writer,
	)

	args := []string{targetPath}
	kustomizeCmd.SetArgs(args)

	return kustomizeCmd.Execute()
}
