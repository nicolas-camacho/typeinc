package game

import (
	"math"
	"math/rand/v2"
	"sort"
	"strings"
)

// Story tuning. The mystery is gated by BestDay (the best day ever
// reached, across runs) so a firing never stalls it, and by keys the
// player finds: both must be met to advance an act.
const (
	StoryMinDay      = 3    // BestDay gate of act 1; acts 2..6 in storyDayGates
	StoryEndAct      = 7    // StoryStage value once the mystery is solved
	CorruptChance    = 0.04 // corrupt-word roll per completed player word
	CorruptMaxPerDay = 1    // frozen-clock interruptions stay rare
	CorruptTimeout   = 15.0 // seconds before an unanswered corrupt word dissolves
	StoryPityDays    = 3    // days without a story beat before one is forced
	PossessChance    = 0.05 // possession roll per intern word boundary (act 4+)
	ClockGlitchTime  = 1.5  // seconds an impossible clock reading stays
	ClockGlitchGap   = 10.0 // seconds between clock-glitch rolls (act 4+)
	ClockGlitchOdds  = 0.06
	StoryToastTime   = 4.0 // seconds a story toast stays visible
	DayEndCodeText   = "4413"
	StoryHintBoost   = 0.15 // extra quip odds while the act's key is unused
	StoryHintOdds    = 0.75 // chance the picked quip is an objective hint then
)

// storyDayGates is the BestDay requirement per act (index = act).
var storyDayGates = [7]int{0, StoryMinDay, 7, 11, 15, 20, 25}

// terminalQuipLine is the act-1 pool index of the line that reveals
// /terminal. Until the player opens the terminal for the first time it is
// the ONLY story quip shown — nobody may miss the door in.
const terminalQuipLine = 4

// storyHintLines marks, per act, the StoryQuipScripts indices that point
// at the act's pending objective (the terminal, the key, the code). They
// get picked with priority while that objective is still unmet.
var storyHintLines = map[int][]int{
	1: {4, 5, 7}, // /terminal, Vera's notes, the name starts with V
	2: {3, 4},    // shift 77 as one word, the old system's key
	3: {2, 4},    // the heavy words, "sotano" opens things
	4: {1, 6},    // let the intern finish, type the purple signals
	5: {3, 4},    // four-four-one-three, the purple digits on the report
	6: {0, 1},    // name+number all together, you already know both
}

// StoryTrigger is one ritual the building watches for. Check reads game
// state only — no RNG, no mutation — so tests drive it deterministically.
// A trigger is live while StoryStage >= Act and fires once per run.
type StoryTrigger struct {
	ID     string
	Act    int
	Check  func(g *Game) bool
	Effect func(g *Game)
}

// StoryTriggers lists the rituals. The glitched quips hint at them.
var StoryTriggers = []StoryTrigger{
	{ID: "fail_even", Act: 1,
		Check:  func(g *Game) bool { return g.Day%2 == 0 && g.DayFails >= 5 },
		Effect: func(g *Game) { g.armCorrupt("cw_help") }},
	{ID: "repeat_fail", Act: 1,
		Check: func(g *Game) bool {
			for _, n := range g.dayFailCounts {
				if n >= 3 {
					return true
				}
			}
			return false
		},
		Effect: func(g *Game) { g.armCorrupt("cw_name") }},
	{ID: "perfect_day", Act: 2,
		Check:  func(g *Game) bool { return g.DayFails == 0 && g.DayWords >= 15 },
		Effect: func(g *Game) { g.armCorrupt("cw_notme") }},
	{ID: "streak_odd", Act: 3,
		Check:  func(g *Game) bool { return g.Day%2 == 1 && g.Streak >= 15 },
		Effect: func(g *Game) { g.armCorrupt("cw_listen") }},
	{ID: "hours_audit", Act: 4,
		Check:  func(g *Game) bool { return g.Day%2 == 0 && g.DayFails >= 3 },
		Effect: func(g *Game) { g.unlockDoc("horas") }},
	{ID: "deep_streak", Act: 5,
		Check:  func(g *Game) bool { return g.Day%5 == 0 && g.DayBestStreak >= 20 },
		Effect: func(g *Game) { g.armCorrupt("cw_alive") }},
}

// StoryActive reports whether the mystery has woken up.
func (g *Game) StoryActive() bool { return g.StoryStage >= 1 }

// storyQuipChance is the odds of a glitched quip replacing a normal one.
func (g *Game) storyQuipChance() float64 {
	switch {
	case g.StoryStage <= 0 || g.StoryStage > 6:
		return 0
	case g.StoryStage <= 2:
		return 0.15
	case g.StoryStage <= 4:
		return 0.25
	default:
		return 0.35
	}
}

// ClueFound reports whether a clue ID has been registered (meta).
func (g *Game) ClueFound(id string) bool { return g.storyClues[id] }

// DocUnlocked reports whether a terminal doc is readable.
func (g *Game) DocUnlocked(id string) bool {
	for _, d := range StoryDocs {
		if d.ID == id {
			return (d.Key == "" && g.StoryStage >= d.Act) || g.storyDocs[id]
		}
	}
	return false
}

// registerClue records a discovered clue (meta, survives firing), resets
// the pity counter and re-checks act progression. Saved immediately: a
// firing seconds later must not lose it.
func (g *Game) registerClue(id string) {
	if id == "" || g.storyClues[id] {
		return
	}
	if g.storyClues == nil {
		g.storyClues = map[string]bool{}
	}
	g.storyClues[id] = true
	g.storyDaysNoBeat = 0
	g.maybeAdvanceStory()
	g.Save()
}

// unlockDoc opens a keyed terminal doc (meta) and counts as a clue.
func (g *Game) unlockDoc(id string) {
	if g.storyDocs[id] {
		return
	}
	if g.storyDocs == nil {
		g.storyDocs = map[string]bool{}
	}
	g.storyDocs[id] = true
	g.registerClue("key_" + id)
}

// maybeAdvanceStory walks the act ladder as far as the current facts
// allow: each step needs its BestDay gate AND its key condition. The only
// place StoryStage ever changes (besides RESET).
func (g *Game) maybeAdvanceStory() {
	for {
		next := g.StoryStage + 1
		if next > StoryEndAct {
			return
		}
		ok := false
		switch next {
		case 1:
			ok = g.BestDay >= storyDayGates[1]
		case 2:
			ok = g.BestDay >= storyDayGates[2] && g.storyDocs["expediente"]
		case 3:
			ok = g.BestDay >= storyDayGates[3] && g.storyDocs["log1"]
		case 4:
			ok = g.BestDay >= storyDayGates[4] && g.storyDocs["log2"] && g.markedCluesFound() >= 1
		case 5:
			ok = g.BestDay >= storyDayGates[5] && (g.storyClues["cw_basement"] || g.storyClues["cw_watching"] || g.storyClues["in_seen"])
		case 6:
			ok = g.BestDay >= storyDayGates[6] && g.storyDocs["acceso"]
		case StoryEndAct:
			ok = g.storyDocs["puerta"]
		}
		if !ok {
			return
		}
		g.StoryStage = next
		g.storyQuipForced = true
	}
}

// storyObjectivePending reports whether the current act's key condition —
// the thing that would advance the story — is still unmet. While pending,
// the glitched quips show up more often and lean on the hint lines.
func (g *Game) storyObjectivePending() bool {
	switch g.StoryStage {
	case 1:
		return !g.storyDocs["expediente"]
	case 2:
		return !g.storyDocs["log1"]
	case 3:
		return !g.storyDocs["log2"] || g.markedCluesFound() < 1
	case 4:
		return !g.storyClues["cw_basement"] && !g.storyClues["cw_watching"] && !g.storyClues["in_seen"]
	case 5:
		return !g.storyDocs["acceso"]
	case 6:
		return !g.storyDocs["puerta"]
	}
	return false
}

// rollHintQuip decides whether a pending objective steers the quip pick.
func (g *Game) rollHintQuip() bool {
	return rand.Float64() < StoryHintOdds
}

// pickStoryQuip picks a line from the act's pool; while the objective is
// pending it prefers the hint lines that point at it. Before the terminal
// has ever been opened, the /terminal reveal is the only line shown.
func (g *Game) pickStoryQuip(act int, pending bool) string {
	code := Languages[g.LangIndex].Code
	if !g.termOpened {
		return StoryQuipScripts[1][code][terminalQuipLine]
	}
	pool := StoryQuipScripts[act]
	if hints := storyHintLines[act]; pending && len(hints) > 0 && g.rollHintQuip() {
		return pool[code][hints[pickQuip(len(hints), &g.lastHint)]]
	}
	return pool[code][pickQuip(len(pool[code]), &g.lastWeird)]
}

// markedCluesFound counts completed marked-word clues.
func (g *Game) markedCluesFound() int {
	n := 0
	for _, mw := range MarkedWords {
		if g.storyClues[mw.ClueID] {
			n++
		}
	}
	return n
}

// checkStoryTriggers fires any live, not-yet-fired ritual whose condition
// holds. Called after typing updates the counters and at day end.
func (g *Game) checkStoryTriggers() {
	if !g.StoryActive() {
		return
	}
	for _, t := range StoryTriggers {
		if g.StoryStage < t.Act || g.ritualDone[t.ID] {
			continue
		}
		if !t.Check(g) {
			continue
		}
		if g.ritualDone == nil {
			g.ritualDone = map[string]bool{}
		}
		g.ritualDone[t.ID] = true
		t.Effect(g)
	}
}

// rollCorrupt decides whether a completed word arms a corrupt one.
func (g *Game) rollCorrupt() bool {
	return rand.Float64() < CorruptChance
}

// armCorrupt queues a corrupt word for the next word boundary. clueID ""
// lets the claim pick any eligible word, preferring undiscovered clues.
// No-ops while the day already had one, another is pending, or an
// inspection owns the word slot.
func (g *Game) armCorrupt(clueID string) {
	if !g.StoryActive() || g.corruptArmed || g.WordCorrupt {
		return
	}
	if g.CorruptToday >= CorruptMaxPerDay {
		return
	}
	if g.InspectionArmed || g.InspectionActive() {
		return
	}
	g.corruptArmed = true
	g.corruptClueID = clueID
}

// pickCorruptWord resolves the armed clue to a concrete word in the active
// language, preferring undiscovered clues when no specific one was asked.
func (g *Game) pickCorruptWord() (string, string) {
	code := Languages[g.LangIndex].Code
	if g.corruptClueID != "" {
		for _, cw := range CorruptWords {
			if cw.ClueID == g.corruptClueID {
				return cw.ClueID, cw.Text[code]
			}
		}
	}
	var eligible, fresh []CorruptWord
	for _, cw := range CorruptWords {
		if cw.Act > g.StoryStage {
			continue
		}
		eligible = append(eligible, cw)
		if !g.storyClues[cw.ClueID] {
			fresh = append(fresh, cw)
		}
	}
	pool := fresh
	if len(pool) == 0 {
		pool = eligible
	}
	if len(pool) == 0 {
		return "", ""
	}
	cw := pool[rand.IntN(len(pool))]
	return cw.ClueID, cw.Text[code]
}

// claimCorrupt turns the next word into the armed corrupt word: violet,
// unpaid, and the day clock freezes until it's typed or it times out.
func (g *Game) claimCorrupt() {
	id, text := g.pickCorruptWord()
	if text == "" {
		g.corruptArmed = false
		return
	}
	g.corruptClueID = id
	g.Word = text
	g.Pos = 0
	g.WordCorrupt = true
	g.WordGolden = false
	g.CorruptTimer = CorruptTimeout
	g.CorruptToday++
}

// corruptComplete banks the clue for a typed corrupt word: no pay, no
// stats — the reward is the clue.
func (g *Game) corruptComplete() {
	g.registerClue(g.corruptClueID)
	code := Languages[g.LangIndex].Code
	g.setStoryToast(CorruptDoneScripts[code][pickQuip(len(CorruptDoneScripts[code]), &g.lastCorruptDone)])
	g.LastGain = 0
	g.Phase = PhaseSuccess
	g.PhaseTimer = SuccessTime
}

// checkMarkedWord counts a completed dictionary word against the watched
// list; hitting a counter fires the clue once, forever.
func (g *Game) checkMarkedWord() {
	if !g.StoryActive() {
		return
	}
	code := Languages[g.LangIndex].Code
	for _, mw := range MarkedWords {
		if g.StoryStage < mw.Act || g.storyClues[mw.ClueID] || g.Word != mw.Word[code] {
			continue
		}
		if g.markedCounts == nil {
			g.markedCounts = map[string]int{}
		}
		g.markedCounts[mw.ClueID]++
		if g.markedCounts[mw.ClueID] >= mw.Count {
			g.registerClue(mw.ClueID)
			g.setStoryToast(MarkedHitScripts[code][pickQuip(len(MarkedHitScripts[code]), &g.lastMarkedHit)])
		}
	}
}

// rollPossession decides whether an intern picks up a phrase of its own.
func (g *Game) rollPossession() bool {
	return rand.Float64() < PossessChance
}

// maybePossessIntern turns intern i's next word into a possessed phrase:
// typed on its own, pays nothing, registers the "in_seen" clue when done.
func (g *Game) maybePossessIntern(i int) bool {
	if g.StoryStage < 4 || g.possessedIdx >= 0 || g.possessedToday || g.WordCorrupt || g.corruptArmed {
		return false
	}
	if !g.rollPossession() {
		return false
	}
	g.possessIntern(i)
	return true
}

// possessIntern applies the possession (split out so tests can force it).
func (g *Game) possessIntern(i int) {
	code := Languages[g.LangIndex].Code
	in := &g.Interns[i]
	in.Word = PossessedScripts[code][pickQuip(len(PossessedScripts[code]), &g.lastPossessed)]
	in.Pos = 0
	in.Possessed = true
	g.possessedIdx = i
	g.possessedToday = true
}

// possessionComplete snaps the intern out of it: clue banked, no pay.
func (g *Game) possessionComplete(in *Intern) {
	in.Possessed = false
	g.possessedIdx = -1
	g.registerClue("in_seen")
	code := Languages[g.LangIndex].Code
	g.setStoryToast(PossessedDoneScripts[code][pickQuip(len(PossessedDoneScripts[code]), &g.lastPossessedDone)])
}

// tickStoryFX advances the cosmetic layer: story toast decay and the
// impossible clock readings of act 4+. Runs with the day clock.
func (g *Game) tickStoryFX(dt float64) {
	if g.ClockGlitch > 0 {
		g.ClockGlitch = math.Max(0, g.ClockGlitch-dt)
	}
	if g.StoryStage < 4 || g.StoryStage >= StoryEndAct {
		return
	}
	g.clockGlitchClock += dt
	for g.clockGlitchClock >= ClockGlitchGap {
		g.clockGlitchClock -= ClockGlitchGap
		if rand.Float64() < ClockGlitchOdds {
			g.ClockGlitch = ClockGlitchTime
			g.clockGlitchText = ClockGlitchTexts[rand.IntN(len(ClockGlitchTexts))]
		}
	}
}

// setStoryToast shows a violet story line above the word for a moment.
func (g *Game) setStoryToast(text string) {
	g.StoryToast = text
	g.StoryToastTimer = StoryToastTime
}

// storyDayBeat runs at every morning compose: pity-arms a corrupt word
// after too many quiet days and decides whether the morning quip glitches.
// Returns the story quip to use, or "" for a normal morning.
func (g *Game) storyDayBeat() string {
	if !g.StoryActive() {
		return ""
	}
	if g.StoryStage <= 6 && g.storyDaysNoBeat >= StoryPityDays {
		g.armCorrupt("")
	}
	act := min(g.StoryStage, StoryEndAct)
	if StoryQuipScripts[act] == nil {
		return ""
	}
	pending := g.storyObjectivePending()
	chance := g.storyQuipChance()
	if pending {
		chance += StoryHintBoost
	}
	if g.storyQuipForced || rand.Float64() < chance {
		g.storyQuipForced = false
		g.storyDaysNoBeat = 0
		return g.pickStoryQuip(act, pending)
	}
	g.storyDaysNoBeat++
	return ""
}

// pickDayEndStory decides whether tonight's summary quip glitches (half
// the morning odds; never over the first-payroll explainer).
func (g *Game) pickDayEndStory() {
	g.DayEndStory = false
	if !g.StoryActive() || g.DayEndFirst {
		return
	}
	act := min(g.StoryStage, StoryEndAct)
	pool := StoryQuipScripts[act]
	if pool == nil {
		return
	}
	pending := g.storyObjectivePending()
	chance := g.storyQuipChance()
	if pending {
		chance += StoryHintBoost
	}
	if rand.Float64() >= chance/2 {
		return
	}
	code := Languages[g.LangIndex].Code
	if !g.termOpened {
		// the /terminal reveal owns every story slot until it's used;
		// pre-terminal the stage is always 1, so the index matches the pool
		g.dayEndStoryIdx = terminalQuipLine
	} else if hints := storyHintLines[act]; pending && len(hints) > 0 && g.rollHintQuip() {
		g.dayEndStoryIdx = hints[pickQuip(len(hints), &g.lastHint)]
	} else {
		g.dayEndStoryIdx = pickQuip(len(pool[code]), &g.lastWeird)
	}
	g.DayEndStory = true
	g.storyDaysNoBeat = 0
}

// --- terminal scene -------------------------------------------------------

// EnterTerminal switches to the hidden terminal (reachable via /terminal
// in play, or from the main menu once discovered) and remembers where to
// return on exit.
func (g *Game) EnterTerminal() {
	g.termFrom = g.Scene
	g.TermBuf = ""
	if len(g.TermLog) == 0 {
		ui := g.UI()
		g.termPrint(ui.TermBanner)
		g.termPrint(ui.TermHelpHint)
	}
	if !g.termOpened {
		// discovery done: the story quips may move past the /terminal reveal
		g.termOpened = true
		g.Save()
	}
	g.Scene = SceneTerminal
}

// TerminalAvailable reports whether /terminal answers at all.
func (g *Game) TerminalAvailable() bool {
	return g.StoryActive() && len(g.storyClues) > 0
}

// TermTypeChar appends one character to the terminal prompt.
func (g *Game) TermTypeChar(r rune) {
	if len(g.TermBuf) < 40 {
		g.TermBuf += string(r)
	}
}

// TermBackspace deletes the last rune of the terminal prompt.
func (g *Game) TermBackspace() {
	if g.TermBuf == "" {
		return
	}
	runes := []rune(g.TermBuf)
	g.TermBuf = string(runes[:len(runes)-1])
}

// TermExit leaves the terminal back to wherever it was opened from (play
// or the main menu), keeping the session log.
func (g *Game) TermExit() {
	g.TermBuf = ""
	g.Scene = g.termFrom
}

// termPrint appends one line to the terminal log, capped.
func (g *Game) termPrint(line string) {
	g.TermLog = append(g.TermLog, line)
	if len(g.TermLog) > 200 {
		g.TermLog = g.TermLog[len(g.TermLog)-200:]
	}
}

// TermSubmit runs the buffered terminal input: built-in commands in both
// languages, doc reads, and loose tokens matched against document keys.
func (g *Game) TermSubmit() {
	input := strings.ToLower(strings.TrimSpace(g.TermBuf))
	g.TermBuf = ""
	if input == "" {
		return
	}
	ui := g.UI()
	g.termPrint("> " + input)
	cmd, arg := input, ""
	if i := strings.IndexByte(input, ' '); i >= 0 {
		cmd, arg = input[:i], strings.TrimSpace(input[i+1:])
	}
	switch cmd {
	case "ayuda", "help":
		for _, l := range ui.TermHelp {
			g.termPrint(l)
		}
	case "docs", "archivos":
		g.termListDocs()
	case "leer", "read":
		g.termReadDoc(arg)
	case "salir", "exit":
		g.TermExit()
	default:
		if g.termTryKey(input) {
			return
		}
		g.termPrint(ui.TermUnknown)
	}
}

// termListDocs prints every doc the current act exposes; keyed ones still
// locked show as a redacted teaser so the player knows there is more.
func (g *Game) termListDocs() {
	ui := g.UI()
	code := Languages[g.LangIndex].Code
	g.termPrint(ui.TermDocsHeader)
	any := false
	for _, d := range StoryDocs {
		if d.Act > g.StoryStage {
			continue
		}
		any = true
		if g.DocUnlocked(d.ID) {
			g.termPrint("  " + d.ID + " — " + d.Title[code])
		} else {
			g.termPrint("  ??? — " + ui.TermLockedTag)
		}
	}
	if !any {
		g.termPrint(ui.TermNoDocs)
	}
}

// termReadDoc prints an unlocked doc's title and body.
func (g *Game) termReadDoc(id string) {
	ui := g.UI()
	code := Languages[g.LangIndex].Code
	for _, d := range StoryDocs {
		if d.ID != id || d.Act > g.StoryStage {
			continue
		}
		if !g.DocUnlocked(d.ID) {
			g.termPrint(ui.TermLocked)
			return
		}
		g.termPrint("── " + d.Title[code] + " ──")
		for _, l := range d.Body[code] {
			g.termPrint(l)
		}
		return
	}
	g.termPrint(ui.TermUnknown)
}

// termTryKey matches a loose token against the keys of stage-visible
// locked docs; a hit unlocks the doc (and may advance the act).
func (g *Game) termTryKey(token string) bool {
	ui := g.UI()
	code := Languages[g.LangIndex].Code
	for _, d := range StoryDocs {
		if d.Key == "" || d.Key != token || d.Act > g.StoryStage || g.storyDocs[d.ID] {
			continue
		}
		g.unlockDoc(d.ID)
		g.termPrint(ui.TermKeyAccepted)
		g.termPrint("── " + d.Title[code] + " ──")
		for _, l := range d.Body[code] {
			g.termPrint(l)
		}
		return true
	}
	return false
}

// --- persistence helpers ---------------------------------------------------

// boolSetToSlice flattens a string-set to a sorted slice for the save file.
func boolSetToSlice(m map[string]bool) []string {
	if len(m) == 0 {
		return nil
	}
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// sliceToBoolSet rebuilds a string-set from its saved slice.
func sliceToBoolSet(s []string) map[string]bool {
	m := map[string]bool{}
	for _, k := range s {
		m[k] = true
	}
	return m
}

// resetStory wipes all mystery progress (menu RESET).
func (g *Game) resetStory() {
	g.StoryStage = 0
	g.storyClues = map[string]bool{}
	g.storyDocs = map[string]bool{}
	g.markedCounts = map[string]int{}
	g.TermLog = nil
	g.termOpened = false
	g.resetStoryRun()
}

// resetStoryRun clears the per-run story state (new run / restore).
func (g *Game) resetStoryRun() {
	g.WordCorrupt = false
	g.corruptArmed = false
	g.corruptClueID = ""
	g.CorruptTimer = 0
	g.CorruptToday = 0
	g.ritualDone = map[string]bool{}
	g.storyDaysNoBeat = 0
	g.possessedIdx = -1
	g.possessedToday = false
	g.ClockGlitch = 0
	g.clockGlitchClock = 0
	g.StoryToast = ""
	g.StoryToastTimer = 0
	g.DayStartStory = false
	g.DayEndStory = false
	g.DayEndCode = ""
}
