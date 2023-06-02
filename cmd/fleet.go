package cmd

import (
	"fmt"
	"os"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"github.com/spf13/cobra"
)

func init() {
	// root
	fleetCmd.PersistentFlags().String("fleet-addr", "127.0.0.1:1338", "the fleet api addr")
	Root.AddCommand(fleetCmd)

	// create
	createDeviceCmd.Flags().String("host", "", "the hostname of the device")
	if err := createDeviceCmd.MarkFlagRequired("host"); err != nil {
		panic(err)
	}
	createDeviceCmd.Flags().String("domain", "", "domain hostname of the device")
	createDeviceCmd.Flags().String("management-ip", "", "management ip of the device")
	createDeviceCmd.Flags().String("serial", "", "serial/mac of the device")
	createDeviceCmd.Flags().String("model", "", "model of the device")
	createDeviceCmd.Flags().String("version", "", "what version the device is running")
	createDeviceCmd.Flags().String("poller-provider", "default_provider", "default provider for the device")
	createDeviceCmd.Flags().String("poller-resource-plugin", "", "resource plugin to use when polling the device (VRP, CTC, etc)")
	fleetCmd.AddCommand(createDeviceCmd)

	updateDeviceCmd.Flags().String("domain", "", "domain hostname of the device")
	updateDeviceCmd.Flags().String("management-ip", "", "management ip of the device")
	updateDeviceCmd.Flags().String("serial", "", "serial/mac of the device")
	updateDeviceCmd.Flags().String("model", "", "model of the device")
	updateDeviceCmd.Flags().String("version", "", "what version the device is running")
	updateDeviceCmd.Flags().String("poller-provider", "default_provider", "default provider for the device")
	updateDeviceCmd.Flags().String("poller-resource-plugin", "", "resource plugin to use when polling the device (VRP, CTC, etc)")
	fleetCmd.AddCommand(updateDeviceCmd)

	fleetCmd.AddCommand(deleteDeviceCmd)

	// list
	listDeviceCmd.Flags().String("search", "", "serach for devices by hostname or serial number")
	fleetCmd.AddCommand(listDeviceCmd)
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

		host, _ := cmd.Flags().GetString("host")
		domain, _ := cmd.Flags().GetString("domain")
		managementIp, _ := cmd.Flags().GetString("management-ip")
		serial, _ := cmd.Flags().GetString("serial")
		model, _ := cmd.Flags().GetString("model")
		version, _ := cmd.Flags().GetString("version")
		pollerProvider, _ := cmd.Flags().GetString("poller-provider")
		pollerResourcePlugin, _ := cmd.Flags().GetString("poller-resource-plugin")

		fleetDeviceClient, err := getDeviceClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		dev, err := fleetDeviceClient.Create(cmd.Context(), &devicepb.CreateParameters{
			Hostname:             host,
			Domain:               domain,
			ManagementIp:         managementIp,
			SerialNumber:         serial,
			Model:                model,
			Version:              version,
			PollerProvider:       pollerProvider,
			PollerResourcePlugin: pollerResourcePlugin,
		})
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
		fmt.Println(prettyPrintJSON(dev))

	},
}

var updateDeviceCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "update a device",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		host, _ := cmd.Flags().GetString("host")
		domain, _ := cmd.Flags().GetString("domain")
		managementIp, _ := cmd.Flags().GetString("management-ip")
		serial, _ := cmd.Flags().GetString("serial")
		model, _ := cmd.Flags().GetString("model")
		version, _ := cmd.Flags().GetString("version")
		pollerProvider, _ := cmd.Flags().GetString("poller-provider")
		pollerResourcePlugin, _ := cmd.Flags().GetString("poller-resource-plugin")

		fleetDeviceClient, err := getDeviceClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
		params := &devicepb.UpdateParameters{
			Id: args[0],
		}
		if host != "" {
			params.Hostname = &host
		}
		if domain != "" {
			params.Domain = &domain
		}
		if managementIp != "" {
			params.ManagementIp = &managementIp
		}
		if serial != "" {
			params.SerialNumber = &serial
		}
		if model != "" {
			params.Model = &model
		}
		if version != "" {
			params.Version = &version
		}
		if pollerProvider != "" {
			params.PollerProvider = &pollerProvider
		}
		if pollerResourcePlugin != "" {
			params.PollerResourcePlugin = &pollerResourcePlugin
		}

		dev, err := fleetDeviceClient.Update(cmd.Context(), params)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
		fmt.Println(prettyPrintJSON(dev))

	},
}

var listDeviceCmd = &cobra.Command{
	Use:   "list",
	Short: "list devices",
	Run: func(cmd *cobra.Command, args []string) {
		search, err := cmd.Flags().GetString("search")
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		fleetDeviceClient, err := getDeviceClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		res, err := fleetDeviceClient.List(cmd.Context(), &devicepb.ListParameters{
			Search: search,
		})

		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
		for _, dev := range res.Devices {
			fmt.Println(prettyPrintJSON(dev))
		}

	},
}

var deleteDeviceCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "delete a device by its id",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		fleetClient, err := getFleetClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		_, err = fleetClient.DeleteDevice(cmd.Context(), &devicepb.DeleteParameters{
			Id: args[0],
		})

		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

	},
}
