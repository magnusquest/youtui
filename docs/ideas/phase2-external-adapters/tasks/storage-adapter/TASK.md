# Task: JSON File Storage Adapter

## Implements

`port.Storage` interface (`internal/port/storage.go`)

## Package

`internal/adapter/storage/`

## Files to Create

- `storage.go` — adapter implementation
- `storage_test.go` — unit tests

## Interface to Satisfy

```go
type Storage interface {
    SaveQueue(queue *domain.Queue) error
    LoadQueue() (*domain.Queue, error)
}
```

## Implementation Details

### Constructor

```go
func New(filePath string) *JSONStorage
```

- Accept the queue file path from config (`StorageConfig.QueueFile`)
- Create parent directories on first `SaveQueue` if they don't exist (`os.MkdirAll`)

### SaveQueue

- Marshal `domain.Queue` to indented JSON (`json.MarshalIndent`)
- Write atomically: write to a temp file in the same directory, then `os.Rename` to target path
- This prevents corruption from interrupted writes

### LoadQueue

- Read file with `os.ReadFile`
- If file doesn't exist (`os.IsNotExist`), return a new empty `domain.Queue` (not an error)
- Unmarshal JSON into `domain.Queue`

### JSON Schema

The serialized queue should include `tracks` and `current` fields. The `Track.StreamURL` field should be excluded from serialization (it's ephemeral — resolved at play time).

## Acceptance Criteria

- [x] `SaveQueue` writes valid JSON to disk
- [x] `LoadQueue` reads it back with identical state
- [x] Missing file returns empty queue, not error
- [x] Parent directories created automatically
- [x] Atomic writes prevent corruption

## Tests

- Round-trip: save queue with tracks, load, assert equality
- Empty queue round-trip
- Missing file returns empty queue
- Parent directory creation
- All tests use `t.TempDir()` — no real `~/.config` touched

## Estimated Effort

~1 hour
