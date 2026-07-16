package game

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// RunData is the persisted state of a run in progress.
type RunData struct {
	Score          int
	MultLevel      int
	StreakCapLevel int
	Streak         int
	Day            int
	DayTimer       float64
	DayWords       int
	DayFails       int
	DayFailCounts  map[string]int
}

// SaveData is the on-disk save file, shared by the desktop and TUI
// frontends. HR fields live outside Run: they survive across runs.
type SaveData struct {
	Version     int
	Run         *RunData // nil when no run is saved
	LangIndex   int
	Volume      float64
	DisplayMode int

	HRPoints     int
	HRMultLevel  int
	HRDayLevel   int
	HRQuotaLevel int

	TotalWords int
	TotalFails int
	FailCounts map[string]int
}

// DefaultSavePath is the shared save location
// (e.g. %APPDATA%\typeinc\save.json). Empty when no config dir exists.
func DefaultSavePath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "typeinc", "save.json")
}

// Save writes the current state to g.SavePath. A game without a SavePath
// (tests) never touches the disk. Errors are ignored: losing an autosave
// must never interrupt play.
func (g *Game) Save() {
	if g.SavePath == "" {
		return
	}
	data := SaveData{
		Version:     2,
		LangIndex:   g.LangIndex,
		Volume:      g.Volume,
		DisplayMode: g.DisplayMode,

		HRPoints:     g.HRPoints,
		HRMultLevel:  g.HRMultLevel,
		HRDayLevel:   g.HRDayLevel,
		HRQuotaLevel: g.HRQuotaLevel,

		TotalWords: g.TotalWords,
		TotalFails: g.TotalFails,
		FailCounts: g.FailCounts,
	}
	if g.RunActive {
		data.Run = &RunData{
			Score:          g.Score,
			MultLevel:      g.MultLevel,
			StreakCapLevel: g.StreakCapLevel,
			Streak:         g.Streak,
			Day:            g.Day,
			DayTimer:       g.DayTimer,
			DayWords:       g.DayWords,
			DayFails:       g.DayFails,
			DayFailCounts:  g.dayFailCounts,
		}
	} else if g.savedRun != nil {
		data.Run = g.savedRun
	}
	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(g.SavePath), 0o755); err != nil {
		return
	}
	_ = os.WriteFile(g.SavePath, raw, 0o644)
}

// Load reads the save file, applying settings immediately and keeping a
// stored run available for CONTINUAR. Missing or corrupt files are
// silently treated as "no save".
func (g *Game) Load() {
	if g.SavePath == "" {
		return
	}
	raw, err := os.ReadFile(g.SavePath)
	if err != nil {
		return
	}
	var data SaveData
	if err := json.Unmarshal(raw, &data); err != nil {
		return
	}
	if data.LangIndex >= 0 && data.LangIndex < len(Languages) {
		g.LangIndex = data.LangIndex
	}
	if data.Volume >= 0 && data.Volume <= 1 {
		g.Volume = data.Volume
	}
	if data.DisplayMode >= 0 && data.DisplayMode < len(uiTables["es"].DisplayModes) {
		g.DisplayMode = data.DisplayMode
	}
	if data.HRPoints > 0 {
		g.HRPoints = data.HRPoints
	}
	if data.HRMultLevel > 0 {
		g.HRMultLevel = data.HRMultLevel
	}
	if data.HRDayLevel > 0 {
		g.HRDayLevel = data.HRDayLevel
	}
	if data.HRQuotaLevel > 0 {
		g.HRQuotaLevel = min(data.HRQuotaLevel, HRQuotaMax)
	}
	if data.TotalWords > 0 {
		g.TotalWords = data.TotalWords
	}
	if data.TotalFails > 0 {
		g.TotalFails = data.TotalFails
	}
	if data.FailCounts != nil {
		g.FailCounts = data.FailCounts
	}
	g.savedRun = data.Run
}
