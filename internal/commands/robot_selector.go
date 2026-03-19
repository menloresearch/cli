package commands

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/menloresearch/menlo-cli/internal/clients/platform"
)

type robotItem struct {
	id   string
	name string
}

func (i robotItem) Title() string       { return i.name }
func (i robotItem) Description() string { return i.id }
func (i robotItem) FilterValue() string { return i.name }

type RobotSelector struct {
	list     list.Model
	selected string
	quitting bool
}

func NewRobotSelector(robots []platform.RobotResponse) *RobotSelector {
	items := make([]list.Item, len(robots))
	for i, r := range robots {
		items[i] = robotItem{id: r.ID, name: r.Name}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select a default robot"

	return &RobotSelector{
		list: l,
	}
}

func (m *RobotSelector) Run() error {
	_, err := tea.NewProgram(m).Run()
	return err
}

func (m *RobotSelector) Selected() string {
	return m.selected
}

func (m *RobotSelector) Init() tea.Cmd {
	return nil
}

func (m *RobotSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if i, ok := m.list.SelectedItem().(robotItem); ok {
				m.selected = i.id
				m.quitting = true
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height-1)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m *RobotSelector) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("\n%s\n\nPress Enter to select, q to quit\n", m.list.View())
}
