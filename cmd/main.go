package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/uemuradevexe/fedorento/internal/ui"
)

func main() {
	m := ui.New()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "fedorento: %v\n", err)
		os.Exit(1)
	}
}
