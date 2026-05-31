package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/giakiet05/uit-hub/apps/agent/internal/storage"
)

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).MarginBottom(1)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true)
	normalStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type pickerModel struct {
	sessions []storage.SessionRecord
	cursor   int
	selected string
	canceled bool
}

func (m pickerModel) Init() tea.Cmd {
	return nil
}

func (m pickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.canceled = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.sessions)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.sessions) > 0 {
				m.selected = m.sessions[m.cursor].ID
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m pickerModel) View() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Select a Session to Resume:"))
	b.WriteString("\n")

	if len(m.sessions) == 0 {
		b.WriteString(normalStyle.Render("No previous sessions found."))
		b.WriteString("\n(Press 'q' or 'esc' to exit)")
		return b.String()
	}

	for i, s := range m.sessions {
		cursor := " "
		style := normalStyle
		if m.cursor == i {
			cursor = ">"
			style = selectedStyle
		}

		timeStr := s.UpdatedAt.Format("2006-01-02 15:04")
		label := fmt.Sprintf("%s %s (Started: %s)", cursor, s.ID, timeStr)
		b.WriteString(style.Render(label))
		b.WriteString("\n")
	}

	b.WriteString("\n(Use up/down arrows to navigate, enter to select, esc to cancel)\n")
	return b.String()
}

// selectSessionTUI runs the picker program and returns the selected session ID.
// If the user cancels, it returns an empty string.
func selectSessionTUI(sessions []storage.SessionRecord) string {
	p := tea.NewProgram(pickerModel{sessions: sessions})
	m, err := p.Run()
	if err != nil {
		return ""
	}
	if model, ok := m.(pickerModel); ok {
		if model.canceled {
			return ""
		}
		return model.selected
	}
	return ""
}
