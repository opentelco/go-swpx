package cmd

import (
	"fmt"

	"git.liero.se/opentelco/go-swpx/proto/go/fleet/stanzapb"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type stanzaItem struct {
	id, title, desc string
	applied         bool
}

func (i stanzaItem) Title() string {
	if i.applied {
		i.title = "🟢 " + i.title
		return colorFg(i.title, "32")
	}

	return i.title
}

func (i stanzaItem) Description() string {
	return i.desc
}
func (i stanzaItem) FilterValue() string { return i.title }

type stanzaModel struct {
	list list.Model
	// choice is the item that was selected by the user
	choice *stanzaItem

	service stanzapb.StanzaServiceClient
}

func (m stanzaModel) Init() tea.Cmd {
	return nil
}

func (m stanzaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			i, ok := m.list.SelectedItem().(stanzaItem)
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

func (m stanzaModel) View() string {
	if m.choice != nil {

		var header = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			PaddingTop(1).
			PaddingLeft(1).
			PaddingBottom(1).
			Width(100)
		var body = lipgloss.NewStyle().
			PaddingTop(1).
			PaddingLeft(1).
			PaddingBottom(1).
			Width(100)

		msg := fmt.Sprintf("%s\n%s\n\n", header.Render(m.choice.title), body.Render(m.choice.desc))

		return wordwrap.String(msg, 100)
	}

	return docStyle.Render(m.list.View())
}
