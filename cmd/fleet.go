package cmd

import (
	"fmt"
	"os"
	"time"

	"git.liero.se/opentelco/go-swpx/proto/go/fleetpb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	Root.AddCommand(fleetCmd)
	fleetCmd.AddCommand(createDeviceCmd)
	fleetCmd.PersistentFlags().String("fleet-api", "127.0.0.1:1338", "the fleet api addr")

	createDeviceCmd.Flags().String("host", "", "the hostname of the device")
	if err := createDeviceCmd.MarkFlagRequired("host"); err != nil {
		panic(err)
	}

}

var fleetCmd = &cobra.Command{
	Use:   "fleet",
	Short: "fleet commands for swpx",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var createDeviceCmd = &cobra.Command{
	Use:   "create",
	Short: "create a device",
	Run: func(cmd *cobra.Command, args []string) {
		addr, _ := cmd.Parent().PersistentFlags().GetString("fleet-api")
		host, err := cmd.Flags().GetString("host")
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		conn, err := grpc.Dial(addr, grpc.WithTimeout(5*time.Second), grpc.WithBlock(), grpc.WithInsecure())
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
		fleetClient := fleetpb.NewFleetClient(conn)

		dev, err := fleetClient.CreateDevice(cmd.Context(), &fleetpb.CreateDeviceParameters{
			Hostname:             host,
			Domain:               "",
			ManagementIp:         "",
			SerialNumber:         "",
			Model:                "",
			Version:              "",
			PollerProvider:       "",
			PollerResourcePlugin: "",
		})
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
		fmt.Println(dev)

	},
}
