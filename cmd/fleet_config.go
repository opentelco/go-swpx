package cmd

import (
	"fmt"
	"os"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {
	fleetRootCmd.AddCommand(fleetConfigRootCmd)

	listConfigCmd.Flags().Int64("limit", 100, "limit")
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

		devClient, err := getDeviceClient(cmd)
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
				items := make([]list.Item, len(res.Configurations))

				for i, n := range res.Configurations {
					d, _ := devClient.GetByID(cmd.Context(), &devicepb.GetByIDParameters{Id: n.DeviceId})
					var title string

					if d != nil {
						if d.Hostname != "" {
							title = d.Hostname
						} else {
							title = d.ManagementIp
						}
					} else {
						title = n.Id
					}

					items[i] = configItem{
						title: title,
						desc:  n.Created.AsTime().String(),
						id:    n.Id,
						c:     n,
					}
				}

				m := configModel{
					list: list.New(items, list.NewDefaultDelegate(), 200, 0), service: client}
				m.list.Title = "Configurations"

				p := tea.NewProgram(m, tea.WithAltScreen())

				if _, err := p.Run(); err != nil {
					fmt.Println("Error running program:", err)
					os.Exit(1)
				}
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
