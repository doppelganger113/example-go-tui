package tui

import "github.com/charmbracelet/lipgloss"

const (
	// In real life situations we'd adjust the document to fit the width we've
	// detected. In the case of this example we're hardcoding the width, and
	// later using the detected width only to truncate in order to avoid jaggy
	// wrapping.
	width = 96

	columnWidth = 30
)

const (
	colorRed    = lipgloss.Color("#f54242")
	colorYellow = lipgloss.Color("#b0ad09")
	colorBlue   = lipgloss.Color("#347aeb")
	colorGray   = lipgloss.Color("#636363")
	colorGreen  = lipgloss.Color("#1fb009")
	colorWhite  = lipgloss.Color("#FFFDF5")
)

var (
	whiteStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWhite)
	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorRed)
	yellowStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorYellow)
	grayStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorGray)
	goodStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorGreen)
	blueStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorBlue)
)
