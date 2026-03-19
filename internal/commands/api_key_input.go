package commands

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

type APIKeyInput struct {
	input textinput.Model
	value string
}

func NewAPIKeyInput() *APIKeyInput {
	ti := textinput.New()
	ti.Placeholder = "Paste your API key here"
	ti.Focus()
	ti.CharLimit = 100

	return &APIKeyInput{
		input: ti,
	}
}

func (m *APIKeyInput) Value() string {
	return m.value
}

func (m *APIKeyInput) Run() error {
	_, err := tea.NewProgram(m).Run()
	return err
}

func (m *APIKeyInput) Init() tea.Cmd {
	return nil
}

func (m *APIKeyInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "enter":
			m.value = m.input.Value()
			if m.value != "" {
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m *APIKeyInput) View() string {
	return fmt.Sprintf(`
Welcome to menlo!

To get your API key, please visit:
https://platform.menlo.ai/account/api-keys

%s

Press Enter to save, or q to cancel
`, m.input.View())
}