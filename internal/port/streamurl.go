package port

import "context"

// StreamURLExtractor extracts playable stream URLs from video IDs
type StreamURLExtractor interface {
	// Extract returns a playable audio stream URL for the given video ID
	Extract(ctx context.Context, videoID string) (string, error)
}
