# Task: YouTube Data API v3 Client Adapter

## Implements

`port.YouTubeClient` interface (`internal/port/youtube.go`)

## Package

`internal/adapter/youtube/`

## Files to Create

- `youtube.go` — adapter implementation
- `youtube_test.go` — unit tests

## Interface to Satisfy

```go
type YouTubeClient interface {
    Search(ctx context.Context, query string, maxResults int) ([]domain.Track, error)
    GetVideoDetails(ctx context.Context, videoID string) (*domain.Track, error)
}
```

## Dependencies

```
google.golang.org/api/youtube/v3
google.golang.org/api/option
```

## Implementation Details

### Constructor

```go
func New(apiKey string, regionCode string) (*Client, error)
```

- Create a `youtube.Service` using `youtube.NewService(ctx, option.WithAPIKey(apiKey))`
- Store the service and region code
- Return error if API key is empty

### Search

1. Call `service.Search.List([]string{"id", "snippet"})` with:
   - `.Q(query)` — search term
   - `.MaxResults(int64(maxResults))` — limit
   - `.Type("video")` — videos only
   - `.RegionCode(regionCode)` — from config
   - `.Context(ctx)` — cancellation support
2. Map results to `[]domain.Track`:
   - `ID` = `item.Id.VideoId`
   - `Title` = `item.Snippet.Title`
   - `Artist` = `item.Snippet.ChannelTitle`
   - `ThumbnailURL` = `item.Snippet.Thumbnails.Default.Url` (or Medium/High if available)
   - `Duration` = zero (Search API doesn't return duration)
3. For duration, make a follow-up `Videos.List` call with the video IDs from results:
   - `service.Videos.List([]string{"contentDetails"}).Id(ids...).Context(ctx)`
   - Parse ISO 8601 duration from `contentDetails.duration` into `time.Duration`
   - Merge durations back into the Track slice

### GetVideoDetails

1. Call `service.Videos.List([]string{"snippet", "contentDetails"})` with:
   - `.Id(videoID)`
   - `.Context(ctx)`
2. Map the single result to `*domain.Track`
3. If no results, return `nil, ErrVideoNotFound`

### ISO 8601 Duration Parsing

YouTube returns durations like `PT4M13S`, `PT1H2M30S`. Write a helper:

```go
func parseISO8601Duration(d string) (time.Duration, error)
```

Parse `PT[nH][nM][nS]` format using regex or manual parsing.

## Acceptance Criteria

- [x] Search returns tracks with title, artist, thumbnail URL, and duration
- [x] GetVideoDetails returns full track metadata
- [x] ISO 8601 duration parsing handles hours, minutes, seconds
- [x] Empty API key returns clear error at construction time
- [x] Context cancellation propagated to API calls

## Tests

- Search for a well-known query (e.g., "Rick Astley Never Gonna Give You Up"), verify results
- GetVideoDetails for a known video ID (e.g., `dQw4w9WgXcQ`), verify metadata
- ISO 8601 parsing: `PT4M13S` → 4m13s, `PT1H2M30S` → 1h2m30s, `PT30S` → 30s
- Empty API key → error at construction
- Skip API tests if `YOUTUBE_API_KEY` env not set (`t.Skip`)

## Estimated Effort

~2 hours
