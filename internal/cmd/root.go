package cmd

import (
	"github.com/spf13/cobra"
)

type FlagPole struct {
	Name   string
	Config string
}

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "vab",
		Version: "vab-v0.0",
		Short:   "A tool for installing the Mia-Platform distro on your clusters",
	}

	rootCmd.AddCommand(NewInitCommand())
	rootCmd.AddCommand(NewValidateCommand())
	rootCmd.AddCommand(NewSyncCommand())
	rootCmd.AddCommand(NewBuildCommand())
	rootCmd.AddCommand(NewApplyCommand())

	return rootCmd
}
