package ui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true).
			Padding(0, 1)

	selectedItem = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true).
			PaddingLeft(2)

	normalItem = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#CCCCCC")).
			PaddingLeft(2)

	dimItem = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555")).
		PaddingLeft(4)

	selectedDimItem = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD93D")).
			Bold(true).
			PaddingLeft(4)

	contentTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6BCB77")).
			Bold(true).
			Underline(true).
			MarginBottom(1)

	contentDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#AAAAAA")).
			MarginBottom(1)

	contentExpl = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#DDDDDD")).
			Italic(true).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444")).
			MarginTop(1)

	paneStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color("#333333")).
			Padding(0, 1)
)
