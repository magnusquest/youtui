package port

import (
	"context"

	"github.com/mganuesquest/youtui/internal/domain"
)

// YouTubeClient defines operations for YouTube API interaction
type YouTubeClient interface {
	// Search searches for videos and returns matching tracks
	Search(ctx context.Context, query string, maxResults int) ([]domain.Track, error)

	// GetVideoDetails fetches detailed info for a video ID
	GetVideoDetails(ctx context.Context, videoID string) (*domain.Track, error)
}
