package game

import (
	"strings"
	"testing"
)

// startPlay puts a game into ScenePlay with a known word, skipping the
// intro story.
func startPlay(t *testing.T, word string) *Game {
	t.Helper()
	g := New()
	g.initRun()
	g.Scene = ScenePlay
	g.Word = word
	g.Pos = 0
	g.WordGolden = false // initRun's word roll must not leak into forced words
	return g
}

func typeWord(g *Game, word string) {
	for _, r := range word {
		g.TypeChar(r)
	}
}

func TestCompleteWordScoresAndEntersSuccess(t *testing.T) {
	g := startPlay(t, "casa")
	typeWord(g, "casa")
	if g.Phase != PhaseSuccess {
		t.Fatalf("Phase = %v, want PhaseSuccess", g.Phase)
	}
	// streak 1 → combo x1.1; 4 × 1 × 1.1 = 4.4 → 4
	if g.Score != 4 || g.LastGain != 4 {
		t.Fatalf("Score = %d, LastGain = %d, want 4, 4", g.Score, g.LastGain)
	}
	if g.Streak != 1 {
		t.Fatalf("Streak = %d, want 1", g.Streak)
	}
}

func TestWrongLetterResetsStreakAndPos(t *testing.T) {
	g := startPlay(t, "casa")
	g.Streak = 3
	g.TypeChar('c')
	g.TypeChar('x')
	if g.Phase != PhaseError {
		t.Fatalf("Phase = %v, want PhaseError", g.Phase)
	}
	if g.Pos != 0 || g.Streak != 0 {
		t.Fatalf("Pos = %d, Streak = %d, want 0, 0", g.Pos, g.Streak)
	}
	if !g.ErrorFlashing() {
		t.Fatal("ErrorFlashing() = false right after mistake, want true")
	}
	// input locked through flash + block
	g.TypeChar('c')
	if g.Pos != 0 {
		t.Fatalf("Pos = %d after typing while locked, want 0", g.Pos)
	}
	// flash ends, block continues
	g.Tick(ErrorFlashTime + 0.01)
	if g.ErrorFlashing() {
		t.Fatal("ErrorFlashing() = true after flash time, want false")
	}
	if g.Phase != PhaseError {
		t.Fatalf("Phase = %v during block second, want PhaseError", g.Phase)
	}
	g.Tick(ErrorBlockTime)
	if g.Phase != PhaseTyping {
		t.Fatalf("Phase = %v after full lock, want PhaseTyping", g.Phase)
	}
}

func TestSuccessPhaseBringsNextWord(t *testing.T) {
	g := startPlay(t, "sol")
	typeWord(g, "sol")
	g.Tick(SuccessTime + 0.01)
	if g.Phase != PhaseTyping {
		t.Fatalf("Phase = %v, want PhaseTyping", g.Phase)
	}
	if g.Word == "sol" || g.Pos != 0 {
		t.Fatalf("Word = %q, Pos = %d: want a fresh word at Pos 0", g.Word, g.Pos)
	}
}

func TestComboMultCapsAtStreakCap(t *testing.T) {
	g := startPlay(t, "casa")
	g.Streak = 99
	if got := g.ComboMult(); got != 1.5 {
		t.Fatalf("ComboMult() = %v with base cap, want 1.5", got)
	}
	g.StreakCapLevel = 1
	if got := g.ComboMult(); got != 2.0 {
		t.Fatalf("ComboMult() = %v with cap level 1, want 2.0", got)
	}
}

func TestCommandModePausesAndRoutes(t *testing.T) {
	g := startPlay(t, "casa")
	g.TypeChar('/')
	if !g.CommandMode || g.CommandBuf != "/" {
		t.Fatalf("CommandMode = %v, CommandBuf = %q, want true, \"/\"", g.CommandMode, g.CommandBuf)
	}
	typeWord(g, "tienda")
	if g.CommandBuf != "/tienda" {
		t.Fatalf("CommandBuf = %q, want \"/tienda\"", g.CommandBuf)
	}
	if g.Pos != 0 {
		t.Fatalf("Pos = %d, command chars leaked into the word", g.Pos)
	}
	g.CommandSubmit()
	if g.Scene != SceneShop {
		t.Fatalf("Scene = %v after /tienda, want SceneShop", g.Scene)
	}
	g.ExitShop()

	// unknown command
	g.TypeChar('/')
	typeWord(g, "nope")
	g.CommandSubmit()
	if g.Scene != ScenePlay || g.CommandErr <= 0 {
		t.Fatalf("Scene = %v, CommandErr = %v: want ScenePlay with error message", g.Scene, g.CommandErr)
	}
}

func TestCommandBackspaceExitsOnSlash(t *testing.T) {
	g := startPlay(t, "casa")
	g.TypeChar('/')
	g.TypeChar('s')
	g.CommandBackspace() // "/s" -> "/"
	if g.CommandBuf != "/" {
		t.Fatalf("CommandBuf = %q, want \"/\"", g.CommandBuf)
	}
	g.CommandBackspace() // deleting the "/" leaves command mode
	if g.CommandMode {
		t.Fatal("CommandMode = true after deleting the slash, want false")
	}
}

func TestIntroTypewriterAndEnterFlow(t *testing.T) {
	g := New()
	g.MenuIndex = 0 // NUEVA PARTIDA (no save yet)
	g.MenuSelect()
	if g.Scene != SceneIntro {
		t.Fatalf("Scene = %v on NUEVA PARTIDA, want SceneIntro", g.Scene)
	}
	if g.IntroVisible() != "" {
		t.Fatalf("IntroVisible() = %q at start, want empty", g.IntroVisible())
	}
	g.Tick(0.1) // 2 chars at 20 cps
	if got := g.IntroVisible(); got != "Bi" {
		t.Fatalf("IntroVisible() = %q after 0.1s, want \"Bi\"", got)
	}
	g.IntroEnter() // completes the half-revealed line
	if !g.IntroLineDone() || g.IntroLine != 0 {
		t.Fatalf("LineDone = %v, IntroLine = %d after first Enter, want true, 0", g.IntroLineDone(), g.IntroLine)
	}
	// step through the remaining lines
	for range len(g.introScript()) - 1 {
		g.IntroEnter() // next line
		g.IntroEnter() // reveal it fully
	}
	g.IntroEnter() // past the last line -> fresh run
	if g.Scene != ScenePlay || !g.RunActive {
		t.Fatalf("Scene = %v, RunActive = %v, want ScenePlay, true", g.Scene, g.RunActive)
	}
}

func TestMenuCommandPausesAndPlayResumes(t *testing.T) {
	g := startPlay(t, "casa")
	g.Score = 123
	g.MultLevel = 2
	g.TypeChar('/')
	typeWord(g, "menu")
	g.CommandSubmit()
	if g.Scene != SceneMenu {
		t.Fatalf("Scene = %v after /menu, want SceneMenu", g.Scene)
	}
	g.MenuIndex = 0 // CONTINUAR resumes (run is active, so it's the first entry)
	if g.MenuItems()[0].ID != "continue" {
		t.Fatalf("MenuItems()[0].ID = %q with an active run, want \"continue\"", g.MenuItems()[0].ID)
	}
	g.MenuSelect()
	// CONTINUAR greets with a one-line HR quip first
	if g.Scene != SceneIntro {
		t.Fatalf("Scene = %v after CONTINUAR, want SceneIntro (return quip)", g.Scene)
	}
	g.IntroEnter() // reveal the quip fully
	g.IntroEnter() // dismiss it
	if g.Scene != ScenePlay || g.Score != 123 || g.MultLevel != 2 {
		t.Fatalf("Scene = %v, Score = %d, MultLevel = %d after resume, want ScenePlay, 123, 2",
			g.Scene, g.Score, g.MultLevel)
	}
}

func TestLanguageChangeKeepsProgress(t *testing.T) {
	g := startPlay(t, "casa")
	g.Score = 77
	g.MenuIndex = 2 // IDIOMA (continue, new, lang, quit)
	g.MenuAdjust(1)
	if g.LangIndex != 1 {
		t.Fatalf("LangIndex = %d, want 1", g.LangIndex)
	}
	if g.Score != 77 {
		t.Fatalf("Score = %d after language change, want 77", g.Score)
	}
	found := false
	for _, w := range WordLists["en"] {
		if w == g.Word {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("Word = %q not from the english list after switching", g.Word)
	}
}

func TestQuitAndOptionsMenuEntries(t *testing.T) {
	g := New()
	items := g.MenuItems()
	if len(items) != 6 { // new, lang, hr, stats, reset, quit — no save, no options support
		t.Fatalf("MenuItems() = %d entries without options, want 6", len(items))
	}
	g.OptionsEnabled = true
	items = g.MenuItems()
	if len(items) != 7 || items[4].ID != "options" || items[5].ID != "reset" {
		t.Fatalf("MenuItems() with options = %d entries, [4] = %q, [5] = %q; want 7, \"options\", \"reset\"",
			len(items), items[4].ID, items[5].ID)
	}
	g.MenuIndex = 6 // SALIR
	g.MenuSelect()
	if !g.QuitRequested {
		t.Fatal("QuitRequested = false after SALIR, want true")
	}
}

func TestSaveLoadRoundtripAndContinue(t *testing.T) {
	path := t.TempDir() + "/save.json"

	g := startPlay(t, "casa")
	g.SavePath = path
	g.Score = 200
	g.MultLevel = 1
	g.StreakCapLevel = 2
	g.LangIndex = 1
	g.Volume = 0.7
	g.DisplayMode = 1
	g.Save()

	// fresh process: load and continue
	g2 := New()
	g2.SavePath = path
	g2.Load()
	if g2.LangIndex != 1 || g2.Volume != 0.7 || g2.DisplayMode != 1 {
		t.Fatalf("settings after Load = lang %d, vol %v, display %d; want 1, 0.7, 1",
			g2.LangIndex, g2.Volume, g2.DisplayMode)
	}
	if !g2.HasSave() {
		t.Fatal("HasSave() = false after loading a run, want true")
	}
	if g2.MenuItems()[0].ID != "continue" {
		t.Fatalf("MenuItems()[0].ID = %q, want \"continue\"", g2.MenuItems()[0].ID)
	}
	g2.MenuIndex = 0
	g2.MenuSelect()
	if g2.Scene != SceneIntro {
		t.Fatalf("Scene = %v after CONTINUAR, want SceneIntro (return quip)", g2.Scene)
	}
	if got := len(g2.introScript()); got != 1 {
		t.Fatalf("introScript() has %d lines on continue, want 1 quip", got)
	}
	g2.IntroEnter()
	g2.IntroEnter()
	if g2.Scene != ScenePlay || g2.Score != 200 || g2.MultLevel != 1 || g2.StreakCapLevel != 2 {
		t.Fatalf("restored run = scene %v, score %d, mult %d, cap %d; want ScenePlay, 200, 1, 2",
			g2.Scene, g2.Score, g2.MultLevel, g2.StreakCapLevel)
	}
}

func TestReturnQuipNeverRepeatsBackToBack(t *testing.T) {
	g := startPlay(t, "casa")
	prev := -1
	for range 30 {
		g.BackToMenu()
		g.MenuIndex = 0 // CONTINUAR
		g.MenuSelect()
		if prev != -1 && g.returnIdx == prev {
			t.Fatalf("returnIdx repeated back-to-back: %d", g.returnIdx)
		}
		prev = g.returnIdx
		g.IntroEnter()
		g.IntroEnter()
	}
	if len(ReturnScripts["es"]) != len(ReturnScripts["en"]) {
		t.Fatalf("return scripts differ in length: es %d, en %d",
			len(ReturnScripts["es"]), len(ReturnScripts["en"]))
	}
}

func TestDictionaryPoolsAreBigAndClean(t *testing.T) {
	for lang, words := range WordLists {
		if len(words) < 10000 {
			t.Fatalf("%s pool = %d words, want a full dictionary", lang, len(words))
		}
		for _, w := range words {
			if len(w) < 3 || len(w) > 14 {
				t.Fatalf("%s word %q outside 3-14 length", lang, w)
			}
			for _, r := range w {
				if r < 'a' || r > 'z' {
					t.Fatalf("%s word %q has non a-z rune %q", lang, w, r)
				}
			}
		}
	}
}

func TestNewGameDiscardsSavedRun(t *testing.T) {
	path := t.TempDir() + "/save.json"

	g := startPlay(t, "casa")
	g.SavePath = path
	g.Score = 500
	g.Save()

	g2 := New()
	g2.SavePath = path
	g2.Load()
	g2.MenuIndex = 1 // NUEVA PARTIDA (continue is first)
	if g2.MenuItems()[1].ID != "new" {
		t.Fatalf("MenuItems()[1].ID = %q, want \"new\"", g2.MenuItems()[1].ID)
	}
	g2.MenuSelect()
	if g2.Scene != SceneIntro {
		t.Fatalf("Scene = %v after NUEVA PARTIDA, want SceneIntro", g2.Scene)
	}
	for !g2.RunActive {
		g2.IntroEnter()
	}
	if g2.Score != 0 || g2.HasSave() && g2.savedRun != nil {
		t.Fatalf("Score = %d, savedRun = %v after new game, want 0 and nil", g2.Score, g2.savedRun)
	}
}

func TestUIStringsFollowLanguage(t *testing.T) {
	g := New()
	if g.UI().MenuNew != "NUEVA PARTIDA" || g.UpgradeName("mult") != "MULTIPLICADOR" {
		t.Fatalf("spanish UI = %q, %q", g.UI().MenuNew, g.UpgradeName("mult"))
	}
	g.LangIndex = 1
	if g.UI().MenuNew != "NEW GAME" || g.UpgradeName("mult") != "MULTIPLIER" {
		t.Fatalf("english UI = %q, %q", g.UI().MenuNew, g.UpgradeName("mult"))
	}
	if len(IntroScripts["es"]) != len(IntroScripts["en"]) {
		t.Fatalf("intro scripts differ in length: es %d, en %d — mid-intro language switch would break",
			len(IntroScripts["es"]), len(IntroScripts["en"]))
	}
}

func TestOptionsAdjust(t *testing.T) {
	g := New()
	g.OptIndex = 0 // display
	g.OptAdjust(1)
	if g.DisplayMode != 1 {
		t.Fatalf("DisplayMode = %d, want 1", g.DisplayMode)
	}
	g.OptIndex = 1 // volume, default 0.4
	g.OptAdjust(1)
	if g.Volume != 0.5 {
		t.Fatalf("Volume = %v, want 0.5", g.Volume)
	}
	for range 10 {
		g.OptAdjust(-1)
	}
	if g.Volume != 0 {
		t.Fatalf("Volume = %v after clamping down, want 0", g.Volume)
	}
	g.OptIndex = 2 // back
	g.Scene = SceneOptions
	g.OptSelect()
	if g.Scene != SceneMenu {
		t.Fatalf("Scene = %v after VOLVER, want SceneMenu", g.Scene)
	}
}

func TestDayQuotaScalesWithDayUpgradesAndHR(t *testing.T) {
	g := startPlay(t, "casa")
	if got := g.computeQuota(); got != 100 {
		t.Fatalf("computeQuota() day 1 = %d, want 100", got)
	}
	g.Day = 3 // 100 × 1.55² = 240.25
	if got := g.computeQuota(); got != 240 {
		t.Fatalf("computeQuota() day 3 = %d, want 240", got)
	}
	g.MultLevel = 2 // +0.70
	g.StreakCapLevel = 1
	if got := g.computeQuota(); got != 444 { // 240.25 × 1.85
		t.Fatalf("computeQuota() with upgrades = %d, want 444", got)
	}
	g.HRQuotaLevel = 5 // −50%
	if got := g.computeQuota(); got != 222 {
		t.Fatalf("computeQuota() with max hrquota = %d, want 222", got)
	}
}

func TestQuotaFreezesDuringDay(t *testing.T) {
	g := startPlay(t, "casa")
	if g.DayQuota() != 100 {
		t.Fatalf("DayQuota() at day start = %d, want 100", g.DayQuota())
	}
	// buying mid-day must not move today's locked quota
	g.Score = 9999
	g.ShopIndex = 0 // mult
	if !g.ShopBuy() {
		t.Fatal("ShopBuy() = false with enough points, want true")
	}
	if g.DayQuota() != 100 {
		t.Fatalf("DayQuota() after mid-day purchase = %d, want 100 (frozen)", g.DayQuota())
	}
	// the payday locks tomorrow's quota WITH the new level
	g.DayTimer = 0.001
	g.Tick(0.01)
	if g.DayEndQuota != 100 {
		t.Fatalf("DayEndQuota = %d, want the frozen 100", g.DayEndQuota)
	}
	want := g.computeQuota() // day 2, mult 1: 155 × 1.35
	if g.DayQuota() != want {
		t.Fatalf("DayQuota() next day = %d, want %d (recomputed with upgrades)", g.DayQuota(), want)
	}
	if g.DayQuota() == 100 {
		t.Fatal("next day's quota still 100: upgrade term never applied")
	}
}

func TestDayTimerRunsOnlyInActivePlay(t *testing.T) {
	g := startPlay(t, "casa")
	if g.Day != 1 || g.DayTimer != DayBaseTime {
		t.Fatalf("Day = %d, DayTimer = %v after initRun, want 1, %v", g.Day, g.DayTimer, DayBaseTime)
	}
	g.Tick(1)
	if g.DayTimer != DayBaseTime-1 {
		t.Fatalf("DayTimer = %v after 1s of play, want %v", g.DayTimer, DayBaseTime-1)
	}
	g.TypeChar('/') // command mode pauses the clock
	g.Tick(5)
	if g.DayTimer != DayBaseTime-1 {
		t.Fatalf("DayTimer = %v after ticking in command mode, want %v", g.DayTimer, DayBaseTime-1)
	}
	g.CommandCancel()
	g.Scene = SceneShop // shop pauses it too
	g.Tick(5)
	if g.DayTimer != DayBaseTime-1 {
		t.Fatalf("DayTimer = %v after ticking in shop, want %v", g.DayTimer, DayBaseTime-1)
	}
}

func TestDayEndPaydayChargesAndAdvances(t *testing.T) {
	g := startPlay(t, "casa")
	g.Score = 130 // covers the 100 quota
	g.DayTimer = 0.01
	g.Tick(0.02)
	if g.Scene != SceneDayEnd || g.DayEndFired {
		t.Fatalf("Scene = %v, fired = %v at day end with funds, want SceneDayEnd payday", g.Scene, g.DayEndFired)
	}
	// day 1 is off-payroll: no HR points yet, only the next-payroll hint
	if g.Score != 30 || g.HRPoints != 0 || g.Day != 2 {
		t.Fatalf("Score = %d, HRPoints = %d, Day = %d, want 30, 0, 2", g.Score, g.HRPoints, g.Day)
	}
	if g.DayTimer != g.DayLength() {
		t.Fatalf("DayTimer = %v after payday, want %v", g.DayTimer, g.DayLength())
	}
	if g.DayEndQuota != 100 || g.DayEndGainHR != 0 || g.DayEndNextPay != HRPayCycle {
		t.Fatalf("summary quota %d, gain %d, next pay %d; want 100, 0, %d",
			g.DayEndQuota, g.DayEndGainHR, g.DayEndNextPay, HRPayCycle)
	}
	if g.DayEndFirst {
		t.Fatal("DayEndFirst = true on day 1, want false (explanation waits for the first payroll)")
	}
	// off-payroll paydays use a single-line quip
	g.DayEndEnter() // reveal it fully
	if !g.DayEndSummary() {
		t.Fatal("DayEndSummary() = false after revealing a single-line quip, want true")
	}
	g.DayEndEnter()
	if g.Scene != SceneDayStart || g.DayStartText == "" {
		t.Fatalf("Scene = %v, DayStartText = %q after payday Enters, want SceneDayStart with a quip", g.Scene, g.DayStartText)
	}
	g.DayStartEnter() // reveal
	g.DayStartEnter() // dismiss
	if g.Scene != ScenePlay || !g.RunActive {
		t.Fatalf("Scene = %v, RunActive = %v after morning Enter, want ScenePlay, true", g.Scene, g.RunActive)
	}
}

func TestHRPayoutSchedule(t *testing.T) {
	cases := map[int]int{1: 0, 2: 0, 3: 1, 4: 0, 5: 0, 6: 3, 9: 4, 12: 6, 15: 7, 18: 9}
	for day, want := range cases {
		if got := hrGainForDay(day); got != want {
			t.Fatalf("hrGainForDay(%d) = %d, want %d", day, got, want)
		}
	}

	// the first payroll (day 3) pays out and triggers the HR explanation
	g := startPlay(t, "casa")
	g.Day = 3
	g.Score = 9999
	g.DayTimer = 0.001
	g.Tick(0.01)
	if g.HRPoints != 1 || g.DayEndGainHR != 1 || g.DayEndNextPay != 0 {
		t.Fatalf("day 3 payroll: points %d, gain %d, next %d; want 1, 1, 0",
			g.HRPoints, g.DayEndGainHR, g.DayEndNextPay)
	}
	if !g.DayEndFirst {
		t.Fatal("DayEndFirst = false on the first payroll, want true (HR points explained)")
	}
	if lines := len(HRExplainScripts["es"]); lines < 2 {
		t.Fatalf("HRExplainScripts has %d lines, want a multi-line script", lines)
	}
	if g.DayEndSummary() {
		t.Fatal("DayEndSummary() = true on line 0 of the explain script, want false")
	}
	for i := 0; g.Scene == SceneDayEnd; i++ {
		if i > 2*len(HRExplainScripts["es"])+2 {
			t.Fatal("DayEndEnter never left SceneDayEnd")
		}
		g.DayEndEnter() // completes each line, then advances
	}
	if g.Scene != SceneDayStart {
		t.Fatalf("Scene = %v after first-payroll Enters, want SceneDayStart", g.Scene)
	}

	// day 6 pays 3 and is no longer "first"
	g.DayStartEnter()
	g.DayStartEnter()
	g.Day = 6
	g.Score = 99999
	g.DayTimer = 0.001
	g.Tick(0.01)
	if g.HRPoints != 4 || g.DayEndGainHR != 3 || g.DayEndFirst {
		t.Fatalf("day 6 payroll: points %d, gain %d, first %v; want 4, 3, false",
			g.HRPoints, g.DayEndGainHR, g.DayEndFirst)
	}
}

func TestDayEndFiringEndsRunKeepsHRPoints(t *testing.T) {
	g := startPlay(t, "casa")
	g.HRPoints = 7
	g.Score = 40 // below the 100 quota
	g.DayTimer = 0.01
	g.Tick(0.02)
	if g.Scene != SceneDayEnd || !g.DayEndFired {
		t.Fatalf("Scene = %v, fired = %v at day end without funds, want SceneDayEnd fired", g.Scene, g.DayEndFired)
	}
	if g.RunActive || g.savedRun != nil {
		t.Fatalf("RunActive = %v, savedRun = %v after firing, want false, nil", g.RunActive, g.savedRun)
	}
	if g.HRPoints != 7 {
		t.Fatalf("HRPoints = %d after firing, want 7 (kept)", g.HRPoints)
	}
	g.DayEndEnter()
	g.DayEndEnter()
	if g.Scene != SceneMenu || g.HasSave() {
		t.Fatalf("Scene = %v, HasSave = %v after firing Enter, want SceneMenu, false", g.Scene, g.HasSave())
	}
}

func TestDayEndWaitsForSuccessPhase(t *testing.T) {
	g := startPlay(t, "sol")
	g.Score = 200
	g.DayTimer = 0.01
	typeWord(g, "sol") // enters PhaseSuccess before the clock runs out
	g.Tick(0.02)       // timer hits 0 mid-success: day must not end yet
	if g.Scene != ScenePlay {
		t.Fatalf("Scene = %v mid-success at timer 0, want ScenePlay", g.Scene)
	}
	g.Tick(SuccessTime) // success finishes, word banked, then the day ends
	if g.Scene != SceneDayEnd {
		t.Fatalf("Scene = %v after success phase, want SceneDayEnd", g.Scene)
	}
	if g.Score != 200+g.LastGain-100 {
		t.Fatalf("Score = %d, want word banked before the 100 quota", g.Score)
	}
}

func TestHRShopBuyAndCaps(t *testing.T) {
	g := New()
	g.HRPoints = 10
	g.HRIndex = 0 // hrmult, base cost 5
	if !g.HRBuy() {
		t.Fatal("HRBuy() = false with enough HR points, want true")
	}
	if g.HRPoints != 5 || g.HRMultLevel != 1 {
		t.Fatalf("HRPoints = %d, HRMultLevel = %d, want 5, 1", g.HRPoints, g.HRMultLevel)
	}
	if g.HRUpgradeCost("hrmult") != 10 {
		t.Fatalf("HRUpgradeCost(hrmult) = %d at level 1, want 10", g.HRUpgradeCost("hrmult"))
	}
	g.HRIndex = 2 // hrquota capped at HRQuotaMax
	g.HRQuotaLevel = HRQuotaMax
	g.HRPoints = 9999
	if g.HRBuy() {
		t.Fatal("HRBuy() = true at hrquota cap, want false")
	}
	if !g.HRUpgradeMaxed("hrquota") {
		t.Fatal("HRUpgradeMaxed(hrquota) = false at cap, want true")
	}
}

func TestHRMetaUpgradesAffectRun(t *testing.T) {
	g := startPlay(t, "casa")
	g.Streak = 0
	base := g.WordGain() // 4 × 1 × 1.0
	g.HRMultLevel = 5    // +50%
	if got := g.WordGain(); got != base+base/2 {
		t.Fatalf("WordGain() = %d with hrmult 5, want %d", got, base+base/2)
	}
	g.HRDayLevel = 2
	if got := g.DayLength(); got != DayBaseTime+2*DayTimeBonus {
		t.Fatalf("DayLength() = %v with hrday 2, want %v", got, DayBaseTime+2*DayTimeBonus)
	}
}

func TestHRCommandAndMenuEntry(t *testing.T) {
	g := startPlay(t, "casa")
	g.TypeChar('/')
	typeWord(g, "rrhh")
	g.CommandSubmit()
	if g.Scene != SceneHR {
		t.Fatalf("Scene = %v after /rrhh, want SceneHR", g.Scene)
	}
	g.ExitHR()
	if g.Scene != ScenePlay {
		t.Fatalf("Scene = %v after ExitHR from play, want ScenePlay", g.Scene)
	}

	g.BackToMenu()
	g.MenuIndex = 3 // RRHH (continue, new, lang, hr)
	if g.MenuItems()[3].ID != "hr" {
		t.Fatalf("MenuItems()[3].ID = %q, want \"hr\"", g.MenuItems()[3].ID)
	}
	g.MenuSelect()
	if g.Scene != SceneHR {
		t.Fatalf("Scene = %v after menu RRHH, want SceneHR", g.Scene)
	}
	g.ExitHR()
	if g.Scene != SceneMenu {
		t.Fatalf("Scene = %v after ExitHR from menu, want SceneMenu", g.Scene)
	}
}

func TestSaveLoadKeepsHRAndDay(t *testing.T) {
	path := t.TempDir() + "/save.json"

	g := startPlay(t, "casa")
	g.SavePath = path
	g.Score = 300
	g.Day = 4
	g.DayTimer = 42.5
	g.HRPoints = 12
	g.HRMultLevel = 1
	g.HRDayLevel = 1
	g.HRQuotaLevel = 2
	g.Save()

	g2 := New()
	g2.SavePath = path
	g2.Load()
	if g2.HRPoints != 12 || g2.HRMultLevel != 1 || g2.HRDayLevel != 1 || g2.HRQuotaLevel != 2 {
		t.Fatalf("HR after Load = %d pts, mult %d, day %d, quota %d; want 12, 1, 1, 2",
			g2.HRPoints, g2.HRMultLevel, g2.HRDayLevel, g2.HRQuotaLevel)
	}
	g2.MenuIndex = 0 // CONTINUAR
	g2.MenuSelect()
	g2.IntroEnter()
	g2.IntroEnter()
	if g2.Day != 4 || g2.DayTimer != 42.5 {
		t.Fatalf("restored day = %d, timer = %v, want 4, 42.5", g2.Day, g2.DayTimer)
	}
}

func TestDayEndQuipsMatchAcrossLanguages(t *testing.T) {
	if len(PaydayScripts["es"]) != len(PaydayScripts["en"]) {
		t.Fatalf("payday scripts differ in length: es %d, en %d",
			len(PaydayScripts["es"]), len(PaydayScripts["en"]))
	}
	if len(FiredScripts["es"]) != len(FiredScripts["en"]) {
		t.Fatalf("fired scripts differ in length: es %d, en %d",
			len(FiredScripts["es"]), len(FiredScripts["en"]))
	}
	if len(DayStartScripts["es"]) != len(DayStartScripts["en"]) {
		t.Fatalf("day-start scripts differ in length: es %d, en %d",
			len(DayStartScripts["es"]), len(DayStartScripts["en"]))
	}
	if len(DayStartStatScripts["es"]) != len(DayStartStatScripts["en"]) {
		t.Fatalf("day-start stat scripts differ in length: es %d, en %d",
			len(DayStartStatScripts["es"]), len(DayStartStatScripts["en"]))
	}
	if len(HRExplainScripts["es"]) != len(HRExplainScripts["en"]) || len(HRExplainScripts["es"]) == 0 {
		t.Fatalf("HR explain scripts differ in length: es %d, en %d",
			len(HRExplainScripts["es"]), len(HRExplainScripts["en"]))
	}
	if len(InternHiredScripts["es"]) != len(InternHiredScripts["en"]) || len(InternHiredScripts["es"]) == 0 {
		t.Fatalf("intern hired scripts differ in length: es %d, en %d",
			len(InternHiredScripts["es"]), len(InternHiredScripts["en"]))
	}
	if len(EventCoffeeScripts["es"]) != len(EventCoffeeScripts["en"]) || len(EventCoffeeScripts["es"]) == 0 {
		t.Fatalf("coffee event scripts differ in length: es %d, en %d",
			len(EventCoffeeScripts["es"]), len(EventCoffeeScripts["en"]))
	}
	if len(EventInspectionScripts["es"]) != len(EventInspectionScripts["en"]) || len(EventInspectionScripts["es"]) == 0 {
		t.Fatalf("inspection event scripts differ in length: es %d, en %d",
			len(EventInspectionScripts["es"]), len(EventInspectionScripts["en"]))
	}
}

func TestStatsCountWordsAndFails(t *testing.T) {
	g := startPlay(t, "sol")
	typeWord(g, "sol")
	g.Tick(SuccessTime + 0.01)
	g.Word, g.Pos = "mar", 0 // force a known word
	g.TypeChar('x')          // fail on "mar"
	g.Tick(ErrorFlashTime + ErrorBlockTime + 0.01)
	g.TypeChar('x') // fail on "mar" again
	if g.DayWords != 1 || g.DayFails != 2 {
		t.Fatalf("DayWords = %d, DayFails = %d, want 1, 2", g.DayWords, g.DayFails)
	}
	if g.TotalWords != 1 || g.TotalFails != 2 {
		t.Fatalf("TotalWords = %d, TotalFails = %d, want 1, 2", g.TotalWords, g.TotalFails)
	}
	if w, n := worstWord(g.dayFailCounts); w != "mar" || n != 2 {
		t.Fatalf("worstWord = %q (%d), want \"mar\" (2)", w, n)
	}

	// day end banks the stats and resets the day counters
	g.Tick(ErrorFlashTime + ErrorBlockTime + 0.01)
	g.Score = 999
	g.DayTimer = 0.001
	g.Tick(0.01)
	if g.LastDayWords != 1 || g.LastDayFails != 2 || g.LastDayWorst != "mar" || g.LastDayWorstN != 2 {
		t.Fatalf("LastDay = %d words, %d fails, worst %q (%d); want 1, 2, \"mar\", 2",
			g.LastDayWords, g.LastDayFails, g.LastDayWorst, g.LastDayWorstN)
	}
	if g.DayWords != 0 || g.DayFails != 0 || len(g.dayFailCounts) != 0 {
		t.Fatalf("day counters not reset after payday: %d words, %d fails, %d map entries",
			g.DayWords, g.DayFails, len(g.dayFailCounts))
	}
	if g.TotalWords != 1 || g.TotalFails != 2 {
		t.Fatalf("global stats changed at day end: %d words, %d fails, want 1, 2", g.TotalWords, g.TotalFails)
	}
}

func TestTopFailWords(t *testing.T) {
	g := New()
	g.FailCounts = map[string]int{"casa": 5, "sol": 2, "mar": 5, "pan": 1}
	top := g.TopFailWords(3)
	if len(top) != 3 {
		t.Fatalf("TopFailWords(3) = %d entries, want 3", len(top))
	}
	// ties broken alphabetically: casa(5), mar(5), sol(2)
	if top[0].Word != "casa" || top[1].Word != "mar" || top[2].Word != "sol" {
		t.Fatalf("TopFailWords order = %q, %q, %q; want casa, mar, sol", top[0].Word, top[1].Word, top[2].Word)
	}
}

func TestGlobalStatsSurviveFiringAndReload(t *testing.T) {
	path := t.TempDir() + "/save.json"

	g := startPlay(t, "casa")
	g.SavePath = path
	g.TypeChar('x') // one global fail on "casa"
	g.Tick(ErrorFlashTime + ErrorBlockTime + 0.01)
	g.Score = 0 // guarantee a firing
	g.DayTimer = 0.001
	g.Tick(0.01)
	if !g.DayEndFired {
		t.Fatal("expected a firing")
	}

	g2 := New()
	g2.SavePath = path
	g2.Load()
	if g2.TotalFails != 1 || g2.FailCounts["casa"] != 1 {
		t.Fatalf("global stats after firing+reload = %d fails, casa %d; want 1, 1",
			g2.TotalFails, g2.FailCounts["casa"])
	}
	if g2.HasSave() {
		t.Fatal("HasSave() = true after a firing, want false (run erased)")
	}
}

func TestDayStartStatQuipFillsPlaceholders(t *testing.T) {
	g := startPlay(t, "casa")
	g.LastDayWords = 12
	g.LastDayFails = 3
	g.LastDayWorst = "burocracia"
	g.LastDayWorstN = 2
	for range 50 {
		g.composeDayStart()
		if strings.Contains(g.DayStartText, "{") {
			t.Fatalf("unfilled placeholder in day-start quip: %q", g.DayStartText)
		}
	}
}

func TestGoldenWordPaysMultiplier(t *testing.T) {
	g := startPlay(t, "casa")
	g.WordGolden = true
	typeWord(g, "casa")
	// streak 1 → combo x1.1; 4 × 1 × 1.1 × 3 = 13.2 → 13
	if g.Score != 13 || g.LastGain != 13 {
		t.Fatalf("Score = %d, LastGain = %d for a golden word, want 13, 13", g.Score, g.LastGain)
	}
}

func TestGoldenScalesWithLevel(t *testing.T) {
	g := startPlay(t, "casa")
	if g.GoldenChance() != GoldenBaseChance || g.GoldenMult() != GoldenBaseMult {
		t.Fatalf("level 0 = %v chance, x%v; want %v, x%v",
			g.GoldenChance(), g.GoldenMult(), GoldenBaseChance, GoldenBaseMult)
	}
	g.GoldenLevel = 1
	if g.GoldenChance() != 0.03 || g.GoldenMult() != 3.5 {
		t.Fatalf("level 1 = %v chance, x%v; want 0.03, x3.5", g.GoldenChance(), g.GoldenMult())
	}
	g.GoldenLevel = GoldenMaxLevel
	if g.GoldenChance() != 0.10 || g.GoldenMult() != 7.0 {
		t.Fatalf("level %d = %v chance, x%v; want 0.10, x7.0",
			GoldenMaxLevel, g.GoldenChance(), g.GoldenMult())
	}
}

func TestGoldenShopBuyCostAndCap(t *testing.T) {
	g := startPlay(t, "casa")
	g.Score = 9999
	g.ShopIndex = 2 // golden, base cost 250
	if !g.ShopBuy() {
		t.Fatal("ShopBuy() = false with enough points, want true")
	}
	if g.GoldenLevel != 1 || g.Score != 9999-250 {
		t.Fatalf("GoldenLevel = %d, Score = %d, want 1, %d", g.GoldenLevel, g.Score, 9999-250)
	}
	if g.UpgradeCost("golden") != 750 {
		t.Fatalf("UpgradeCost(golden) = %d at level 1, want 750", g.UpgradeCost("golden"))
	}
	g.GoldenLevel = GoldenMaxLevel
	g.Score = 1 << 30
	if !g.UpgradeMaxed("golden") {
		t.Fatal("UpgradeMaxed(golden) = false at cap, want true")
	}
	if g.ShopBuy() {
		t.Fatal("ShopBuy() = true at golden cap, want false")
	}
}

func TestGoldenLevelPersistsAcrossSaveLoad(t *testing.T) {
	path := t.TempDir() + "/save.json"

	g := startPlay(t, "casa")
	g.SavePath = path
	g.GoldenLevel = 3
	g.Save()

	g2 := New()
	g2.SavePath = path
	g2.Load()
	g2.MenuIndex = 0 // CONTINUAR
	g2.MenuSelect()
	g2.IntroEnter()
	g2.IntroEnter()
	if g2.GoldenLevel != 3 {
		t.Fatalf("GoldenLevel after reload = %d, want 3", g2.GoldenLevel)
	}
}

func TestQuotaIncludesGoldenTerm(t *testing.T) {
	g := startPlay(t, "casa")
	g.GoldenLevel = 4 // 100 × (1 + 0.08×4) = 132
	if got := g.computeQuota(); got != 132 {
		t.Fatalf("computeQuota() with golden 4 = %d, want 132", got)
	}
}

// internIndex is the shop index of the "intern" entry within ShopItems().
func internIndex(t *testing.T, g *Game) int {
	t.Helper()
	for i, up := range g.ShopItems() {
		if up.ID == "intern" {
			return i
		}
	}
	t.Fatal("intern entry not found in ShopItems()")
	return -1
}

func TestInternBuyAddsInternAndCapsPerRun(t *testing.T) {
	g := startPlay(t, "casa")
	g.Score = 9999
	g.ShopIndex = internIndex(t, g)
	if !g.ShopBuy() {
		t.Fatal("ShopBuy() = false with enough points, want true")
	}
	if len(g.Interns) != 1 || g.Score != 9999-InternBaseCost {
		t.Fatalf("interns = %d, Score = %d, want 1, %d", len(g.Interns), g.Score, 9999-InternBaseCost)
	}
	if g.Interns[0].Word == "" {
		t.Fatal("hired intern has no word")
	}
	if g.ShopQuipText == "" {
		t.Fatal("ShopQuipText empty after hiring, want an HR quip")
	}
	// default headcount is 1: a second intern can't be hired
	if !g.UpgradeMaxed("intern") || g.ShopBuy() {
		t.Fatal("second intern bought at MaxInterns()=1, want refusal")
	}
	// the meta upgrade lifts the cap; the second intern costs ×3
	g.HRInternLevel = 1
	if g.UpgradeMaxed("intern") {
		t.Fatal("UpgradeMaxed(intern) = true with hrintern 1 and 1 intern, want false")
	}
	if got := g.UpgradeCost("intern"); got != InternBaseCost*3 {
		t.Fatalf("UpgradeCost(intern) with 1 owned = %d, want %d", got, InternBaseCost*3)
	}
	if !g.ShopBuy() {
		t.Fatal("ShopBuy() = false for the second intern with the cap lifted, want true")
	}
	if len(g.Interns) != 2 {
		t.Fatalf("interns = %d after second hire, want 2", len(g.Interns))
	}
	g.ExitShop()
	if g.ShopQuipText != "" {
		t.Fatal("ShopQuipText kept after leaving the shop, want cleared")
	}
}

func TestInternTypesAndEarnsWithoutCombo(t *testing.T) {
	g := startPlay(t, "casa")
	g.Streak = 99 // combo must NOT apply to intern gains
	g.Interns = []Intern{{Word: "sol"}}
	g.Tick(2.0) // 1.5 cps × 2s = 3 letters: completes "sol"
	if g.Score != 3 {
		t.Fatalf("Score = %d after intern word, want 3 (len × 1 × 1, no combo)", g.Score)
	}
	if g.DayWords != 1 || g.TotalWords != 1 {
		t.Fatalf("DayWords = %d, TotalWords = %d, want 1, 1", g.DayWords, g.TotalWords)
	}
	in := g.Interns[0]
	if in.Word == "sol" || in.Pos != 0 {
		t.Fatalf("intern word = %q, pos = %d after completing, want a fresh word at 0", in.Word, in.Pos)
	}
	if in.LastGain != 3 || in.GainTimer <= 0 {
		t.Fatalf("LastGain = %d, GainTimer = %v, want 3 and a running flash", in.LastGain, in.GainTimer)
	}
	// the post-word rest holds typing before the next word starts
	g.Tick(InternWordRest / 2)
	if g.Interns[0].Pos != 0 {
		t.Fatalf("intern typed during its rest, pos = %d", g.Interns[0].Pos)
	}
}

func TestInternSpeedUpgradeSpeedsTyping(t *testing.T) {
	g := startPlay(t, "casa")
	if g.InternCPS() != InternBaseCPS {
		t.Fatalf("InternCPS() = %v at level 0, want %v", g.InternCPS(), InternBaseCPS)
	}
	g.InternSpeedLevel = 2
	if g.InternCPS() != InternBaseCPS+2*InternCPSPerLevel {
		t.Fatalf("InternCPS() = %v at level 2, want %v", g.InternCPS(), InternBaseCPS+2*InternCPSPerLevel)
	}
}

func TestInternGainUsesMultAndHRMultOnly(t *testing.T) {
	g := startPlay(t, "casa")
	g.MultLevel = 1                           // ×2
	g.HRMultLevel = 5                         // ×1.5
	g.Streak = 99                             // must be ignored
	if got := g.InternGain("sol"); got != 9 { // 3 × 2 × 1.5
		t.Fatalf("InternGain(sol) = %d, want 9", got)
	}
}

func TestInternPausesWithDayClock(t *testing.T) {
	g := startPlay(t, "casa")
	g.Interns = []Intern{{Word: "burocracia"}}
	g.TypeChar('/') // command mode pauses interns
	g.Tick(5)
	if g.Interns[0].Pos != 0 {
		t.Fatalf("intern pos = %d after ticking in command mode, want 0", g.Interns[0].Pos)
	}
	g.CommandCancel()
	g.Scene = SceneShop // the shop pauses them too
	g.Tick(5)
	if g.Interns[0].Pos != 0 {
		t.Fatalf("intern pos = %d after ticking in the shop, want 0", g.Interns[0].Pos)
	}
}

func TestInternPersistence(t *testing.T) {
	path := t.TempDir() + "/save.json"

	g := startPlay(t, "casa")
	g.SavePath = path
	g.HRInternLevel = 2
	g.Interns = []Intern{{Word: "sol", Pos: 2}, {Word: "mar"}}
	g.InternSpeedLevel = 3
	g.Save()

	g2 := New()
	g2.SavePath = path
	g2.Load()
	if g2.HRInternLevel != 2 {
		t.Fatalf("HRInternLevel after Load = %d, want 2", g2.HRInternLevel)
	}
	g2.MenuIndex = 0 // CONTINUAR
	g2.MenuSelect()
	g2.IntroEnter()
	g2.IntroEnter()
	if len(g2.Interns) != 2 || g2.InternSpeedLevel != 3 {
		t.Fatalf("restored = %d interns, speed %d; want 2, 3", len(g2.Interns), g2.InternSpeedLevel)
	}
	// words/progress are rebuilt fresh, like the player word
	for i, in := range g2.Interns {
		if in.Word == "" || in.Pos != 0 {
			t.Fatalf("intern %d = %q at pos %d, want a fresh word at 0", i, in.Word, in.Pos)
		}
	}
}

func TestQuotaIncludesInternTerms(t *testing.T) {
	g := startPlay(t, "casa")
	g.Interns = []Intern{{Word: "sol"}, {Word: "mar"}}
	g.InternSpeedLevel = 1 // 100 × (1 + 0.25×2 + 0.08) = 158
	if got := g.computeQuota(); got != 158 {
		t.Fatalf("computeQuota() with 2 interns + speed 1 = %d, want 158", got)
	}
}

func TestHRInternBuyAndCap(t *testing.T) {
	g := New()
	g.HRPoints = 30
	g.HRIndex = 3 // hrintern, base cost 25
	if HRUpgrades[g.HRIndex].ID != "hrintern" {
		t.Fatalf("HRUpgrades[3].ID = %q, want \"hrintern\"", HRUpgrades[3].ID)
	}
	if !g.HRBuy() {
		t.Fatal("HRBuy() = false with enough HR points, want true")
	}
	if g.HRInternLevel != 1 || g.HRPoints != 5 || g.MaxInterns() != 2 {
		t.Fatalf("level %d, points %d, max %d; want 1, 5, 2", g.HRInternLevel, g.HRPoints, g.MaxInterns())
	}
	g.HRInternLevel = HRInternMax
	g.HRPoints = 9999
	if !g.HRUpgradeMaxed("hrintern") || g.HRBuy() {
		t.Fatal("hrintern bought past its cap, want refusal")
	}
}

func TestShopItemsHideInternSpeedUntilOwned(t *testing.T) {
	g := startPlay(t, "casa")
	for _, up := range g.ShopItems() {
		if up.ID == "internspeed" {
			t.Fatal("internspeed listed with no intern hired")
		}
	}
	g.Interns = []Intern{{Word: "sol"}}
	found := false
	for _, up := range g.ShopItems() {
		if up.ID == "internspeed" {
			found = true
		}
	}
	if !found {
		t.Fatal("internspeed not listed with an intern hired")
	}
}

func TestFrenzyDoublesPlayerAndInternGain(t *testing.T) {
	g := startPlay(t, "casa")
	base := g.WordGain()          // streak 0: 4
	baseIn := g.InternGain("sol") // 3
	g.startFrenzy()
	if got := g.WordGain(); got != base*2 {
		t.Fatalf("WordGain() in frenzy = %d, want %d", got, base*2)
	}
	if got := g.InternGain("sol"); got != baseIn*2 {
		t.Fatalf("InternGain() in frenzy = %d, want %d", got, baseIn*2)
	}
	if g.EventToast == "" || g.EventsToday != 1 {
		t.Fatalf("toast = %q, EventsToday = %d after frenzy, want announce + 1", g.EventToast, g.EventsToday)
	}
	g.Tick(FrenzyDuration + 0.01)
	if g.FrenzyActive() {
		t.Fatal("FrenzyActive() = true after its duration, want false")
	}
	if got := g.WordGain(); got != base {
		t.Fatalf("WordGain() after frenzy = %d, want %d", got, base)
	}
}

func TestGoldenStacksWithFrenzy(t *testing.T) {
	g := startPlay(t, "casa")
	g.WordGolden = true
	g.startFrenzy()
	if got := g.WordGain(); got != 24 { // 4 × 3 (golden) × 2 (frenzy)
		t.Fatalf("WordGain() golden+frenzy = %d, want 24", got)
	}
}

func TestInspectionArmsAtWordBoundaryAndPaysBonus(t *testing.T) {
	g := startPlay(t, "sol")
	g.armInspection()
	if !g.InspectionArmed || g.InspectionActive() {
		t.Fatalf("armed = %v, active = %v right after arming, want true, false (word in progress untouched)",
			g.InspectionArmed, g.InspectionActive())
	}
	typeWord(g, "sol")   // current word pays normal
	if g.LastGain != 3 { // 3 × 1.1 = 3.3 → 3
		t.Fatalf("LastGain = %d for the word in progress, want 3 (no bonus)", g.LastGain)
	}
	g.Tick(SuccessTime + 0.01) // word boundary: the next word turns urgent
	if g.InspectionArmed || !g.InspectionActive() || g.WordGolden {
		t.Fatalf("armed = %v, active = %v, golden = %v at the boundary, want false, true, false",
			g.InspectionArmed, g.InspectionActive(), g.WordGolden)
	}
	g.Word, g.Pos = "mar", 0 // force a known urgent word
	typeWord(g, "mar")
	// streak 2 → combo 1.2; 3 × 1.2 × 3 (inspection) = 10.8 → 11
	if g.LastGain != 11 {
		t.Fatalf("LastGain = %d for the urgent word, want 11", g.LastGain)
	}
	if g.InspectionActive() {
		t.Fatal("InspectionActive() = true after delivering the urgent word, want false")
	}
}

func TestInspectionExpiryIsHarmless(t *testing.T) {
	g := startPlay(t, "mar")
	g.InspectionTimer = 0.5
	g.TypeChar('m')
	g.TypeChar('a')
	g.Tick(1.0) // the inspection expires mid-word
	if g.InspectionActive() {
		t.Fatal("InspectionActive() = true after expiry, want false")
	}
	if g.Pos != 2 || g.Streak != 0 || g.Phase != PhaseTyping {
		t.Fatalf("Pos = %d, Streak = %d, Phase = %v after expiry, want progress untouched", g.Pos, g.Streak, g.Phase)
	}
	g.TypeChar('r')      // completing it pays normal, no punishment anywhere
	if g.LastGain != 3 { // 3 × 1.1 = 3.3 → 3
		t.Fatalf("LastGain = %d after expiry, want 3 (normal pay)", g.LastGain)
	}
}

func TestEventGuardsSpacingAndDailyCap(t *testing.T) {
	g := startPlay(t, "casa")
	g.EventsToday = EventMaxPerDay
	g.eventSpacing = 999
	for range 500 {
		g.maybeTriggerEvent()
	}
	if g.FrenzyActive() || g.InspectionArmed || g.EventsToday != EventMaxPerDay {
		t.Fatal("an event triggered past the daily cap")
	}
	g.EventsToday = 0
	g.eventSpacing = EventMinSpacing - 1
	for range 500 {
		g.maybeTriggerEvent()
	}
	if g.EventsToday != 0 {
		t.Fatal("an event triggered inside the min spacing window")
	}
	g.eventSpacing = EventMinSpacing
	for i := 0; g.EventsToday == 0; i++ {
		if i > 100000 {
			t.Fatal("no event ever triggered with every guard satisfied")
		}
		g.maybeTriggerEvent()
	}
	if !g.FrenzyActive() && !g.InspectionArmed {
		t.Fatal("EventsToday bumped but no event is active/armed")
	}
	if g.eventSpacing != 0 {
		t.Fatalf("eventSpacing = %v after an event, want 0", g.eventSpacing)
	}
}

func TestEventTimersPauseInCommandMode(t *testing.T) {
	g := startPlay(t, "casa")
	g.startFrenzy()
	g.TypeChar('/')
	g.Tick(5)
	if g.FrenzyTimer != FrenzyDuration {
		t.Fatalf("FrenzyTimer = %v after ticking in command mode, want %v", g.FrenzyTimer, FrenzyDuration)
	}
}

func TestEventsClearAtDayEnd(t *testing.T) {
	g := startPlay(t, "casa")
	g.startFrenzy()
	g.Score = 999
	g.DayTimer = 0.001
	g.Tick(0.01)
	if g.Scene != SceneDayEnd {
		t.Fatalf("Scene = %v, want SceneDayEnd", g.Scene)
	}
	if g.FrenzyActive() || g.EventsToday != 0 || g.EventToast != "" {
		t.Fatalf("frenzy = %v, EventsToday = %d, toast = %q after day end, want all cleared",
			g.FrenzyActive(), g.EventsToday, g.EventToast)
	}
}

func TestEventsTodayPersists(t *testing.T) {
	path := t.TempDir() + "/save.json"

	g := startPlay(t, "casa")
	g.SavePath = path
	g.EventsToday = 2
	g.Save()

	g2 := New()
	g2.SavePath = path
	g2.Load()
	g2.MenuIndex = 0 // CONTINUAR
	g2.MenuSelect()
	g2.IntroEnter()
	g2.IntroEnter()
	if g2.EventsToday != 2 {
		t.Fatalf("EventsToday after reload = %d, want 2 (no save-quit event farming)", g2.EventsToday)
	}
}

func TestEndDayCommandNeedsUnlock(t *testing.T) {
	g := startPlay(t, "casa")
	g.Score = 999
	g.TypeChar('/')
	typeWord(g, "terminar")
	g.CommandSubmit()
	if g.Scene != ScenePlay || g.CommandErr <= 0 {
		t.Fatalf("Scene = %v, CommandErr = %v without hrend, want ScenePlay + error", g.Scene, g.CommandErr)
	}
	g.HREndLevel = 1
	g.TypeChar('/')
	typeWord(g, "terminar")
	g.CommandSubmit()
	if g.Scene != SceneDayEnd || g.DayEndFired {
		t.Fatalf("Scene = %v, fired = %v after /terminar with funds, want payday", g.Scene, g.DayEndFired)
	}
	if g.Day != 2 {
		t.Fatalf("Day = %d after /terminar payday, want 2", g.Day)
	}
}

func TestEndDayCommandCanFire(t *testing.T) {
	g := startPlay(t, "casa")
	g.HREndLevel = 1
	g.Score = 10 // below the 100 quota: /terminar is a real gamble
	g.TypeChar('/')
	typeWord(g, "endday")
	g.CommandSubmit()
	if g.Scene != SceneDayEnd || !g.DayEndFired {
		t.Fatalf("Scene = %v, fired = %v after a short /endday, want fired", g.Scene, g.DayEndFired)
	}
}

func TestHREndBuyAndCap(t *testing.T) {
	g := New()
	g.HRPoints = 40
	g.HRIndex = 4 // hrend, base cost 40
	if HRUpgrades[g.HRIndex].ID != "hrend" {
		t.Fatalf("HRUpgrades[4].ID = %q, want \"hrend\"", HRUpgrades[4].ID)
	}
	if !g.HRBuy() {
		t.Fatal("HRBuy() = false with enough HR points, want true")
	}
	if g.HREndLevel != 1 || g.HRPoints != 0 {
		t.Fatalf("HREndLevel = %d, HRPoints = %d, want 1, 0", g.HREndLevel, g.HRPoints)
	}
	g.HRPoints = 9999
	if !g.HRUpgradeMaxed("hrend") || g.HRBuy() {
		t.Fatal("hrend bought past its single level, want refusal")
	}
}

func TestResetNeedsConfirmAndWipes(t *testing.T) {
	path := t.TempDir() + "/save.json"

	g := startPlay(t, "casa")
	g.SavePath = path
	g.Score = 500
	g.HRPoints = 42
	g.HRMultLevel = 2
	g.TotalWords = 100
	g.BestStreak = 9
	g.BestDay = 6
	g.FailCounts["casa"] = 3
	g.Scene = SceneMenu

	// select RESET: first Enter arms, moving away disarms
	items := g.MenuItems()
	for i, it := range items {
		if it.ID == "reset" {
			g.MenuIndex = i
		}
	}
	g.MenuSelect()
	if !g.resetArmed || g.HRPoints != 42 {
		t.Fatalf("armed = %v, HRPoints = %d after first Enter, want armed and untouched", g.resetArmed, g.HRPoints)
	}
	if g.menuItem().Label != g.UI().ResetConfirm {
		t.Fatalf("armed label = %q, want %q", g.menuItem().Label, g.UI().ResetConfirm)
	}
	g.MenuMove(1)
	if g.resetArmed {
		t.Fatal("resetArmed = true after moving away, want disarmed")
	}

	// arm again and confirm
	g.MenuMove(-1)
	g.MenuSelect()
	g.MenuSelect()
	if g.HRPoints != 0 || g.HRMultLevel != 0 || g.TotalWords != 0 ||
		g.BestStreak != 0 || g.BestDay != 0 || len(g.FailCounts) != 0 {
		t.Fatal("progress survived the confirmed reset")
	}
	if g.RunActive || g.HasSave() {
		t.Fatalf("RunActive = %v, HasSave = %v after reset, want false, false", g.RunActive, g.HasSave())
	}

	// the wiped state hit the disk too
	g2 := New()
	g2.SavePath = path
	g2.Load()
	if g2.HRPoints != 0 || g2.TotalWords != 0 || g2.HasSave() {
		t.Fatalf("reloaded after reset = %d pts, %d words, save %v; want all zero",
			g2.HRPoints, g2.TotalWords, g2.HasSave())
	}
}

func TestEventToastVisibleForWholeEvent(t *testing.T) {
	g := startPlay(t, "casa")
	g.startFrenzy()
	if !g.EventToastVisible() || g.EventToast == "" {
		t.Fatal("toast not visible during an active frenzy")
	}
	g.Tick(FrenzyDuration - 1)
	if !g.EventToastVisible() {
		t.Fatal("toast gone before the frenzy ended")
	}
	g.Tick(1.01)
	if g.EventToastVisible() {
		t.Fatal("toast still visible after the frenzy expired")
	}

	// inspection: visible from armed through the urgent word
	g.armInspection()
	if !g.EventToastVisible() {
		t.Fatal("toast not visible while the inspection is armed")
	}
	typeWord(g, g.Word)
	g.Tick(SuccessTime + 0.01) // boundary: the next word turns urgent
	if !g.InspectionActive() || !g.EventToastVisible() {
		t.Fatal("toast not visible during the urgent word")
	}
	g.Word, g.Pos = "sol", 0
	typeWord(g, "sol") // delivered: inspection over
	if g.EventToastVisible() {
		t.Fatal("toast still visible after delivering the urgent word")
	}
}

func TestBestStreakTrackedAndBanked(t *testing.T) {
	g := startPlay(t, "sol")
	g.Streak = 4
	typeWord(g, "sol") // streak 5
	if g.DayBestStreak != 5 || g.BestStreak != 5 {
		t.Fatalf("DayBestStreak = %d, BestStreak = %d, want 5, 5", g.DayBestStreak, g.BestStreak)
	}
	// a fail resets the streak but not the records
	g.Tick(SuccessTime + 0.01)
	g.TypeChar('!')
	if g.DayBestStreak != 5 || g.BestStreak != 5 {
		t.Fatalf("records dropped after a fail: day %d, global %d", g.DayBestStreak, g.BestStreak)
	}
	// day end banks the day record and resets the day counter
	g.Tick(ErrorFlashTime + ErrorBlockTime + 0.01)
	g.Score = 999
	g.DayTimer = 0.001
	g.Tick(0.01)
	if g.LastDayBestStreak != 5 || g.DayBestStreak != 0 {
		t.Fatalf("LastDayBestStreak = %d, DayBestStreak = %d after payday, want 5, 0",
			g.LastDayBestStreak, g.DayBestStreak)
	}
	if g.BestStreak != 5 {
		t.Fatalf("BestStreak = %d after payday, want 5 (global record kept)", g.BestStreak)
	}
}

func TestBestStreakAndBestDaySurviveFiringAndReload(t *testing.T) {
	path := t.TempDir() + "/save.json"

	g := startPlay(t, "casa")
	g.SavePath = path
	g.BestStreak = 12
	g.Day = 7
	g.Score = 999
	g.DayTimer = 0.001
	g.Tick(0.01) // payday: day 8 reached
	if g.BestDay != 8 {
		t.Fatalf("BestDay = %d after reaching day 8, want 8", g.BestDay)
	}
	g.Score = 0
	g.DayTimer = 0.001
	g.Phase = PhaseTyping
	g.Scene = ScenePlay
	g.Tick(0.01) // fired
	if !g.DayEndFired {
		t.Fatal("expected a firing")
	}

	g2 := New()
	g2.SavePath = path
	g2.Load()
	if g2.BestStreak != 12 || g2.BestDay != 8 {
		t.Fatalf("after firing+reload: BestStreak = %d, BestDay = %d; want 12, 8", g2.BestStreak, g2.BestDay)
	}
}

func TestEveryPurchaseSetsAQuip(t *testing.T) {
	g := startPlay(t, "casa")
	buy := func(id string) {
		t.Helper()
		g.Score = 1 << 30
		idx := -1
		for i, up := range g.ShopItems() {
			if up.ID == id {
				idx = i
			}
		}
		if idx < 0 {
			t.Fatalf("%q not listed in ShopItems()", id)
		}
		g.ShopIndex = idx
		g.ShopQuipText = ""
		if !g.ShopBuy() {
			t.Fatalf("ShopBuy(%s) = false with infinite points", id)
		}
		if g.ShopQuipText == "" {
			t.Fatalf("no quip after buying %q", id)
		}
	}
	// intern before internspeed: the speed entry hides until one exists
	for _, id := range []string{"mult", "streakcap", "golden", "intern", "internspeed"} {
		buy(id)
	}

	g.HRPoints = 1 << 30
	for i, up := range HRUpgrades {
		if g.HRUpgradeMaxed(up.ID) {
			continue
		}
		g.HRIndex = i
		g.HRQuipText = ""
		if !g.HRBuy() {
			t.Fatalf("HRBuy(%s) = false with infinite points", up.ID)
		}
		if g.HRQuipText == "" {
			t.Fatalf("no quip after buying %q", up.ID)
		}
	}
	g.ExitHR()
	if g.HRQuipText != "" {
		t.Fatal("HRQuipText kept after leaving the HR office, want cleared")
	}
}

func TestUpgradeBoughtScriptsCoverEveryUpgrade(t *testing.T) {
	ids := []string{"mult", "streakcap", "golden", "internspeed",
		"hrmult", "hrday", "hrquota", "hrintern", "hrend"}
	for _, id := range ids {
		pool := UpgradeBoughtScripts[id]
		if len(pool["es"]) == 0 || len(pool["es"]) != len(pool["en"]) {
			t.Fatalf("%s purchase scripts: es %d, en %d; want equal and non-empty",
				id, len(pool["es"]), len(pool["en"]))
		}
	}
}

func TestShopBuyAppliesCostAndEffect(t *testing.T) {
	g := startPlay(t, "casa")
	g.Score = 80
	g.ShopIndex = 0 // mult, base cost 75
	if !g.ShopBuy() {
		t.Fatal("ShopBuy() = false with enough points, want true")
	}
	if g.Score != 5 || g.MultLevel != 1 {
		t.Fatalf("Score = %d, MultLevel = %d, want 5, 1", g.Score, g.MultLevel)
	}
	if g.UpgradeCost("mult") != 225 {
		t.Fatalf("UpgradeCost(mult) = %d at level 1, want 225", g.UpgradeCost("mult"))
	}
	if g.ShopBuy() {
		t.Fatal("ShopBuy() = true without enough points, want false")
	}
}
