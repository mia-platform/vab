package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	vabUtils "github.com/mia-platform/vab/internal/utils"
	vabConfig "github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/spf13/cobra"
)

var flags = &FlagPole{}

func NewInitCommand() *cobra.Command {
	initCmd := &cobra.Command{

		Use:   "init",
		Short: "Initialize the vab project",
		Long: `Creates the project folder with a preliminary directory structure, together with the skeleton of the configuration file.
The project directory will contain the clusters directory (including the all-clusters folder with a minimal kustomize
configuration), and the configuration file.`,

		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Initializing...")

			currentPath, err := os.Getwd()
			if err != nil {
				return err
			}

			configPath, err := vabUtils.GetProjectRelativePath(currentPath, flags.Name)
			if err != nil {
				return err
			}

			if err := vabUtils.WriteConfig(vabConfig.EmptyConfig(filepath.Base(configPath)), configPath); err != nil {
				return err
			}

			if err := vabUtils.CreateClusterOverride(configPath, "all-clusters"); err != nil {
				return err
			}

			return nil
		},
	}
	initCmd.Flags().StringVarP(&flags.Name, "name", "n", "", "project name, defaults to current directory name")
	return initCmd
}
