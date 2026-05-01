package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/uemuradevexe/fedorento/content/laravel"
	"github.com/uemuradevexe/fedorento/internal/data"
	"github.com/uemuradevexe/fedorento/internal/highlight"
)

type screen int

const (
	screenSplash screen = iota
	screenMenu
	screenNav
	screenContent
)

// Model is the root Bubble Tea model.
type Model struct {
	screen screen

	// menu
	menuCursor int
	menuItems  []string

	// chapters loaded
	chapters []data.Chapter

	// nav state
	chapterIdx int
	sectionIdx int
	topicIdx   int
	navPane    int // 0 = sections, 1 = topics

	// terminal size
	width  int
	height int
}

func New() Model {
	return Model{
		screen:    screenSplash,
		menuItems: []string{"Laravel 13", "Sair"},
		chapters:  []data.Chapter{laravel.Laravel13},
	}
}

// --- Init ---

func (m Model) Init() tea.Cmd {
	return nil
}

// --- Update ---

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch m.screen {
		case screenSplash:
			if msg.String() == "enter" || msg.String() == " " {
				m.screen = screenMenu
			}
		case screenMenu:
			m = m.updateMenu(msg)
		case screenNav:
			m = m.updateNav(msg)
		case screenContent:
			if msg.String() == "esc" || msg.String() == "q" {
				m.screen = screenNav
			}
		}

		// global quit
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) updateMenu(msg tea.KeyMsg) Model {
	switch msg.String() {
	case "up", "k":
		if m.menuCursor > 0 {
			m.menuCursor--
		}
	case "down", "j":
		if m.menuCursor < len(m.menuItems)-1 {
			m.menuCursor++
		}
	case "enter":
		switch m.menuCursor {
		case 0: // Laravel 13
			m.chapterIdx = 0
			m.sectionIdx = 0
			m.topicIdx = 0
			m.navPane = 0
			m.screen = screenNav
		case 1: // Sair
			// handled below via tea.Quit
		}
		if m.menuCursor == len(m.menuItems)-1 {
			// last item = exit — signal via returning quit in outer Update
		}
	case "q":
		// will be caught by global quit check — noop here
	}
	return m
}

func (m Model) updateNav(msg tea.KeyMsg) Model {
	chapter := m.chapters[m.chapterIdx]
	switch msg.String() {
	case "esc", "q":
		m.screen = screenMenu
	case "tab":
		if m.navPane == 0 {
			m.navPane = 1
		} else {
			m.navPane = 0
		}
	case "up", "k":
		if m.navPane == 0 {
			if m.sectionIdx > 0 {
				m.sectionIdx--
				m.topicIdx = 0
			}
		} else {
			if m.topicIdx > 0 {
				m.topicIdx--
			}
		}
	case "down", "j":
		if m.navPane == 0 {
			if m.sectionIdx < len(chapter.Sections)-1 {
				m.sectionIdx++
				m.topicIdx = 0
			}
		} else {
			section := chapter.Sections[m.sectionIdx]
			if m.topicIdx < len(section.Topics)-1 {
				m.topicIdx++
			}
		}
	case "enter", "l":
		if m.navPane == 1 {
			m.screen = screenContent
		} else {
			m.navPane = 1
		}
	case "h":
		m.navPane = 0
	}
	return m
}

// --- View ---

func (m Model) View() string {
	switch m.screen {
	case screenSplash:
		return SplashView()
	case screenMenu:
		return m.viewMenu()
	case screenNav:
		return m.viewNav()
	case screenContent:
		return m.viewContent()
	}
	return ""
}

func (m Model) viewMenu() string {
	var sb strings.Builder
	sb.WriteString(titleStyle.Render("fedorento") + "\n\n")
	for i, item := range m.menuItems {
		if i == m.menuCursor {
			sb.WriteString(selectedItem.Render("▶ " + item))
		} else {
			sb.WriteString(normalItem.Render("  " + item))
		}
		sb.WriteString("\n")
	}
	sb.WriteString(helpStyle.Render("\n↑/↓ navegar • enter selecionar • ctrl+c sair"))
	return sb.String()
}

func (m Model) viewNav() string {
	chapter := m.chapters[m.chapterIdx]

	// left pane — sections
	var left strings.Builder
	left.WriteString(titleStyle.Render(chapter.Title) + "\n\n")
	for i, sec := range chapter.Sections {
		if i == m.sectionIdx {
			left.WriteString(selectedItem.Render("▶ " + sec.Title))
		} else {
			left.WriteString(normalItem.Render("  " + sec.Title))
		}
		left.WriteString("\n")
	}

	// right pane — topics
	var right strings.Builder
	section := chapter.Sections[m.sectionIdx]
	right.WriteString(titleStyle.Render(section.Title) + "\n\n")
	for i, topic := range section.Topics {
		if m.navPane == 1 && i == m.topicIdx {
			right.WriteString(selectedDimItem.Render("▶ " + topic.Title))
		} else {
			right.WriteString(dimItem.Render("  " + topic.Title))
		}
		right.WriteString("\n")
	}

	leftW := 26
	leftPane := paneStyle.Width(leftW).Render(left.String())
	rightPane := lipgloss.NewStyle().Padding(0, 1).Render(right.String())

	body := lipgloss.JoinHorizontal(lipgloss.Top, leftPane, rightPane)

	paneLabel := "seções"
	if m.navPane == 1 {
		paneLabel = "tópicos"
	}
	help := helpStyle.Render(
		fmt.Sprintf("tab trocar painel [%s] • ↑/↓ navegar • enter abrir • esc voltar", paneLabel),
	)
	return body + "\n" + help
}

func (m Model) viewContent() string {
	chapter := m.chapters[m.chapterIdx]
	section := chapter.Sections[m.sectionIdx]
	topic := section.Topics[m.topicIdx]

	lang := topic.Language
	if lang == "" {
		lang = "php"
	}

	code := highlight.RenderCode(topic.Code, lang)

	var sb strings.Builder
	sb.WriteString(contentTitle.Render(topic.Title) + "\n")
	sb.WriteString(contentDesc.Render(topic.Description) + "\n\n")
	sb.WriteString(code + "\n")
	sb.WriteString(contentExpl.Render(topic.Explanation) + "\n")
	sb.WriteString(helpStyle.Render("\nesc/q voltar"))
	return sb.String()
}
