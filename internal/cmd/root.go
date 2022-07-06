package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

type FlagPole struct {
	Name   string
	Config string
}

var rootCmd = &cobra.Command{
	Use:   "vab",
	Short: "A tool for installing the Mia-Platform distro on your clusters",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of vab",
	Long:  `All software has versions. This is vab's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("vab-v0.0")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(InitCmd)
	rootCmd.AddCommand(ValidateCmd)
	rootCmd.AddCommand(SyncCmd)
	rootCmd.AddCommand(BuildCmd)
	rootCmd.AddCommand(ApplyCmd)
}
