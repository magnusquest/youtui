# Task: Mosaic Thumbnail Renderer Adapter

## Implements

`port.ThumbnailRenderer` interface (`internal/port/thumbnail.go`)

## Package

`internal/adapter/mosaic/`

## Files to Create

- `mosaic.go` — adapter implementation
- `mosaic_test.go` — unit tests

## Interface to Satisfy

```go
type ThumbnailRenderer interface {
    Render(url string, width, height int) (string, error)
}
```

## Dependencies

```
github.com/charmbracelet/x/mosaic
```

## Implementation Details

### Constructor

```go
func New() *Renderer
```

No config needed — the width/height come per-call from the `Render` method.

### Render

1. **Download the image** from the URL:
   - `http.Get(url)` with a reasonable timeout (10 seconds)
   - Verify response status is 200
   - Read the response body
2. **Decode the image** using Go's `image` package:
   - Import `_ "image/jpeg"` and `_ "image/png"` for format registration
   - `image.Decode(resp.Body)` to get `image.Image`
3. **Render to terminal string** using charmbracelet/x/mosaic:
   ```go
   output := mosaic.New().
       Width(width).
       Height(height).
       Render(img)
   ```
4. Return the rendered string

### Error Handling

- Network errors → wrap with "failed to download thumbnail"
- Non-200 status → return error with status code
- Image decode errors → wrap with "failed to decode image"
- Empty URL → return empty string (no thumbnail to render), not an error

### HTTP Client

Use a package-level `http.Client` with a 10-second timeout rather than `http.DefaultClient`:

```go
var httpClient = &http.Client{Timeout: 10 * time.Second}
```

## Acceptance Criteria

- [x] Downloads a JPEG/PNG thumbnail from a URL
- [x] Renders it as Unicode block art using charmbracelet/x/mosaic
- [x] Output fits within the specified width/height character cells
- [x] Handles network errors and invalid image data gracefully
- [x] Empty URL returns empty string

## Tests

- Render a known YouTube thumbnail URL (e.g., `https://i.ytimg.com/vi/dQw4w9WgXcQ/default.jpg`)
- Verify output is non-empty string
- Verify output contains ANSI escape sequences (color data)
- Test with empty URL → empty string
- Test with invalid URL → error
- Skip network tests if offline (check with a quick connectivity probe or use `testing.Short()`)

## Estimated Effort

~1 hour
