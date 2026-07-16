// Package game holds the shared typing-game logic, free of any rendering
// or input-backend dependencies so both the desktop and TUI frontends can
// drive it.
package game

import (
	"fmt"
	"math"
	"math/rand/v2"
	"sort"
	"strconv"
	"strings"
)

// Scene identifies which screen the game is on.
type Scene int

const (
	SceneMenu Scene = iota
	SceneOptions
	SceneIntro
	ScenePlay
	SceneShop
	SceneDayEnd   // HR end-of-day summary: quota paid or fired
	SceneHR       // meta shop, paid with HR points
	SceneDayStart // HR morning quip before a new day begins
	SceneStats    // global stats screen, from the main menu
)

// Phase timings, in seconds.
const (
	SuccessTime    = 0.5 // green highlight + input lock after completing a word
	ErrorFlashTime = 0.5 // pink highlight after a mistake
	ErrorBlockTime = 0.5 // extra input lock after the error flash ends
	CommandErrTime = 1.5 // "unknown command" message duration
)

// BaseStreakCap is how many streak words count toward the combo
// multiplier before upgrades.
const BaseStreakCap = 5

// Day/quota tuning: a day lasts DayBaseTime seconds of active play; at its
// end HR charges DayQuota() — not covering it means getting fired.
const (
	DayBaseTime  = 90.0 // seconds per day before upgrades (1.5 min)
	DayTimeBonus = 15.0 // extra seconds per "hrday" level
	QuotaBase    = 100  // day-1 quota before scaling
	QuotaGrowth  = 1.5  // per-day exponential growth
	HRQuotaMax   = 5    // "hrquota" level cap (−10% each, max −50%)
)

// Upgrade describes one shop entry; the display name is localized via
// UpgradeName. Cost at level n is BaseCost × 3ⁿ.
type Upgrade struct {
	ID       string
	BaseCost int
}

// Upgrades lists the shop entries in display order.
var Upgrades = []Upgrade{
	{ID: "mult", BaseCost: 50},
	{ID: "streakcap", BaseCost: 40},
}

// HRUpgrades lists the meta-shop entries, paid with HR points and kept
// across runs. Cost at level n is BaseCost × 2ⁿ.
var HRUpgrades = []Upgrade{
	{ID: "hrmult", BaseCost: 3},
	{ID: "hrday", BaseCost: 2},
	{ID: "hrquota", BaseCost: 3},
}

// Phase is the sub-state within ScenePlay.
type Phase int

const (
	PhaseTyping  Phase = iota
	PhaseSuccess       // word completed: highlight + "+x", input locked
	PhaseError         // mistake: flash, then a plain locked second
)

// Languages lists the word languages in display order.
var Languages = []struct {
	Label string
	Code  string
}{
	{Label: "ESPANOL", Code: "es"},
	{Label: "ENGLISH", Code: "en"},
}

// IntroCPS is the typewriter speed of the intro story, in characters per
// second.
const IntroCPS = 20.0

// MenuItem is one main-menu entry.
type MenuItem struct {
	ID    string
	Label string
}

// OptItemIDs lists the options screen entries in display order; labels
// come from UI().
var OptItemIDs = []string{"display", "volume", "back"}

// Game is the full game state. Frontends mutate it only through methods.
type Game struct {
	Scene      Scene
	MenuIndex  int      // selected entry in Languages
	Words      []string // active word list
	Word       string   // word currently being typed
	Pos        int      // number of letters typed correctly so far
	Score      int
	Phase      Phase
	PhaseTimer float64 // seconds left in the current non-typing phase
	LastGain   int     // points awarded for the last completed word

	// incremental state
	Streak         int // consecutive completed words without a mistake
	MultLevel      int // "mult" upgrade level
	StreakCapLevel int // "streakcap" upgrade level

	// day cycle (part of the run)
	Day      int     // current day, 1-based
	DayTimer float64 // active-play seconds left in the day

	// current-day stats (part of the run)
	DayWords      int            // words completed today
	DayFails      int            // mistakes today
	dayFailCounts map[string]int // per-word mistakes today

	// stats of the last settled day, for the summary and morning quips
	LastDayWords  int
	LastDayFails  int
	LastDayWorst  string // most-failed word ("" if no fails)
	LastDayWorstN int

	// global stats, kept across runs
	TotalWords int
	TotalFails int
	FailCounts map[string]int // per-word mistakes, all-time

	// meta currency, kept across runs
	HRPoints     int
	HRMultLevel  int // "hrmult": +10% gain per level, permanent
	HRDayLevel   int // "hrday": +15s day length per level
	HRQuotaLevel int // "hrquota": −10% quota per level, capped

	// end-of-day summary state
	DayEndFired  bool    // this summary is a firing, not a payday
	DayEndFirst  bool    // this payday closed day 1: HR points get explained
	DayEndQuota  int     // quota HR charged (or demanded) this day
	DayEndGainHR int     // HR points earned (0 when fired)
	DayEndQuip   int     // picked quip index
	DayEndLine   int     // current line of a multi-line summary script
	QuipShown    float64 // typewriter progress of the quip, in runes
	lastPayday   int     // previous payday quip, to never repeat
	lastFired    int     // previous firing quip, to never repeat

	// day-start morning quip state
	DayStartText   string // composed quip (may embed yesterday's stats)
	lastStartPlain int    // previous plain morning quip
	lastStartStat  int    // previous stat-based morning quip

	HRIndex int   // selected meta-shop entry
	hrFrom  Scene // scene to return to when leaving the meta shop

	// command mode ("/..." pauses play)
	CommandMode bool
	CommandBuf  string
	CommandErr  float64 // seconds the "unknown command" message stays visible

	ShopIndex int // selected shop entry

	// menu / session state
	OptionsEnabled bool // frontend support for the options screen (desktop)
	QuitRequested  bool // SALIR chosen; the frontend reads this and exits
	LangIndex      int  // selected language (words + interface)
	RunActive      bool // a run exists in memory, CONTINUAR resumes it

	// persistence
	SavePath string   // "" disables saving (tests)
	savedRun *RunData // run loaded from disk, restored by CONTINUAR

	// intro story (plays on every new game); CONTINUAR reuses the scene
	// for a one-line HR quip
	IntroLine   int
	IntroShown  float64 // runes of the current line revealed so far
	introReturn bool    // showing a return quip, not the full story
	returnIdx   int     // quip picked for this return
	lastReturn  int     // previous quip, to never repeat back-to-back

	// options (applied by the desktop frontend)
	DisplayMode int     // index into DisplayModeLabels
	Volume      float64 // master volume 0..1
	OptIndex    int     // selected options entry
}

// New returns a game sitting on the main menu.
func New() *Game {
	return &Game{Scene: SceneMenu, Volume: 0.4, FailCounts: map[string]int{}}
}

// HasSave reports whether CONTINUAR has something to resume: a run in
// memory or one loaded from disk.
func (g *Game) HasSave() bool {
	return g.RunActive || g.savedRun != nil
}

// MenuItems builds the main menu in the interface language; CONTINUAR
// only shows with a resumable run, OPCIONES only when the frontend
// supports it.
func (g *Game) MenuItems() []MenuItem {
	ui := g.UI()
	var items []MenuItem
	if g.HasSave() {
		items = append(items, MenuItem{ID: "continue", Label: ui.MenuContinue})
	}
	items = append(items,
		MenuItem{ID: "new", Label: ui.MenuNew},
		MenuItem{ID: "lang", Label: ui.MenuLang + ": " + Languages[g.LangIndex].Label},
		MenuItem{ID: "hr", Label: fmt.Sprintf("%s: %d", ui.MenuHR, g.HRPoints)},
		MenuItem{ID: "stats", Label: ui.MenuStats},
	)
	if g.OptionsEnabled {
		items = append(items, MenuItem{ID: "options", Label: ui.MenuOptions})
	}
	return append(items, MenuItem{ID: "quit", Label: ui.MenuQuit})
}

// MenuMove moves the menu selection by delta, wrapping around.
func (g *Game) MenuMove(delta int) {
	n := len(g.MenuItems())
	g.MenuIndex = ((g.MenuIndex+delta)%n + n) % n
}

// menuItem is the highlighted entry, clamped in case the item list shrank.
func (g *Game) menuItem() MenuItem {
	items := g.MenuItems()
	if g.MenuIndex >= len(items) {
		g.MenuIndex = len(items) - 1
	}
	return items[g.MenuIndex]
}

// MenuAdjust handles left/right on the highlighted menu entry.
func (g *Game) MenuAdjust(delta int) {
	if g.menuItem().ID == "lang" {
		g.cycleLang(delta)
	}
}

// MenuSelect activates the highlighted menu entry.
func (g *Game) MenuSelect() {
	switch g.menuItem().ID {
	case "continue":
		g.continueRun()
	case "new":
		g.introReturn = false
		g.IntroLine = 0
		g.IntroShown = 0
		g.Scene = SceneIntro
	case "lang":
		g.cycleLang(1)
	case "hr":
		g.HRIndex = 0
		g.hrFrom = SceneMenu
		g.Scene = SceneHR
	case "stats":
		g.Scene = SceneStats
	case "options":
		g.OptIndex = 0
		g.Scene = SceneOptions
	case "quit":
		g.QuitRequested = true
		g.Save()
	}
}

// continueRun resumes the in-memory run (or restores the one from disk)
// and greets the player with a short random HR quip before play.
func (g *Game) continueRun() {
	if !g.RunActive && g.savedRun != nil {
		g.initRun()
		g.Score = g.savedRun.Score
		g.MultLevel = g.savedRun.MultLevel
		g.StreakCapLevel = g.savedRun.StreakCapLevel
		g.Streak = g.savedRun.Streak
		g.Day = g.savedRun.Day
		g.DayTimer = g.savedRun.DayTimer
		g.DayWords = g.savedRun.DayWords
		g.DayFails = g.savedRun.DayFails
		g.dayFailCounts = g.savedRun.DayFailCounts
		if g.dayFailCounts == nil {
			g.dayFailCounts = map[string]int{}
		}
		g.savedRun = nil
		// older saves have no day data; start the day fresh
		if g.Day < 1 {
			g.Day = 1
		}
		if g.DayTimer <= 0 || g.DayTimer > g.DayLength() {
			g.DayTimer = g.DayLength()
		}
	}
	if !g.RunActive {
		g.initRun()
	}

	g.returnIdx = pickQuip(len(ReturnScripts[Languages[g.LangIndex].Code]), &g.lastReturn)
	g.introReturn = true
	g.IntroLine = 0
	g.IntroShown = 0
	g.Scene = SceneIntro
}

// cycleLang switches the language for words and interface. A running game
// keeps its score and upgrades; only the word list (and current word)
// changes.
func (g *Game) cycleLang(delta int) {
	n := len(Languages)
	g.LangIndex = ((g.LangIndex+delta)%n + n) % n
	if g.RunActive {
		g.Words = WordLists[Languages[g.LangIndex].Code]
		g.Word = ""
		g.nextWord()
	}
	g.Save()
}

// initRun resets all run state for a brand-new game.
func (g *Game) initRun() {
	g.Words = WordLists[Languages[g.LangIndex].Code]
	g.Score = 0
	g.Phase = PhaseTyping
	g.PhaseTimer = 0
	g.Streak = 0
	g.MultLevel = 0
	g.StreakCapLevel = 0
	g.CommandMode = false
	g.CommandBuf = ""
	g.CommandErr = 0
	g.ShopIndex = 0
	g.Word = ""
	g.nextWord()
	g.Day = 1
	g.DayTimer = g.DayLength()
	g.resetDayStats()
	g.RunActive = true
}

// resetDayStats clears the current-day counters.
func (g *Game) resetDayStats() {
	g.DayWords = 0
	g.DayFails = 0
	g.dayFailCounts = map[string]int{}
}

// pickQuip returns a random index in [0,n) different from *last (when
// possible) and records it as the new last pick.
func pickQuip(n int, last *int) int {
	idx := rand.IntN(n)
	for idx == *last && n > 1 {
		idx = rand.IntN(n)
	}
	*last = idx
	return idx
}

// introScript is what the intro scene is currently showing: the full
// story on a new game, or a single return quip on CONTINUAR.
func (g *Game) introScript() []string {
	code := Languages[g.LangIndex].Code
	if g.introReturn {
		return []string{ReturnScripts[code][g.returnIdx]}
	}
	return IntroScripts[code]
}

// introLine returns the current script line as runes (the script contains
// multibyte characters).
func (g *Game) introLine() []rune {
	return []rune(g.introScript()[g.IntroLine])
}

// IntroVisible is the typewriter-revealed prefix of the current line.
func (g *Game) IntroVisible() string {
	line := g.introLine()
	return string(line[:min(int(g.IntroShown), len(line))])
}

// IntroLineDone reports whether the current line is fully revealed.
func (g *Game) IntroLineDone() bool {
	return int(g.IntroShown) >= len(g.introLine())
}

// IntroEnter advances the intro: it completes a half-revealed line,
// otherwise moves to the next one; after the last line a fresh run starts
// (the intro only plays on NUEVA PARTIDA).
func (g *Game) IntroEnter() {
	if !g.IntroLineDone() {
		g.IntroShown = float64(len(g.introLine()))
		return
	}
	if g.IntroLine+1 < len(g.introScript()) {
		g.IntroLine++
		g.IntroShown = 0
		return
	}
	if g.introReturn {
		// the run was already restored by continueRun
		g.introReturn = false
		g.Scene = ScenePlay
		return
	}
	g.initRun()
	g.savedRun = nil
	g.Scene = ScenePlay
	g.Save()
}

// OptMove moves the options selection by delta, wrapping around.
func (g *Game) OptMove(delta int) {
	n := len(OptItemIDs)
	g.OptIndex = ((g.OptIndex+delta)%n + n) % n
}

// OptAdjust handles left/right on the highlighted options entry.
func (g *Game) OptAdjust(delta int) {
	switch OptItemIDs[g.OptIndex] {
	case "display":
		n := len(g.UI().DisplayModes)
		g.DisplayMode = ((g.DisplayMode+delta)%n + n) % n
	case "volume":
		g.Volume = math.Round(math.Min(1, math.Max(0, g.Volume+0.1*float64(delta)))*10) / 10
	default:
		return
	}
	g.Save()
}

// OptSelect activates the highlighted options entry.
func (g *Game) OptSelect() {
	switch OptItemIDs[g.OptIndex] {
	case "display":
		g.OptAdjust(1)
	case "back":
		g.Scene = SceneMenu
	}
}

// OptLabel renders one options row in the interface language, including
// its current value.
func (g *Game) OptLabel(id string) string {
	ui := g.UI()
	switch id {
	case "display":
		return ui.OptDisplay + ": " + ui.DisplayModes[g.DisplayMode]
	case "volume":
		return fmt.Sprintf("%s: %d%%", ui.OptVolume, int(g.Volume*100+0.5))
	case "back":
		return ui.OptBack
	}
	return ""
}

// BackToMenu returns to the main menu.
func (g *Game) BackToMenu() {
	g.CommandMode = false
	g.CommandBuf = ""
	g.Scene = SceneMenu
	g.Save()
}

// StreakCap is the number of streak words that count toward the combo.
func (g *Game) StreakCap() int {
	return BaseStreakCap + BaseStreakCap*g.StreakCapLevel
}

// ComboMult is the current streak multiplier: +10% per streak word, capped.
func (g *Game) ComboMult() float64 {
	return 1 + 0.1*float64(min(g.Streak, g.StreakCap()))
}

// WordGain is what completing the current word pays right now:
// len × (1 + mult level) × combo multiplier × HR bonus, rounded.
func (g *Game) WordGain() int {
	return int(math.Round(float64(len(g.Word)) * float64(1+g.MultLevel) * g.ComboMult() * g.HRMult()))
}

// HRMult is the permanent gain bonus from the "hrmult" meta upgrade.
func (g *Game) HRMult() float64 {
	return 1 + 0.1*float64(g.HRMultLevel)
}

// DayLength is how many active-play seconds a day lasts.
func (g *Game) DayLength() float64 {
	return DayBaseTime + DayTimeBonus*float64(g.HRDayLevel)
}

// DayQuota is what HR charges at the end of the current day: exponential
// in the day number, scaled up by owned run upgrades and trimmed by the
// "hrquota" meta upgrade.
func (g *Game) DayQuota() int {
	quota := float64(QuotaBase) * math.Pow(QuotaGrowth, float64(g.Day-1))
	quota *= 1 + 0.25*float64(g.MultLevel) + 0.10*float64(g.StreakCapLevel)
	quota *= 1 - 0.10*float64(g.HRQuotaLevel)
	return int(math.Round(quota))
}

// DayClock renders the remaining day time as M:SS.
func (g *Game) DayClock() string {
	t := int(math.Ceil(math.Max(0, g.DayTimer)))
	return fmt.Sprintf("%d:%02d", t/60, t%60)
}

// DayProgress is how much of the day has elapsed, 0..1, for the bottom
// progress bar.
func (g *Game) DayProgress() float64 {
	l := g.DayLength()
	if l <= 0 {
		return 0
	}
	return math.Min(1, math.Max(0, 1-g.DayTimer/l))
}

// TypeChar feeds one typed character into the game. In command mode the
// character goes to the command buffer. Otherwise, "/" opens command mode,
// a correct letter advances the cursor, completing the word scores
// WordGain() and enters PhaseSuccess, and a wrong letter enters PhaseError,
// resetting both the cursor and the streak. Input is ignored outside
// PhaseTyping.
func (g *Game) TypeChar(r rune) {
	if g.Scene != ScenePlay {
		return
	}
	if g.CommandMode {
		g.CommandBuf += string(r)
		return
	}
	if g.Phase != PhaseTyping {
		return
	}
	if r == '/' {
		g.CommandMode = true
		g.CommandBuf = "/"
		g.CommandErr = 0
		return
	}
	if g.Pos < len(g.Word) && rune(g.Word[g.Pos]) == r {
		g.Pos++
		if g.Pos == len(g.Word) {
			g.Streak++
			g.LastGain = g.WordGain()
			g.Score += g.LastGain
			g.DayWords++
			g.TotalWords++
			g.Phase = PhaseSuccess
			g.PhaseTimer = SuccessTime
		}
		return
	}
	g.Phase = PhaseError
	g.PhaseTimer = ErrorFlashTime + ErrorBlockTime
	g.Pos = 0
	g.Streak = 0
	g.DayFails++
	g.TotalFails++
	g.dayFailCounts[g.Word]++
	g.FailCounts[g.Word]++
}

// Tick advances time-based state; dt is in seconds. The intro reveals
// characters at IntroCPS; in play, command mode pauses the phases, and
// when a phase expires play returns to PhaseTyping (a finished success
// phase also brings the next word).
func (g *Game) Tick(dt float64) {
	if g.CommandErr > 0 {
		g.CommandErr = max(0, g.CommandErr-dt)
	}
	switch g.Scene {
	case SceneIntro:
		if !g.IntroLineDone() {
			g.IntroShown = math.Min(g.IntroShown+dt*IntroCPS, float64(len(g.introLine())))
		}
	case ScenePlay:
		if g.CommandMode {
			return // command mode pauses phases and the day clock
		}
		g.DayTimer -= dt
		if g.Phase != PhaseTyping {
			g.PhaseTimer -= dt
			if g.PhaseTimer <= 0 {
				if g.Phase == PhaseSuccess {
					g.nextWord()
					g.Save() // autosave with the completed word banked
				}
				g.Phase = PhaseTyping
				g.PhaseTimer = 0
			}
		}
		// only end the day between phases so a word mid-success still
		// banks its points first
		if g.DayTimer <= 0 && g.Phase == PhaseTyping {
			g.endDay()
		}
	case SceneDayEnd:
		if !g.DayEndQuipDone() {
			g.QuipShown = math.Min(g.QuipShown+dt*IntroCPS, float64(len(g.dayEndQuip())))
		}
	case SceneDayStart:
		if !g.DayStartDone() {
			g.QuipShown = math.Min(g.QuipShown+dt*IntroCPS, float64(len([]rune(g.DayStartText))))
		}
	}
}

// CommandBackspace deletes the last rune of the command buffer; deleting
// the leading "/" leaves command mode.
func (g *Game) CommandBackspace() {
	if !g.CommandMode {
		return
	}
	if len(g.CommandBuf) <= 1 {
		g.CommandCancel()
		return
	}
	runes := []rune(g.CommandBuf)
	g.CommandBuf = string(runes[:len(runes)-1])
}

// CommandCancel leaves command mode discarding the buffer.
func (g *Game) CommandCancel() {
	g.CommandMode = false
	g.CommandBuf = ""
}

// CommandSubmit executes the buffered command. Unknown commands show a
// brief error message.
func (g *Game) CommandSubmit() {
	cmd := strings.ToLower(strings.TrimSpace(g.CommandBuf))
	g.CommandCancel()
	switch cmd {
	case "/shop", "/tienda":
		g.ShopIndex = 0
		g.Scene = SceneShop
	case "/rrhh", "/hr":
		g.HRIndex = 0
		g.hrFrom = ScenePlay
		g.Scene = SceneHR
	case "/menu":
		g.Scene = SceneMenu
		g.Save()
	default:
		g.CommandErr = CommandErrTime
	}
}

// ShopMove moves the shop selection by delta, wrapping around.
func (g *Game) ShopMove(delta int) {
	n := len(Upgrades)
	g.ShopIndex = ((g.ShopIndex+delta)%n + n) % n
}

// ShopBuy purchases the selected upgrade if the score covers its cost.
func (g *Game) ShopBuy() bool {
	up := Upgrades[g.ShopIndex]
	cost := g.UpgradeCost(up.ID)
	if g.Score < cost {
		return false
	}
	g.Score -= cost
	switch up.ID {
	case "mult":
		g.MultLevel++
	case "streakcap":
		g.StreakCapLevel++
	}
	g.Save()
	return true
}

// ExitShop returns from the shop to play.
func (g *Game) ExitShop() {
	g.Scene = ScenePlay
}

// endDay settles the day with HR: covering the quota pays it, earns HR
// points equal to the day number and starts the next day; falling short
// gets the player fired and the run erased (HR points survive).
func (g *Game) endDay() {
	code := Languages[g.LangIndex].Code
	g.DayEndQuota = g.DayQuota()
	g.QuipShown = 0
	g.DayEndLine = 0

	// bank today's stats for the summary and tomorrow's morning quip
	g.LastDayWords = g.DayWords
	g.LastDayFails = g.DayFails
	g.LastDayWorst, g.LastDayWorstN = worstWord(g.dayFailCounts)

	if g.Score >= g.DayEndQuota {
		g.Score -= g.DayEndQuota
		g.DayEndGainHR = g.Day
		g.HRPoints += g.Day
		g.DayEndFirst = g.Day == 1
		g.Day++
		g.DayTimer = g.DayLength()
		g.resetDayStats()
		g.DayEndFired = false
		g.DayEndQuip = pickQuip(len(PaydayScripts[code]), &g.lastPayday)
	} else {
		g.DayEndFired = true
		g.DayEndFirst = false
		g.DayEndGainHR = 0
		g.DayEndQuip = pickQuip(len(FiredScripts[code]), &g.lastFired)
		g.RunActive = false
		g.savedRun = nil
	}
	g.Scene = SceneDayEnd
	g.Save()
}

// worstWord returns the most-failed word of a count map (ties broken
// alphabetically for determinism) or "" when empty.
func worstWord(counts map[string]int) (string, int) {
	word, n := "", 0
	for w, c := range counts {
		if c > n || (c == n && (word == "" || w < word)) {
			word, n = w, c
		}
	}
	return word, n
}

// dayEndLines is the summary script: normally a single quip, but the
// first payday explains HR points across several lines, intro-style.
func (g *Game) dayEndLines() []string {
	code := Languages[g.LangIndex].Code
	if g.DayEndFired {
		return []string{FiredScripts[code][g.DayEndQuip]}
	}
	if g.DayEndFirst {
		return HRExplainScripts[code]
	}
	return []string{PaydayScripts[code][g.DayEndQuip]}
}

// dayEndQuip is the current summary line, as runes.
func (g *Game) dayEndQuip() []rune {
	return []rune(g.dayEndLines()[g.DayEndLine])
}

// DayEndQuipVisible is the typewriter-revealed prefix of the summary quip.
func (g *Game) DayEndQuipVisible() string {
	line := g.dayEndQuip()
	return string(line[:min(int(g.QuipShown), len(line))])
}

// DayEndQuipDone reports whether the summary quip is fully revealed.
func (g *Game) DayEndQuipDone() bool {
	return int(g.QuipShown) >= len(g.dayEndQuip())
}

// DayEndSummary reports whether the numbers (stats + settlement) should
// show: the script's last line is fully revealed.
func (g *Game) DayEndSummary() bool {
	return g.DayEndQuipDone() && g.DayEndLine == len(g.dayEndLines())-1
}

// DayEndEnter advances the end-of-day summary: it completes a
// half-revealed line, then moves to the next one; after the last, a paid
// day moves on to the morning quip, a firing goes back to the main menu.
func (g *Game) DayEndEnter() {
	if !g.DayEndQuipDone() {
		g.QuipShown = float64(len(g.dayEndQuip()))
		return
	}
	if g.DayEndLine+1 < len(g.dayEndLines()) {
		g.DayEndLine++
		g.QuipShown = 0
		return
	}
	if g.DayEndFired {
		g.MenuIndex = 0
		g.Scene = SceneMenu
		return
	}
	g.composeDayStart()
	g.QuipShown = 0
	g.Scene = SceneDayStart
}

// composeDayStart picks the morning quip for the new day: half the time
// (when yesterday left stats worth mocking) a stat-aware line, otherwise a
// plain one. Stat templates use {words}, {fails}, {worst} and {worstn}.
func (g *Game) composeDayStart() {
	code := Languages[g.LangIndex].Code
	if g.LastDayWorst != "" && rand.IntN(2) == 0 {
		t := DayStartStatScripts[code][pickQuip(len(DayStartStatScripts[code]), &g.lastStartStat)]
		t = strings.ReplaceAll(t, "{words}", strconv.Itoa(g.LastDayWords))
		t = strings.ReplaceAll(t, "{fails}", strconv.Itoa(g.LastDayFails))
		t = strings.ReplaceAll(t, "{worst}", g.LastDayWorst)
		t = strings.ReplaceAll(t, "{worstn}", strconv.Itoa(g.LastDayWorstN))
		g.DayStartText = t
		return
	}
	g.DayStartText = DayStartScripts[code][pickQuip(len(DayStartScripts[code]), &g.lastStartPlain)]
}

// DayStartVisible is the typewriter-revealed prefix of the morning quip.
func (g *Game) DayStartVisible() string {
	line := []rune(g.DayStartText)
	return string(line[:min(int(g.QuipShown), len(line))])
}

// DayStartDone reports whether the morning quip is fully revealed.
func (g *Game) DayStartDone() bool {
	return int(g.QuipShown) >= len([]rune(g.DayStartText))
}

// DayStartEnter leaves the morning quip into the new day's play, with a
// fresh word (a half-revealed quip is completed first).
func (g *Game) DayStartEnter() {
	if !g.DayStartDone() {
		g.QuipShown = float64(len([]rune(g.DayStartText)))
		return
	}
	g.Word = ""
	g.nextWord()
	g.Phase = PhaseTyping
	g.PhaseTimer = 0
	g.Scene = ScenePlay
}

// FailStat is one entry of the most-failed-words ranking.
type FailStat struct {
	Word string
	N    int
}

// TopFailWords returns the n all-time most-failed words, most failed
// first (ties broken alphabetically).
func (g *Game) TopFailWords(n int) []FailStat {
	stats := make([]FailStat, 0, len(g.FailCounts))
	for w, c := range g.FailCounts {
		stats = append(stats, FailStat{Word: w, N: c})
	}
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].N != stats[j].N {
			return stats[i].N > stats[j].N
		}
		return stats[i].Word < stats[j].Word
	})
	if len(stats) > n {
		stats = stats[:n]
	}
	return stats
}

// StatsBack leaves the stats screen back to the main menu.
func (g *Game) StatsBack() {
	g.Scene = SceneMenu
}

// HRMove moves the meta-shop selection by delta, wrapping around.
func (g *Game) HRMove(delta int) {
	n := len(HRUpgrades)
	g.HRIndex = ((g.HRIndex+delta)%n + n) % n
}

// HRBuy purchases the selected meta upgrade with HR points; capped
// upgrades at their limit can't be bought.
func (g *Game) HRBuy() bool {
	up := HRUpgrades[g.HRIndex]
	if up.ID == "hrquota" && g.HRQuotaLevel >= HRQuotaMax {
		return false
	}
	cost := g.HRUpgradeCost(up.ID)
	if g.HRPoints < cost {
		return false
	}
	g.HRPoints -= cost
	switch up.ID {
	case "hrmult":
		g.HRMultLevel++
	case "hrday":
		g.HRDayLevel++
	case "hrquota":
		g.HRQuotaLevel++
	}
	g.Save()
	return true
}

// ExitHR leaves the meta shop, back to wherever it was opened from.
func (g *Game) ExitHR() {
	g.Scene = g.hrFrom
}

// HRUpgradeLevel returns the owned level of a meta upgrade.
func (g *Game) HRUpgradeLevel(id string) int {
	switch id {
	case "hrmult":
		return g.HRMultLevel
	case "hrday":
		return g.HRDayLevel
	case "hrquota":
		return g.HRQuotaLevel
	}
	return 0
}

// HRUpgradeCost is the price of the next meta level: BaseCost × 2^level.
func (g *Game) HRUpgradeCost(id string) int {
	cost := 0
	for _, up := range HRUpgrades {
		if up.ID == id {
			cost = up.BaseCost
		}
	}
	for range g.HRUpgradeLevel(id) {
		cost *= 2
	}
	return cost
}

// HRUpgradeMaxed reports whether a meta upgrade hit its level cap.
func (g *Game) HRUpgradeMaxed(id string) bool {
	return id == "hrquota" && g.HRQuotaLevel >= HRQuotaMax
}

// HRUpgradeEffect describes what buying the next meta level changes.
func (g *Game) HRUpgradeEffect(id string) string {
	switch id {
	case "hrmult":
		return fmt.Sprintf("+%d%% → +%d%%", 10*g.HRMultLevel, 10*(g.HRMultLevel+1))
	case "hrday":
		return fmt.Sprintf("%ds → %ds",
			int(DayBaseTime+DayTimeBonus*float64(g.HRDayLevel)),
			int(DayBaseTime+DayTimeBonus*float64(g.HRDayLevel+1)))
	case "hrquota":
		if g.HRQuotaLevel >= HRQuotaMax {
			return fmt.Sprintf("-%d%%", 10*g.HRQuotaLevel)
		}
		return fmt.Sprintf("-%d%% → -%d%%", 10*g.HRQuotaLevel, 10*(g.HRQuotaLevel+1))
	}
	return ""
}

// UpgradeLevel returns the owned level of an upgrade.
func (g *Game) UpgradeLevel(id string) int {
	switch id {
	case "mult":
		return g.MultLevel
	case "streakcap":
		return g.StreakCapLevel
	}
	return 0
}

// UpgradeCost is the price of the next level: BaseCost × 3^level.
func (g *Game) UpgradeCost(id string) int {
	cost := 0
	for _, up := range Upgrades {
		if up.ID == id {
			cost = up.BaseCost
		}
	}
	for range g.UpgradeLevel(id) {
		cost *= 3
	}
	return cost
}

// UpgradeEffect describes what buying the next level changes, e.g.
// "x1 → x2" for the point multiplier or "5 → 10" for the streak cap.
func (g *Game) UpgradeEffect(id string) string {
	switch id {
	case "mult":
		return fmt.Sprintf("x%d → x%d", 1+g.MultLevel, 2+g.MultLevel)
	case "streakcap":
		return fmt.Sprintf("%d → %d", g.StreakCap(), g.StreakCap()+BaseStreakCap)
	}
	return ""
}

// ErrorFlashing reports whether the pink error highlight is currently
// visible (the first second of PhaseError).
func (g *Game) ErrorFlashing() bool {
	return g.Phase == PhaseError && g.PhaseTimer > ErrorBlockTime
}

// nextWord picks a random word, avoiding an immediate repeat.
func (g *Game) nextWord() {
	next := g.Words[rand.IntN(len(g.Words))]
	for next == g.Word && len(g.Words) > 1 {
		next = g.Words[rand.IntN(len(g.Words))]
	}
	g.Word = next
	g.Pos = 0
}
