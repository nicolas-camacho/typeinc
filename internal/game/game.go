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
	SceneTerminal // hidden story terminal, reached via /terminal
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
// end HR charges DayQuota() — not covering it means getting fired. The
// quota is locked at day start: upgrades bought mid-day only inflate the
// NEXT day's quota.
const (
	DayBaseTime     = 60.0 // seconds per day before upgrades (1 min)
	DayTimeBonus    = 15.0 // extra seconds per "hrday" level
	QuotaBase       = 100  // day-1 quota before scaling
	QuotaGrowth     = 1.55 // per-day exponential growth
	QuotaMultTerm   = 0.25 // quota inflation per "mult" level
	QuotaStreakTerm = 0.08 // quota inflation per "streakcap" level
	HRQuotaMax      = 5    // "hrquota" level cap (−10% each, max −50%)
)

// HR payroll tuning: HR points are only paid on days multiple of
// HRPayCycle, and the payout is day/HRPayDivisor (integer division):
// day 3 → 1, day 6 → 3, day 9 → 4, day 12 → 6… Getting fired mid-cycle
// forfeits that cycle.
const (
	HRPayCycle   = 3 // a payroll lands every Nth day
	HRPayDivisor = 2 // payout on a payroll day = day / HRPayDivisor
)

// hrGainForDay is how many HR points closing the given day pays.
func hrGainForDay(day int) int {
	if day%HRPayCycle != 0 {
		return 0
	}
	return day / HRPayDivisor
}

// Golden-word tuning: each new player word has a chance of coming out
// golden, paying a multiplied gain. The "golden" shop upgrade raises both
// the chance and the multiplier, and inflates the quota in return.
const (
	GoldenBaseChance     = 0.02 // base chance the next player word is golden
	GoldenChancePerLevel = 0.01 // +1% chance per "golden" shop level
	GoldenBaseMult       = 3.0  // golden payout multiplier at level 0
	GoldenMultPerLevel   = 0.5  // +0.5x payout per level
	GoldenMaxLevel       = 8    // level cap: 10% chance, ×7 payout
	QuotaGoldenTerm      = 0.04 // quota inflation per golden level
)

// Intern tuning: an intern types its own words in the HUD's left column,
// banking WordGain without the combo or golden bonuses. One intern per run
// by default; the "hrintern" meta upgrade raises the per-run headcount.
const (
	InternBaseCost       = 1500 // first intern; each next one costs ×3
	InternBaseCPS        = 1.5  // letters per second at speed level 0
	InternCPSPerLevel    = 0.75 // extra letters/sec per "internspeed" level
	InternSpeedBaseCost  = 750  // "internspeed" base cost (×3 per level)
	InternWordRest       = 0.5  // pause after an intern finishes a word (HUD readability)
	InternGainFlash      = 0.8  // seconds the intern "+n" flash stays visible
	HRInternMax          = 3    // "hrintern" level cap: up to 4 interns per run
	QuotaInternTerm      = 0.15 // quota inflation per intern owned
	QuotaInternSpeedTerm = 0.04 // quota inflation per "internspeed" level
)

// Intern is one hired auto-typer: it types Word letter by letter and banks
// the gain on completion. Frontends only read the exported fields.
type Intern struct {
	Word      string
	Pos       int     // letters typed so far
	LastGain  int     // last payout, for the HUD flash
	GainTimer float64 // seconds the payout flash stays visible
	Possessed bool    // typing a story phrase: renders violet, pays nothing
	acc       float64 // fractional letters accumulated
	rest      float64 // post-word pause left
}

// Office-event tuning: during active play the game occasionally rolls an
// event — a coffee frenzy (everything pays double for a while) or an HR
// inspection (the next word is urgent: typed in time it pays a bonus,
// expired it just goes back to normal — never a punishment).
// EventMinSpacing > FrenzyDuration guarantees events never overlap.
const (
	EventCheckInterval   = 5.0  // seconds of active play between event rolls
	EventChance          = 0.20 // trigger chance per roll (~2.4 expected per 60s day)
	EventMinSpacing      = 20.0 // min active seconds between events
	EventMaxPerDay       = 5
	FrenzyDuration       = 10.0
	FrenzyMult           = 2.0 // ×2 on ALL income: player words and interns
	InspectionTime       = 8.0 // seconds to finish the urgent word
	InspectionRewardMult = 3.0 // urgent word completed in time pays ×3
)

// Upgrade describes one shop entry; the display name is localized via
// UpgradeName. Cost at level n is BaseCost × 3ⁿ.
type Upgrade struct {
	ID       string
	BaseCost int
}

// Upgrades lists the shop entries in display order.
var Upgrades = []Upgrade{
	{ID: "mult", BaseCost: 75},
	{ID: "streakcap", BaseCost: 60},
	{ID: "golden", BaseCost: 250},
	{ID: "intern", BaseCost: InternBaseCost},
	{ID: "internspeed", BaseCost: InternSpeedBaseCost},
}

// HRUpgrades lists the meta-shop entries, paid with HR points and kept
// across runs. Cost at level n is BaseCost × 2ⁿ.
var HRUpgrades = []Upgrade{
	{ID: "hrmult", BaseCost: 5},
	{ID: "hrday", BaseCost: 4},
	{ID: "hrquota", BaseCost: 6},
	{ID: "hrintern", BaseCost: 25},
	{ID: "hrend", BaseCost: 40},
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
	WordGolden bool     // current word is golden: it pays GoldenMult()
	Score      int
	Phase      Phase
	PhaseTimer float64 // seconds left in the current non-typing phase
	LastGain   int     // points awarded for the last completed word

	// incremental state
	Streak         int // consecutive completed words without a mistake
	MultLevel      int // "mult" upgrade level
	StreakCapLevel int // "streakcap" upgrade level
	GoldenLevel    int // "golden" upgrade level

	// interns (part of the run); words/progress are not persisted
	Interns          []Intern
	InternSpeedLevel int    // "internspeed" upgrade level
	ShopQuipText     string // HR quip shown in the shop after hiring
	lastHired        int    // previous hiring quip, to never repeat

	// office events; only EventsToday persists — active timers reset on
	// load, like the word itself
	FrenzyTimer     float64 // >0: everything pays ×FrenzyMult
	InspectionArmed bool    // inspection waiting for the next word boundary
	InspectionTimer float64 // >0: the current word is urgent
	EventsToday     int
	EventToast      string  // announce line, visible for the whole event
	eventClock      float64 // active seconds since the last event roll
	eventSpacing    float64 // active seconds since the last event
	lastCoffeeQuip  int
	lastInspectQuip int

	// day cycle (part of the run)
	Day      int     // current day, 1-based
	DayTimer float64 // active-play seconds left in the day
	dayQuota int     // quota locked at day start; upgrades hit the NEXT day

	// current-day stats (part of the run)
	DayWords      int            // words completed today
	DayFails      int            // mistakes today
	DayBestStreak int            // longest streak today
	dayFailCounts map[string]int // per-word mistakes today

	// stats of the last settled day, for the summary and morning quips
	LastDayWords      int
	LastDayFails      int
	LastDayWorst      string // most-failed word ("" if no fails)
	LastDayWorstN     int
	LastDayBestStreak int

	// global stats, kept across runs
	TotalWords int
	TotalFails int
	BestStreak int            // longest streak ever
	BestDay    int            // highest day ever reached
	FailCounts map[string]int // per-word mistakes, all-time

	// meta currency, kept across runs
	HRPoints      int
	HRMultLevel   int // "hrmult": +10% gain per level, permanent
	HRDayLevel    int // "hrday": +15s day length per level
	HRQuotaLevel  int // "hrquota": −10% quota per level, capped
	HRInternLevel int // "hrintern": +1 intern per run per level, capped
	HREndLevel    int // "hrend": unlocks the /terminar command (0 or 1)

	// end-of-day summary state
	DayEndFired   bool    // this summary is a firing, not a payday
	DayEndFirst   bool    // first payroll (day HRPayCycle): HR points get explained
	DayEndQuota   int     // quota HR charged (or demanded) this day
	DayEndGainHR  int     // HR points earned (0 off payroll days or when fired)
	DayEndNextPay int     // next payroll day (0 when this day paid out)
	DayEndQuip    int     // picked quip index
	DayEndLine    int     // current line of a multi-line summary script
	QuipShown     float64 // typewriter progress of the quip, in runes
	lastPayday    int     // previous payday quip, to never repeat
	lastFired     int     // previous firing quip, to never repeat

	// day-start morning quip state
	DayStartText   string // composed quip (may embed yesterday's stats)
	lastStartPlain int    // previous plain morning quip
	lastStartStat  int    // previous stat-based morning quip

	HRIndex    int    // selected meta-shop entry
	hrFrom     Scene  // scene to return to when leaving the meta shop
	HRQuipText string // HR quip shown in the meta shop after a purchase

	lastBought map[string]int // previous purchase quip per upgrade, to never repeat
	resetArmed bool           // RESET selected once, waiting for the confirm Enter

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

	// background mystery — meta layer, survives firings (like HRPoints)
	StoryStage   int             // 0 dormant, 1..6 acts, StoryEndAct solved
	storyClues   map[string]bool // discovered clue IDs
	storyDocs    map[string]bool // keyed terminal docs unlocked
	markedCounts map[string]int  // marked-word progress, accumulates across runs

	// background mystery — run layer, reset by initRun
	WordCorrupt     bool            // current word is corrupt: violet, unpaid, clock frozen
	CorruptTimer    float64         // seconds before the corrupt word dissolves
	corruptArmed    bool            // corrupt word waiting for the next word boundary
	corruptClueID   string          // clue the armed/active corrupt word registers
	CorruptToday    int             // per-day corrupt cap counter
	ritualDone      map[string]bool // triggers already fired this run
	storyDaysNoBeat int             // mornings without a story beat (pity)
	possessedIdx    int             // intern currently possessed (-1 none)
	possessedToday  bool

	// background mystery — presentation, not persisted
	ClockGlitch       float64 // >0: DayClock shows an impossible time
	clockGlitchClock  float64
	clockGlitchText   string
	StoryToast        string  // violet toast (corrupt/marked/possession payoffs)
	StoryToastTimer   float64 // seconds the story toast stays visible
	DayStartStory     bool    // this morning's quip is a glitched one
	DayEndStory       bool    // this evening's quip is a glitched one
	dayEndStoryIdx    int
	DayEndCode        string // act-5 code embedded in the day-end summary
	storyQuipForced   bool   // act just advanced: next morning quip glitches
	lastWeird         int
	lastHint          int   // no-repeat cursor for objective-hint quips
	termFrom          Scene // scene to return to when leaving the terminal
	termOpened        bool  // meta: terminal entered at least once (persisted)
	lastCorruptDone   int
	lastMarkedHit     int
	lastPossessed     int
	lastPossessedDone int

	// hidden terminal scene
	TermBuf string   // prompt buffer
	TermLog []string // session log, capped
}

// New returns a game sitting on the main menu.
func New() *Game {
	g := &Game{Scene: SceneMenu, Volume: 0.4, FailCounts: map[string]int{}}
	g.storyClues = map[string]bool{}
	g.storyDocs = map[string]bool{}
	g.markedCounts = map[string]int{}
	g.resetStoryRun()
	return g
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
	if g.TerminalAvailable() {
		items = append(items, MenuItem{ID: "terminal", Label: ui.MenuTerminal})
	}
	if g.OptionsEnabled {
		items = append(items, MenuItem{ID: "options", Label: ui.MenuOptions})
	}
	resetLabel := ui.MenuReset
	if g.resetArmed {
		resetLabel = ui.ResetConfirm
	}
	items = append(items, MenuItem{ID: "reset", Label: resetLabel})
	return append(items, MenuItem{ID: "quit", Label: ui.MenuQuit})
}

// MenuMove moves the menu selection by delta, wrapping around. Moving
// away from an armed RESET disarms it.
func (g *Game) MenuMove(delta int) {
	g.resetArmed = false
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

// MenuSelect activates the highlighted menu entry. RESET needs a second
// Enter to confirm; selecting anything else disarms it.
func (g *Game) MenuSelect() {
	id := g.menuItem().ID
	if id != "reset" {
		g.resetArmed = false
	}
	switch id {
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
	case "terminal":
		g.EnterTerminal()
	case "options":
		g.OptIndex = 0
		g.Scene = SceneOptions
	case "reset":
		if !g.resetArmed {
			g.resetArmed = true
			return
		}
		g.resetArmed = false
		g.resetSave()
	case "quit":
		g.QuitRequested = true
		g.Save()
	}
}

// resetSave wipes every bit of progress — run, HR points and levels,
// global stats — and overwrites the save file. Settings (language, volume,
// display) survive.
func (g *Game) resetSave() {
	g.RunActive = false
	g.savedRun = nil
	g.HRPoints = 0
	g.HRMultLevel = 0
	g.HRDayLevel = 0
	g.HRQuotaLevel = 0
	g.HRInternLevel = 0
	g.HREndLevel = 0
	g.TotalWords = 0
	g.TotalFails = 0
	g.BestStreak = 0
	g.BestDay = 0
	g.FailCounts = map[string]int{}
	g.resetStory()
	g.MenuIndex = 0
	g.Save()
}

// continueRun resumes the in-memory run (or restores the one from disk)
// and greets the player with a short random HR quip before play.
func (g *Game) continueRun() {
	if !g.RunActive && g.savedRun != nil {
		g.initRun()
		g.Score = g.savedRun.Score
		g.MultLevel = g.savedRun.MultLevel
		g.StreakCapLevel = g.savedRun.StreakCapLevel
		g.GoldenLevel = g.savedRun.GoldenLevel
		g.InternSpeedLevel = g.savedRun.InternSpeedLevel
		// intern words/progress aren't persisted: rebuild them fresh, like
		// the player word
		g.Interns = nil
		for range g.savedRun.InternCount {
			g.Interns = append(g.Interns, Intern{Word: g.pickWord("")})
		}
		g.EventsToday = g.savedRun.EventsToday
		g.CorruptToday = g.savedRun.CorruptToday
		g.ritualDone = sliceToBoolSet(g.savedRun.RitualDone)
		g.storyDaysNoBeat = g.savedRun.StoryDaysNoBeat
		g.Streak = g.savedRun.Streak
		g.Day = g.savedRun.Day
		g.DayTimer = g.savedRun.DayTimer
		g.dayQuota = g.savedRun.DayQuota
		g.DayWords = g.savedRun.DayWords
		g.DayFails = g.savedRun.DayFails
		g.DayBestStreak = g.savedRun.DayBestStreak
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
		// older saves have no locked quota
		if g.dayQuota <= 0 {
			g.dayQuota = g.computeQuota()
		}
		g.BestDay = max(g.BestDay, g.Day)
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
	g.GoldenLevel = 0
	g.Interns = nil
	g.InternSpeedLevel = 0
	g.ShopQuipText = ""
	g.CommandMode = false
	g.CommandBuf = ""
	g.CommandErr = 0
	g.ShopIndex = 0
	g.resetStoryRun()
	g.Word = ""
	g.nextWord()
	g.Day = 1
	g.DayTimer = g.DayLength()
	g.dayQuota = g.computeQuota()
	g.BestDay = max(g.BestDay, g.Day)
	g.resetDayStats()
	g.clearEvents()
	g.RunActive = true
}

// resetDayStats clears the current-day counters.
func (g *Game) resetDayStats() {
	g.DayWords = 0
	g.DayFails = 0
	g.DayBestStreak = 0
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

// IntroSkip jumps past the whole intro straight into play — the full
// new-game story and the CONTINUAR return quip alike.
func (g *Game) IntroSkip() {
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
// len × (1 + mult level) × combo multiplier × HR bonus, times the golden,
// frenzy and inspection multipliers when they apply, rounded.
func (g *Game) WordGain() int {
	gain := float64(len(g.Word)) * float64(1+g.MultLevel) * g.ComboMult() * g.HRMult()
	if g.WordGolden {
		gain *= g.GoldenMult()
	}
	if g.FrenzyActive() {
		gain *= FrenzyMult
	}
	if g.InspectionActive() {
		gain *= InspectionRewardMult
	}
	return int(math.Round(gain))
}

// GoldenChance is the probability of the next player word being golden.
func (g *Game) GoldenChance() float64 {
	return GoldenBaseChance + GoldenChancePerLevel*float64(g.GoldenLevel)
}

// GoldenMult is what a golden word multiplies its gain by.
func (g *Game) GoldenMult() float64 {
	return GoldenBaseMult + GoldenMultPerLevel*float64(g.GoldenLevel)
}

// rollGolden decides whether a freshly picked word comes out golden.
func (g *Game) rollGolden() bool {
	return rand.Float64() < g.GoldenChance()
}

// tickEvents advances the event timers and rolls for a new event every
// EventCheckInterval. It runs exactly when the day clock runs, so paused
// scenes freeze events too.
func (g *Game) tickEvents(dt float64) {
	g.FrenzyTimer = math.Max(0, g.FrenzyTimer-dt)
	g.InspectionTimer = math.Max(0, g.InspectionTimer-dt)
	g.eventSpacing += dt
	g.eventClock += dt
	for g.eventClock >= EventCheckInterval {
		g.eventClock -= EventCheckInterval
		g.maybeTriggerEvent()
	}
}

// maybeTriggerEvent rolls for a new event, skipping when the daily cap is
// hit, the last event is too recent, or one is still active.
func (g *Game) maybeTriggerEvent() {
	if g.EventsToday >= EventMaxPerDay || g.eventSpacing < EventMinSpacing {
		return
	}
	if g.FrenzyActive() || g.InspectionArmed || g.InspectionActive() {
		return
	}
	if g.corruptArmed || g.WordCorrupt {
		return // the corrupt signal owns the word slot
	}
	if rand.Float64() >= EventChance {
		return
	}
	if rand.IntN(2) == 0 {
		g.startFrenzy()
	} else {
		g.armInspection()
	}
}

// startFrenzy begins a coffee frenzy: everything pays ×FrenzyMult while
// the timer runs.
func (g *Game) startFrenzy() {
	g.FrenzyTimer = FrenzyDuration
	g.EventsToday++
	g.eventSpacing = 0
	code := Languages[g.LangIndex].Code
	g.EventToast = EventCoffeeScripts[code][pickQuip(len(EventCoffeeScripts[code]), &g.lastCoffeeQuip)]
}

// armInspection queues an HR inspection: the NEXT word becomes urgent (a
// word in progress is never interrupted).
func (g *Game) armInspection() {
	g.InspectionArmed = true
	g.EventsToday++
	g.eventSpacing = 0
	code := Languages[g.LangIndex].Code
	g.EventToast = EventInspectionScripts[code][pickQuip(len(EventInspectionScripts[code]), &g.lastInspectQuip)]
}

// clearEvents drops all event state; a new day starts calm.
func (g *Game) clearEvents() {
	g.FrenzyTimer = 0
	g.InspectionArmed = false
	g.InspectionTimer = 0
	g.EventToast = ""
	g.eventClock = 0
	g.eventSpacing = 0
	g.EventsToday = 0
}

// FrenzyActive reports whether a coffee frenzy is running.
func (g *Game) FrenzyActive() bool { return g.FrenzyTimer > 0 }

// FrenzyLeft is the frenzy's remaining seconds.
func (g *Game) FrenzyLeft() float64 { return g.FrenzyTimer }

// InspectionActive reports whether the current word is urgent.
func (g *Game) InspectionActive() bool { return g.InspectionTimer > 0 }

// InspectionLeft is the urgent word's remaining seconds.
func (g *Game) InspectionLeft() float64 { return g.InspectionTimer }

// EventToastVisible reports whether the event announce line should show:
// for the whole life of the event, not on a timer.
func (g *Game) EventToastVisible() bool {
	return g.FrenzyActive() || g.InspectionArmed || g.InspectionActive()
}

// MaxInterns is how many interns a run may hire, raised by "hrintern".
func (g *Game) MaxInterns() int {
	return 1 + g.HRInternLevel
}

// InternCPS is how many letters per second each intern types.
func (g *Game) InternCPS() float64 {
	return InternBaseCPS + InternCPSPerLevel*float64(g.InternSpeedLevel)
}

// InternGain is what an intern banks for completing word: like WordGain
// but without the combo and golden bonuses — those stay player-skill. A
// coffee frenzy is company-wide, so it doubles interns too.
func (g *Game) InternGain(word string) int {
	gain := float64(len(word)) * float64(1+g.MultLevel) * g.HRMult()
	if g.FrenzyActive() {
		gain *= FrenzyMult
	}
	return int(math.Round(gain))
}

// tickInterns advances every intern's typing. It runs exactly when the day
// clock runs, so pausing the day (command mode, shop, menus) freezes the
// interns too — idling in a paused scene must never print free money.
func (g *Game) tickInterns(dt float64) {
	for i := range g.Interns {
		in := &g.Interns[i]
		in.GainTimer = math.Max(0, in.GainTimer-dt)
		left := dt
		if in.rest > 0 {
			if in.rest >= left {
				in.rest -= left
				continue
			}
			left -= in.rest
			in.rest = 0
		}
		in.acc += left * g.InternCPS()
		for in.acc >= 1 && in.rest == 0 {
			in.acc--
			in.Pos++
			if in.Pos >= len(in.Word) {
				if in.Possessed {
					// the phrase pays nothing; the intern snaps out of it
					g.possessionComplete(in)
				} else {
					in.LastGain = g.InternGain(in.Word)
					in.GainTimer = InternGainFlash
					g.Score += in.LastGain
					g.DayWords++
					g.TotalWords++
				}
				in.Word = g.pickWord(in.Word)
				in.Pos = 0
				in.rest = InternWordRest
				in.acc = 0
				g.maybePossessIntern(i)
			}
		}
	}
}

// HRMult is the permanent gain bonus from the "hrmult" meta upgrade.
func (g *Game) HRMult() float64 {
	return 1 + 0.1*float64(g.HRMultLevel)
}

// DayLength is how many active-play seconds a day lasts.
func (g *Game) DayLength() float64 {
	return DayBaseTime + DayTimeBonus*float64(g.HRDayLevel)
}

// DayQuota is what HR charges at the end of the current day. It's locked
// at day start (see computeQuota): buying upgrades mid-day only inflates
// the next day's quota.
func (g *Game) DayQuota() int {
	return g.dayQuota
}

// computeQuota evaluates the quota formula for the current day and levels:
// exponential in the day number, scaled up by owned run upgrades and
// trimmed by the "hrquota" meta upgrade. Called only at day boundaries.
func (g *Game) computeQuota() int {
	quota := float64(QuotaBase) * math.Pow(QuotaGrowth, float64(g.Day-1))
	quota *= 1 + QuotaMultTerm*float64(g.MultLevel) + QuotaStreakTerm*float64(g.StreakCapLevel) +
		QuotaGoldenTerm*float64(g.GoldenLevel) +
		QuotaInternTerm*float64(len(g.Interns)) +
		QuotaInternSpeedTerm*float64(g.InternSpeedLevel)
	quota *= 1 - 0.10*float64(g.HRQuotaLevel)
	return int(math.Round(quota))
}

// DayClock renders the remaining day time as M:SS — unless the building
// briefly shows the local time of somewhere else (act 4+).
func (g *Game) DayClock() string {
	if g.ClockGlitch > 0 && g.clockGlitchText != "" {
		return g.clockGlitchText
	}
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
			if g.WordCorrupt {
				// corrupt words pay nothing and touch no stats: the
				// reward is the clue
				g.corruptComplete()
				return
			}
			g.Streak++
			g.DayBestStreak = max(g.DayBestStreak, g.Streak)
			g.BestStreak = max(g.BestStreak, g.Streak)
			g.LastGain = g.WordGain()
			g.Score += g.LastGain
			g.DayWords++
			g.TotalWords++
			g.InspectionTimer = 0 // urgent word delivered: bonus banked above
			g.Phase = PhaseSuccess
			g.PhaseTimer = SuccessTime
			g.checkMarkedWord()
			g.checkStoryTriggers()
			if g.StoryActive() && g.rollCorrupt() {
				g.armCorrupt("")
			}
		}
		return
	}
	g.Phase = PhaseError
	g.PhaseTimer = ErrorFlashTime + ErrorBlockTime
	g.Pos = 0
	if g.WordCorrupt {
		// a corrupt word is off the company's books: the streak survives
		// and "ayudame" must never pollute the fail stats
		return
	}
	g.Streak = 0
	g.DayFails++
	g.TotalFails++
	g.dayFailCounts[g.Word]++
	g.FailCounts[g.Word]++
	g.checkStoryTriggers()
}

// Tick advances time-based state; dt is in seconds. The intro reveals
// characters at IntroCPS; in play, command mode pauses the phases, and
// when a phase expires play returns to PhaseTyping (a finished success
// phase also brings the next word).
func (g *Game) Tick(dt float64) {
	if g.CommandErr > 0 {
		g.CommandErr = max(0, g.CommandErr-dt)
	}
	if g.StoryToastTimer > 0 {
		g.StoryToastTimer = max(0, g.StoryToastTimer-dt)
	}
	switch g.Scene {
	case SceneIntro:
		if !g.IntroLineDone() {
			g.IntroShown = math.Min(g.IntroShown+dt*IntroCPS, float64(len(g.introLine())))
		}
	case ScenePlay:
		if g.CommandMode {
			return // command mode pauses phases, interns, events and the day clock
		}
		if g.WordCorrupt {
			// a corrupt word freezes the whole office — day clock, events
			// and interns — so the interruption can never cost or farm
			// anything; only the signal's own timer runs
			if g.Phase == PhaseTyping {
				g.CorruptTimer -= dt
				if g.CorruptTimer <= 0 {
					g.WordCorrupt = false
					g.nextWord()
				}
			}
		} else {
			g.DayTimer -= dt
			g.tickEvents(dt)
			g.tickInterns(dt)
			g.tickStoryFX(dt)
		}
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
		// banks its points first; an active corrupt word plays out first
		if g.DayTimer <= 0 && g.Phase == PhaseTyping && !g.WordCorrupt {
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
	case "/terminar", "/endday":
		if g.HREndLevel > 0 {
			g.endDay() // settles the quota now: short score means fired
		} else {
			g.CommandErr = CommandErrTime
		}
	case "/terminal":
		// the old system only answers once the story is awake and the
		// player holds at least one clue; before that it stays a secret
		if g.TerminalAvailable() {
			g.EnterTerminal()
		} else {
			g.CommandErr = CommandErrTime
		}
	default:
		g.CommandErr = CommandErrTime
	}
}

// ShopItems is what the shop actually lists: every upgrade, except that
// the intern speed entry stays hidden until an intern exists.
func (g *Game) ShopItems() []Upgrade {
	if len(g.Interns) > 0 {
		return Upgrades
	}
	items := make([]Upgrade, 0, len(Upgrades)-1)
	for _, up := range Upgrades {
		if up.ID != "internspeed" {
			items = append(items, up)
		}
	}
	return items
}

// ShopMove moves the shop selection by delta, wrapping around.
func (g *Game) ShopMove(delta int) {
	n := len(g.ShopItems())
	g.ShopIndex = ((g.ShopIndex+delta)%n + n) % n
}

// ShopBuy purchases the selected upgrade if the score covers its cost;
// capped upgrades at their limit can't be bought. Hiring an intern also
// picks an HR quip for the shop screen.
func (g *Game) ShopBuy() bool {
	up := g.ShopItems()[g.ShopIndex]
	if g.UpgradeMaxed(up.ID) {
		return false
	}
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
	case "golden":
		g.GoldenLevel++
	case "intern":
		g.Interns = append(g.Interns, Intern{Word: g.pickWord("")})
	case "internspeed":
		g.InternSpeedLevel++
	}
	if up.ID == "intern" {
		code := Languages[g.LangIndex].Code
		g.ShopQuipText = InternHiredScripts[code][pickQuip(len(InternHiredScripts[code]), &g.lastHired)]
	} else {
		g.ShopQuipText = g.pickBoughtQuip(up.ID)
	}
	g.Save()
	return true
}

// ExitShop returns from the shop to play, dropping any hiring quip.
func (g *Game) ExitShop() {
	g.ShopQuipText = ""
	g.Scene = ScenePlay
}

// endDay settles the day with HR: covering the quota pays it, earns HR
// points equal to the day number and starts the next day; falling short
// gets the player fired and the run erased (HR points survive).
func (g *Game) endDay() {
	g.checkStoryTriggers() // day-boundary rituals see today's final numbers
	code := Languages[g.LangIndex].Code
	g.DayEndQuota = g.DayQuota()
	g.QuipShown = 0
	g.DayEndLine = 0

	// bank today's stats for the summary and tomorrow's morning quip
	g.LastDayWords = g.DayWords
	g.LastDayFails = g.DayFails
	g.LastDayBestStreak = g.DayBestStreak
	g.LastDayWorst, g.LastDayWorstN = worstWord(g.dayFailCounts)
	g.clearEvents()

	if g.Score >= g.DayEndQuota {
		g.Score -= g.DayEndQuota
		g.DayEndGainHR = hrGainForDay(g.Day)
		g.HRPoints += g.DayEndGainHR
		g.DayEndFirst = g.Day == HRPayCycle // first payroll: explain HR points
		g.DayEndNextPay = 0
		if g.DayEndGainHR == 0 {
			g.DayEndNextPay = (g.Day/HRPayCycle + 1) * HRPayCycle
		}
		g.Day++
		g.BestDay = max(g.BestDay, g.Day)
		g.maybeAdvanceStory() // BestDay gates re-check at every day boundary
		g.DayTimer = g.DayLength()
		g.dayQuota = g.computeQuota() // lock tomorrow's quota with today's upgrades
		g.resetDayStats()
		g.CorruptToday = 0
		g.possessedToday = false
		g.DayEndFired = false
		g.DayEndQuip = pickQuip(len(PaydayScripts[code]), &g.lastPayday)
		g.pickDayEndStory()
		g.DayEndCode = ""
		if g.StoryStage == 5 && !g.storyDocs["acceso"] {
			g.DayEndCode = DayEndCodeText // the payroll anomaly shows its digits
		}
	} else {
		g.DayEndFired = true
		g.DayEndFirst = false
		g.DayEndGainHR = 0
		g.DayEndNextPay = 0
		g.DayEndStory = false
		g.DayEndCode = ""
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
	if g.DayEndStory {
		act := min(g.StoryStage, StoryEndAct)
		if pool := StoryQuipScripts[act]; pool != nil && g.dayEndStoryIdx < len(pool[code]) {
			return []string{pool[code][g.dayEndStoryIdx]}
		}
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

// DayEndSkip jumps to the summary's last line fully revealed — numbers on
// screen at once. The first-payroll HR explanation is the long one; ESC
// skips it like the intro.
func (g *Game) DayEndSkip() {
	g.DayEndLine = len(g.dayEndLines()) - 1
	g.QuipShown = float64(len(g.dayEndQuip()))
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
	g.DayStartStory = false
	if q := g.storyDayBeat(); q != "" {
		g.DayStartText = q
		g.DayStartStory = true
		return
	}
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
	if g.HRUpgradeMaxed(up.ID) {
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
	case "hrintern":
		g.HRInternLevel++
	case "hrend":
		g.HREndLevel++
	}
	g.HRQuipText = g.pickBoughtQuip(up.ID)
	g.Save()
	return true
}

// pickBoughtQuip returns an HR quip for buying the given upgrade, never
// repeating the previous pick for that same upgrade.
func (g *Game) pickBoughtQuip(id string) string {
	code := Languages[g.LangIndex].Code
	pool := UpgradeBoughtScripts[id][code]
	if len(pool) == 0 {
		return ""
	}
	if g.lastBought == nil {
		g.lastBought = map[string]int{}
	}
	last := g.lastBought[id]
	quip := pool[pickQuip(len(pool), &last)]
	g.lastBought[id] = last
	return quip
}

// ExitHR leaves the meta shop, back to wherever it was opened from,
// dropping any purchase quip.
func (g *Game) ExitHR() {
	g.HRQuipText = ""
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
	case "hrintern":
		return g.HRInternLevel
	case "hrend":
		return g.HREndLevel
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
	switch id {
	case "hrquota":
		return g.HRQuotaLevel >= HRQuotaMax
	case "hrintern":
		return g.HRInternLevel >= HRInternMax
	case "hrend":
		return g.HREndLevel >= 1
	}
	return false
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
	case "hrintern":
		if g.HRInternLevel >= HRInternMax {
			return fmt.Sprintf("%d", g.MaxInterns())
		}
		return fmt.Sprintf("%d → %d", g.MaxInterns(), g.MaxInterns()+1)
	case "hrend":
		if Languages[g.LangIndex].Code == "es" {
			return "/terminar"
		}
		return "/endday"
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
	case "golden":
		return g.GoldenLevel
	case "intern":
		return len(g.Interns)
	case "internspeed":
		return g.InternSpeedLevel
	}
	return 0
}

// UpgradeMaxed reports whether a shop upgrade hit its level cap.
func (g *Game) UpgradeMaxed(id string) bool {
	switch id {
	case "golden":
		return g.GoldenLevel >= GoldenMaxLevel
	case "intern":
		return len(g.Interns) >= g.MaxInterns()
	}
	return false
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
	case "golden":
		if g.GoldenLevel >= GoldenMaxLevel {
			return fmt.Sprintf("%.0f%% x%.1f", 100*g.GoldenChance(), g.GoldenMult())
		}
		return fmt.Sprintf("%.0f%% x%.1f → %.0f%% x%.1f",
			100*g.GoldenChance(), g.GoldenMult(),
			100*(g.GoldenChance()+GoldenChancePerLevel), g.GoldenMult()+GoldenMultPerLevel)
	case "intern":
		if len(g.Interns) >= g.MaxInterns() {
			return fmt.Sprintf("%d/%d", len(g.Interns), g.MaxInterns())
		}
		return fmt.Sprintf("%d/%d → %d/%d",
			len(g.Interns), g.MaxInterns(), len(g.Interns)+1, g.MaxInterns())
	case "internspeed":
		return fmt.Sprintf("%.1f → %.1f l/s", g.InternCPS(), g.InternCPS()+InternCPSPerLevel)
	}
	return ""
}

// ErrorFlashing reports whether the pink error highlight is currently
// visible (the first second of PhaseError).
func (g *Game) ErrorFlashing() bool {
	return g.Phase == PhaseError && g.PhaseTimer > ErrorBlockTime
}

// nextWord picks a random word for the player, avoiding an immediate
// repeat, and rolls whether it comes out golden. An armed inspection
// claims the word here — at the word boundary — making it urgent instead
// of golden.
func (g *Game) nextWord() {
	g.WordCorrupt = false
	if g.corruptArmed {
		// the corrupt signal claims the boundary first; an armed
		// inspection stays armed for the following word
		g.corruptArmed = false
		g.claimCorrupt()
		if g.WordCorrupt {
			return
		}
	}
	g.Word = g.pickWord(g.Word)
	g.Pos = 0
	if g.InspectionArmed {
		g.InspectionArmed = false
		g.InspectionTimer = InspectionTime
		g.WordGolden = false
		return
	}
	g.WordGolden = g.rollGolden()
}

// pickWord returns a random word from the active list, avoiding an
// immediate repeat of avoid.
func (g *Game) pickWord(avoid string) string {
	next := g.Words[rand.IntN(len(g.Words))]
	for next == avoid && len(g.Words) > 1 {
		next = g.Words[rand.IntN(len(g.Words))]
	}
	return next
}
