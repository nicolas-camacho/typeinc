package game

import (
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// docByID returns the StoryDoc with the given ID (order-independent).
func docByID(t *testing.T, id string) StoryDoc {
	t.Helper()
	for _, d := range StoryDocs {
		if d.ID == id {
			return d
		}
	}
	t.Fatalf("no StoryDoc with ID %q", id)
	return StoryDoc{}
}

// storyPlay is startPlay with the mystery already awake at the given act.
func storyPlay(t *testing.T, word string, stage int) *Game {
	t.Helper()
	g := startPlay(t, word)
	g.StoryStage = stage
	g.BestDay = storyDayGates[min(stage, 6)]
	return g
}

// --- content shape ----------------------------------------------------------

func TestStoryScriptsMatchAcrossLanguages(t *testing.T) {
	for act, pool := range StoryQuipScripts {
		if len(pool["es"]) == 0 || len(pool["es"]) != len(pool["en"]) {
			t.Fatalf("StoryQuipScripts[%d]: es %d lines, en %d lines", act, len(pool["es"]), len(pool["en"]))
		}
	}
	for act := 1; act <= 6; act++ {
		if StoryQuipScripts[act] == nil {
			t.Fatalf("StoryQuipScripts missing act %d", act)
		}
	}
	if StoryQuipScripts[StoryEndAct] == nil {
		t.Fatal("StoryQuipScripts missing the epilogue pool")
	}
	for name, pool := range map[string]map[string][]string{
		"CorruptDone": CorruptDoneScripts, "MarkedHit": MarkedHitScripts,
		"Possessed": PossessedScripts, "PossessedDone": PossessedDoneScripts,
	} {
		if len(pool["es"]) == 0 || len(pool["es"]) != len(pool["en"]) {
			t.Fatalf("%sScripts: es %d lines, en %d lines", name, len(pool["es"]), len(pool["en"]))
		}
	}
}

func TestCorruptWordsAndDocsWellFormed(t *testing.T) {
	typeable := regexp.MustCompile(`^[a-z ]{3,24}$`)
	seen := map[string]bool{}
	for _, cw := range CorruptWords {
		if seen[cw.ClueID] {
			t.Fatalf("duplicate corrupt ClueID %q", cw.ClueID)
		}
		seen[cw.ClueID] = true
		if cw.Act < 1 || cw.Act > 6 {
			t.Fatalf("corrupt word %q has act %d", cw.ClueID, cw.Act)
		}
		for _, code := range []string{"es", "en"} {
			if !typeable.MatchString(cw.Text[code]) {
				t.Fatalf("corrupt word %q (%s) = %q: not typeable ASCII", cw.ClueID, code, cw.Text[code])
			}
		}
	}
	// possessed phrases go through the same byte-indexed typing renderer
	for _, code := range []string{"es", "en"} {
		for _, p := range PossessedScripts[code] {
			if !typeable.MatchString(p) {
				t.Fatalf("possessed phrase %q (%s): not typeable ASCII", p, code)
			}
		}
	}
	token := regexp.MustCompile(`^[a-z0-9]+$`)
	ids, keys := map[string]bool{}, map[string]bool{}
	for _, d := range StoryDocs {
		if ids[d.ID] || !token.MatchString(d.ID) {
			t.Fatalf("doc ID %q duplicated or not a lowercase token", d.ID)
		}
		ids[d.ID] = true
		if d.Key != "" {
			if keys[d.Key] || !token.MatchString(d.Key) {
				t.Fatalf("doc key %q duplicated or not a lowercase token", d.Key)
			}
			keys[d.Key] = true
		}
		if len(d.Body["es"]) == 0 || len(d.Body["es"]) != len(d.Body["en"]) {
			t.Fatalf("doc %q: es %d body lines, en %d", d.ID, len(d.Body["es"]), len(d.Body["en"]))
		}
		if d.Title["es"] == "" || d.Title["en"] == "" {
			t.Fatalf("doc %q: missing a title", d.ID)
		}
	}
}

func TestMarkedWordsExistInDictionaries(t *testing.T) {
	for _, mw := range MarkedWords {
		for _, code := range []string{"es", "en"} {
			if !slices.Contains(WordLists[code], mw.Word[code]) {
				t.Fatalf("marked word %q (%s) not in the %s dictionary", mw.Word[code], mw.ClueID, code)
			}
		}
	}
}

// --- corrupt word mechanic --------------------------------------------------

func TestCorruptArmClaimsNextWordBoundary(t *testing.T) {
	g := storyPlay(t, "casa", 1)
	g.armCorrupt("cw_help")
	if g.WordCorrupt || g.Word != "casa" {
		t.Fatal("arming must not interrupt the word in progress")
	}
	typeWord(g, "casa")
	g.Tick(SuccessTime + 0.01) // success expires → nextWord claims
	if !g.WordCorrupt || g.Word != "ayudame" {
		t.Fatalf("Word = %q, WordCorrupt = %v; want the armed corrupt word", g.Word, g.WordCorrupt)
	}
	if g.WordGolden {
		t.Fatal("a corrupt word must never be golden")
	}
	if g.CorruptToday != 1 || g.CorruptTimer != CorruptTimeout {
		t.Fatalf("CorruptToday = %d, CorruptTimer = %v; want 1, %v", g.CorruptToday, g.CorruptTimer, CorruptTimeout)
	}
}

func TestCorruptFreezesClockEventsAndInterns(t *testing.T) {
	g := storyPlay(t, "casa", 1)
	g.Interns = []Intern{{Word: "burocracia"}}
	g.FrenzyTimer = 5
	g.armCorrupt("cw_help")
	g.Word = "x" // finish instantly
	g.Pos = 0
	typeWord(g, "x")
	g.Tick(SuccessTime + 0.01)
	if !g.WordCorrupt {
		t.Fatal("corrupt word not claimed")
	}
	day, frenzy, pos, score := g.DayTimer, g.FrenzyTimer, g.Interns[0].Pos, g.Score
	g.Tick(5)
	if g.DayTimer != day {
		t.Fatalf("DayTimer moved during a corrupt word: %v → %v", day, g.DayTimer)
	}
	if g.FrenzyTimer != frenzy {
		t.Fatalf("events ticked during a corrupt word: %v → %v", frenzy, g.FrenzyTimer)
	}
	if g.Interns[0].Pos != pos || g.Score != score {
		t.Fatal("interns earned during a corrupt word")
	}
	if g.CorruptTimer >= CorruptTimeout {
		t.Fatal("CorruptTimer did not run")
	}
}

func TestCorruptTimeoutRestoresNormalPlay(t *testing.T) {
	g := storyPlay(t, "casa", 1)
	g.armCorrupt("cw_help")
	g.Word = "x"
	typeWord(g, "x")
	g.Tick(SuccessTime + 0.01)
	g.Tick(CorruptTimeout + 0.01)
	if g.WordCorrupt {
		t.Fatal("corrupt word still active past its timeout")
	}
	if g.ClueFound("cw_help") {
		t.Fatal("an unanswered corrupt word must not register its clue")
	}
	day := g.DayTimer
	g.Tick(1)
	if g.DayTimer >= day {
		t.Fatal("day clock still frozen after the corrupt word dissolved")
	}
}

func TestCorruptCompletionRegistersClueNoPay(t *testing.T) {
	g := storyPlay(t, "casa", 1)
	g.armCorrupt("cw_help")
	g.Word = "x"
	typeWord(g, "x")
	score, words, streak := g.Score, g.DayWords, g.Streak
	g.Tick(SuccessTime + 0.01)
	typeWord(g, "ayudame")
	if g.Phase != PhaseSuccess {
		t.Fatalf("Phase = %v after typing the corrupt word, want PhaseSuccess", g.Phase)
	}
	if !g.ClueFound("cw_help") {
		t.Fatal("completing the corrupt word must register its clue")
	}
	if g.Score != score || g.DayWords != words || g.Streak != streak {
		t.Fatalf("corrupt word touched the books: score %d→%d, words %d→%d, streak %d→%d",
			score, g.Score, words, g.DayWords, streak, g.Streak)
	}
	if g.StoryToastTimer <= 0 || g.StoryToast == "" {
		t.Fatal("no story toast after the corrupt word")
	}
	g.Tick(SuccessTime + 0.01)
	if g.WordCorrupt {
		t.Fatal("WordCorrupt still set on the next word")
	}
}

func TestCorruptMistakesDontTouchStats(t *testing.T) {
	g := storyPlay(t, "casa", 1)
	g.Streak = 7
	g.armCorrupt("cw_help")
	g.Word = "x"
	typeWord(g, "x")
	g.Tick(SuccessTime + 0.01)
	g.TypeChar('z') // wrong
	if g.Phase != PhaseError || g.Pos != 0 {
		t.Fatalf("Phase = %v, Pos = %d; want the usual error flash", g.Phase, g.Pos)
	}
	if g.Streak != 8 { // 7 + the "x" word
		t.Fatalf("Streak = %d, want 8: corrupt mistakes must not reset it", g.Streak)
	}
	if g.DayFails != 0 || len(g.dayFailCounts) != 0 || len(g.FailCounts) != 0 {
		t.Fatal("corrupt mistakes leaked into the fail stats")
	}
}

func TestCorruptInspectionMutualExclusion(t *testing.T) {
	g := storyPlay(t, "casa", 1)
	g.InspectionArmed = true
	g.armCorrupt("cw_help")
	if g.corruptArmed {
		t.Fatal("armCorrupt must refuse while an inspection owns the word slot")
	}
	g.InspectionArmed = false
	g.armCorrupt("cw_help")
	g.EventsToday = 0
	g.eventSpacing = EventMinSpacing + 1
	for range 200 {
		g.maybeTriggerEvent()
	}
	if g.EventsToday != 0 {
		t.Fatal("events rolled while a corrupt word was armed")
	}
}

func TestCorruptDailyCap(t *testing.T) {
	g := storyPlay(t, "casa", 1)
	g.CorruptToday = CorruptMaxPerDay
	g.armCorrupt("cw_help")
	if g.corruptArmed {
		t.Fatal("armCorrupt ignored the daily cap")
	}
}

func TestCycleLangClearsActiveCorrupt(t *testing.T) {
	g := storyPlay(t, "casa", 1)
	g.armCorrupt("cw_help")
	g.Word = "x"
	typeWord(g, "x")
	g.Tick(SuccessTime + 0.01)
	if !g.WordCorrupt {
		t.Fatal("corrupt word not claimed")
	}
	g.cycleLang(1)
	if g.WordCorrupt {
		t.Fatal("language switch left a corrupt word active on a fresh dictionary word")
	}
}

// --- companion mechanics ----------------------------------------------------

func TestMarkedWordCountsAndFiresOnce(t *testing.T) {
	g := storyPlay(t, "sombra", 3)
	for i := 0; i < MarkedWords[0].Count; i++ {
		g.Word = "sombra"
		g.Pos = 0
		g.Phase = PhaseTyping
		g.WordCorrupt = false
		g.corruptArmed = false
		typeWord(g, "sombra")
	}
	if !g.ClueFound("mk_shadow") {
		t.Fatalf("markedCounts = %v: clue not registered after %d completions", g.markedCounts, MarkedWords[0].Count)
	}
	if g.StoryToastTimer <= 0 {
		t.Fatal("no toast when the marked counter completed")
	}
	// once fired it never counts again
	n := g.markedCounts["mk_shadow"]
	g.Word = "sombra"
	g.Pos = 0
	g.Phase = PhaseTyping
	typeWord(g, "sombra")
	if g.markedCounts["mk_shadow"] != n {
		t.Fatal("marked counter kept counting after its clue fired")
	}
}

func TestMarkedCountsSurviveReload(t *testing.T) {
	path := t.TempDir() + "/save.json"
	g := storyPlay(t, "sombra", 3)
	g.SavePath = path
	g.markedCounts["mk_shadow"] = 2
	g.Save()

	g2 := New()
	g2.SavePath = path
	g2.Load()
	if g2.markedCounts["mk_shadow"] != 2 {
		t.Fatalf("markedCounts after reload = %v, want mk_shadow:2", g2.markedCounts)
	}
}

func TestPossessedInternPaysNothingAndRegistersClue(t *testing.T) {
	g := storyPlay(t, "casa", 4)
	g.Interns = []Intern{{Word: "casa"}}
	g.possessIntern(0)
	if !g.Interns[0].Possessed || g.Interns[0].Pos != 0 {
		t.Fatal("possession did not take the intern's word")
	}
	phrase := g.Interns[0].Word
	score := g.Score
	// enough active seconds to type the phrase at base CPS
	need := float64(len(phrase))/g.InternCPS() + 2
	for i := 0.0; i < need; i += 0.05 {
		g.Tick(0.05)
	}
	if g.Score != score {
		t.Fatalf("Score %d → %d: a possessed intern must not earn", score, g.Score)
	}
	if !g.ClueFound("in_seen") {
		t.Fatal("possession completion did not register its clue")
	}
	if g.Interns[0].Possessed {
		t.Fatal("intern still possessed after finishing the phrase")
	}
}

func TestDayEndCodeShowsOnlyWhileAccesoLocked(t *testing.T) {
	g := storyPlay(t, "casa", 5)
	g.Score = 999999
	g.DayTimer = 0.001
	g.Tick(0.01)
	if g.DayEndCode != DayEndCodeText {
		t.Fatalf("DayEndCode = %q at act 5 with acceso locked, want %q", g.DayEndCode, DayEndCodeText)
	}
	g.storyDocs["acceso"] = true
	g.Scene = ScenePlay
	g.Phase = PhaseTyping
	g.Score = 999999
	g.DayTimer = 0.001
	g.Tick(0.01)
	if g.DayEndCode != "" {
		t.Fatalf("DayEndCode = %q after acceso unlocked, want empty", g.DayEndCode)
	}
}

func TestClockGlitchOverridesDayClock(t *testing.T) {
	g := storyPlay(t, "casa", 4)
	g.ClockGlitch = 1
	g.clockGlitchText = "66:99"
	if g.DayClock() != "66:99" {
		t.Fatalf("DayClock() = %q during a glitch, want 66:99", g.DayClock())
	}
	g.Tick(1.5)
	if g.DayClock() == "66:99" {
		t.Fatal("clock glitch never expired")
	}
}

// --- rituals ----------------------------------------------------------------

func TestRitualFailFiveOnEvenDayArmsCorrupt(t *testing.T) {
	g := storyPlay(t, "casa", 1)
	g.Day = 3 // odd: must not fire
	g.DayFails = 4
	g.TypeChar('z')
	if g.corruptArmed {
		t.Fatal("ritual fired on an odd day")
	}
	g.ritualDone = map[string]bool{}
	g.Day = 4
	g.Phase = PhaseTyping
	g.DayFails = 4
	g.TypeChar('z') // 5th fail on an even day
	if !g.corruptArmed || g.corruptClueID != "cw_help" {
		t.Fatalf("armed = %v, clue = %q; want cw_help armed", g.corruptArmed, g.corruptClueID)
	}
	// once per run (dayFailCounts cleared so the repeat-fail ritual can't
	// fire instead)
	g.corruptArmed = false
	g.dayFailCounts = map[string]int{}
	g.Phase = PhaseTyping
	g.TypeChar('z')
	if g.corruptArmed {
		t.Fatal("ritual re-fired within the same run")
	}
}

func TestRitualDonePersistsAndResetsOnNewRun(t *testing.T) {
	path := t.TempDir() + "/save.json"
	g := storyPlay(t, "casa", 1)
	g.SavePath = path
	g.ritualDone["fail_even"] = true
	g.Save()

	g2 := New()
	g2.SavePath = path
	g2.Load()
	g2.continueRun()
	if !g2.ritualDone["fail_even"] {
		t.Fatal("ritualDone lost on save-quit-resume")
	}
	g2.initRun()
	if len(g2.ritualDone) != 0 {
		t.Fatal("ritualDone survived a brand-new run")
	}
}

// --- act progression --------------------------------------------------------

func TestStoryWakesAtBestDayGate(t *testing.T) {
	g := startPlay(t, "casa")
	g.Day = 2
	g.Score = 999999
	g.DayTimer = 0.001
	g.Tick(0.01) // day 3 reached
	if g.StoryStage != 1 {
		t.Fatalf("StoryStage = %d after reaching day %d, want 1", g.StoryStage, g.Day)
	}
	if !g.storyQuipForced {
		t.Fatal("waking up must force the next morning quip")
	}
	// idempotent
	g.maybeAdvanceStory()
	if g.StoryStage != 1 {
		t.Fatalf("StoryStage = %d after re-check, want 1", g.StoryStage)
	}
}

func TestActGatesNeedBestDayAndKeys(t *testing.T) {
	g := New()
	// all the keys, none of the days: nothing advances past act 1's needs
	g.storyDocs = map[string]bool{"expediente": true, "log1": true, "log2": true, "acceso": true, "puerta": true}
	g.storyClues = map[string]bool{"mk_shadow": true, "cw_basement": true}
	g.BestDay = 3
	g.maybeAdvanceStory()
	if g.StoryStage != 1 {
		t.Fatalf("StoryStage = %d with BestDay 3, want 1 (act 2 needs BestDay 7)", g.StoryStage)
	}
	// now the days too: the ladder climbs to the end
	g.BestDay = 30
	g.maybeAdvanceStory()
	if g.StoryStage != StoryEndAct {
		t.Fatalf("StoryStage = %d with everything unlocked, want %d", g.StoryStage, StoryEndAct)
	}
}

func TestForcedStoryQuipOnDayStart(t *testing.T) {
	g := storyPlay(t, "casa", 1)
	g.storyQuipForced = true
	g.composeDayStart()
	if !g.DayStartStory {
		t.Fatal("forced story beat did not glitch the morning quip")
	}
	if !slices.Contains(StoryQuipScripts[1]["es"], g.DayStartText) {
		t.Fatalf("DayStartText = %q: not from the act-1 pool", g.DayStartText)
	}
	if g.storyQuipForced {
		t.Fatal("forced flag not consumed")
	}
}

func TestTerminalQuipRepeatsUntilOpened(t *testing.T) {
	g := storyPlay(t, "casa", 1)
	want := StoryQuipScripts[1]["es"][terminalQuipLine]
	for range 2 {
		g.storyQuipForced = true
		g.composeDayStart()
		if g.DayStartText != want {
			t.Fatalf("DayStartText = %q before opening the terminal, want the /terminal reveal", g.DayStartText)
		}
	}
	g.storyClues["cw_help"] = true
	g.EnterTerminal()
	if !g.termOpened {
		t.Fatal("EnterTerminal did not mark the terminal as opened")
	}

	// the flag survives a save/load roundtrip
	path := t.TempDir() + "/save.json"
	g.SavePath = path
	g.Save()
	g2 := New()
	g2.SavePath = path
	g2.Load()
	if !g2.termOpened {
		t.Fatal("TermOpened lost in the save roundtrip")
	}
}

func TestPityArmsCorruptAfterQuietDays(t *testing.T) {
	g := storyPlay(t, "casa", 1)
	g.storyDaysNoBeat = StoryPityDays
	g.composeDayStart()
	if !g.corruptArmed {
		t.Fatal("pity did not arm a corrupt word")
	}
}

// --- terminal ---------------------------------------------------------------

func TestTerminalCommandGatedByStage(t *testing.T) {
	g := startPlay(t, "casa")
	g.TypeChar('/')
	typeWord(g, "terminal")
	g.CommandSubmit()
	if g.Scene != ScenePlay || g.CommandErr <= 0 {
		t.Fatal("/terminal answered before the story woke up")
	}
	g.StoryStage = 1
	g.storyClues["cw_help"] = true
	g.CommandErr = 0
	g.TypeChar('/')
	typeWord(g, "terminal")
	g.CommandSubmit()
	if g.Scene != SceneTerminal {
		t.Fatalf("Scene = %v after /terminal with a clue, want SceneTerminal", g.Scene)
	}
	if len(g.TermLog) == 0 {
		t.Fatal("terminal opened without its banner")
	}
}

func TestTerminalKeyUnlocksDocAndAdvances(t *testing.T) {
	g := storyPlay(t, "casa", 1)
	g.BestDay = 7 // act 2's day gate already reached
	g.storyClues["cw_name"] = true
	g.Scene = SceneTerminal

	g.TermBuf = "nope"
	g.TermSubmit()
	if g.DocUnlocked("expediente") {
		t.Fatal("a wrong token unlocked a doc")
	}
	g.TermBuf = "vera"
	g.TermSubmit()
	if !g.DocUnlocked("expediente") {
		t.Fatal("the right key did not unlock its doc")
	}
	if g.StoryStage != 2 {
		t.Fatalf("StoryStage = %d after the act-1 key with BestDay 7, want 2", g.StoryStage)
	}
	joined := strings.Join(g.TermLog, "\n")
	if !strings.Contains(joined, docByID(t, "expediente").Title["es"]) {
		t.Fatal("unlocking a doc did not print it")
	}
}

func TestTerminalReadAndLocks(t *testing.T) {
	g := storyPlay(t, "casa", 1)
	g.Scene = SceneTerminal
	g.TermBuf = "leer aviso" // act-1 keyless doc: readable
	g.TermSubmit()
	joined := strings.Join(g.TermLog, "\n")
	if !strings.Contains(joined, docByID(t, "aviso").Title["es"]) {
		t.Fatal("keyless act-1 doc not readable")
	}
	g.TermBuf = "leer expediente" // keyed, locked
	g.TermSubmit()
	if !strings.Contains(strings.Join(g.TermLog, "\n"), g.UI().TermLocked) {
		t.Fatal("locked doc did not answer with the locked message")
	}
	g.TermBuf = "leer puerta" // act 6: not even visible
	g.TermSubmit()
	if !strings.Contains(g.TermLog[len(g.TermLog)-1], g.UI().TermUnknown) {
		t.Fatal("an out-of-act doc leaked")
	}
}

func TestTerminalDocsSurviveFiring(t *testing.T) {
	path := t.TempDir() + "/save.json"
	g := storyPlay(t, "casa", 2)
	g.SavePath = path
	g.unlockDoc("expediente")

	// get fired
	g.Score = 0
	g.DayTimer = 0.001
	g.Tick(0.01)
	if !g.DayEndFired {
		t.Fatal("expected a firing")
	}

	g2 := New()
	g2.SavePath = path
	g2.Load()
	g2.initRun() // new run after the firing
	if !g2.DocUnlocked("expediente") {
		t.Fatal("an unlocked doc did not survive the firing")
	}
	if g2.StoryStage < 2 {
		t.Fatalf("StoryStage = %d after reload, want ≥2", g2.StoryStage)
	}
}

func TestTerminalPausesDayClock(t *testing.T) {
	g := storyPlay(t, "casa", 1)
	g.Interns = []Intern{{Word: "burocracia"}}
	g.Scene = SceneTerminal
	day, pos := g.DayTimer, g.Interns[0].Pos
	g.Tick(5)
	if g.DayTimer != day || g.Interns[0].Pos != pos {
		t.Fatal("the terminal did not pause the day")
	}
}

// --- persistence ------------------------------------------------------------

func TestStorySaveRoundtrip(t *testing.T) {
	path := t.TempDir() + "/save.json"
	g := storyPlay(t, "casa", 3)
	g.SavePath = path
	g.storyClues["cw_help"] = true
	g.storyDocs["expediente"] = true
	g.markedCounts["mk_door"] = 2
	g.Save()

	g2 := New()
	g2.SavePath = path
	g2.Load()
	if g2.StoryStage != 3 || !g2.ClueFound("cw_help") || !g2.storyDocs["expediente"] || g2.markedCounts["mk_door"] != 2 {
		t.Fatalf("story state lost in the roundtrip: stage %d, clues %v, docs %v, marked %v",
			g2.StoryStage, g2.storyClues, g2.storyDocs, g2.markedCounts)
	}
}

func TestSaveVersionGate(t *testing.T) {
	// a v2 save full of progress: everything wiped except settings
	raw := `{"Version":2,"Run":{"Score":500,"Day":6},"LangIndex":1,"Volume":0.7,"DisplayMode":1,
		"HRPoints":5,"HRMultLevel":2,"TotalWords":100,"BestStreak":9,"BestDay":8,"FailCounts":{"casa":3},
		"StoryStage":3,"StoryClues":["cw_help"]}`
	path := t.TempDir() + "/save.json"
	if err := os.WriteFile(path, []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}
	g := New()
	g.SavePath = path
	g.Load()
	if g.HRPoints != 0 || g.BestDay != 0 || g.TotalWords != 0 || g.savedRun != nil || g.StoryStage != 0 {
		t.Fatalf("v2 save not wiped: HR %d, BestDay %d, words %d, run %v, stage %d",
			g.HRPoints, g.BestDay, g.TotalWords, g.savedRun, g.StoryStage)
	}
	if g.LangIndex != 1 || g.Volume != 0.7 || g.DisplayMode != 1 {
		t.Fatalf("settings dropped by the wipe: lang %d, vol %v, display %d", g.LangIndex, g.Volume, g.DisplayMode)
	}

	// an explicit older GameVersion gets the same wipe
	older := strings.Replace(raw, `"Version":2,`, `"GameVersion":"0.2.0",`, 1)
	if err := os.WriteFile(path, []byte(older), 0o644); err != nil {
		t.Fatal(err)
	}
	g2 := New()
	g2.SavePath = path
	g2.Load()
	if g2.HRPoints != 0 || g2.savedRun != nil || g2.StoryStage != 0 {
		t.Fatal("0.2.0 save not wiped")
	}

	// dotted compare is numeric per part, not lexicographic
	if saveVersionLess("0.10.0", "0.3.0") {
		t.Fatal("saveVersionLess: 0.10.0 counted as older than 0.3.0")
	}
	if !saveVersionLess("", SaveVersion) || !saveVersionLess("0.2.9", SaveVersion) {
		t.Fatal("saveVersionLess: empty or 0.2.9 not counted as older")
	}

	// a current-version roundtrip keeps everything
	g3 := storyPlay(t, "casa", 2)
	g3.SavePath = path
	g3.HRPoints = 7
	g3.Save()
	g4 := New()
	g4.SavePath = path
	g4.Load()
	if g4.HRPoints != 7 || g4.StoryStage != 2 || g4.savedRun == nil {
		t.Fatalf("v%s roundtrip lost progress: HR %d, stage %d, run %v",
			SaveVersion, g4.HRPoints, g4.StoryStage, g4.savedRun)
	}
}

func TestVeraNotesCoverEveryRitualAct(t *testing.T) {
	for act := 1; act <= 5; act++ {
		id := "nota" + string(rune('0'+act))
		d := docByID(t, id)
		if d.Act != act || d.Key != "" {
			t.Fatalf("doc %q: act %d key %q, want act %d and no key", id, d.Act, d.Key, act)
		}
	}
}

func TestStoryHintLinesValid(t *testing.T) {
	for act := 1; act <= 6; act++ {
		hints := storyHintLines[act]
		if len(hints) < 2 {
			t.Fatalf("act %d has %d hint lines, want at least 2", act, len(hints))
		}
		for _, code := range []string{"es", "en"} {
			pool := StoryQuipScripts[act][code]
			for _, i := range hints {
				if i < 0 || i >= len(pool) {
					t.Fatalf("act %d hint index %d out of range (%s pool has %d lines)", act, i, code, len(pool))
				}
			}
		}
	}
}

func TestStoryObjectivePendingLadder(t *testing.T) {
	g := New()
	g.StoryStage = 1
	if !g.storyObjectivePending() {
		t.Fatal("act 1 objective not pending without expediente")
	}
	g.storyDocs["expediente"] = true
	if g.storyObjectivePending() {
		t.Fatal("act 1 objective still pending with expediente unlocked")
	}
	g.StoryStage = 4
	if !g.storyObjectivePending() {
		t.Fatal("act 4 objective not pending without any act-4 clue")
	}
	g.storyClues["in_seen"] = true
	if g.storyObjectivePending() {
		t.Fatal("act 4 objective still pending with the possession clue")
	}
	g.StoryStage = StoryEndAct
	if g.storyObjectivePending() {
		t.Fatal("a solved story must have nothing pending")
	}
}

func TestTerminalMenuEntryGated(t *testing.T) {
	g := New()
	for _, it := range g.MenuItems() {
		if it.ID == "terminal" {
			t.Fatal("terminal menu entry visible before the story woke up")
		}
	}
	g.StoryStage = 1
	g.storyClues["cw_help"] = true
	idx := -1
	for i, it := range g.MenuItems() {
		if it.ID == "terminal" {
			idx = i
		}
	}
	if idx < 0 {
		t.Fatal("terminal menu entry missing once discovered")
	}
	g.MenuIndex = idx
	g.MenuSelect()
	if g.Scene != SceneTerminal {
		t.Fatalf("Scene = %v after selecting the menu terminal, want SceneTerminal", g.Scene)
	}
	g.TermExit()
	if g.Scene != SceneMenu {
		t.Fatalf("Scene = %v after TermExit from the menu, want SceneMenu", g.Scene)
	}

	// entered from play, it returns to play
	g2 := storyPlay(t, "casa", 1)
	g2.storyClues["cw_help"] = true
	g2.TypeChar('/')
	typeWord(g2, "terminal")
	g2.CommandSubmit()
	if g2.Scene != SceneTerminal {
		t.Fatalf("Scene = %v after /terminal, want SceneTerminal", g2.Scene)
	}
	g2.TermExit()
	if g2.Scene != ScenePlay {
		t.Fatalf("Scene = %v after TermExit from play, want ScenePlay", g2.Scene)
	}
}

func TestResetWipesStory(t *testing.T) {
	g := storyPlay(t, "casa", 3)
	g.storyClues["cw_help"] = true
	g.storyDocs["expediente"] = true
	g.TermLog = []string{"x"}
	g.resetSave()
	if g.StoryStage != 0 || len(g.storyClues) != 0 || len(g.storyDocs) != 0 || g.TermLog != nil {
		t.Fatal("RESET left story progress behind")
	}
}
