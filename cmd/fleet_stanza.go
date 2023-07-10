package cmd

import (
	"fmt"
	"os"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func init() {

	listStanzaCmd.Flags().Bool("include-applied", false, "include applied")
	listStanzaCmd.Flags().Bool("include-not-applied", false, "include not applied")
	listStanzaCmd.Flags().String("device-type", "", "list stanzas for a specific device_type")
	listStanzaCmd.Flags().StringP("device-id", "d", "", "show stanzas for a specific device_id")
	// limit
	listStanzaCmd.Flags().Int("limit", 100, "limit")
	// offset
	listStanzaCmd.Flags().Int("offset", 0, "offset")

	createStanzaCmd.Flags().String("name", "", "the name of the stanza, e.g 'set hostname to foo'")
	_ = createStanzaCmd.MarkFlagRequired("name")

	createStanzaCmd.Flags().String("description", "", "the description of the stanza")
	createStanzaCmd.Flags().String("device-type", "", "device type")
	createStanzaCmd.Flags().String("content", "", "content of the stanza, e.g 'sysname foo'")
	_ = createStanzaCmd.MarkFlagRequired("content")

	createStanzaCmd.Flags().String("revert-content", "", "revert-content of the stanza, e.g 'sysname old'")

	applyStanzaCmd.Flags().StringP("device-id", "d", "", "device id")
	applyStanzaCmd.MarkFlagRequired("device-id")
	applyStanzaCmd.Flags().Bool("non-blocking", false, "blocking")

	fleetStanzaRootCmd.AddCommand(listStanzaCmd)
	fleetStanzaRootCmd.AddCommand(getStanzaCmd)
	fleetStanzaRootCmd.AddCommand(deleteStanzaCmd)
	fleetStanzaRootCmd.AddCommand(createStanzaCmd)
	fleetStanzaRootCmd.AddCommand(applyStanzaCmd)

	fleetRootCmd.AddCommand(fleetStanzaRootCmd)

}

var fleetStanzaRootCmd = &cobra.Command{
	Use:     "stanza",
	Aliases: []string{"st"},
	Short:   "feet configuration stanzas",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

var listStanzaCmd = &cobra.Command{
	Use:   "list ",
	Short: "list stanzas",
	Run: func(cmd *cobra.Command, args []string) {

		includeApplied, _ := cmd.Flags().GetBool("include-applied")
		excludeApplied, _ := cmd.Flags().GetBool("include-not-applied")
		deviceId, _ := cmd.Flags().GetString("device-id")
		limit, _ := cmd.Flags().GetInt64("limit")
		offset, _ := cmd.Flags().GetInt64("offset")

		stanzaClient, err := getStanzaClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		params := &stanzapb.ListRequest{
			Limit:  &limit,
			Offset: &offset,
		}

		if includeApplied {
			params.Filters = []stanzapb.ListRequest_Filter{stanzapb.ListRequest_FILTER_APPLIED}
		}
		if excludeApplied {
			params.Filters = []stanzapb.ListRequest_Filter{stanzapb.ListRequest_FILTER_NOT_APPLIED}
		}

		if deviceId != "" {
			params.DeviceId = &deviceId
		}

		res, err := stanzaClient.List(cmd.Context(), params)

		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		items := make([]list.Item, len(res.Stanzas))
		for i, n := range res.Stanzas {
			items[i] = stanzaItem{
				title:   n.Name,
				desc:    n.Description,
				id:      n.Id,
				applied: n.AppliedAt != nil,
			}
		}

		m := stanzaModel{
			list: list.New(items, list.NewDefaultDelegate(), 200, 0), service: stanzaClient}
		m.list.Title = "Stanzas"

		p := tea.NewProgram(m, tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

	},
}

// create
var createStanzaCmd = &cobra.Command{
	Use:     "create",
	Short:   "create stanza",
	Aliases: []string{"new"},
	Run: func(cmd *cobra.Command, args []string) {

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		deviceType, _ := cmd.Flags().GetString("device-type")
		content, _ := cmd.Flags().GetString("content")
		revertContent, _ := cmd.Flags().GetString("revert-content")

		stanzaClient, err := getStanzaClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		res, err := stanzaClient.Create(cmd.Context(), &stanzapb.CreateRequest{
			Name:          name,
			Description:   &description,
			DeviceType:    deviceType,
			Content:       content,
			RevertContent: &revertContent,
		})

		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		fmt.Println(res)

	},
}

// get by id
var getStanzaCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "get stanza",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		id := args[0]
		stanzaClient, err := getStanzaClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		res, err := stanzaClient.GetByID(cmd.Context(), &stanzapb.GetByIDRequest{
			Id: id,
		})

		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		fmt.Println(res)

	},
}

// delete by id
var deleteStanzaCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "delete stanza",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		id := args[0]
		stanzaClient, err := getStanzaClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		res, err := stanzaClient.Delete(cmd.Context(), &stanzapb.DeleteRequest{
			Id: id,
		})

		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		fmt.Println(res)

	},
}

// apply stanza
var applyStanzaCmd = &cobra.Command{
	Use:   "apply [id]",
	Short: "apply stanza",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		deviceId, _ := cmd.Flags().GetString("device-id")
		nonBlocking, _ := cmd.Flags().GetBool("non-blocking")

		id := args[0]
		stanzaClient, err := getStanzaClient(cmd)
		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		applyRes, err := stanzaClient.Apply(cmd.Context(), &stanzapb.ApplyRequest{
			Id:       id,
			DeviceId: deviceId,
			Blocking: !nonBlocking,
		})

		if err != nil {
			cmd.PrintErr(err)
			os.Exit(1)
		}

		fmt.Println(applyRes)
	},
}
