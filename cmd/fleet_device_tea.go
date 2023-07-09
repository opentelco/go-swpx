package cmd

import (
	"fmt"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/devicepb"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var term = termenv.EnvColorProfile()
var docStyle = lipgloss.NewStyle().Margin(1, 2)

type deviceItem struct {
	id, title, desc string
	dev             *devicepb.Device
}

func (i deviceItem) Title() string {
	switch i.dev.Status {
	case devicepb.Device_DEVICE_STATUS_UNREACHABLE:
		i.title = "🔴 " + i.title
	case devicepb.Device_DEVICE_STATUS_REACHABLE:
		i.title = "🟢 " + i.title
	case devicepb.Device_DEVICE_STATUS_NEW:
		i.title = "⚪ " + i.title

	}

	return fmt.Sprintf("%s (%s)", i.title, i.dev.NetworkRegion)
}
func (i deviceItem) Description() string { return i.desc }
func (i deviceItem) FilterValue() string { return i.title }

type deviceModel struct {
	list list.Model
	// choice is the item that was selected by the user
	choice *deviceItem
}

func (m deviceModel) Init() tea.Cmd {
	return nil
}

func (m deviceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case tea.KeyLeft.String(), "b":
			m.choice = nil
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(deviceItem)
			if ok {
				m.choice = &i
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m deviceModel) View() string {
	if m.choice != nil {
		// render the device

		return renderDevice(m.choice)
	}
	return docStyle.Render(m.list.View())
}

func renderDevice(i *deviceItem) string {
	msg := fmt.Sprintf(`%s

Schedule:
%s
	`, i.Title(), i.dev.Schedules)
	return msg
}
