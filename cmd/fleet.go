package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	// root
	fleetRootCmd.PersistentFlags().String("fleet-addr", "127.0.0.1:1338", "the fleet api addr")
	Root.AddCommand(fleetRootCmd)

}

var fleetRootCmd = &cobra.Command{
	Use:   "fleet",
	Short: "fleet commands for swpx",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}
