package cmd

import (
	"git.liero.se/opentelco/go-swpx/proto/go/fleet/configurationpb"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	tea "github.com/charmbracelet/bubbletea"
)

type configItem struct {
	id, title, desc string

	c *configurationpb.Configuration
}

func (i configItem) Title() string {
	return i.title
}

func (i configItem) Description() string {
	return i.desc
}
func (i configItem) FilterValue() string { return i.title }

type configModel struct {
	list list.Model
	// choice is the item that was selected by the user
	choice *configItem

	configVp viewport.Model

	service configurationpb.ConfigurationServiceClient
}

func (m configModel) Init() tea.Cmd {
	return nil
}

func (m configModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case *tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case tea.KeyLeft.String(), "b":
			m.choice = nil

		case "up", "down":
			if m.choice != nil {
				var cmd tea.Cmd
				m.configVp, cmd = m.configVp.Update(msg)
				return m, cmd
			} else {
				var cmd tea.Cmd
				m.list, cmd = m.list.Update(msg)
				return m, cmd
			}

		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter", tea.KeyRight.String():
			i, ok := m.list.SelectedItem().(configItem)
			if ok {
				m.choice = &i
				vp, err := newConfigVP(i.c.Configuration)
				if err != nil {
					panic(err)
				}
				m.configVp = vp

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

func (m configModel) View() string {
	if m.choice != nil {
		style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))

		out := style.Render("Device: ") + m.choice.title + "\n"
		out += style.Render("ID: ") + m.choice.c.Id + "\n"
		out += style.Render("Feched: ") + m.choice.c.Created.AsTime().String() + "\n"
		out += "\n"
		out += m.configVp.View() + m.helpView()

		if m.choice.c.Changes != "" {
			chg := style.Render("Changes") + "\n"
			chg += "... in a tab? "
			out += "\n" + chg
		}

		return out

	}

	return docStyle.Render(m.list.View())
}

const useHighPerformanceRenderer = false

func newConfigVP(content string) (viewport.Model, error) {
	const width = 300

	vp := viewport.New(width, 30)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	vp.SetContent(content)

	return vp, nil
}

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

func (m configModel) helpView() string {
	return helpStyle("\n  ↑/↓: Navigate • q: Quit\n")
}
