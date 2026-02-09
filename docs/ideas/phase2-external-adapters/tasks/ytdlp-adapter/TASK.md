# Task: yt-dlp Stream URL Extractor Adapter

## Implements

`port.StreamURLExtractor` interface (`internal/port/streamurl.go`)

## Package

`internal/adapter/ytdlp/`

## Files to Create

- `ytdlp.go` — adapter implementation
- `ytdlp_test.go` — unit tests

## Interface to Satisfy

```go
type StreamURLExtractor interface {
    Extract(ctx context.Context, videoID string) (string, error)
}
```

## Implementation Details

### Constructor

```go
func New() *Extractor
```

No config needed — yt-dlp is expected to be in `$PATH`.

### Extract

1. Build the YouTube URL: `https://www.youtube.com/watch?v=VIDEO_ID`
2. Execute yt-dlp as a subprocess with context support:
   ```
   yt-dlp -f bestaudio --get-url --no-playlist <url>
   ```
   - `-f bestaudio` — select best audio-only format
   - `--get-url` — output only the direct stream URL (no download)
   - `--no-playlist` — ensure single video, not playlist expansion
3. Capture stdout, trim whitespace
4. Use `exec.CommandContext(ctx, ...)` for cancellation support
5. If yt-dlp exits non-zero, return stderr as the error message
6. Validate that the returned URL is non-empty

### Error Handling

- yt-dlp not found in PATH → clear error: "yt-dlp not found: install from https://github.com/yt-dlp/yt-dlp"
- Non-zero exit → wrap stderr in error
- Empty output → return error "no stream URL returned"
- Context cancelled → propagated via `exec.CommandContext`

## Acceptance Criteria

- [x] Extracts a valid, playable audio stream URL from a known video ID
- [x] Supports context cancellation
- [x] Returns clear error when yt-dlp is not installed
- [x] Handles invalid video IDs gracefully

## Tests

- Extract URL from a known, stable YouTube video (e.g., YouTube's own test video)
- Verify returned URL is a valid HTTP(S) URL
- Test with invalid video ID → expect error
- Skip tests if yt-dlp not installed (`exec.LookPath` check + `t.Skip`)

## Estimated Effort

~1 hour
