package tui

import (
	"github.com/charmbracelet/lipgloss"
	"slices"
)

var selected = func(s string) string {
	return lipgloss.NewStyle().Foreground(colorBlue).Render(s)
}
var disabled = func(s string) string {
	return lipgloss.NewStyle().Foreground(colorGray).Render(s)
}

var subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}

var list = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder(), false, true, false, false).
	BorderForeground(subtle).
	MarginRight(2).
	Height(8).
	Width(columnWidth + 1)

var listHeader = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderBottom(true).
	BorderForeground(subtle).
	MarginRight(2).
	Render

func RenderListDisplay(header string, items []string) string {
	return list.Width(columnWidth).
		Border(lipgloss.NormalBorder(), false, false, false, false).
		Render(
			lipgloss.JoinVertical(lipgloss.Left,
				slices.Insert(items, 0, listHeader(header))...,
			),
		)
}
