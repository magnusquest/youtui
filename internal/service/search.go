package service

import (
	"context"

	"github.com/mganuesquest/youtui/internal/domain"
	"github.com/mganuesquest/youtui/internal/port"
)

// SearchResult wraps a track with its rendered thumbnail.
type SearchResult struct {
	Track     domain.Track
	Thumbnail string // rendered terminal art, may be empty
}

// SearchService orchestrates YouTube search and thumbnail rendering.
type SearchService struct {
	youtube   port.YouTubeClient
	thumbnail port.ThumbnailRenderer

	thumbWidth  int
	thumbHeight int
}

// NewSearchService creates a SearchService with the given dependencies.
func NewSearchService(yt port.YouTubeClient, thumb port.ThumbnailRenderer, thumbW, thumbH int) *SearchService {
	return &SearchService{
		youtube:     yt,
		thumbnail:   thumb,
		thumbWidth:  thumbW,
		thumbHeight: thumbH,
	}
}

// Search queries YouTube and returns results with optional thumbnails.
func (s *SearchService) Search(ctx context.Context, query string, maxResults int) ([]SearchResult, error) {
	tracks, err := s.youtube.Search(ctx, query, maxResults)
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, len(tracks))
	for i, t := range tracks {
		results[i] = SearchResult{Track: t}

		// Best-effort thumbnail rendering — never fail the search on it.
		if s.thumbnail != nil && t.ThumbnailURL != "" {
			rendered, err := s.thumbnail.Render(t.ThumbnailURL, s.thumbWidth, s.thumbHeight)
			if err == nil {
				results[i].Thumbnail = rendered
			}
		}
	}

	return results, nil
}

// GetVideoDetails fetches full details for a single video.
func (s *SearchService) GetVideoDetails(ctx context.Context, videoID string) (*domain.Track, error) {
	return s.youtube.GetVideoDetails(ctx, videoID)
}
