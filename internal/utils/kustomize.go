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
	"io"
	"os"

	"sigs.k8s.io/kustomize/kustomize/v4/commands/build"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// RunKustomizeBuild runs the kustomize build command in the target path
func RunKustomizeBuild(targetPath string, writer io.Writer) error {
	if writer == nil {
		writer = os.Stdout
	}

	kustomizeCmd := build.NewCmdBuild(
		filesys.MakeFsOnDisk(),
		&build.Help{},
		writer,
	)

	args := []string{targetPath}
	kustomizeCmd.SetArgs(args)

	return kustomizeCmd.Execute()
}
