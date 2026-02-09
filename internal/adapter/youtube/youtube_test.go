package youtube

import (
	"context"
	"os"
	"testing"
	"time"
)

func skipIfNoAPIKey(t *testing.T) string {
	t.Helper()
	key := os.Getenv("YOUTUBE_API_KEY")
	if key == "" {
		t.Skip("YOUTUBE_API_KEY not set, skipping")
	}
	return key
}

func TestSearch(t *testing.T) {
	key := skipIfNoAPIKey(t)
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}

	client, err := New(context.Background(), key, "US")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	tracks, err := client.Search(ctx, "Rick Astley Never Gonna Give You Up", 5)
	if err != nil {
		t.Fatalf("Search: %v", err)
	}
	if len(tracks) == 0 {
		t.Fatal("expected at least one result")
	}

	tr := tracks[0]
	if tr.ID == "" {
		t.Error("expected non-empty video ID")
	}
	if tr.Title == "" {
		t.Error("expected non-empty title")
	}
	if tr.Artist == "" {
		t.Error("expected non-empty artist")
	}
}

func TestGetVideoDetails(t *testing.T) {
	key := skipIfNoAPIKey(t)
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}

	client, err := New(context.Background(), key, "US")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	track, err := client.GetVideoDetails(ctx, "dQw4w9WgXcQ")
	if err != nil {
		t.Fatalf("GetVideoDetails: %v", err)
	}
	if track.ID != "dQw4w9WgXcQ" {
		t.Errorf("expected ID dQw4w9WgXcQ, got %s", track.ID)
	}
	if track.Duration == 0 {
		t.Error("expected non-zero duration")
	}
}

func TestNewEmptyAPIKey(t *testing.T) {
	_, err := New(context.Background(), "", "US")
	if err == nil {
		t.Fatal("expected error for empty API key")
	}
}

func TestParseISO8601Duration(t *testing.T) {
	tests := []struct {
		input string
		want  time.Duration
	}{
		{"PT4M13S", 4*time.Minute + 13*time.Second},
		{"PT1H2M30S", 1*time.Hour + 2*time.Minute + 30*time.Second},
		{"PT30S", 30 * time.Second},
		{"PT5M", 5 * time.Minute},
		{"PT2H", 2 * time.Hour},
		{"PT1H30S", 1*time.Hour + 30*time.Second},
	}

	for _, tt := range tests {
		got, err := parseISO8601Duration(tt.input)
		if err != nil {
			t.Errorf("parseISO8601Duration(%q): %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("parseISO8601Duration(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseISO8601DurationInvalid(t *testing.T) {
	invalid := []string{"", "P4M13S", "4M13S", "PTXYZ", "hello"}
	for _, s := range invalid {
		if _, err := parseISO8601Duration(s); err == nil {
			t.Errorf("expected error for %q", s)
		}
	}
}
