package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gotui/internal/commands"
	"gotui/internal/storage"
	"gotui/internal/tui"
	"log"
	"os"
	"strings"
)

const (
	// In real life situations we'd adjust the document to fit the width we've
	// detected. In the case of this example we're hardcoding the width, and
	// later using the detected width only to truncate in order to avoid jaggy
	// wrapping.
	width = 96

	columnWidth = 30
)

type command struct {
	disabled bool
	name     string
}

type appMode string

const (
	appModeInput   appMode = "input"
	appModeLoading appMode = "loading"
	appModeDefault appMode = ""
)

type model struct {
	stateDescription string
	stateStatus      tui.StatusBarState
	commands         []command // items on the to-do list
	cursor           int       // which to-do list item our cursor is pointing at
	secondListHeader string
	secondListValues []string
	dbConnection     *commands.DbConnection
	user             *storage.User
	loading          bool
	spinner          spinner.Model
	textInput        textinput.Model // <- text input component
	mode             appMode         // <- mode in which our app is input entering or not
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	ti := textinput.New()
	ti.Placeholder = "john@email.com"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput:        ti,
		spinner:          s,
		loading:          true,
		stateDescription: "Initializing...",
		stateStatus:      tui.StatusBarStateBlue,
		commands: []command{
			{name: "Set user"},
			{name: "Fetch token", disabled: true},
			{name: "Other..."},
		},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		commands.InitDatabase,
		m.spinner.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case commands.DbConnection:
		m.stateDescription = ""
		m.dbConnection = &msg
		if m.dbConnection != nil {
			if m.dbConnection.Err != nil {
				m.stateStatus = tui.StatusBarStateRed
				m.stateDescription = "Failed to connect to database: " + shortenErr(m.dbConnection.Err, 35)
			} else {
				m.stateStatus = tui.StatusBarStateGreen
				m.stateDescription = "Connected to database"
			}
		}
		m.loading = false
		return m, nil

	case commands.GetUserByEmailMsg:
		m.user = msg.User
		if msg.Err != nil {
			m.stateDescription = msg.Err.Error()
		}
		if m.user == nil {
			m.stateDescription = "User not found"
			m.stateStatus = tui.StatusBarStateYellow
			m.commands[1].disabled = true
			m.secondListHeader = "User"
			m.secondListValues = []string{"Not Found"}
		} else {
			m.stateDescription = "User set"
			m.stateStatus = tui.StatusBarStateBlue
			m.commands[1].disabled = false
			m.secondListHeader = "User"
			m.secondListValues = userFieldsToArray(m.user)
		}
		return m, nil

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {

		case tea.KeyTab.String():
			if m.mode == appModeInput {
				m.mode = appModeDefault
			}
			return m, nil

		case tea.KeyEnter.String():
			if m.mode == appModeInput {
				m.mode = ""
				email := m.textInput.Value()
				m.textInput.SetValue("")
				if email == "" {
					return m, nil
				}

				return m, commands.GetUserByEmail(m.dbConnection.UserRepository, email)
			}

			if m.cursor == 0 {
				m.mode = appModeInput
			}

			return m, nil

		// These keys should exit the program.
		case "ctrl+c", "ctrl+q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				if m.commands[m.cursor-1].disabled {
					m.cursor--
				}
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.commands)-1 {
				if m.commands[m.cursor+1].disabled {
					m.cursor++
				}
				m.cursor++
			}
		}
	}

	var cmd tea.Cmd
	if m.mode == appModeInput {
		m.textInput, cmd = m.textInput.Update(msg)
	}
	m.spinner, cmd = m.spinner.Update(msg)

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, cmd
}

func shortenErr(err error, length int) string {
	if len(err.Error()) < length {
		return err.Error()
	}

	return err.Error()[:length] + "..."
}

func (m model) View() string {
	doc := &strings.Builder{}

	tui.RenderTitleRow(width, doc, tui.TitleRowProps{Title: "GO TUI example"})
	doc.WriteString("\n\n")

	var stateDescription string
	if !m.loading {
		stateDescription = m.stateDescription

		if m.mode == appModeInput {
			doc.WriteString(m.textInput.View())
			doc.WriteString("\n\n")
			doc.WriteString("Press tab to return")
			doc.WriteString("\n\n")
		} else {
			// Lists
			renderLists(doc, m)
		}
	} else {
		stateDescription = m.spinner.View()
	}

	tui.RenderStatusBar(doc, tui.NewStatusBarProps(&tui.StatusBarProps{
		Description: stateDescription,
		User:        "NONE",
		StatusState: m.stateStatus,
		Width:       width,
	}))

	// Footer
	doc.WriteString("Press ctrl+q to quit.")
	doc.WriteString("\n")

	// Send the UI for rendering
	return doc.String()
}

func renderLists(doc *strings.Builder, m model) {
	var items []tui.Item
	for _, c := range m.commands {
		items = append(items, tui.Item{
			Value:    c.name,
			Disabled: c.disabled,
		})
	}

	lists := lipgloss.JoinHorizontal(lipgloss.Top,
		tui.RenderListCommands(doc, &tui.ListProps{
			Items:    items,
			Selected: m.cursor,
		}),
		tui.RenderListDisplay(m.secondListHeader, m.secondListValues),
	)

	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lists))
	doc.WriteString("\n\n")
}

func userFieldsToArray(user *storage.User) []string {
	if user == nil {
		return []string{}
	}

	return []string{
		fmt.Sprintf("id: %s", user.Id.Hex()),
		fmt.Sprintf("email: %s", user.Email),
	}
}

func main() {
	// Set DEBUG=true and watch the file for logs: 'tail -f debug.log'
	if len(os.Getenv("DEBUG")) > 0 {
		f, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			log.Fatalf("failed setting the debug log file: %v", err)
		}
		defer f.Close()
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatalf("TUI run error: %v", err)
	}
}
