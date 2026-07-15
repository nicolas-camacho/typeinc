// Terminal frontend for typeinc, built on Bubble Tea + Lipgloss.
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"typeinc/internal/game"
)

const tickInterval = 50 * time.Millisecond

// The TUI keeps the terminal's own background color; only text is styled.
var (
	titleStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	ghostStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	typedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	cursorStyle  = lipgloss.NewStyle().Reverse(true)
	errorStyle   = lipgloss.NewStyle().Background(lipgloss.Color("#D6336C")).Foreground(lipgloss.Color("#FFFFFF"))
	successStyle = lipgloss.NewStyle().Background(lipgloss.Color("#2FA084")).Foreground(lipgloss.Color("#FFFFFF"))
	goldStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	dimStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Faint(true)
	errorTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D6336C")).Bold(true)
	comboStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B35")).Bold(true)
)

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(tickInterval, func(t time.Time) tea.Msg { return tickMsg(t) })
}

type model struct {
	g      *game.Game
	width  int
	height int
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tickMsg:
		m.g.Tick(tickInterval.Seconds())
		return m, tick()
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.g.Save()
			return m, tea.Quit
		}
		switch m.g.Scene {
		case game.SceneMenu:
			switch msg.String() {
			case "esc":
				m.g.Save()
				return m, tea.Quit
			case "up":
				m.g.MenuMove(-1)
			case "down":
				m.g.MenuMove(1)
			case "left":
				m.g.MenuAdjust(-1)
			case "right":
				m.g.MenuAdjust(1)
			case "enter":
				m.g.MenuSelect()
				if m.g.QuitRequested {
					return m, tea.Quit
				}
			}
		case game.SceneIntro:
			if msg.String() == "enter" {
				m.g.IntroEnter()
			}
		case game.ScenePlay:
			switch msg.String() {
			case "esc":
				if m.g.CommandMode {
					m.g.CommandCancel()
				} else {
					m.g.BackToMenu()
				}
			case "backspace":
				m.g.CommandBackspace()
			case "enter":
				if m.g.CommandMode {
					m.g.CommandSubmit()
				}
			default:
				if msg.Type == tea.KeyRunes {
					for _, r := range msg.Runes {
						m.g.TypeChar(r)
					}
				}
			}
		case game.SceneShop:
			switch msg.String() {
			case "esc":
				m.g.ExitShop()
			case "up":
				m.g.ShopMove(-1)
			case "down":
				m.g.ShopMove(1)
			case "enter":
				m.g.ShopBuy()
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return ""
	}
	switch m.g.Scene {
	case game.SceneMenu:
		return m.viewMenu()
	case game.SceneIntro:
		return m.viewIntro()
	case game.SceneShop:
		return m.viewShop()
	}
	return m.viewPlay()
}

func (m model) viewMenu() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("T Y P E I N C"))
	b.WriteString("\n\n\n")
	for i, item := range m.g.MenuItems() {
		if i == m.g.MenuIndex {
			b.WriteString(typedStyle.Render("> " + item.Label))
		} else {
			b.WriteString(ghostStyle.Render("  " + item.Label))
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(dimStyle.Render(strings.ToLower(m.g.UI().MenuHint)))
	// menu block sits left-aligned with a small margin
	block := lipgloss.NewStyle().PaddingLeft(6).Render(
		lipgloss.JoinVertical(lipgloss.Left, strings.Split(b.String(), "\n")...))
	return lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Center, block)
}

func (m model) viewIntro() string {
	text := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Width(min(70, max(20, m.width-8))).
		Align(lipgloss.Center).
		Render(m.g.IntroVisible())
	hint := " "
	if m.g.IntroLineDone() {
		hint = dimStyle.Render(strings.ToLower(m.g.UI().IntroHint))
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, text, "", hint))
}

// styledWord renders the word with the typed prefix bright, the current
// letter as a block cursor, and the rest as a faint ghost. During the
// error flash the whole word gets a pink background; during the success
// phase, a green one.
func (m model) styledWord() string {
	w := m.g.Word
	if m.g.ErrorFlashing() {
		return errorStyle.Render(" " + w + " ")
	}
	if m.g.Phase == game.PhaseSuccess {
		return successStyle.Render(" " + w + " ")
	}
	var b strings.Builder
	b.WriteString(" ")
	if m.g.Pos > 0 {
		b.WriteString(typedStyle.Render(w[:m.g.Pos]))
	}
	if m.g.Pos < len(w) {
		if m.g.Phase == game.PhaseTyping {
			b.WriteString(cursorStyle.Render(string(w[m.g.Pos])))
		} else {
			b.WriteString(ghostStyle.Render(string(w[m.g.Pos])))
		}
		b.WriteString(ghostStyle.Render(w[m.g.Pos+1:]))
	}
	b.WriteString(" ")
	return b.String()
}

func (m model) viewPlay() string {
	// points gained floats above the word during the success phase
	gain := " "
	if m.g.Phase == game.PhaseSuccess {
		gain = goldStyle.Render(fmt.Sprintf("+%d", m.g.LastGain))
	}
	center := lipgloss.JoinVertical(lipgloss.Center, gain, "", m.styledWord())
	word := lipgloss.Place(m.width, m.height-2, lipgloss.Center, lipgloss.Center, center)

	// score and combo multiplier, top-right
	comboSt := dimStyle
	if m.g.ComboMult() > 1 {
		comboSt = comboStyle
	}
	top := goldStyle.Render(fmt.Sprintf("%d ", m.g.Score)) +
		comboSt.Render(fmt.Sprintf("x%.1f ", m.g.ComboMult()))
	topLine := lipgloss.PlaceHorizontal(m.width, lipgloss.Right, top)

	// command bar / unknown-command message, bottom
	bottom := " "
	if m.g.CommandMode {
		bottom = typedStyle.Render(m.g.CommandBuf + "█")
	} else if m.g.CommandErr > 0 {
		bottom = errorTextStyle.Render(strings.ToLower(m.g.UI().UnknownCmd))
	}
	bottomLine := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, bottom)

	return topLine + "\n" + word + "\n" + bottomLine
}

func (m model) viewShop() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(spaced(m.g.UI().ShopTitle)))
	b.WriteString("\n\n\n")
	for i, up := range game.Upgrades {
		level := m.g.UpgradeLevel(up.ID)
		cost := m.g.UpgradeCost(up.ID)
		line := fmt.Sprintf("%s  nv.%d  %s   ", m.g.UpgradeName(up.ID), level, m.g.UpgradeEffect(up.ID))

		priceStyle := goldStyle
		if m.g.Score < cost {
			priceStyle = errorTextStyle
		}
		price := priceStyle.Render(fmt.Sprintf("%d", cost))

		if i == m.g.ShopIndex {
			b.WriteString(typedStyle.Render("> "+line) + price)
		} else {
			b.WriteString(ghostStyle.Render("  "+line) + price)
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(dimStyle.Render(strings.ToLower(m.g.UI().ShopHint)))

	body := lipgloss.Place(m.width, m.height-1, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, strings.Split(b.String(), "\n")...))
	scoreLine := lipgloss.PlaceHorizontal(m.width, lipgloss.Right,
		goldStyle.Render(fmt.Sprintf("%d ", m.g.Score)))
	return scoreLine + "\n" + body
}

// spaced renders a title with letter spacing ("SHOP" -> "S H O P").
func spaced(s string) string {
	return strings.Join(strings.Split(s, ""), " ")
}

func main() {
	g := game.New()
	g.SavePath = game.DefaultSavePath()
	g.Load()
	p := tea.NewProgram(model{g: g}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	g.Save()
}
