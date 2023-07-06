package cmd

import (
	"fmt"
	"os"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/notificationpb"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {

	fleetNotificationsCmd.AddCommand(listNotificationsCmd)

	fleetRootCmd.AddCommand(fleetNotificationsCmd)

}

var fleetNotificationsCmd = &cobra.Command{
	Use:     "notifications",
	Aliases: []string{"not"},
	Short:   "feet notifications",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var listNotificationsCmd = &cobra.Command{
	Use:   "list ",
	Short: "list notifications",
	Run: func(cmd *cobra.Command, args []string) {

		fleetDeviceClient, err := getNotificationClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		res, err := fleetDeviceClient.List(cmd.Context(), &notificationpb.ListRequest{})

		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		items := make([]list.Item, len(res.Notifications))
		for i, n := range res.Notifications {
			items[i] = notificationItem{
				title: n.Title,
				desc:  n.Message,
				id:    n.Id,
			}
		}

		m := notificationModel{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
		m.list.Title = "Notifications"

		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

	},
}
