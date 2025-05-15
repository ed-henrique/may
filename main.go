package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#5533ff"))
	workStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#0000ff"))
	academicStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00ff00"))
	personalStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff9900"))
)

type model struct {
	cursor int

	editMode     bool
	selectedMode string
	mode         map[string]lipgloss.Style
	modeTasks    map[string][]string
	textInput    textinput.Model
}

func initialModel() model {
	ti := textinput.New()
	ti.TextStyle = inputStyle
	ti.PromptStyle = inputStyle
	ti.Cursor.Style = inputStyle

	return model{
		mode: map[string]lipgloss.Style{
			"work":     workStyle,
			"academic": academicStyle,
			"personal": personalStyle,
		},
		selectedMode: "work",
		modeTasks: map[string][]string{
			"work": {
				"Buy carrots",
				"Buy celery",
				"Buy kohlrabi",
			},
			"academic": {
				"Buy carrots",
			},
			"personal": {
				"Buy kohlrabi",
			},
		},
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			if !m.editMode {
				m.editMode = true
				m.textInput.SetValue(m.modeTasks[m.selectedMode][m.cursor])
				m.textInput.Focus()
			}

		case "w":
			if !m.editMode {
				m.selectedMode = "work"
				if m.cursor > len(m.modeTasks[m.selectedMode])-1 {
					m.cursor = len(m.modeTasks[m.selectedMode]) - 1
				}
			}

		case "e":
			if !m.editMode {
				m.selectedMode = "academic"
				if m.cursor > len(m.modeTasks[m.selectedMode])-1 {
					m.cursor = len(m.modeTasks[m.selectedMode]) - 1
				}
			}

		case "r":
			if !m.editMode {
				m.selectedMode = "personal"
				if m.cursor > len(m.modeTasks[m.selectedMode])-1 {
					m.cursor = len(m.modeTasks[m.selectedMode]) - 1
				}
			}

		case "q":
			if !m.editMode {
				return m, tea.Quit
			}

		case "up", "k":
			if !m.editMode {
				if m.cursor > 0 {
					m.cursor--
				}
			}

		case "down", "j":
			if !m.editMode {
				if m.cursor < len(m.modeTasks[m.selectedMode])-1 {
					m.cursor++
				}
			}

		case "enter":
			if m.editMode {
				m.editMode = false
				m.modeTasks[m.selectedMode][m.cursor] = strings.TrimSpace(m.textInput.Value())
				m.textInput.Blur()
			}

		case "ctrl+c":
			return m, tea.Quit

		}
	}

	return m, cmd
}

func (m model) View() string {
	s := "Tasks ["

	modes := make([]string, 0, len(m.mode))
	for k := range m.mode {
		modes = append(modes, k)
	}

	slices.Sort(modes)
	for i := range len(modes) {
		modes[i] = m.mode[modes[i]].Render(modes[i])
	}

	s += strings.Join(modes, "|")
	s += "]\n\n"

	for i, choice := range m.modeTasks[m.selectedMode] {
		if m.editMode && m.cursor == i {
			s += m.textInput.View() + "\n"
		} else {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			checked := " "
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}
	}

	s += "\nPress q to quit\n"

	return m.mode[m.selectedMode].Render(s)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
