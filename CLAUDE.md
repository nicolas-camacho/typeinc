# TYPE.Inc (typeinc)

Incremental typing game. Shared core `internal/game` (zero render deps) +
two frontends: `cmd/typeinc` (raylib desktop) and `cmd/typeinc-tui`
(Bubble Tea). Both share one save file.

## Critical build notes
- raylib-go runs via purego. NEVER call rl.SetWindowSize: its purego
  binding has a broken ffi signature and panics at runtime; resize only
  through ToggleFullscreen/ToggleBorderlessWindowed (see applyDisplayMode).
- Desktop assets are embedded (go:embed), so they must live inside
  cmd/typeinc/. Nicolas drops new audio/fonts into assets/ (gitignored);
  copy them into cmd/typeinc/ to use them.

## Game direction (agreed with Nicolas)
- More mechanics are planned. Balance numbers are tuned by Nicolas
  playtesting, but actively help: flag how any new mechanic shifts the
  quota-vs-income curve, and watch for boring, frustrating or too-easy
  outcomes.
- The sarcastic-HR narrator aesthetic is permanent. New content should add
  more HR quips; quips never repeat back-to-back (use pickQuip) and every
  script table must keep identical line counts in "es" and "en".
- The color palette is settled — change a color only if Nicolas asks.

## Running
- TUI is interactive: never launch it, Nicolas tests it himself.
- Desktop run: $env:CGO_ENABLED='0'; go run ./cmd/typeinc
- Release: git tag vX.Y.Z && git push origin vX.Y.Z → GitHub Actions
  publishes TYPE.Inc.exe. Still v0.x — not the final version.
