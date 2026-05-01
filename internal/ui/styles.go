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

	badgeBaseStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1)

	badgeWebStyle = badgeBaseStyle.Copy().
			Foreground(lipgloss.Color("#0B1F14")).
			Background(lipgloss.Color("#6BCB77"))

	badgeAPIStyle = badgeBaseStyle.Copy().
			Foreground(lipgloss.Color("#0D1B2A")).
			Background(lipgloss.Color("#4D96FF"))

	badgeSharedStyle = badgeBaseStyle.Copy().
				Foreground(lipgloss.Color("#2A1A00")).
				Background(lipgloss.Color("#FFD93D"))

	legendStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

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

	scrollTrackStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#333333"))

	scrollFillStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD93D")).
			Bold(true)

	scrollStatusStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFD93D")).
				Bold(true)

	paneStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, true, false, false).
			BorderForeground(lipgloss.Color("#333333")).
			Padding(0, 1)
)
