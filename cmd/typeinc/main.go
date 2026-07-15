// Desktop frontend for typeinc, rendered with raylib.
package main

import (
	_ "embed"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"typeinc/internal/game"
)

//go:embed font.ttf
var fontData []byte

//go:embed key.wav
var keySoundData []byte

//go:embed wrong.mp3
var wrongSoundData []byte

const (
	wordSize  = 48 // px, the word being typed
	uiSize    = 32 // px, menu options / score / hints
	titleSize = 96 // px, menu title
)

var (
	bgColor      = rl.NewColor(0x2C, 0x2C, 0x2C, 0xFF)
	ghostColor   = rl.NewColor(0xFF, 0xFF, 0xFF, 0x55) // faint word to type over
	dimColor     = rl.NewColor(0xFF, 0xFF, 0xFF, 0x66) // unselected menu entries
	goldColor    = rl.NewColor(0xFF, 0xD7, 0x00, 0xFF)
	errorColor   = rl.NewColor(0xD6, 0x33, 0x6C, 0xFF) // wrong-key word highlight
	successColor = rl.NewColor(0x2F, 0xA0, 0x84, 0xFF) // completed-word highlight
	comboColor   = rl.NewColor(0xFF, 0x6B, 0x35, 0xFF) // active combo multiplier
)

// window size in px, refreshed every frame so windowed mode lays out right
var screenW, screenH float32

var (
	keySound   rl.Sound
	wrongSound rl.Sound
)

// prevIntroChars tracks the typewriter to play a key sound per new char.
var prevIntroChars int

type fonts struct {
	word  rl.Font
	ui    rl.Font
	title rl.Font
}

// fontChars is the baked charset: full Latin-1 so the intro story renders
// accents, ñ and ¿¡ correctly.
func fontChars() []rune {
	chars := make([]rune, 0, 224)
	for r := rune(32); r <= 255; r++ {
		chars = append(chars, r)
	}
	return chars
}

func loadFont(size int32) rl.Font {
	return rl.LoadFontFromMemory(".ttf", fontData, size, fontChars())
}

// applyDisplayMode switches between fullscreen, a 1280x720 window, and a
// borderless full-monitor window. It first leaves whatever special state
// the window is in.
//
// NOTE: rl.SetWindowSize is unusable here — raylib-go's purego binding
// preps it with a wrong ffi signature and panics ("too many arguments"),
// so window sizes only change through the toggles, which manage size
// themselves. The recover keeps a failing native call from killing the
// game (the bad mode would be saved and brick every startup).
func applyDisplayMode(mode int) {
	defer func() { _ = recover() }()

	if rl.IsWindowFullscreen() {
		rl.ToggleFullscreen()
	}
	if rl.IsWindowState(rl.FlagBorderlessWindowedMode) {
		rl.ToggleBorderlessWindowed()
	}
	// the window is now decorated at its original 1280x720
	switch mode {
	case 0: // COMPLETA: grow to monitor size via borderless, then go exclusive
		rl.ToggleBorderlessWindowed()
		rl.ToggleFullscreen()
	case 1: // VENTANA: recenter the restored 1280x720 window
		monitor := rl.GetCurrentMonitor()
		mw, mh := rl.GetMonitorWidth(monitor), rl.GetMonitorHeight(monitor)
		rl.SetWindowPosition((mw-1280)/2, (mh-720)/2)
	case 2: // SIN BORDE: borderless window covering the monitor
		rl.ToggleBorderlessWindowed()
	}
}

func main() {
	rl.SetConfigFlags(rl.FlagVsyncHint)
	rl.InitWindow(1280, 720, "typeinc")
	defer rl.CloseWindow()
	rl.SetExitKey(0) // Esc is a game key, not "quit"
	rl.SetTargetFPS(60)

	rl.InitAudioDevice()
	defer rl.CloseAudioDevice()
	keySound = rl.LoadSoundFromWave(rl.LoadWaveFromMemory(".wav", keySoundData, int32(len(keySoundData))))
	defer rl.UnloadSound(keySound)
	wrongSound = rl.LoadSoundFromWave(rl.LoadWaveFromMemory(".mp3", wrongSoundData, int32(len(wrongSoundData))))
	defer rl.UnloadSound(wrongSound)

	f := fonts{word: loadFont(wordSize), ui: loadFont(uiSize), title: loadFont(titleSize)}
	defer rl.UnloadFont(f.word)
	defer rl.UnloadFont(f.ui)
	defer rl.UnloadFont(f.title)

	g := game.New()
	g.OptionsEnabled = true
	g.DisplayMode = 2 // default borderless full-monitor; Load may override
	g.SavePath = game.DefaultSavePath()
	g.Load()

	appliedDisplay := -1
	appliedVolume := -1.0

	for !rl.WindowShouldClose() && !g.QuitRequested {
		if g.DisplayMode != appliedDisplay {
			applyDisplayMode(g.DisplayMode)
			appliedDisplay = g.DisplayMode
		}
		if g.Volume != appliedVolume {
			rl.SetMasterVolume(float32(g.Volume))
			appliedVolume = g.Volume
		}
		screenW = float32(rl.GetScreenWidth())
		screenH = float32(rl.GetScreenHeight())

		update(g)

		rl.BeginDrawing()
		rl.ClearBackground(bgColor)
		switch g.Scene {
		case game.SceneMenu:
			drawMenu(g, f)
		case game.SceneOptions:
			drawOptions(g, f)
		case game.SceneIntro:
			drawIntro(g, f)
		case game.ScenePlay:
			drawPlay(g, f)
		case game.SceneShop:
			drawShop(g, f)
		}
		rl.EndDrawing()

		if g.Scene == game.SceneMenu && rl.IsKeyPressed(rl.KeyEscape) {
			break
		}
	}
	g.Save()
}

func update(g *game.Game) {
	g.Tick(float64(rl.GetFrameTime()))

	switch g.Scene {
	case game.SceneMenu:
		if rl.IsKeyPressed(rl.KeyUp) {
			g.MenuMove(-1)
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			g.MenuMove(1)
		}
		if rl.IsKeyPressed(rl.KeyLeft) {
			g.MenuAdjust(-1)
		}
		if rl.IsKeyPressed(rl.KeyRight) {
			g.MenuAdjust(1)
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			g.MenuSelect()
		}
	case game.SceneOptions:
		if rl.IsKeyPressed(rl.KeyEscape) {
			g.Scene = game.SceneMenu
		}
		if rl.IsKeyPressed(rl.KeyUp) {
			g.OptMove(-1)
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			g.OptMove(1)
		}
		if rl.IsKeyPressed(rl.KeyLeft) {
			g.OptAdjust(-1)
		}
		if rl.IsKeyPressed(rl.KeyRight) {
			g.OptAdjust(1)
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			g.OptSelect()
		}
	case game.SceneIntro:
		if rl.IsKeyPressed(rl.KeyEnter) {
			g.IntroEnter()
		}
		// typewriter key sound: one tap per newly revealed character
		cur := len([]rune(g.IntroVisible()))
		if cur > prevIntroChars {
			rl.PlaySound(keySound)
		}
		prevIntroChars = cur
	case game.ScenePlay:
		if g.CommandMode {
			if rl.IsKeyPressed(rl.KeyEscape) {
				g.CommandCancel()
			}
			if rl.IsKeyPressed(rl.KeyBackspace) {
				g.CommandBackspace()
			}
			if rl.IsKeyPressed(rl.KeyEnter) {
				g.CommandSubmit()
			}
			for r := rl.GetCharPressed(); r != 0; r = rl.GetCharPressed() {
				rl.PlaySound(keySound)
				g.TypeChar(r)
			}
			return
		}
		if rl.IsKeyPressed(rl.KeyEscape) {
			g.BackToMenu()
			return
		}
		for r := rl.GetCharPressed(); r != 0; r = rl.GetCharPressed() {
			wasTyping := g.Phase == game.PhaseTyping
			g.TypeChar(r)
			if !wasTyping {
				continue
			}
			// a key that just failed the word plays the error sound instead
			if g.Phase == game.PhaseError {
				rl.PlaySound(wrongSound)
			} else {
				rl.PlaySound(keySound)
			}
		}
	case game.SceneShop:
		if rl.IsKeyPressed(rl.KeyEscape) {
			g.ExitShop()
		}
		if rl.IsKeyPressed(rl.KeyUp) {
			g.ShopMove(-1)
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			g.ShopMove(1)
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			g.ShopBuy()
		}
	}
}

// drawCentered draws text horizontally centered at the given y.
func drawCentered(f rl.Font, text string, y float32, size float32, color rl.Color) {
	m := rl.MeasureTextEx(f, text, size, 0)
	rl.DrawTextEx(f, text, rl.Vector2{X: (screenW - m.X) / 2, Y: y}, size, 0, color)
}

func drawMenu(g *game.Game, f fonts) {
	x := screenW * 0.12 // menu block is left-aligned
	rl.DrawTextEx(f.title, "TYPEINC", rl.Vector2{X: x, Y: screenH * 0.22}, titleSize, 0, rl.White)

	for i, item := range g.MenuItems() {
		label := "  " + item.Label
		color := dimColor
		if i == g.MenuIndex {
			label = "> " + item.Label
			color = rl.White
		}
		rl.DrawTextEx(f.ui, label, rl.Vector2{X: x, Y: screenH*0.48 + float32(i)*(uiSize*1.6)}, uiSize, 0, color)
	}

	rl.DrawTextEx(f.ui, g.UI().MenuHint, rl.Vector2{X: x, Y: screenH - uiSize*3}, uiSize, 0, ghostColor)
}

func drawOptions(g *game.Game, f fonts) {
	x := screenW * 0.12 // options share the menu's left alignment
	rl.DrawTextEx(f.title, g.UI().MenuOptions, rl.Vector2{X: x, Y: screenH * 0.22}, titleSize, 0, rl.White)

	for i, id := range game.OptItemIDs {
		label := "  " + g.OptLabel(id)
		color := dimColor
		if i == g.OptIndex {
			label = "> " + g.OptLabel(id)
			color = rl.White
		}
		rl.DrawTextEx(f.ui, label, rl.Vector2{X: x, Y: screenH*0.48 + float32(i)*(uiSize*1.6)}, uiSize, 0, color)
	}

	rl.DrawTextEx(f.ui, g.UI().OptionsHint, rl.Vector2{X: x, Y: screenH - uiSize*3}, uiSize, 0, ghostColor)
}

// wrapText greedily wraps words so each line fits maxW at the given size.
func wrapText(f rl.Font, text string, size, maxW float32) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}
	var lines []string
	line := words[0]
	for _, w := range words[1:] {
		if rl.MeasureTextEx(f, line+" "+w, size, 0).X > maxW {
			lines = append(lines, line)
			line = w
		} else {
			line += " " + w
		}
	}
	return append(lines, line)
}

func drawIntro(g *game.Game, f fonts) {
	lines := wrapText(f.ui, g.IntroVisible(), uiSize, screenW*0.8)
	lineH := float32(uiSize) * 1.4
	y := (screenH - lineH*float32(len(lines))) / 2
	for i, line := range lines {
		drawCentered(f.ui, line, y+float32(i)*lineH, uiSize, rl.White)
	}
	if g.IntroLineDone() {
		drawCentered(f.ui, g.UI().IntroHint, screenH-uiSize*3, uiSize, ghostColor)
	}
}

func drawPlay(g *game.Game, f fonts) {
	m := rl.MeasureTextEx(f.word, g.Word, wordSize, 0)
	pos := rl.Vector2{X: (screenW - m.X) / 2, Y: (screenH - m.Y) / 2}

	// word highlight: pink flash on error, green on success
	if g.ErrorFlashing() {
		rl.DrawRectangle(int32(pos.X)-16, int32(pos.Y)-8, int32(m.X)+32, int32(m.Y)+16, errorColor)
	} else if g.Phase == game.PhaseSuccess {
		rl.DrawRectangle(int32(pos.X)-16, int32(pos.Y)-8, int32(m.X)+32, int32(m.Y)+16, successColor)
		// points gained, gold, floating above the word
		gain := "+" + strconv.Itoa(g.LastGain)
		gm := rl.MeasureTextEx(f.word, gain, wordSize, 0)
		rl.DrawTextEx(f.word, gain, rl.Vector2{X: (screenW - gm.X) / 2, Y: pos.Y - gm.Y - 24}, wordSize, 0, goldColor)
	}

	// ghost word with the correctly typed prefix drawn on top; monospace
	// font keeps both perfectly aligned
	rl.DrawTextEx(f.word, g.Word, pos, wordSize, 0, ghostColor)
	prefixW := float32(0)
	if g.Pos > 0 {
		prefix := g.Word[:g.Pos]
		rl.DrawTextEx(f.word, prefix, pos, wordSize, 0, rl.White)
		prefixW = rl.MeasureTextEx(f.word, prefix, wordSize, 0).X
	}

	// cursor: thin bar in front of the current letter, only while typing
	if g.Phase == game.PhaseTyping {
		rl.DrawRectangle(int32(pos.X+prefixW), int32(pos.Y), 3, int32(m.Y), rl.White)
	}

	// score, gold, top-right corner; combo multiplier right below it
	score := strconv.Itoa(g.Score)
	sm := rl.MeasureTextEx(f.ui, score, uiSize, 0)
	rl.DrawTextEx(f.ui, score, rl.Vector2{X: screenW - sm.X - 32, Y: 24}, uiSize, 0, goldColor)

	combo := "x" + strconv.FormatFloat(g.ComboMult(), 'f', 1, 64)
	comboCol := dimColor
	if g.ComboMult() > 1 {
		comboCol = comboColor
	}
	cm := rl.MeasureTextEx(f.ui, combo, uiSize, 0)
	rl.DrawTextEx(f.ui, combo, rl.Vector2{X: screenW - cm.X - 32, Y: 24 + sm.Y + 8}, uiSize, 0, comboCol)

	// command bar / unknown-command message, bottom center
	if g.CommandMode {
		drawCentered(f.ui, g.CommandBuf+"_", screenH-uiSize*3, uiSize, rl.White)
	} else if g.CommandErr > 0 {
		drawCentered(f.ui, g.UI().UnknownCmd, screenH-uiSize*3, uiSize, errorColor)
	}
}

func drawShop(g *game.Game, f fonts) {
	drawCentered(f.title, g.UI().ShopTitle, screenH*0.2, titleSize, rl.White)

	for i, up := range game.Upgrades {
		level := g.UpgradeLevel(up.ID)
		cost := g.UpgradeCost(up.ID)
		line := g.UpgradeName(up.ID) + "  NV." + strconv.Itoa(level) + "  " + g.UpgradeEffect(up.ID) + "   "
		price := strconv.Itoa(cost)

		color := dimColor
		if i == g.ShopIndex {
			line = "> " + line
			color = rl.White
		} else {
			line = "  " + line
		}

		y := screenH*0.45 + float32(i)*(uiSize*2)
		full := rl.MeasureTextEx(f.ui, line+price, uiSize, 0)
		pos := rl.Vector2{X: (screenW - full.X) / 2, Y: y}
		rl.DrawTextEx(f.ui, line, pos, uiSize, 0, color)

		priceColor := goldColor
		if g.Score < cost {
			priceColor = errorColor
		}
		lineW := rl.MeasureTextEx(f.ui, line, uiSize, 0).X
		rl.DrawTextEx(f.ui, price, rl.Vector2{X: pos.X + lineW, Y: y}, uiSize, 0, priceColor)
	}

	// score, gold, top-right corner
	score := strconv.Itoa(g.Score)
	sm := rl.MeasureTextEx(f.ui, score, uiSize, 0)
	rl.DrawTextEx(f.ui, score, rl.Vector2{X: screenW - sm.X - 32, Y: 24}, uiSize, 0, goldColor)

	drawCentered(f.ui, g.UI().ShopHint, screenH-uiSize*3, uiSize, ghostColor)
}
