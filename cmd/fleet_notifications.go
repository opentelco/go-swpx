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

	listNotificationsCmd.Flags().BoolP("include-read", "r", false, "include read notifications")

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

		includeread, _ := cmd.Flags().GetBool("include-read")

		notificationClient, err := getNotificationClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		params := &notificationpb.ListRequest{}

		if includeread {
			params.Filter = []notificationpb.ListRequest_Filter{notificationpb.ListRequest_INCLUDE_READ}
		}

		res, err := notificationClient.List(cmd.Context(), params)

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
				read:  n.Read,
			}
		}

		m := notificationModel{
			list: list.New(items, list.NewDefaultDelegate(), 200, 0), notification: notificationClient}
		m.list.Title = "Notifications"

		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

	},
}
