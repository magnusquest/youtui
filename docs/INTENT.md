# YouTui — Project Intent

## Vision

A terminal-based music player that streams YouTube audio with queue management and an integrated Pomodoro timer. Built as a single Go binary using Charm Bubbletea, mpv, and yt-dlp.

## Values

- **Simplicity**: Single binary, minimal config, works out of the box (given mpv + yt-dlp installed)
- **Terminal-native**: Beautiful TUI using Charm ecosystem — no browser, no Electron
- **Hexagonal architecture**: Domain logic decoupled from external dependencies via ports/adapters
- **Testability**: Each adapter testable in isolation against real external tools

## Constraints

- Go 1.25+
- External runtime dependencies: mpv, yt-dlp
- YouTube Data API v3 key required for search
- Charm Bubbletea for TUI framework
- `charmbracelet/x/mosaic` for thumbnail rendering

## Success Metrics

- All 5 port interfaces have working adapter implementations
- Each adapter has unit tests that exercise the real external dependency
- Clean `go build` with no warnings
- Application can search YouTube, extract stream URLs, and play audio through mpv

## Decision Log

| Date | Decision | Reference |
|------|----------|-----------|
| 2025-01-15 | Phase 1 complete — domain entities, ports, config, TUI skeleton | commit `a180d34` |
| 2025-02-08 | Phase 2 ideation started — external adapters | Issue #1 |
| 2025-02-08 | Phase 2 complete — all 5 adapters implemented and tested | Issue #1 closed, commit `c63ca90` |
