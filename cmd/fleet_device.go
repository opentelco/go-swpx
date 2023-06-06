package cmd

import (
	"fmt"
	"os"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/fleetpb"
	"github.com/spf13/cobra"
)

func init() {
	// root

	// create
	createDeviceCmd.Flags().String("host", "", "the hostname of the device")
	createDeviceCmd.Flags().String("domain", "", "domain hostname of the device")
	createDeviceCmd.Flags().String("management-ip", "", "management ip of the device")
	createDeviceCmd.Flags().String("serial", "", "serial/mac of the device")
	createDeviceCmd.Flags().String("model", "", "model of the device")
	createDeviceCmd.Flags().String("version", "", "what version the device is running")
	createDeviceCmd.Flags().String("network-region", "", "network region of the device,used by the poller to determine which network to use")
	createDeviceCmd.Flags().String("poller-provider-plugin", "default_provider", "default provider for the device")
	createDeviceCmd.Flags().String("poller-resource-plugin", "", "resource plugin to use when polling the device (VRP, CTC, etc)")
	fleetDeviceCmd.AddCommand(createDeviceCmd)

	discoverAndCreateDeviceCmd.Flags().String("host", "", "the hostname of the device")
	discoverAndCreateDeviceCmd.Flags().String("domain", "", "domain hostname of the device")
	discoverAndCreateDeviceCmd.Flags().String("management-ip", "", "management ip of the device")
	discoverAndCreateDeviceCmd.Flags().String("serial", "", "serial/mac of the device")
	discoverAndCreateDeviceCmd.Flags().String("model", "", "model of the device")
	discoverAndCreateDeviceCmd.Flags().String("version", "", "what version the device is running")
	discoverAndCreateDeviceCmd.Flags().String("network-region", "", "network region of the device,used by the poller to determine which network to use")
	discoverAndCreateDeviceCmd.Flags().String("poller-provider-plugin", "default_provider", "default provider for the device")
	discoverAndCreateDeviceCmd.Flags().String("poller-resource-plugin", "", "resource plugin to use when polling the device (VRP, CTC, etc)")
	fleetDeviceCmd.AddCommand(discoverAndCreateDeviceCmd)

	updateDeviceCmd.Flags().String("domain", "", "domain hostname of the device")
	updateDeviceCmd.Flags().String("management-ip", "", "management ip of the device")
	updateDeviceCmd.Flags().String("serial", "", "serial/mac of the device")
	updateDeviceCmd.Flags().String("model", "", "model of the device")
	updateDeviceCmd.Flags().String("version", "", "what version the device is running")
	updateDeviceCmd.Flags().String("poller-provider-plugin", "default_provider", "default provider for the device")
	updateDeviceCmd.Flags().String("poller-resource-plugin", "", "resource plugin to use when polling the device (VRP, CTC, etc)")

	fleetDeviceCmd.AddCommand(updateDeviceCmd)

	fleetDeviceCmd.AddCommand(listDeviceChangesCmd)

	fleetDeviceCmd.AddCommand(deleteDeviceCmd)

	fleetDeviceCmd.AddCommand(collectDeviceConfig)

	fleetDeviceCmd.AddCommand(collectDeviceCmd)

	fleetDeviceCmd.AddCommand(listDeviceEventsCmd)

	// list
	listDeviceCmd.Flags().String("search", "", "serach for devices by hostname or serial number")
	fleetDeviceCmd.AddCommand(listDeviceCmd)

	fleetRootCmd.AddCommand(fleetDeviceCmd)
}

var fleetDeviceCmd = &cobra.Command{
	Use:     "device",
	Aliases: []string{"dev"},
	Short:   "fleet device commands for swpx",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var createDeviceCmd = &cobra.Command{
	Use:   "create",
	Short: "create a device without any discovery",
	Run: func(cmd *cobra.Command, args []string) {

		host, _ := cmd.Flags().GetString("host")
		domain, _ := cmd.Flags().GetString("domain")
		managementIp, _ := cmd.Flags().GetString("management-ip")
		serial, _ := cmd.Flags().GetString("serial")
		model, _ := cmd.Flags().GetString("model")
		version, _ := cmd.Flags().GetString("version")
		networkRegion, _ := cmd.Flags().GetString("network-region")
		pollerProviderPlugin, _ := cmd.Flags().GetString("poller-provider-plugin")
		pollerResourcePlugin, _ := cmd.Flags().GetString("poller-resource-plugin")

		fleetDeviceClient, err := getDeviceClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		dev, err := fleetDeviceClient.Create(cmd.Context(), &devicepb.CreateParameters{
			Hostname:             &host,
			Domain:               &domain,
			ManagementIp:         &managementIp,
			SerialNumber:         &serial,
			NetworkRegion:        &networkRegion,
			Model:                &model,
			Version:              &version,
			PollerProviderPlugin: &pollerProviderPlugin,
			PollerResourcePlugin: &pollerResourcePlugin,
		})
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
		fmt.Println(prettyPrintJSON(dev))

	},
}

var discoverAndCreateDeviceCmd = &cobra.Command{
	Use:   "discover",
	Short: "discover rand create a device",
	Run: func(cmd *cobra.Command, args []string) {

		host, _ := cmd.Flags().GetString("host")
		domain, _ := cmd.Flags().GetString("domain")
		managementIp, _ := cmd.Flags().GetString("management-ip")
		serial, _ := cmd.Flags().GetString("serial")
		model, _ := cmd.Flags().GetString("model")
		version, _ := cmd.Flags().GetString("version")
		networkRegion, _ := cmd.Flags().GetString("network-region")
		pollerProviderPlugin, _ := cmd.Flags().GetString("poller-provider-plugin")
		pollerResourcePlugin, _ := cmd.Flags().GetString("poller-resource-plugin")

		fleetClient, err := getFleetClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		dev, err := fleetClient.DiscoverDevice(cmd.Context(), &devicepb.CreateParameters{
			Hostname:             &host,
			Domain:               &domain,
			ManagementIp:         &managementIp,
			SerialNumber:         &serial,
			NetworkRegion:        &networkRegion,
			Model:                &model,
			Version:              &version,
			PollerProviderPlugin: &pollerProviderPlugin,
			PollerResourcePlugin: &pollerResourcePlugin,
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
		pollerProviderPlugin, _ := cmd.Flags().GetString("poller-provider-plugin")
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
		if pollerProviderPlugin != "" {
			params.PollerProviderPlugin = &pollerProviderPlugin
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

var listDeviceChangesCmd = &cobra.Command{
	Use:   "changes [id]",
	Short: "list changes for a device",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		fleetDeviceClient, err := getDeviceClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		res, err := fleetDeviceClient.ListChanges(cmd.Context(), &devicepb.ListChangesParameters{
			DeviceId: args[0],
		})

		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
		if len(res.Changes) == 0 {
			fmt.Println("no changes found for device: ", args[0])
			return
		}
		for _, change := range res.Changes {
			fmt.Println(prettyPrintJSON(change))
		}

	},
}

var listDeviceEventsCmd = &cobra.Command{
	Use:   "events [id]",
	Args:  cobra.MinimumNArgs(1),
	Short: "list events for a device",
	Run: func(cmd *cobra.Command, args []string) {

		fleetDeviceClient, err := getDeviceClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		res, err := fleetDeviceClient.ListEvents(cmd.Context(), &devicepb.ListEventsParameters{
			DeviceId: args[0],
		})

		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
		for _, evnt := range res.Events {
			fmt.Println(prettyPrintJSON(evnt))
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

var collectDeviceConfig = &cobra.Command{
	Use:   "collect-config [deviceId]",
	Short: "collect config for a device",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		deviceId := args[0]

		fleetClient, err := getFleetClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		dev, err := fleetClient.CollectConfig(cmd.Context(), &fleetpb.CollectConfigParameters{
			DeviceId: deviceId,
		})
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
		fmt.Println(prettyPrintJSON(dev))

	},
}
var collectDeviceCmd = &cobra.Command{
	Use:   "collect [deviceId]",
	Short: "collect info for a device",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		deviceId := args[0]

		fleetClient, err := getFleetClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		dev, err := fleetClient.CollectDevice(cmd.Context(), &fleetpb.CollectDeviceParameters{
			DeviceId: deviceId,
		})
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
		fmt.Println(prettyPrintJSON(dev))

	},
}