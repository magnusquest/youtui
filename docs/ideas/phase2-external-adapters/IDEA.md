# Phase 2: External Adapters

## Summary

Implement the 5 adapter packages that connect YouTui's domain layer to external services: YouTube API, yt-dlp, mpv, Charm mosaic, and JSON file storage.

## Approach

Each adapter lives in `internal/adapter/<name>/` and implements its corresponding `port` interface. Adapters are the only packages that import external dependencies (Google API client, exec.Command for yt-dlp, net.Conn for mpv IPC, charmbracelet/x/mosaic for thumbnails).

## Implementation Order

Dependency-driven — simplest and least-coupled first:

1. **Storage** — JSON file read/write, zero external deps beyond stdlib
2. **yt-dlp** — subprocess invocation, no network state to manage
3. **YouTube** — API client, needs API key from config
4. **mpv** — IPC socket, most complex protocol (JSON-based commands)
5. **Mosaic** — image download + charmbracelet/x/mosaic rendering

## Testing Strategy

Each adapter gets unit tests that test against the **real external tool**, not mocks:

| Adapter | Test Approach |
|---------|--------------|
| Storage | Write/read temp files, verify round-trip |
| yt-dlp | Call real yt-dlp binary with a known video ID, verify URL returned |
| YouTube | Call real API with test key (from `YOUTUBE_API_KEY` env), verify search results |
| mpv | Launch real mpv in idle mode with IPC socket, send commands, verify responses |
| Mosaic | Download a real thumbnail URL, render to string, verify non-empty output |

Tests requiring external tools/network should use `testing.Short()` skip guards so `go test -short` skips them.

## New Dependencies Required

```
google.golang.org/api          # YouTube Data API v3 client
github.com/charmbracelet/x/mosaic  # Terminal image rendering
```

## Architecture Invariants

- Adapters MUST NOT import other adapters
- Adapters MUST NOT import the service layer
- Adapters ONLY import: their port interface, domain entities, config, and external libraries
- All adapter constructors accept config values (not the whole Config struct)

## Risks & Mitigations

| Risk | Mitigation |
|------|-----------|
| yt-dlp output format changes | Pin to known output format flags (`-j` JSON), parse defensively |
| mpv IPC protocol quirks | Use `--input-ipc-server` with JSON IPC, handle async events |
| YouTube API rate limits | Config-driven `maxResults`, respect quota in tests |
| mosaic library API instability (x/ experimental) | Pin version in go.mod |
