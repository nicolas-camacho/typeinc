package game

import "testing"

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
	if len(items) != 3 { // new, lang, quit — no save, no options support
		t.Fatalf("MenuItems() = %d entries without options, want 3", len(items))
	}
	g.OptionsEnabled = true
	items = g.MenuItems()
	if len(items) != 4 || items[2].ID != "options" {
		t.Fatalf("MenuItems() with options = %d entries, [2] = %q, want 4, \"options\"", len(items), items[2].ID)
	}
	g.MenuIndex = 3 // SALIR
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
