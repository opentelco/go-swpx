package cmd

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type notificationItem struct {
	id, title, desc string
}

func (i notificationItem) Title() string       { return i.title }
func (i notificationItem) Description() string { return i.desc }
func (i notificationItem) FilterValue() string { return i.title }

type notificationModel struct {
	list list.Model
	// choice is the item that was selected by the user
	choice *notificationItem
}

func (m notificationModel) Init() tea.Cmd {
	return nil
}

func (m notificationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
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
		return m.choice.desc
	}
	return docStyle.Render(m.list.View())
}
