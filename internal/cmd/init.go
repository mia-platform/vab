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

		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Initializing...")

			configPath, err := vabUtils.GetProjectRelativePath(".", flags.Name)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if err := vabUtils.WriteConfig(vabConfig.EmptyConfig(filepath.Base(configPath)), configPath); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if err := vabUtils.CreateClusterOverride(configPath, "all-clusters"); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	initCmd.Flags().StringVarP(&flags.Name, "name", "n", "", "project name, defaults to current directory name")
	return initCmd
}
