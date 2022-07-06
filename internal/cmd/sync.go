package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Fetches new and updated vendor versions and updates the clusters configuration locally.",
	Long: `Fetches new and updated vendor versions and updates the clusters configuration locally
					to the latest changes of the configuration file. After the execution, the vendors folder
					will include the new and updated modules/add-ons (if not already present), and the
					directory structure inside the clusters folder will be updated according
					to the current configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Synchronizing...")
	},
}
