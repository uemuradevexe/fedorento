package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	splashTitle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B6B")).
			Bold(true)

	splashSub = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true)

	splashBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6B6B")).
			Padding(1, 4)
)

const asciiArt = `
  __         _                        _
 / _|       | |                      | |
| |_ ___  __| | ___  _ __ ___ _ __  | |_ ___
|  _/ _ \/ _` + "`" + ` |/ _ \| '__/ _ \ '_ \ | __/ _ \
| ||  __/ (_| | (_) | | |  __/ | | || ||  __/
|_| \___|\__,_|\___/|_|  \___|_| |_| \__\___|`

// SplashView renders the full splash screen string.
func SplashView() string {
	art := splashTitle.Render(asciiArt)
	sub := splashSub.Render("A smelly CLI for clean Laravel notes.")
	hint := lipgloss.NewStyle().Foreground(lipgloss.Color("#555555")).Render("\npress enter to continue")

	content := art + "\n\n" + sub + hint
	return lipgloss.Place(80, 20, lipgloss.Center, lipgloss.Center, splashBox.Render(content))
}
