package service

import (
	"sync"

	"github.com/mganuesquest/youtui/internal/domain"
	"github.com/mganuesquest/youtui/internal/port"
)

// QueueService manages the playback queue with persistence.
type QueueService struct {
	mu      sync.Mutex
	queue   *domain.Queue
	storage port.Storage
}

// NewQueueService creates a QueueService. If storage is non-nil, the queue is
// loaded from disk on creation (falling back to an empty queue on error).
func NewQueueService(storage port.Storage) *QueueService {
	q := domain.NewQueue()
	if storage != nil {
		if loaded, err := storage.LoadQueue(); err == nil {
			q = loaded
		}
	}
	return &QueueService{queue: q, storage: storage}
}

// Add appends a track to the queue and persists.
func (s *QueueService) Add(track domain.Track) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queue.Add(track)
	s.persist()
}

// Remove removes the track at the given index and persists.
func (s *QueueService) Remove(index int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ok := s.queue.Remove(index); !ok {
		return false
	}
	s.persist()
	return true
}

// MoveUp moves a track one position up and persists.
func (s *QueueService) MoveUp(index int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ok := s.queue.MoveUp(index); !ok {
		return false
	}
	s.persist()
	return true
}

// MoveDown moves a track one position down and persists.
func (s *QueueService) MoveDown(index int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if ok := s.queue.MoveDown(index); !ok {
		return false
	}
	s.persist()
	return true
}

// Clear removes all tracks and persists.
func (s *QueueService) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queue.Clear()
	s.persist()
}

// CurrentTrack returns the currently selected track, or nil.
func (s *QueueService) CurrentTrack() *domain.Track {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queue.CurrentTrack()
}

// CurrentIndex returns the current track index.
func (s *QueueService) CurrentIndex() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queue.Current
}

// Next advances to the next track. Returns false if at end.
func (s *QueueService) Next() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queue.Next()
}

// Previous goes to the previous track. Returns false if at start.
func (s *QueueService) Previous() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queue.Previous()
}

// Tracks returns a copy of the track list.
func (s *QueueService) Tracks() []domain.Track {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]domain.Track, len(s.queue.Tracks))
	copy(out, s.queue.Tracks)
	return out
}

// Len returns the number of tracks.
func (s *QueueService) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queue.Len()
}

// IsEmpty returns true if the queue is empty.
func (s *QueueService) IsEmpty() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queue.IsEmpty()
}

// SetCurrent sets the current track index directly.
func (s *QueueService) SetCurrent(index int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.queue.Tracks) {
		return false
	}
	s.queue.Current = index
	return true
}

// persist saves the queue to storage, ignoring errors (best-effort).
func (s *QueueService) persist() {
	if s.storage != nil {
		_ = s.storage.SaveQueue(s.queue)
	}
}
