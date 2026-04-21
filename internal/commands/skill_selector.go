package commands

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type skillTargetItem struct {
	target string
	label  string
}

func (i skillTargetItem) Title() string       { return i.label }
func (i skillTargetItem) Description() string { return i.target }
func (i skillTargetItem) FilterValue() string { return i.label }

type SkillTargetSelector struct {
	list     list.Model
	selected string
	quitting bool
}

func NewSkillTargetSelector(targets []skillTargetItem) *SkillTargetSelector {
	items := make([]list.Item, len(targets))
	for i, t := range targets {
		items[i] = t
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select skill install target"

	return &SkillTargetSelector{list: l}
}

func (m *SkillTargetSelector) Run() error {
	_, err := tea.NewProgram(m).Run()
	return err
}

func (m *SkillTargetSelector) Selected() string {
	return m.selected
}

func (m *SkillTargetSelector) Init() tea.Cmd {
	return nil
}

func (m *SkillTargetSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if i, ok := m.list.SelectedItem().(skillTargetItem); ok {
				m.selected = i.target
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

func (m *SkillTargetSelector) View() string {
	if m.quitting {
		return ""
	}
	return fmt.Sprintf("\n%s\n\nPress Enter to select, q to quit\n", m.list.View())
}
