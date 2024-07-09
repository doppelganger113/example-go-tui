package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"gotui/internal/storage"
	"gotui/internal/tui"
	"log"
	"os"
	"strings"
)

type appMode string

const (
	appModeInput   appMode = "input"
	appModeLoading appMode = "loading"
	appModeDefault appMode = ""
)

type command struct {
	disabled bool
	name     string
}

type model struct {
	spinner          spinner.Model
	mode             appMode
	stateDescription string
	stateStatus      tui.StatusBarState
	commands         []command // items on the to-do list
	cursor           int       // which to-do list item our cursor is pointing at
	textInput        textinput.Model
	dbConnection     *dbConnection
	user             *storage.User
	loading          bool
	secondListHeader string
	secondListValues []string
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

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "john@email.com"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		spinner:          s,
		textInput:        ti,
		stateStatus:      tui.StatusBarStateBlue,
		stateDescription: "Initializing...",
		commands: []command{
			{name: "Set user"},
			{name: "Fetch token", disabled: true},
			{name: "Other..."},
		},
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return tea.Batch(
		initDatabase,
		m.spinner.Tick,
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case dbConnection:
		m.stateDescription = ""
		m.dbConnection = &msg
		if m.dbConnection != nil {
			if m.dbConnection.err != nil {
				m.stateStatus = tui.StatusBarStateRed
				m.stateDescription = "Failed to connect to database."
			} else {
				m.stateStatus = tui.StatusBarStateGreen
				m.stateDescription = "Connected to database"
			}
		}
		return m, nil
	case getUserByEmailMsg:
		m.user = msg.user
		if msg.err != nil {
			m.stateDescription = msg.err.Error()
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

	case getTokenByUserEmail:
		m.loading = false
		m.mode = appModeDefault
		if msg.err != nil {
			m.stateDescription = msg.err.Error()
			m.stateStatus = tui.StatusBarStateRed
		} else {
			m.stateDescription = "Retrieved token"
			m.stateStatus = tui.StatusBarStateGreen
			m.secondListHeader = "Token"
			m.secondListValues = []string{"1234"}
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

				return m, getUserByEmail(m.dbConnection.userRepository, email)
			}

			if m.cursor == 0 {
				m.mode = appModeInput
			} else if m.cursor == 1 && m.user != nil {
				m.stateDescription = fmt.Sprintf("Fetching %s token...", m.user.Email)
				m.loading = true
				m.mode = appModeLoading
				m.stateStatus = tui.StatusBarStateBlue
				return m, getLatestTokenByUserEmail(m.user.Email)
			}

			return m, nil

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.mode == appModeLoading {
				return m, nil
			}
			if m.cursor > 0 {
				if m.commands[1].disabled {
					m.cursor--
				}
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.mode == appModeLoading {
				return m, nil
			}
			if m.cursor < len(m.commands)-1 {
				if m.commands[1].disabled {
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

func (m model) View() string {
	doc := &strings.Builder{}

	tui.RenderTitleRow(doc, tui.TitleRowProps{Title: "GO TUI example"})
	doc.WriteString("\n")

	if m.mode == appModeInput {
		doc.WriteString(m.textInput.View())
		doc.WriteString("\n\n")
		doc.WriteString("Press tab to return")
		doc.WriteString("\n\n")
	} else {

		// Lists
		renderLists(doc, m)
	}

	renderStatusBar(doc, m)

	// Footer
	doc.WriteString("Press q to quit.")
	doc.WriteString("\n")

	// Send the UI for rendering
	return doc.String()
}

func renderStatusBar(doc *strings.Builder, m model) {
	statusBarProps := tui.StatusBarProps{
		StatusState: m.stateStatus,
		Description: m.stateDescription,
	}

	if m.user != nil {
		statusBarProps.User = m.user.Email
	}

	tui.RenderStatusBar(doc, tui.NewStatusBarProps(&statusBarProps))
}

func renderLists(doc *strings.Builder, m model) {
	var items []tui.Item
	for _, c := range m.commands {
		items = append(items, tui.Item{
			Value:    c.name,
			Disabled: c.disabled,
		})
	}

	var secondList string
	if m.loading {
		secondList = tui.RenderRightListView(m.spinner.View(), nil)
	} else {
		secondList = tui.RenderRightListView(m.secondListHeader, m.secondListValues)
	}

	lists := lipgloss.JoinHorizontal(lipgloss.Top,
		tui.RenderCommandsList(doc, &tui.ListProps{
			Items:    items,
			Selected: m.cursor,
		}),
		secondList,
	)

	doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, lists))
	doc.WriteString("\n\n")
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
