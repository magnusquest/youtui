package port

// ThumbnailRenderer renders thumbnails as ASCII/block art for the terminal
type ThumbnailRenderer interface {
	// Render downloads and renders a thumbnail URL as terminal-compatible art
	// width and height are in character cells
	Render(url string, width, height int) (string, error)
}
