package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type TypeMeta struct {
	APIVersion string `yaml:"apiVersion,omitempty"`
	Kind       string `yaml:"kind,omitempty"`
}

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the vab project",
	Long: `Creates the project folder with a preliminary directory structure,
					together with the skeleton of the configuration file.
					The project directory will include the clusters directory
					(including the all-clusters folder with a minimal
					kustomize configuration), and the configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing...")
		viper.SetConfigFile("./config.yaml")
		viper.SetConfigType("yaml")
		if err := viper.ReadInConfig(); err != nil {
			panic(fmt.Errorf("fatal error config file: %w", err))
		}
		c := &TypeMeta{}
		viper.Unmarshal(c)
	},
}
