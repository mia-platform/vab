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
