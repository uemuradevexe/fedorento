package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
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

const contentHelpText = "esc/q voltar • ↑/↓ ou j/k rolar • space/pgdn descer página • b/pgup subir página"

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

	content viewport.Model
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
		m.syncContentViewport()

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

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
			} else {
				var cmd tea.Cmd
				m.content, cmd = m.content.Update(msg)
				return m, cmd
			}
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
			m.setContentViewport()
		} else {
			m.navPane = 1
		}
	case "h":
		m.navPane = 0
	}
	return m
}

func (m *Model) syncContentViewport() {
	if m.width <= 0 || m.height <= 0 {
		return
	}

	helpHeight := lipgloss.Height(helpStyle.Render(contentHelpText))
	contentHeight := m.height - helpHeight - 1
	if contentHeight < 1 {
		contentHeight = 1
	}

	contentWidth := m.width - 4
	if contentWidth < 20 {
		contentWidth = m.width
	}
	if contentWidth < 1 {
		contentWidth = 1
	}

	m.content.Width = contentWidth
	m.content.Height = contentHeight
	if m.screen == screenContent {
		m.setContentViewport()
	}
}

func (m *Model) setContentViewport() {
	chapter := m.chapters[m.chapterIdx]
	section := chapter.Sections[m.sectionIdx]
	topic := section.Topics[m.topicIdx]

	lang := topic.Language
	if lang == "" {
		lang = "php"
	}

	if m.content.Width == 0 || m.content.Height == 0 {
		m.syncContentViewport()
		if m.content.Width == 0 {
			m.content.Width = 80
		}
		if m.content.Height == 0 {
			m.content.Height = 20
		}
	}

	code := highlight.RenderCode(topic.Code, lang)
	title := contentTitle.Render(topic.Title)
	if badge := topicAudienceBadge(topic.Audience); badge != "" {
		title = lipgloss.JoinHorizontal(lipgloss.Center, title, " ", badge)
	}

	var sb strings.Builder
	sb.WriteString(title + "\n")
	sb.WriteString(contentDesc.Width(m.content.Width).Render(topic.Description) + "\n\n")
	sb.WriteString(code + "\n")
	sb.WriteString(contentExpl.Width(m.content.Width).Render(topic.Explanation) + "\n")

	m.content.SetContent(sb.String())
	m.content.GotoTop()
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
		label := topic.Title
		if short := topicAudienceShort(topic.Audience); short != "" {
			label += " [" + short + "]"
		}
		if m.navPane == 1 && i == m.topicIdx {
			right.WriteString(selectedDimItem.Render("▶ " + label))
		} else {
			right.WriteString(dimItem.Render("  " + label))
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
	legend := legendStyle.Render(
		topicAudienceBadge("web") + " " +
			topicAudienceBadge("api") + " " +
			topicAudienceBadge("shared"),
	)
	help := helpStyle.Render(
		fmt.Sprintf("tab trocar painel [%s] • ↑/↓ navegar • enter abrir • esc voltar", paneLabel),
	)
	footer := lipgloss.JoinHorizontal(lipgloss.Top, help, "  ", legend)
	footer = lipgloss.NewStyle().Width(m.width).Render(footer)
	return body + "\n" + footer
}

func (m Model) viewContent() string {
	progress := m.content.ScrollPercent()
	status := fmt.Sprintf("scroll %d%%", int(m.content.ScrollPercent()*100))
	if m.content.AtTop() {
		status = "topo"
	}
	if m.content.AtBottom() {
		status = "fim"
	}

	help := helpStyle.Render(contentHelpText)
	indicator := scrollStatusStyle.Render(status)
	rightWidth := max(18, min(40, m.width/3))
	barWidth := max(10, rightWidth-lipgloss.Width(indicator)-2)
	filled := int(progress * float64(barWidth))
	if m.content.AtBottom() {
		filled = barWidth
	}
	if filled > barWidth {
		filled = barWidth
	}
	bar := scrollFillStyle.Render(strings.Repeat("█", filled)) + scrollTrackStyle.Render(strings.Repeat("░", barWidth-filled))
	right := lipgloss.JoinHorizontal(lipgloss.Top, bar, "  ", indicator)
	right = lipgloss.NewStyle().Width(rightWidth).Align(lipgloss.Right).Render(right)
	leftWidth := max(1, m.width-rightWidth)
	left := lipgloss.NewStyle().Width(leftWidth).Render(help)
	footer := lipgloss.JoinHorizontal(lipgloss.Top, left, right)

	return m.content.View() + "\n" + footer
}

func topicAudienceShort(audience string) string {
	switch audience {
	case "web":
		return "W"
	case "api":
		return "A"
	default:
		return "*"
	}
}

func topicAudienceBadge(audience string) string {
	switch audience {
	case "web":
		return badgeWebStyle.Render("W WEB")
	case "api":
		return badgeAPIStyle.Render("A API")
	default:
		return badgeSharedStyle.Render("* ALL")
	}
}
