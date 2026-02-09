package mosaic

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"time"

	"github.com/charmbracelet/x/mosaic"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

// Renderer implements port.ThumbnailRenderer using charmbracelet/x/mosaic.
type Renderer struct{}

// New creates a new mosaic thumbnail renderer.
func New() *Renderer {
	return &Renderer{}
}

// Render downloads the image at the given URL and renders it as Unicode block art.
func (r *Renderer) Render(url string, width, height int) (string, error) {
	if url == "" {
		return "", nil
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return "", fmt.Errorf("mosaic: failed to download thumbnail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("mosaic: thumbnail download returned status %d", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return "", fmt.Errorf("mosaic: failed to decode image: %w", err)
	}

	m := mosaic.New().
		Width(width).
		Height(height)
	output := m.Render(img)

	return output, nil
}
