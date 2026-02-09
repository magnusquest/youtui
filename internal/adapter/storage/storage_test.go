package storage

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/mganuesquest/youtui/internal/domain"
)

func TestRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "queue.json")
	s := New(path)

	q := domain.NewQueue()
	q.Add(domain.Track{
		ID:           "abc123",
		Title:        "Test Track",
		Artist:       "Test Artist",
		Duration:     3*time.Minute + 45*time.Second,
		ThumbnailURL: "https://example.com/thumb.jpg",
		StreamURL:    "https://should-not-persist.com/stream",
	})
	q.Add(domain.Track{
		ID:     "def456",
		Title:  "Second Track",
		Artist: "Another Artist",
	})

	if err := s.SaveQueue(q); err != nil {
		t.Fatalf("SaveQueue: %v", err)
	}

	loaded, err := s.LoadQueue()
	if err != nil {
		t.Fatalf("LoadQueue: %v", err)
	}

	if loaded.Len() != 2 {
		t.Fatalf("expected 2 tracks, got %d", loaded.Len())
	}
	if loaded.Current != 0 {
		t.Fatalf("expected current=0, got %d", loaded.Current)
	}

	tr := loaded.Tracks[0]
	if tr.ID != "abc123" {
		t.Errorf("expected ID abc123, got %s", tr.ID)
	}
	if tr.Title != "Test Track" {
		t.Errorf("expected title Test Track, got %s", tr.Title)
	}
	if tr.Duration != 3*time.Minute+45*time.Second {
		t.Errorf("expected duration 3m45s, got %v", tr.Duration)
	}
	if tr.StreamURL != "" {
		t.Errorf("StreamURL should not be persisted, got %s", tr.StreamURL)
	}
}

func TestEmptyQueueRoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "queue.json")
	s := New(path)

	q := domain.NewQueue()
	if err := s.SaveQueue(q); err != nil {
		t.Fatalf("SaveQueue: %v", err)
	}

	loaded, err := s.LoadQueue()
	if err != nil {
		t.Fatalf("LoadQueue: %v", err)
	}
	if !loaded.IsEmpty() {
		t.Fatalf("expected empty queue, got %d tracks", loaded.Len())
	}
	if loaded.Current != -1 {
		t.Fatalf("expected current=-1, got %d", loaded.Current)
	}
}

func TestMissingFileReturnsEmptyQueue(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent", "queue.json")
	s := New(path)

	loaded, err := s.LoadQueue()
	if err != nil {
		t.Fatalf("LoadQueue on missing file should not error: %v", err)
	}
	if !loaded.IsEmpty() {
		t.Fatalf("expected empty queue for missing file")
	}
}

func TestCreatesParentDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "deep", "queue.json")
	s := New(path)

	q := domain.NewQueue()
	q.Add(domain.Track{ID: "xyz", Title: "Nested"})

	if err := s.SaveQueue(q); err != nil {
		t.Fatalf("SaveQueue should create parent dirs: %v", err)
	}

	loaded, err := s.LoadQueue()
	if err != nil {
		t.Fatalf("LoadQueue: %v", err)
	}
	if loaded.Len() != 1 {
		t.Fatalf("expected 1 track, got %d", loaded.Len())
	}
}
