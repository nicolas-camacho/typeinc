// Terminal frontend for typeinc, built on Bubble Tea + Lipgloss.
package main

import (
	"fmt"
	"math"
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
	titleStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	ghostStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	typedStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	cursorStyle    = lipgloss.NewStyle().Reverse(true)
	errorStyle     = lipgloss.NewStyle().Background(lipgloss.Color("#D6336C")).Foreground(lipgloss.Color("#FFFFFF"))
	successStyle   = lipgloss.NewStyle().Background(lipgloss.Color("#2FA084")).Foreground(lipgloss.Color("#FFFFFF"))
	goldStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	goldGhostStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Faint(true)
	dimStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("245")).Faint(true)
	errorTextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D6336C")).Bold(true)
	comboStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B35")).Bold(true)

	// day progress bar: blue start, purple past halfway, red in the last 10%
	barBlueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#1B4EF5"))
	barPurpleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8140DC"))
	barRedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#D0311E"))
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
		case game.SceneDayEnd:
			if msg.String() == "enter" {
				m.g.DayEndEnter()
			}
		case game.SceneHR:
			switch msg.String() {
			case "esc":
				m.g.ExitHR()
			case "up":
				m.g.HRMove(-1)
			case "down":
				m.g.HRMove(1)
			case "enter":
				m.g.HRBuy()
			}
		case game.SceneDayStart:
			if msg.String() == "enter" {
				m.g.DayStartEnter()
			}
		case game.SceneStats:
			if msg.String() == "esc" {
				m.g.StatsBack()
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
	case game.SceneDayEnd:
		return m.viewDayEnd()
	case game.SceneHR:
		return m.viewHR()
	case game.SceneDayStart:
		return m.viewDayStart()
	case game.SceneStats:
		return m.viewStats()
	}
	return m.viewPlay()
}

func (m model) viewMenu() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("TYPE.Inc"))
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
	typedSt, ghostSt := typedStyle, ghostStyle
	if m.g.WordGolden {
		typedSt, ghostSt = goldStyle, goldGhostStyle
	}
	var b strings.Builder
	b.WriteString(" ")
	if m.g.Pos > 0 {
		b.WriteString(typedSt.Render(w[:m.g.Pos]))
	}
	if m.g.Pos < len(w) {
		if m.g.Phase == game.PhaseTyping {
			b.WriteString(cursorStyle.Render(string(w[m.g.Pos])))
		} else {
			b.WriteString(ghostSt.Render(string(w[m.g.Pos])))
		}
		b.WriteString(ghostSt.Render(w[m.g.Pos+1:]))
	}
	b.WriteString(" ")
	return b.String()
}

func (m model) viewPlay() string {
	ui := m.g.UI()

	// points gained floats above the word (a golden or urgent word
	// announces itself there while typing); streak + combo sit below it
	gain := " "
	if m.g.Phase == game.PhaseSuccess {
		gain = goldStyle.Render(fmt.Sprintf("+%d", m.g.LastGain))
	} else if m.g.Phase == game.PhaseTyping && m.g.InspectionActive() {
		gain = errorTextStyle.Render(fmt.Sprintf("%s %d", ui.InspectionLabel, int(math.Ceil(m.g.InspectionLeft()))))
	} else if m.g.Phase == game.PhaseTyping && m.g.WordGolden {
		gain = goldStyle.Render(ui.GoldenLabel)
	}
	streakSt := dimStyle
	if m.g.ComboMult() > 1 {
		streakSt = comboStyle
	}
	streak := streakSt.Render(fmt.Sprintf("%s %d · x%.1f", strings.ToLower(ui.StreakLabel), m.g.Streak, m.g.ComboMult()))
	center := lipgloss.JoinVertical(lipgloss.Center, gain, "", m.styledWord(), "", streak)

	// interns column on the far left; the word keeps the remaining width
	wordH := m.height - 4
	word := lipgloss.Place(m.width, wordH, lipgloss.Center, lipgloss.Center, center)
	if interns := m.internColumn(); interns != "" {
		internW := lipgloss.Width(interns) + 2
		col := lipgloss.Place(internW, wordH, lipgloss.Left, lipgloss.Center,
			lipgloss.NewStyle().PaddingLeft(1).Render(interns))
		rest := lipgloss.Place(max(1, m.width-internW), wordH, lipgloss.Center, lipgloss.Center, center)
		word = lipgloss.JoinHorizontal(lipgloss.Top, col, rest)
	}

	// score, day and quota, top-right; the coffee-frenzy banner, top-left
	top := dimStyle.Render(fmt.Sprintf("%s %d  ", ui.DayLabel, m.g.Day)) +
		dimStyle.Render(fmt.Sprintf("%s %d  ", ui.QuotaLabel, m.g.DayQuota())) +
		goldStyle.Render(fmt.Sprintf("%d ", m.g.Score))
	topLine := lipgloss.PlaceHorizontal(m.width, lipgloss.Right, top)
	if m.g.FrenzyActive() {
		banner := goldStyle.Render(fmt.Sprintf(" %s %d", ui.FrenzyLabel, int(math.Ceil(m.g.FrenzyLeft()))))
		topLine = banner + lipgloss.PlaceHorizontal(max(1, m.width-lipgloss.Width(banner)), lipgloss.Right, top)
	}

	// event announce line under the top row, visible for the whole event
	toast := " "
	if m.g.EventToastVisible() {
		toast = dimStyle.Render(m.g.EventToast)
	}
	toastLine := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, toast)

	// command bar / unknown-command message above the day bar
	bottom := " "
	if m.g.CommandMode {
		bottom = typedStyle.Render(m.g.CommandBuf + "█")
	} else if m.g.CommandErr > 0 {
		bottom = errorTextStyle.Render(strings.ToLower(m.g.UI().UnknownCmd))
	}
	bottomLine := lipgloss.PlaceHorizontal(m.width, lipgloss.Center, bottom)

	return topLine + "\n" + toastLine + "\n" + word + "\n" + bottomLine + "\n" + m.dayBar()
}

// internColumn renders one line per intern: typed prefix bright, the rest
// as a ghost, plus a brief gold "+n" flash on completion. Empty string
// when no interns are hired.
func (m model) internColumn() string {
	if len(m.g.Interns) == 0 {
		return ""
	}
	lines := make([]string, 0, len(m.g.Interns))
	for _, in := range m.g.Interns {
		var b strings.Builder
		if in.Pos > 0 {
			b.WriteString(typedStyle.Render(in.Word[:in.Pos]))
		}
		if in.Pos < len(in.Word) {
			b.WriteString(ghostStyle.Render(in.Word[in.Pos:]))
		}
		if in.GainTimer > 0 {
			b.WriteString(goldStyle.Render(fmt.Sprintf(" +%d", in.LastGain)))
		}
		lines = append(lines, b.String())
	}
	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// dayBar renders the very-bottom day progress bar, filling as the day
// burns, with the remaining-time clock at its left. Blue at first, purple
// past halfway, red in the last 10%.
func (m model) dayBar() string {
	p := m.g.DayProgress()
	barSt := barBlueStyle
	switch {
	case p >= 0.9:
		barSt = barRedStyle
	case p >= 0.5:
		barSt = barPurpleStyle
	}

	clock := m.g.DayClock()
	width := max(1, m.width-len(clock)-3)
	filled := min(width, int(p*float64(width)+0.5))
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return " " + barSt.Render(clock+" "+bar)
}

func (m model) viewShop() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(spaced(m.g.UI().ShopTitle)))
	b.WriteString("\n\n\n")
	for i, up := range m.g.ShopItems() {
		level := m.g.UpgradeLevel(up.ID)
		cost := m.g.UpgradeCost(up.ID)
		line := fmt.Sprintf("%s  nv.%d  %s   ", m.g.UpgradeName(up.ID), level, m.g.UpgradeEffect(up.ID))

		var price string
		switch {
		case m.g.UpgradeMaxed(up.ID):
			price = dimStyle.Render(strings.ToLower(m.g.UI().MaxedLabel))
		case m.g.Score < cost:
			price = errorTextStyle.Render(fmt.Sprintf("%d", cost))
		default:
			price = goldStyle.Render(fmt.Sprintf("%d", cost))
		}

		if i == m.g.ShopIndex {
			b.WriteString(typedStyle.Render("> "+line) + price)
		} else {
			b.WriteString(ghostStyle.Render("  "+line) + price)
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
	if m.g.ShopQuipText != "" {
		b.WriteString(dimStyle.Render(m.g.ShopQuipText))
		b.WriteString("\n\n")
	}
	b.WriteString(dimStyle.Render(strings.ToLower(m.g.UI().ShopHint)))

	body := lipgloss.Place(m.width, m.height-1, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, strings.Split(b.String(), "\n")...))
	scoreLine := lipgloss.PlaceHorizontal(m.width, lipgloss.Right,
		goldStyle.Render(fmt.Sprintf("%d ", m.g.Score)))
	return scoreLine + "\n" + body
}

func (m model) viewDayEnd() string {
	ui := m.g.UI()
	title := goldStyle.Render(spaced(ui.PaydayTitle))
	if m.g.DayEndFired {
		title = errorTextStyle.Render(spaced(ui.FiredTitle))
	}

	quip := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Width(min(70, max(20, m.width-8))).
		Align(lipgloss.Center).
		Render(m.g.DayEndQuipVisible())

	stats, breakdown, hint := " ", " ", " "
	if m.g.DayEndQuipDone() {
		hint = dimStyle.Render(strings.ToLower(ui.DayEndHint))
	}
	if m.g.DayEndSummary() {
		worst := ui.NoneLabel
		if m.g.LastDayWorst != "" {
			worst = fmt.Sprintf("%s (%d)", m.g.LastDayWorst, m.g.LastDayWorstN)
		}
		stats = dimStyle.Render(fmt.Sprintf("%s: %d · %s: %d · %s: %d · %s: %s",
			ui.StatsWords, m.g.LastDayWords, ui.StatsFails, m.g.LastDayFails,
			ui.StatsStreak, m.g.LastDayBestStreak, ui.StatsWorst, worst))
		breakdown = errorTextStyle.Render(fmt.Sprintf("%s: -%d", ui.QuotaLabel, m.g.DayEndQuota))
		if !m.g.DayEndFired {
			if m.g.DayEndGainHR > 0 {
				breakdown += "   " + goldStyle.Render(fmt.Sprintf("%s: +%d", ui.HRPointsLbl, m.g.DayEndGainHR))
			} else {
				breakdown += "   " + dimStyle.Render(fmt.Sprintf("%s %d", ui.NextPayLabel, m.g.DayEndNextPay))
			}
		}
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, "", quip, "", stats, breakdown, "", hint))
}

func (m model) viewDayStart() string {
	ui := m.g.UI()
	title := titleStyle.Render(spaced(fmt.Sprintf("%s %d", ui.DayLabel, m.g.Day)))
	quip := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Width(min(70, max(20, m.width-8))).
		Align(lipgloss.Center).
		Render(m.g.DayStartVisible())
	hint := " "
	if m.g.DayStartDone() {
		hint = dimStyle.Render(strings.ToLower(ui.DayEndHint))
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, "", quip, "", hint))
}

func (m model) viewStats() string {
	ui := m.g.UI()
	var b strings.Builder
	b.WriteString(titleStyle.Render(spaced(ui.StatsTitle)))
	b.WriteString("\n\n\n")
	b.WriteString(typedStyle.Render(fmt.Sprintf("%s: %d", ui.StatsWords, m.g.TotalWords)))
	b.WriteString("\n")
	b.WriteString(typedStyle.Render(fmt.Sprintf("%s: %d", ui.StatsFails, m.g.TotalFails)))
	b.WriteString("\n")
	b.WriteString(typedStyle.Render(fmt.Sprintf("%s: %d", ui.StatsStreak, m.g.BestStreak)))
	b.WriteString("\n")
	b.WriteString(typedStyle.Render(fmt.Sprintf("%s: %d", ui.StatsBestDay, m.g.BestDay)))
	b.WriteString("\n\n")
	b.WriteString(dimStyle.Render(ui.StatsTop))
	b.WriteString("\n")
	top := m.g.TopFailWords(3)
	if len(top) == 0 {
		b.WriteString(ghostStyle.Render(ui.NoneLabel))
		b.WriteString("\n")
	}
	for i, s := range top {
		b.WriteString(typedStyle.Render(fmt.Sprintf("%d. %s (%d)", i+1, s.Word, s.N)))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(dimStyle.Render(strings.ToLower(ui.StatsHint)))

	block := lipgloss.NewStyle().PaddingLeft(6).Render(
		lipgloss.JoinVertical(lipgloss.Left, strings.Split(b.String(), "\n")...))
	return lipgloss.Place(m.width, m.height, lipgloss.Left, lipgloss.Center, block)
}

func (m model) viewHR() string {
	ui := m.g.UI()
	var b strings.Builder
	b.WriteString(titleStyle.Render(spaced(ui.HRTitle)))
	b.WriteString("\n\n\n")
	for i, up := range game.HRUpgrades {
		cost := m.g.HRUpgradeCost(up.ID)
		line := fmt.Sprintf("%s  nv.%d  %s   ", m.g.UpgradeName(up.ID), m.g.HRUpgradeLevel(up.ID), m.g.HRUpgradeEffect(up.ID))

		var price string
		switch {
		case m.g.HRUpgradeMaxed(up.ID):
			price = dimStyle.Render(strings.ToLower(ui.MaxedLabel))
		case m.g.HRPoints < cost:
			price = errorTextStyle.Render(fmt.Sprintf("%d", cost))
		default:
			price = goldStyle.Render(fmt.Sprintf("%d", cost))
		}

		if i == m.g.HRIndex {
			b.WriteString(typedStyle.Render("> "+line) + price)
		} else {
			b.WriteString(ghostStyle.Render("  "+line) + price)
		}
		b.WriteString("\n")
	}
	b.WriteString("\n")
	if m.g.HRQuipText != "" {
		b.WriteString(dimStyle.Render(m.g.HRQuipText))
		b.WriteString("\n\n")
	}
	b.WriteString(dimStyle.Render(strings.ToLower(ui.HRHint)))

	body := lipgloss.Place(m.width, m.height-1, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, strings.Split(b.String(), "\n")...))
	ptsLine := lipgloss.PlaceHorizontal(m.width, lipgloss.Right,
		goldStyle.Render(fmt.Sprintf("%s: %d ", ui.HRPointsLbl, m.g.HRPoints)))
	return ptsLine + "\n" + body
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
