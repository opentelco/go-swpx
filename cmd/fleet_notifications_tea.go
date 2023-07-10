package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/glamour"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/notificationpb"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/reflow/wordwrap"
	"github.com/muesli/termenv"
)

type notificationItem struct {
	id, title, desc string
	read            bool
}

func colorFg(val, color string) string {
	return termenv.String(val).Foreground(term.Color(color)).String()
}

func colorBg(val, color string) string {
	return termenv.String(val).Background(term.Color(color)).String()
}

func (i notificationItem) Title() string {
	if !i.read {
		i.title = "🟢 " + i.title
		return colorFg(i.title, "32")
	} else {
		i.title = "⚪ " + i.title
	}

	return i.title
}

func (i notificationItem) Description() string {
	return i.desc
}
func (i notificationItem) FilterValue() string { return i.title }

type notificationModel struct {
	list list.Model
	// choice is the item that was selected by the user
	choice *notificationItem

	notification notificationpb.NotificationServiceClient
}

func (m notificationModel) Init() tea.Cmd {
	return nil
}

func (m notificationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case *tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case tea.KeyLeft.String(), "b":
			m.choice = nil

		case "ctrl+c":
			return m, tea.Quit

		case "enter", tea.KeyRight.String():
			i, ok := m.list.SelectedItem().(notificationItem)
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

func (m notificationModel) View() string {
	if m.choice != nil {
		if !m.choice.read {
			m.notification.MarkAsRead(context.Background(), &notificationpb.MarkAsReadRequest{Ids: []string{m.choice.id}})
		}

		in := fmt.Sprintf(`# %s (ID: %s)

## Description

%s
`, m.choice.Title(), m.choice.id, m.choice.desc)

		out, err := glamour.Render(in, "dark")
		if err != nil {
			panic(err)
		}
		return wordwrap.String(out, 100)
	}

	return docStyle.Render(m.list.View())
}
