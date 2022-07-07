package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	vabConfig "github.com/mia-platform/vab/pkg/apis/vab.mia-platform.eu/v1alpha1"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type flagpole struct {
	Name string
}

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

		var b bytes.Buffer
		yamlEncoder := yaml.NewEncoder(&b)
		yamlEncoder.SetIndent(2)

		if flags.Name != "" {
			prjName := flags.Name
			initConfig.Name = prjName
			os.Mkdir(prjName, os.ModePerm)
			os.Chdir(prjName)
		}

		yamlEncoder.Encode(&initConfig)

		if writeErr := ioutil.WriteFile("config.yaml", b.Bytes(), 0666); writeErr != nil {
			fmt.Println(writeErr)
			os.Exit(1)
		}
	},
}

func init() {

	dir, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	InitCmd.Flags().StringVarP(&flags.Name, "name", "n", filepath.Base(dir), "project name, defaults to current directory name")
}
