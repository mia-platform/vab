package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type flagpole struct {
	Name string
}

var flags = &FlagPole{}

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
		viper.SetConfigFile("config.yaml")
		viper.SetConfigType("yaml")
		viper.SetDefault("apiVersion", "vab.mia-platform.eu/v1alpha1")
		viper.WriteConfig()
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
