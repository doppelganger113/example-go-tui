package tui

import (
	"github.com/charmbracelet/lipgloss"
	"slices"
	"strings"
)

type Item struct {
	Value    string
	Disabled bool
}

type ListProps struct {
	Items    []Item
	Selected int
}

func RenderListCommands(doc *strings.Builder, props *ListProps) string {
	var list = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, true, false, false).
		BorderForeground(subtle).
		MarginRight(2).
		Height(8).
		Width(columnWidth + 1)

	var processedItems []string
	for i, item := range props.Items {
		if i == props.Selected && !item.Disabled {
			processedItems = append(processedItems, selected(item.Value))
		} else {
			if item.Disabled {
				processedItems = append(processedItems, disabled(item.Value))
			} else {
				processedItems = append(processedItems, item.Value)
			}
		}
	}

	return list.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			slices.Insert(processedItems, 0, listHeader("Commands"))...,
		),
	)
}
