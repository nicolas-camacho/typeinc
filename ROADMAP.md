# TYPE.Inc — Roadmap

Planned features, roughly in suggested priority order within each block.
Reorder freely; nothing here is committed until it ships (and lands in
CHANGELOG.md).

## 1. Debug & tuning tools (dev-only)

Gated behind a build tag (`//go:build debug`) or a `-debug` CLI flag so
release builds never carry them.

### Story inspector
- Panel/command that dumps the live mystery state: current act, BestDay
  gate met or not, every clue/doc with its unlock condition and whether it
  is satisfied, every ritual with its per-run fired flag.
- Cheats: force-arm any corrupt word, force intern possession, unlock any
  doc, jump to an act (respecting or overriding gates), reset story only.
- Goal: verify the full puzzle chain (vera → turno77 → sotano → 4413 →
  vera4413) end to end in minutes instead of a 25-day playthrough.

### Balance simulator
- Headless day-by-day simulation over the shared core: input = average
  WPM (and error rate), a purchase strategy (greedy, none, custom), and
  the meta levels; output = per-day table of expected income vs locked
  quota, flagged the first unreachable day.
- "Words per day needed" table: for each day N, minimum completed words
  (at current multipliers) to cover the quota — spots impossible walls
  after balance changes without playtesting.
- Could run as a `go test` helper or a tiny `cmd/typeinc-sim` CLI printing
  the table; reuses `computeQuota`, `WordGain`, `hrGainForDay` directly so
  it can never drift from the game.

## 2. New upgrades

Every run upgrade that adds income must add a matching term to
`computeQuota()` (see README conventions). Candidates:

| Upgrade | Type | Effect knob | Quota note |
|---|---|---|---|
| Event magnet | shop | +event roll chance per level | small term (~4%) |
| Golden polish | shop | +golden multiplier only (chance untouched) | reuse golden term |
| Overdrive | shop | combo keeps growing past the cap at half rate | needs its own term |
| Streak insurance | shop | first fail of the day keeps the streak (1/day) | small term |
| Coffee machine | HR meta | frenzy lasts longer per level | none (meta) |
| Union rep | HR meta | one quota miss forgiven per run (demotion instead of firing) | none, expensive |
| Discount badge | HR meta | shop prices −5% per level | indirect — watch it |
| Overtime | HR meta | +1 payroll point on every payday | none (meta) |

Balance flags to keep in mind: discounts compound with everything;
insurance-style upgrades remove tension if cheap — price them as luxuries.

## 3. UI, sound & screen effects

### Sound (desktop)
- Distinct sounds: golden word completed, corrupt word completed (violet
  channel needs its own voice), event start, firing stinger, payday chime,
  terminal keystrokes (lower pitch than the game keys).
- Drop new files into `assets/`, copy into `cmd/typeinc/` for go:embed
  (see CLAUDE.md).

### Screen effects (desktop)
- Screen shake on firing and on corrupt-word timeout.
- Brief white/gold flash on golden completion; violet vignette while a
  corrupt word is active.
- CRT scanlines + slight flicker in the terminal scene (sells the "old
  system" fiction).
- Word-completion particles (letters scattering) — cheap with raylib.

### TUI
- Sound: agreed route is oto v3 + go-mp3 (pure Go on Windows; cgo on
  Linux/macOS) — deferred until TUI binaries are distributed, needs an
  Actions build matrix.
- Renders: brief inverse-video flash on golden/corrupt completion, softer
  than desktop but same language.

### Both frontends
- Day-end summary: small count-up animation for the settlement numbers.
- Optional reduced-effects toggle in options (accessibility).
