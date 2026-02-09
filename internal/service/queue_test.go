package service

import (
	"testing"
	"time"

	"github.com/mganuesquest/youtui/internal/domain"
)

// --- stub storage ---

type stubStorage struct {
	saved *domain.Queue
}

func (s *stubStorage) SaveQueue(q *domain.Queue) error {
	s.saved = q
	return nil
}

func (s *stubStorage) LoadQueue() (*domain.Queue, error) {
	if s.saved != nil {
		return s.saved, nil
	}
	return domain.NewQueue(), nil
}

func makeTrack(id, title string) domain.Track {
	return domain.Track{ID: id, Title: title, Duration: 3 * time.Minute}
}

// --- tests ---

func TestQueueService_AddAndLen(t *testing.T) {
	svc := NewQueueService(nil)
	if svc.Len() != 0 {
		t.Fatalf("expected empty queue")
	}

	svc.Add(makeTrack("a", "A"))
	svc.Add(makeTrack("b", "B"))

	if svc.Len() != 2 {
		t.Fatalf("expected 2 tracks, got %d", svc.Len())
	}
	if svc.CurrentIndex() != 0 {
		t.Errorf("expected current=0, got %d", svc.CurrentIndex())
	}
}

func TestQueueService_Remove(t *testing.T) {
	svc := NewQueueService(nil)
	svc.Add(makeTrack("a", "A"))
	svc.Add(makeTrack("b", "B"))
	svc.Add(makeTrack("c", "C"))

	if !svc.Remove(1) {
		t.Fatal("remove should succeed")
	}
	if svc.Len() != 2 {
		t.Fatalf("expected 2 tracks, got %d", svc.Len())
	}
	tracks := svc.Tracks()
	if tracks[1].ID != "c" {
		t.Errorf("expected second track to be 'c', got %q", tracks[1].ID)
	}
}

func TestQueueService_RemoveOutOfBounds(t *testing.T) {
	svc := NewQueueService(nil)
	if svc.Remove(0) {
		t.Fatal("remove on empty queue should return false")
	}
	svc.Add(makeTrack("a", "A"))
	if svc.Remove(5) {
		t.Fatal("remove at bad index should return false")
	}
}

func TestQueueService_MoveUpDown(t *testing.T) {
	svc := NewQueueService(nil)
	svc.Add(makeTrack("a", "A"))
	svc.Add(makeTrack("b", "B"))
	svc.Add(makeTrack("c", "C"))

	// Move B (index 1) up to index 0.
	if !svc.MoveUp(1) {
		t.Fatal("MoveUp should succeed")
	}
	tracks := svc.Tracks()
	if tracks[0].ID != "b" {
		t.Errorf("expected first track 'b', got %q", tracks[0].ID)
	}

	// Move B (now index 0) down to index 1.
	if !svc.MoveDown(0) {
		t.Fatal("MoveDown should succeed")
	}
	tracks = svc.Tracks()
	if tracks[1].ID != "b" {
		t.Errorf("expected second track 'b', got %q", tracks[1].ID)
	}

	// Boundary cases.
	if svc.MoveUp(0) {
		t.Error("MoveUp at 0 should return false")
	}
	if svc.MoveDown(2) {
		t.Error("MoveDown at last should return false")
	}
}

func TestQueueService_Clear(t *testing.T) {
	svc := NewQueueService(nil)
	svc.Add(makeTrack("a", "A"))
	svc.Add(makeTrack("b", "B"))
	svc.Clear()

	if !svc.IsEmpty() {
		t.Fatal("expected empty queue after clear")
	}
	if svc.CurrentTrack() != nil {
		t.Fatal("expected nil current track after clear")
	}
}

func TestQueueService_Navigation(t *testing.T) {
	svc := NewQueueService(nil)
	svc.Add(makeTrack("a", "A"))
	svc.Add(makeTrack("b", "B"))
	svc.Add(makeTrack("c", "C"))

	if !svc.Next() {
		t.Fatal("Next should succeed")
	}
	if svc.CurrentTrack().ID != "b" {
		t.Errorf("expected current=b, got %s", svc.CurrentTrack().ID)
	}

	if !svc.Next() {
		t.Fatal("Next should succeed")
	}
	if svc.Next() {
		t.Fatal("Next at end should return false")
	}

	if !svc.Previous() {
		t.Fatal("Previous should succeed")
	}
	if svc.CurrentTrack().ID != "b" {
		t.Errorf("expected current=b, got %s", svc.CurrentTrack().ID)
	}
}

func TestQueueService_SetCurrent(t *testing.T) {
	svc := NewQueueService(nil)
	svc.Add(makeTrack("a", "A"))
	svc.Add(makeTrack("b", "B"))

	if !svc.SetCurrent(1) {
		t.Fatal("SetCurrent(1) should succeed")
	}
	if svc.CurrentTrack().ID != "b" {
		t.Errorf("expected current=b, got %s", svc.CurrentTrack().ID)
	}
	if svc.SetCurrent(5) {
		t.Fatal("SetCurrent(5) should return false")
	}
}

func TestQueueService_Persistence(t *testing.T) {
	store := &stubStorage{}

	svc := NewQueueService(store)
	svc.Add(makeTrack("a", "A"))
	svc.Add(makeTrack("b", "B"))

	if store.saved == nil {
		t.Fatal("expected storage to be written")
	}
	if len(store.saved.Tracks) != 2 {
		t.Fatalf("expected 2 saved tracks, got %d", len(store.saved.Tracks))
	}

	// Create a new service from the same storage — should load the queue.
	svc2 := NewQueueService(store)
	if svc2.Len() != 2 {
		t.Fatalf("expected loaded queue to have 2 tracks, got %d", svc2.Len())
	}
}

func TestQueueService_Tracks_IsCopy(t *testing.T) {
	svc := NewQueueService(nil)
	svc.Add(makeTrack("a", "A"))

	tracks := svc.Tracks()
	tracks[0].ID = "modified"

	if svc.CurrentTrack().ID == "modified" {
		t.Fatal("Tracks() should return a copy, not a reference to internal state")
	}
}
