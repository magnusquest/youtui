# Task: mpv IPC Audio Player Adapter

## Implements

`port.AudioPlayer` interface (`internal/port/player.go`)

## Package

`internal/adapter/mpv/`

## Files to Create

- `mpv.go` — adapter implementation
- `mpv_test.go` — unit tests

## Interface to Satisfy

```go
type AudioPlayer interface {
    Play(ctx context.Context, streamURL string) error
    Pause(ctx context.Context) error
    Resume(ctx context.Context) error
    Stop(ctx context.Context) error
    Seek(ctx context.Context, position time.Duration) error
    SeekRelative(ctx context.Context, delta time.Duration) error
    SetVolume(ctx context.Context, volume int) error
    GetPosition(ctx context.Context) (time.Duration, error)
    GetDuration(ctx context.Context) (time.Duration, error)
    IsPlaying(ctx context.Context) (bool, error)
    Close() error
}
```

## Implementation Details

### Constructor

```go
func New(socketPath string) (*Player, error)
```

1. Launch mpv in idle mode as a subprocess:
   ```
   mpv --idle --no-video --no-terminal --input-ipc-server=<socketPath>
   ```
   - `--idle` — start without playing anything
   - `--no-video` — audio only
   - `--no-terminal` — suppress terminal output
   - `--input-ipc-server` — enable JSON IPC on the given Unix socket
2. Wait briefly for the socket to become available (poll with short backoff)
3. Connect to the Unix socket (`net.Dial("unix", socketPath)`)
4. Start a background goroutine to read responses from the socket (request/response multiplexing)
5. Store the `*exec.Cmd`, `net.Conn`, and a `sync.Mutex` for command serialization

### IPC Protocol

mpv uses a JSON-based IPC protocol over Unix sockets:

**Sending commands:**
```json
{"command": ["loadfile", "URL"], "request_id": 1}
{"command": ["set_property", "pause", true], "request_id": 2}
{"command": ["get_property", "time-pos"], "request_id": 3}
```

**Receiving responses:**
```json
{"error": "success", "data": null, "request_id": 1}
{"error": "success", "data": 42.5, "request_id": 3}
```

Implement a `sendCommand(ctx context.Context, args ...interface{}) (interface{}, error)` method that:
- Assigns an incrementing `request_id`
- Marshals and writes the JSON command + newline
- Waits for the matching response (by `request_id`) or context cancellation
- Returns the `data` field or wraps the `error` field

### Method Mapping

| Method | mpv Command |
|--------|------------|
| `Play(url)` | `["loadfile", url]` |
| `Pause()` | `["set_property", "pause", true]` |
| `Resume()` | `["set_property", "pause", false]` |
| `Stop()` | `["stop"]` |
| `Seek(pos)` | `["seek", seconds, "absolute"]` |
| `SeekRelative(delta)` | `["seek", seconds, "relative"]` |
| `SetVolume(vol)` | `["set_property", "volume", vol]` |
| `GetPosition()` | `["get_property", "time-pos"]` → float64 seconds |
| `GetDuration()` | `["get_property", "duration"]` → float64 seconds |
| `IsPlaying()` | `["get_property", "pause"]` → invert bool |

### Close

1. Send `["quit"]` command
2. Close the socket connection
3. Wait for the mpv process to exit (`cmd.Wait()`)
4. Clean up the socket file

### Concurrency

- Use `sync.Mutex` to serialize writes to the socket
- Use a `map[int]chan json.RawMessage` for pending response channels keyed by `request_id`
- Background reader goroutine dispatches responses to the correct channel

## Acceptance Criteria

- [x] Launches mpv subprocess and connects via IPC socket
- [x] Play loads a URL and starts playback
- [x] Pause/Resume toggle playback state
- [x] Seek and SeekRelative position correctly
- [x] Volume control works (0-100)
- [x] GetPosition/GetDuration return valid durations
- [x] IsPlaying reflects actual state
- [x] Close cleanly shuts down mpv process and socket
- [x] Context cancellation supported on all methods

## Tests

- Launch mpv, play a short audio URL (can use a public domain clip or yt-dlp extracted URL)
- Verify pause/resume toggles state
- Verify volume changes
- Verify position/duration return reasonable values
- Verify Close() exits cleanly
- Skip tests if mpv not installed (`exec.LookPath` check + `t.Skip`)

## Estimated Effort

~3 hours (most complex adapter due to IPC protocol)
