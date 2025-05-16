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

type baseModel struct {
	cursor int

	editMode       bool
	insertMode     bool
	selectedBucket string
	bucket         map[string]lipgloss.Style
	bucketTasks    map[string][]string
	textInput      textinput.Model
}

func initialModel() baseModel {
	ti := textinput.New()
	ti.TextStyle = inputStyle
	ti.PromptStyle = inputStyle
	ti.Cursor.Style = inputStyle

	return baseModel{
		bucket: map[string]lipgloss.Style{
			"work":     workStyle,
			"academic": academicStyle,
			"personal": personalStyle,
		},
		selectedBucket: "work",
		bucketTasks: map[string][]string{
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

func (m baseModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m baseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			if !m.editMode && !m.insertMode {
				m.editMode = true
				m.textInput.SetValue(m.bucketTasks[m.selectedBucket][m.cursor])
				m.textInput.Focus()
			}

		case "n":
			if !m.editMode && !m.insertMode {
				m.insertMode = true
				m.textInput.Placeholder = "Do something..."
				m.textInput.Width = 14
				m.textInput.Focus()
			}

		case "w":
			if !m.editMode && !m.insertMode {
				m.selectedBucket = "work"
				if m.cursor > len(m.bucketTasks[m.selectedBucket])-1 {
					m.cursor = len(m.bucketTasks[m.selectedBucket]) - 1
				}
			}

		case "e":
			if !m.editMode && !m.insertMode {
				m.selectedBucket = "academic"
				if m.cursor > len(m.bucketTasks[m.selectedBucket])-1 {
					m.cursor = len(m.bucketTasks[m.selectedBucket]) - 1
				}
			}

		case "r":
			if !m.editMode && !m.insertMode {
				m.selectedBucket = "personal"
				if m.cursor > len(m.bucketTasks[m.selectedBucket])-1 {
					m.cursor = len(m.bucketTasks[m.selectedBucket]) - 1
				}
			}

		case "q":
			if !m.editMode && !m.insertMode {
				return m, tea.Quit
			}

		case "up", "k":
			if !m.editMode && !m.insertMode {
				if m.cursor > 0 {
					m.cursor--
				}
			}

		case "down", "j":
			if !m.editMode && !m.insertMode {
				if m.cursor < len(m.bucketTasks[m.selectedBucket])-1 {
					m.cursor++
				}
			}

		case "enter":
			if m.editMode {
				m.editMode = false
				m.bucketTasks[m.selectedBucket][m.cursor] = strings.TrimSpace(m.textInput.Value())
				m.textInput.Blur()
			} else if m.insertMode {
				m.insertMode = false
				m.bucketTasks[m.selectedBucket] = append(m.bucketTasks[m.selectedBucket], strings.TrimSpace(m.textInput.Value()))
				m.textInput.Blur()
			}

		case "ctrl+c":
			return m, tea.Quit

		}
	}

	return m, cmd
}

func (m baseModel) View() string {
	s := "Tasks ["

	buckets := make([]string, 0, len(m.bucket))
	for k := range m.bucket {
		buckets = append(buckets, k)
	}

	slices.Sort(buckets)
	for i := range len(buckets) {
		buckets[i] = m.bucket[buckets[i]].Render(buckets[i])
	}

	s += strings.Join(buckets, "|")
	s += "]\n\n"

	if m.insertMode {
		s += m.textInput.View() + "\n\n"
	}

	for i, choice := range m.bucketTasks[m.selectedBucket] {
		if m.editMode && m.cursor == i {
			s += m.textInput.View() + "\n"
		} else {
			cursor := " "
			if m.cursor == i && !m.insertMode {
				cursor = ">"
			}

			checked := " "
			s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
		}
	}

	s += "\nPress q to quit\n"

	return m.bucket[m.selectedBucket].Render(s)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
