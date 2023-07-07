package cmd

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var term = termenv.EnvColorProfile()
var docStyle = lipgloss.NewStyle().Margin(1, 2)

type deviceItem struct {
	id, title, desc string
	networkRegion   string
}

func (i deviceItem) Title() string {
	return fmt.Sprintf("%s (%s)", i.title, i.networkRegion)
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
		return m.choice.desc
	}
	return docStyle.Render(m.list.View())
}
