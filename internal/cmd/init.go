package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	vabUtils "github.com/mia-platform/vab/internal/utils"
	vabConfig "github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/spf13/cobra"
	kustomizeTypes "sigs.k8s.io/kustomize/api/types"
)

var flags = &FlagPole{}
var emptyConfig = &vabConfig.ClustersConfiguration{
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
var emptyKustomization = &kustomizeTypes.Kustomization{
	TypeMeta: kustomizeTypes.TypeMeta{
		Kind:       "Kustomization",
		APIVersion: "kustomize.config.k8s.io/v1beta1",
	},
}

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

		configPath, err := vabUtils.GetProjectRelativePath(".", flags.Name)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		emptyConfig.Name = filepath.Base(configPath)
		if err := vabUtils.WriteConfig(*emptyConfig, configPath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := os.MkdirAll(path.Join(configPath, "clusters", "all-clusters"), os.ModePerm); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := vabUtils.WriteKustomization(*emptyKustomization, path.Join(configPath, "clusters", "all-clusters")); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	InitCmd.Flags().StringVarP(&flags.Name, "name", "n", "", "project name, defaults to current directory name")
}
