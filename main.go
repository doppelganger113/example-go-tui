package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
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
}

func initialModel() model {
	return model{
		stateDescription: "Initializing...",
		commands: []command{
			{name: "Set user"},
			{name: "Fetch token", disabled: true},
			{name: "Other..."},
		},
	}
}

func (m model) Init() tea.Cmd {
	// No I/O
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

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
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.commands)-1 {
				m.cursor++
			}
		}
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	doc := &strings.Builder{}

	tui.RenderTitleRow(width, doc, tui.TitleRowProps{Title: "GO TUI example"})
	doc.WriteString("\n\n")

	doc.WriteString(fmt.Sprintf("Cursor: %d", m.cursor))
	doc.WriteString("\n\n")

	tui.RenderStatusBar(doc, tui.NewStatusBarProps(&tui.StatusBarProps{
		Description: m.stateDescription,
		User:        "NONE",
		StatusState: tui.StatusBarStateBlue,
		Width:       width,
	}))

	// Footer
	doc.WriteString("Press q to quit.")
	doc.WriteString("\n")

	// Send the UI for rendering
	return doc.String()
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
