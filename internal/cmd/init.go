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
var initConfig = &vabConfig.ClustersConfiguration{

	TypeMeta: vabConfig.TypeMeta{
		Kind:       "ClustersConfiguration",
		APIVersion: "vab.mia-platform.eu/v1alpha1",
	},

	Spec: vabConfig.ConfigSpec{
		Modules: make(map[string]vabConfig.Module),
		AddOns:  make(map[string]vabConfig.AddOn),
		Groups:  make([]vabConfig.Group, 0),
	},
}
var kustomization = `apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization`

var InitCmd = &cobra.Command{

	Use:   "init",
	Short: "Initialize the vab project",
	Long: `Creates the project folder with a preliminary directory structure,
together with the skeleton of the configuration file.
The project directory will contain the clusters directory
(including the all-clusters folder with a minimal
kustomize configuration), and the configuration file.`,

	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("Initializing...")

		if flags.Name != "" {
			os.Mkdir(flags.Name, os.ModePerm)
			os.Chdir(flags.Name)
		}

		configPath, _ := os.Getwd()
		initConfig.Name = filepath.Base(configPath)
		vabUtils.WriteConfig(*initConfig, configPath)

		os.Mkdir("clusters", os.ModePerm)
		os.Mkdir("clusters/all-clusters", os.ModePerm)

		if writeErr := os.WriteFile("clusters/all-clusters/kustomization.yaml", []byte(kustomization), 0644); writeErr != nil {
			fmt.Println(writeErr)
			os.Exit(1)
		}
	},
}

func init() {
	InitCmd.Flags().StringVarP(&flags.Name, "name", "n", "", "project name, defaults to current directory name")
}
