package utils

import (
	"io"
	"os"

	"sigs.k8s.io/kustomize/api/filesys"
	"sigs.k8s.io/kustomize/kustomize/v4/commands/build"
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

	if err := kustomizeCmd.Execute(); err != nil {
		return err
	}

	return nil
}
