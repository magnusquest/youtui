package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mganuesquest/youtui/internal/domain"
)

// --- stubs ---

type stubYouTube struct {
	tracks []domain.Track
	err    error
	detail *domain.Track
}

func (s *stubYouTube) Search(_ context.Context, _ string, _ int) ([]domain.Track, error) {
	return s.tracks, s.err
}

func (s *stubYouTube) GetVideoDetails(_ context.Context, _ string) (*domain.Track, error) {
	if s.detail != nil {
		return s.detail, nil
	}
	return nil, fmt.Errorf("not found")
}

type stubThumbnail struct {
	output string
	err    error
}

func (s *stubThumbnail) Render(_ string, _, _ int) (string, error) {
	return s.output, s.err
}

// --- tests ---

func TestSearchService_Search(t *testing.T) {
	tracks := []domain.Track{
		{ID: "a1", Title: "Song A", Artist: "Artist", Duration: 3 * time.Minute, ThumbnailURL: "http://thumb/a"},
		{ID: "b2", Title: "Song B", Artist: "Artist", Duration: 4 * time.Minute, ThumbnailURL: "http://thumb/b"},
	}

	svc := NewSearchService(
		&stubYouTube{tracks: tracks},
		&stubThumbnail{output: "[thumb]"},
		12, 6,
	)

	results, err := svc.Search(context.Background(), "test", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Track.ID != "a1" {
		t.Errorf("expected track ID a1, got %s", results[0].Track.ID)
	}
	if results[0].Thumbnail != "[thumb]" {
		t.Errorf("expected thumbnail '[thumb]', got %q", results[0].Thumbnail)
	}
}

func TestSearchService_SearchYouTubeError(t *testing.T) {
	svc := NewSearchService(
		&stubYouTube{err: fmt.Errorf("api error")},
		nil, 12, 6,
	)

	_, err := svc.Search(context.Background(), "test", 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSearchService_ThumbnailErrorNonFatal(t *testing.T) {
	tracks := []domain.Track{
		{ID: "a1", Title: "Song", ThumbnailURL: "http://thumb/a"},
	}
	svc := NewSearchService(
		&stubYouTube{tracks: tracks},
		&stubThumbnail{err: fmt.Errorf("render failed")},
		12, 6,
	)

	results, err := svc.Search(context.Background(), "test", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Thumbnail != "" {
		t.Errorf("expected empty thumbnail on render error, got %q", results[0].Thumbnail)
	}
}

func TestSearchService_NilThumbnailRenderer(t *testing.T) {
	tracks := []domain.Track{
		{ID: "a1", Title: "Song", ThumbnailURL: "http://thumb/a"},
	}
	svc := NewSearchService(
		&stubYouTube{tracks: tracks},
		nil, 12, 6,
	)

	results, err := svc.Search(context.Background(), "test", 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Thumbnail != "" {
		t.Errorf("expected empty thumbnail with nil renderer, got %q", results[0].Thumbnail)
	}
}

func TestSearchService_GetVideoDetails(t *testing.T) {
	detail := &domain.Track{ID: "x1", Title: "Detail Track"}
	svc := NewSearchService(
		&stubYouTube{detail: detail},
		nil, 12, 6,
	)

	track, err := svc.GetVideoDetails(context.Background(), "x1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if track.ID != "x1" {
		t.Errorf("expected ID x1, got %s", track.ID)
	}
}
