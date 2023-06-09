package cmd

import (
	"fmt"
	"os"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"github.com/spf13/cobra"
)

func init() {
	fleetRootCmd.AddCommand(fleetConfigRootCmd)

	listConfigCmd.Flags().Int64("limit", 2, "limit")
	listConfigCmd.Flags().Int64("offset", 0, "offset")

	listConfigCmd.Flags().BoolP("quite", "q", false, "only print ids")

	fleetConfigRootCmd.AddCommand(listConfigCmd)

	getConfigCmd.Flags().BoolP("diff-only", "d", false, "show diff")

	fleetConfigRootCmd.AddCommand(getConfigCmd)

}

var fleetConfigRootCmd = &cobra.Command{
	Use:     "configuration",
	Aliases: []string{"config", "conf"},
	Short:   "fleet configuration commands for swpx",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var listConfigCmd = &cobra.Command{
	Use:   "list [device-id]",
	Short: "list configurations",
	Run: func(cmd *cobra.Command, args []string) {

		var deviceId string
		if len(args) > 0 {
			deviceId = args[0]
		}
		limit, _ := cmd.Flags().GetInt64("limit")
		offset, _ := cmd.Flags().GetInt64("offset")
		quite, _ := cmd.Flags().GetBool("quite")

		client, err := getConfigClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		res, err := client.List(cmd.Context(), &configurationpb.ListParameters{
			DeviceId: deviceId,
			Limit:    &limit,
			Offset:   &offset,
		})

		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}
		for _, c := range res.Configurations {
			if quite {
				fmt.Printf("%s ", c.Id)
			} else {
				fmt.Println(prettyPrintJSON(c))
			}
		}

	},
}

var getConfigCmd = &cobra.Command{
	Use:   "get [config-id]",
	Short: "get specific configuration by id",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		configId := args[0]

		diffonly, _ := cmd.Flags().GetBool("diff-only")

		client, err := getConfigClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		res, err := client.GetByID(cmd.Context(), &configurationpb.GetByIDParameters{
			Id: configId,
		})

		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		if diffonly {
			if res.Changes == "" {
				fmt.Println("no changes")
				return
			} else {
				fmt.Println(res.Changes)
			}
			return
		} else {
			fmt.Println(prettyPrintJSON(res))
		}

	},
}
