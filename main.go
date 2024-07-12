package main

import (
	"github.com/charmbracelet/bubbles/spinner"
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
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
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

	// Is it a key press?
	case tea.KeyMsg:
		// Cool, what was the actual key pressed?
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
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
		renderLists(doc, m)
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
	doc.WriteString("Press q to quit.")
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
