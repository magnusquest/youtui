package youtube

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mganuesquest/youtui/internal/domain"
	yt "google.golang.org/api/youtube/v3"

	"google.golang.org/api/option"
)

// Client implements port.YouTubeClient using the YouTube Data API v3.
type Client struct {
	service    *yt.Service
	regionCode string
}

// New creates a YouTube API client with the given API key and region code.
func New(ctx context.Context, apiKey string, regionCode string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("youtube: API key is required")
	}

	service, err := yt.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("youtube: failed to create service: %w", err)
	}

	return &Client{
		service:    service,
		regionCode: regionCode,
	}, nil
}

// Search searches for videos and returns matching tracks with durations.
func (c *Client) Search(ctx context.Context, query string, maxResults int) ([]domain.Track, error) {
	call := c.service.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(int64(maxResults)).
		Type("video").
		RegionCode(c.regionCode).
		Context(ctx)

	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("youtube: search failed: %w", err)
	}

	if len(resp.Items) == 0 {
		return nil, nil
	}

	// Collect video IDs for the duration lookup.
	ids := make([]string, 0, len(resp.Items))
	tracks := make([]domain.Track, 0, len(resp.Items))
	for _, item := range resp.Items {
		thumb := ""
		if item.Snippet.Thumbnails != nil {
			if item.Snippet.Thumbnails.Medium != nil {
				thumb = item.Snippet.Thumbnails.Medium.Url
			} else if item.Snippet.Thumbnails.Default != nil {
				thumb = item.Snippet.Thumbnails.Default.Url
			}
		}
		tracks = append(tracks, domain.Track{
			ID:           item.Id.VideoId,
			Title:        item.Snippet.Title,
			Artist:       item.Snippet.ChannelTitle,
			ThumbnailURL: thumb,
		})
		ids = append(ids, item.Id.VideoId)
	}

	// Fetch durations in a single batch call.
	durations, err := c.fetchDurations(ctx, ids)
	if err != nil {
		// Non-fatal: return tracks without durations rather than failing.
		return tracks, nil
	}
	for i, t := range tracks {
		if d, ok := durations[t.ID]; ok {
			tracks[i].Duration = d
		}
	}

	return tracks, nil
}

// GetVideoDetails fetches detailed info for a single video ID.
func (c *Client) GetVideoDetails(ctx context.Context, videoID string) (*domain.Track, error) {
	call := c.service.Videos.List([]string{"snippet", "contentDetails"}).
		Id(videoID).
		Context(ctx)

	resp, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("youtube: video details failed: %w", err)
	}

	if len(resp.Items) == 0 {
		return nil, fmt.Errorf("youtube: video not found: %s", videoID)
	}

	item := resp.Items[0]
	thumb := ""
	if item.Snippet.Thumbnails != nil {
		if item.Snippet.Thumbnails.Medium != nil {
			thumb = item.Snippet.Thumbnails.Medium.Url
		} else if item.Snippet.Thumbnails.Default != nil {
			thumb = item.Snippet.Thumbnails.Default.Url
		}
	}

	dur, _ := parseISO8601Duration(item.ContentDetails.Duration)

	return &domain.Track{
		ID:           item.Id,
		Title:        item.Snippet.Title,
		Artist:       item.Snippet.ChannelTitle,
		ThumbnailURL: thumb,
		Duration:     dur,
	}, nil
}

// fetchDurations fetches durations for a batch of video IDs.
func (c *Client) fetchDurations(ctx context.Context, ids []string) (map[string]time.Duration, error) {
	call := c.service.Videos.List([]string{"contentDetails"}).
		Id(strings.Join(ids, ",")).
		Context(ctx)

	resp, err := call.Do()
	if err != nil {
		return nil, err
	}

	durations := make(map[string]time.Duration, len(resp.Items))
	for _, item := range resp.Items {
		if d, err := parseISO8601Duration(item.ContentDetails.Duration); err == nil {
			durations[item.Id] = d
		}
	}
	return durations, nil
}

var iso8601Re = regexp.MustCompile(`^PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?$`)

// parseISO8601Duration parses a YouTube ISO 8601 duration string like "PT4M13S".
func parseISO8601Duration(s string) (time.Duration, error) {
	matches := iso8601Re.FindStringSubmatch(s)
	if matches == nil {
		return 0, fmt.Errorf("invalid ISO 8601 duration: %s", s)
	}

	var hours, minutes, seconds int
	if matches[1] != "" {
		hours, _ = strconv.Atoi(matches[1])
	}
	if matches[2] != "" {
		minutes, _ = strconv.Atoi(matches[2])
	}
	if matches[3] != "" {
		seconds, _ = strconv.Atoi(matches[3])
	}

	return time.Duration(hours)*time.Hour +
		time.Duration(minutes)*time.Minute +
		time.Duration(seconds)*time.Second, nil
}
