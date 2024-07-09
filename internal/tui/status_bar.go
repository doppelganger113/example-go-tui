package tui

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type StatusBarState string

const (
	StatusBarStateGreen  StatusBarState = "green"
	StatusBarStateYellow StatusBarState = "yellow"
	StatusBarStateBlue   StatusBarState = "blue"
	StatusBarStateGray   StatusBarState = "gray"
	StatusBarStateRed    StatusBarState = "red"
)

type StatusBarProps struct {
	Status      string
	Description string
	User        string
	StatusState StatusBarState
	Width       int
}

func NewStatusBarProps(props *StatusBarProps) StatusBarProps {
	defaultProps := StatusBarProps{
		Status:      "STATUS",
		Description: "",
		User:        "NONE",
		StatusState: StatusBarStateGreen,
		Width:       98,
	}
	if props == nil {
		return defaultProps
	}

	if props.User != "" {
		defaultProps.User = props.User
	}
	if props.Status != "" {
		defaultProps.Status = props.Status
	}
	if props.Description != "" {
		defaultProps.Description = props.Description
	}
	if props.Width > 0 {
		defaultProps.Width = props.Width
	}
	if props.StatusState != "" {
		defaultProps.StatusState = props.StatusState
	}

	return defaultProps
}

var (
	statusNugget = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Padding(0, 1)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#343433", Dark: "#C1C6B2"}).
			Background(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#353533"})

	statusStyle = lipgloss.NewStyle().
			Inherit(statusBarStyle).
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color(colorGreen)).
			Padding(0, 1).
			MarginRight(1)

	statusStyleErr    = statusStyle.Background(colorRed)
	statusStyleGray   = statusStyle.Background(colorGray)
	statusStyleYellow = statusStyle.Background(colorYellow)
	statusStyleBlue   = statusStyle.Background(colorBlue)

	encodingStyle = statusNugget.
			Background(lipgloss.Color("#A550DF")).
			Align(lipgloss.Right)

	statusText = lipgloss.NewStyle().Inherit(statusBarStyle)

	fishCakeStyle = statusNugget.Background(lipgloss.Color("#6124DF"))
)

func RenderStatusBar(doc *strings.Builder, props StatusBarProps) {

	w := lipgloss.Width

	statusKey := getStatusStyle(props)

	encoding := encodingStyle.Render("USER")
	fishCake := fishCakeStyle.Render(props.User)
	statusVal := statusText.
		Width(props.Width - w(statusKey) - w(encoding) - w(fishCake)).
		Render(getStatusDescription(props))

	bar := lipgloss.JoinHorizontal(lipgloss.Top,
		statusKey,
		statusVal,
		encoding,
		fishCake,
	)

	doc.WriteString(statusBarStyle.Width(props.Width).Render(bar))

	doc.WriteString("\n\n")
}

func getStatusStyle(props StatusBarProps) string {
	if props.StatusState == StatusBarStateRed {
		return statusStyleErr.Render(props.Status)
	} else if props.StatusState == StatusBarStateGray {
		return statusStyleGray.Render(props.Status)
	} else if props.StatusState == StatusBarStateYellow {
		return statusStyleYellow.Render(props.Status)
	} else if props.StatusState == StatusBarStateBlue {
		return statusStyleBlue.Render(props.Status)
	}

	return statusStyle.Render(props.Status)
}

func getStatusDescription(props StatusBarProps) string {
	if props.StatusState == StatusBarStateRed {
		return whiteStyle.Render(props.Description)
	} else if props.StatusState == StatusBarStateYellow {
		return whiteStyle.Render(props.Description)
	} else if props.StatusState == StatusBarStateGray {
		return whiteStyle.Render(props.Description)
	} else if props.StatusState == StatusBarStateBlue {
		return whiteStyle.Render(props.Description)
	}

	return whiteStyle.Render(props.Description)
}
