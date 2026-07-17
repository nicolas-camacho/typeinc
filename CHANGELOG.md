# Changelog

All notable changes to TYPE.Inc are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/); versions are the
`vX.Y.Z` release tags. Add entries under **Unreleased** as you go and cut
them into a new section when tagging.

## [Unreleased]

## [0.3.1] - 2026-07-17

### Added
- ESC skips the end-of-day summary text (including the long first-payroll
  HR explanation), revealing the numbers at once.

### Fixed
- The terminal unlocks as soon as its reveal quip can appear — no longer
  requires holding a clue first.
- Re-typing an already-used key in the terminal re-reads its document
  instead of answering "not recognized".

## [0.3.0] - 2026-07-17

### Added
- Background mystery storyline: six acts unlocking over in-game days, a
  hidden terminal system with keyed documents, corrupt words that freeze
  the whole office, rituals, watched words, possessed interns and clock
  glitches. Spoiler-free hint: mind anything violet.
- Terminal entry in the main menu once discovered; the quip revealing the
  terminal repeats until it is opened for the first time.
- ESC skips the intro story and the CONTINUE return quip.

### Changed
- Days last 60s of active play (was 90s).
- Quota inflation per upgrade level lowered to 25/15/8/4/4 percent
  (multiplier / intern / max streak / golden / intern coffee).
- Office events roll at 20% chance (was 10%) with a cap of 5 per day
  (was 3).

### Breaking
- Saves are now gated by `GameVersion`: any save older than 0.3.0 (or
  without the field) is reset on load, keeping only settings.

## [0.2.0] - 2026-07-16

### Added
- Golden words: rare words paying a multiplier; shop upgrade raises both
  chance and payout.
- Intern auto-typer with its own HUD word, speed upgrade, and an HR meta
  upgrade raising the per-run headcount.
- Office events: coffee frenzy (everything ×2) and HR inspection (urgent
  word, bonus on time, harmless on expiry).
- `/terminar` · `/endday` command unlocked by an HR meta upgrade (settles
  the quota immediately, firing included).
- Two-step RESET entry in the main menu.
- Best streak (day and global) and best day stats.
- Per-upgrade purchase quips and larger HR script pools.

### Changed
- Quota locked at day start: upgrades bought mid-day only inflate the
  next day; steeper daily growth.
- HR points paid only on days multiple of 3 (payout = day / 2).
- Shop and HR meta upgrade costs raised.

## [0.1.0] - 2026-07-15

### Added
- Incremental typing core: word loop, streak combo, score, shop
  (multiplier, max streak).
- Day cycle with HR quota and firing; sarcastic HR narrator (es/en).
- HR meta progression paid with HR points; permanent upgrades.
- Global stats screen and shared save file for both frontends.
- Two frontends over one headless core: raylib desktop and Bubble Tea TUI.
- Windows release workflow (tag-triggered, `TYPE.Inc.exe`).

[Unreleased]: https://github.com/nicolas-camacho/typeinc/compare/v0.3.1...HEAD
[0.3.1]: https://github.com/nicolas-camacho/typeinc/compare/v0.3.0...v0.3.1
[0.3.0]: https://github.com/nicolas-camacho/typeinc/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/nicolas-camacho/typeinc/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/nicolas-camacho/typeinc/releases/tag/v0.1.0
