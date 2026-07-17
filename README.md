# TYPE.Inc

Incremental typing game narrated by a sarcastic HR department. Type words,
earn points, survive the daily quota, and climb a meta progression that is
openly rigged against you.

Still v0.x — mechanics and balance change often.

## How the game works

- **Intro** — Enter advances the HR welcome story; ESC skips it entirely.
- **Typing loop** — a word appears; type it letter by letter. A completed
  word pays `len × (1 + mult) × combo × HR bonus`, times the golden, coffee
  and inspection multipliers when they apply. A wrong key resets the word
  and the streak.
- **Combo** — consecutive words without a mistake build a streak: +10% per
  word, capped (`MAX STREAK` upgrade raises the cap).
- **Day cycle** — a day lasts 60s of active play (pauses in menus, shop and
  command mode). At day end HR charges the **quota**, locked at day start:
  exponential in the day number and inflated by the upgrades you owned when
  the day began. Cover it and the next day starts; fall short and you are
  fired — the run is erased, HR points survive.
- **Payroll** — HR points are only paid on days multiple of 3, at
  `day / 2` points (day 3 → 1, day 6 → 3, day 9 → 4…). Fired mid-cycle?
  That cycle is forfeited.
- **Shop** (in-run, paid with score, via `/tienda` · `/shop`): point
  multiplier, max streak, golden words, intern, intern coffee.
- **HR office** (meta, paid with HR points, via `/rrhh` · `/hr`): permanent
  bonus, extended shift, reduced quota, intern headcount, and the
  `/terminar` · `/endday` command (ends the day on the spot — quota settles
  immediately, firing included).
- **Golden words** — small chance a word comes out golden and pays a
  multiplier; the shop upgrade raises both chance and payout.
- **Intern** — an auto-typer with its own word on the left of the HUD.
  Earns like you but without combo or golden bonuses. Freezes whenever the
  day clock freezes. More interns unlock through the HR office.
- **Office events** — random during play: coffee frenzy (everything pays
  ×2 for a while) and HR inspection (the next word is urgent; typed in
  time it pays ×3, expired it just goes back to normal).
- **Menu extras** — global stats (words, fails, best streak, best day, top
  failed words) and a two-step RESET that wipes all progress.
- **Something else** — after a few days, the game starts behaving oddly:
  purple words, glitched messages, a system that should not exist. It is a
  long background mystery that unlocks over many in-game days. No spoilers
  here — pay attention to anything violet.

Everything is localized: Spanish and English, switchable from the menu.

## Architecture

```
internal/game/        shared core — ALL game logic, zero render deps
├── game.go           state machine (Scene/Phase), Tick(dt), balance constants
├── save.go           JSON save/load, GameVersion-gated (old saves reset)
├── strings.go        every player-facing string + HR quip pools (es/en)
├── story.go          background-mystery engine (acts, triggers, terminal)
├── storystrings.go   story content: quips, documents, keys (es/en)
└── words.go          embedded dictionaries (a-z words, 3-14 letters)

cmd/typeinc/          desktop frontend — raylib via purego (no cgo)
                      embeds font + sounds; audio, display modes, options
cmd/typeinc-tui/      terminal frontend — Bubble Tea + Lipgloss
                      no audio, no options screen
```

The core is a headless state machine: frontends construct `game.New()`,
call `Load()`, then every frame feed input through methods (`TypeChar`,
`MenuSelect`, `ShopBuy`…), advance time with `Tick(dt)` and render from
exported state and getters. They never reach into game internals.

Both frontends share **one save file** at `%APPDATA%\typeinc\save.json`
(`os.UserConfigDir`). Run state lives under `Run` and dies with a firing;
HR points, meta levels and global stats live outside it and persist.

Planned features live in [ROADMAP.md](ROADMAP.md).

## Working on the code

Conventions the codebase relies on:

- **Balance is named constants** at the top of `internal/game/game.go`
  (quota growth and per-upgrade inflation terms, costs, event timings…).
  Tune there, never inline numbers.
- **Quota freeze**: `DayQuota()` returns the value locked at day start;
  `computeQuota()` runs only at day boundaries. New run upgrades that add
  income should add an inflation term to `computeQuota()`.
- **Strings**: every player-facing label exists in both `"es"` and `"en"`
  tables in `strings.go`. Quip script tables must keep **identical line
  counts** in both languages (tests enforce it). Random quips go through
  `pickQuip`, which never repeats the previous pick of the same pool.
- **New scenes** need a `Scene` constant plus input and render wiring in
  BOTH frontends (`update`/`drawX` in the desktop, `Update`/`viewX` in the
  TUI). Shop and HR office render generically from the upgrade slices —
  new upgrades usually need zero frontend work.
- **Save changes are additive** within a save generation: new fields with
  zero-value defaults, guard on `Load`. `SaveVersion` (a dotted game
  version) gates the whole file — loading an older or versionless save
  keeps only the settings. Bump it ONLY when a release must invalidate old
  saves on purpose.
- **Story content** lives in `storystrings.go` (quips per act, corrupt
  words, terminal documents). Keys and document IDs are language-neutral
  lowercase tokens; corrupt-word texts are typeable ASCII (a-z + space).
  Structural tests enforce all of it.
- **Every mechanic ships with tests** in `internal/game/game_test.go`.
  Tests force words/state directly to avoid RNG; chance rolls live in tiny
  helpers so effects can be tested deterministically.

Critical, learned the hard way:

- **Build with `CGO_ENABLED=0`.** raylib-go runs through purego; no C
  compiler is needed or wanted.
- **Never call `rl.SetWindowSize`** — its purego binding has a broken ffi
  signature and panics at runtime. Window sizes only change through
  `ToggleFullscreen` / `ToggleBorderlessWindowed` (see `applyDisplayMode`).

## Commands

PowerShell:

```powershell
# run the desktop version
$env:CGO_ENABLED='0'; go run ./cmd/typeinc

# run the TUI (interactive — run it in a real terminal)
$env:CGO_ENABLED='0'; go run ./cmd/typeinc-tui

# tests, vet, build everything
$env:CGO_ENABLED='0'; go test ./...
$env:CGO_ENABLED='0'; go vet ./...
$env:CGO_ENABLED='0'; go build ./...
```

## Release

Before tagging: move the `Unreleased` entries in `CHANGELOG.md` into a new
version section.

```powershell
git tag v0.X.Y
git push origin v0.X.Y
```

The `release` GitHub Actions workflow builds `TYPE.Inc.exe` (windows/amd64,
icon + version info via go-winres, `-H windowsgui`), zips it and publishes
a GitHub Release with generated notes. The TUI is not distributed yet —
build it from source.
