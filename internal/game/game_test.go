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
	if len(items) != 5 { // new, lang, hr, stats, quit — no save, no options support
		t.Fatalf("MenuItems() = %d entries without options, want 5", len(items))
	}
	g.OptionsEnabled = true
	items = g.MenuItems()
	if len(items) != 6 || items[4].ID != "options" {
		t.Fatalf("MenuItems() with options = %d entries, [4] = %q, want 6, \"options\"", len(items), items[4].ID)
	}
	g.MenuIndex = 5 // SALIR
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
	if got := g.DayQuota(); got != 100 {
		t.Fatalf("DayQuota() day 1 = %d, want 100", got)
	}
	g.Day = 3 // 100 × 1.5² = 225
	if got := g.DayQuota(); got != 225 {
		t.Fatalf("DayQuota() day 3 = %d, want 225", got)
	}
	g.MultLevel = 2 // ×1.5
	g.StreakCapLevel = 1
	if got := g.DayQuota(); got != 360 { // 225 × 1.6
		t.Fatalf("DayQuota() with upgrades = %d, want 360", got)
	}
	g.HRQuotaLevel = 5 // −50%
	if got := g.DayQuota(); got != 180 {
		t.Fatalf("DayQuota() with max hrquota = %d, want 180", got)
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
	if g.Score != 30 || g.HRPoints != 1 || g.Day != 2 {
		t.Fatalf("Score = %d, HRPoints = %d, Day = %d, want 30, 1, 2", g.Score, g.HRPoints, g.Day)
	}
	if g.DayTimer != g.DayLength() {
		t.Fatalf("DayTimer = %v after payday, want %v", g.DayTimer, g.DayLength())
	}
	if g.DayEndQuota != 100 || g.DayEndGainHR != 1 {
		t.Fatalf("summary quota %d, gain %d, want 100, 1", g.DayEndQuota, g.DayEndGainHR)
	}
	if !g.DayEndFirst {
		t.Fatal("DayEndFirst = false on the first payday, want true (HR points explained)")
	}
	// the first payday explains HR points across several lines: numbers
	// only show on the last one
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
	if g.Scene != SceneDayStart || g.DayStartText == "" {
		t.Fatalf("Scene = %v, DayStartText = %q after payday Enters, want SceneDayStart with a quip", g.Scene, g.DayStartText)
	}
	g.DayStartEnter() // reveal
	g.DayStartEnter() // dismiss
	if g.Scene != ScenePlay || !g.RunActive {
		t.Fatalf("Scene = %v, RunActive = %v after morning Enter, want ScenePlay, true", g.Scene, g.RunActive)
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
	g.HRIndex = 0 // hrmult, base cost 3
	if !g.HRBuy() {
		t.Fatal("HRBuy() = false with enough HR points, want true")
	}
	if g.HRPoints != 7 || g.HRMultLevel != 1 {
		t.Fatalf("HRPoints = %d, HRMultLevel = %d, want 7, 1", g.HRPoints, g.HRMultLevel)
	}
	if g.HRUpgradeCost("hrmult") != 6 {
		t.Fatalf("HRUpgradeCost(hrmult) = %d at level 1, want 6", g.HRUpgradeCost("hrmult"))
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

func TestShopBuyAppliesCostAndEffect(t *testing.T) {
	g := startPlay(t, "casa")
	g.Score = 60
	g.ShopIndex = 0 // mult, base cost 50
	if !g.ShopBuy() {
		t.Fatal("ShopBuy() = false with enough points, want true")
	}
	if g.Score != 10 || g.MultLevel != 1 {
		t.Fatalf("Score = %d, MultLevel = %d, want 10, 1", g.Score, g.MultLevel)
	}
	if g.UpgradeCost("mult") != 150 {
		t.Fatalf("UpgradeCost(mult) = %d at level 1, want 150", g.UpgradeCost("mult"))
	}
	if g.ShopBuy() {
		t.Fatal("ShopBuy() = true without enough points, want false")
	}
}
