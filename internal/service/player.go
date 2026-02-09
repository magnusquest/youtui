package service

import (
	"context"
	"sync"
	"time"

	"github.com/mganuesquest/youtui/internal/domain"
	"github.com/mganuesquest/youtui/internal/port"
)

// PlayerEvent is emitted when playback state changes.
type PlayerEvent int

const (
	EventTrackStarted PlayerEvent = iota
	EventTrackFinished
	EventPaused
	EventResumed
	EventStopped
	EventPositionUpdate
	EventVolumeChanged
	EventError
)

// PlayerStatus holds a snapshot of the current player state.
type PlayerStatus struct {
	Event    PlayerEvent
	Track    *domain.Track
	Playback *domain.Playback
	Err      error
}

// PlayerService orchestrates audio playback with queue management.
// It bridges the AudioPlayer and StreamURLExtractor ports with
// the QueueService, managing the full lifecycle of track playback.
type PlayerService struct {
	player    port.AudioPlayer
	extractor port.StreamURLExtractor
	queue     *QueueService

	mu       sync.Mutex
	playback *domain.Playback
	current  *domain.Track // the track currently loaded in the player

	// Subscribers receive player events. Callbacks are invoked synchronously
	// from the service goroutine, so they should be non-blocking.
	subMu       sync.Mutex
	subscribers []func(PlayerStatus)

	stopPoll chan struct{}
}

// NewPlayerService creates a PlayerService with the given dependencies.
func NewPlayerService(
	player port.AudioPlayer,
	extractor port.StreamURLExtractor,
	queue *QueueService,
	defaultVolume int,
) *PlayerService {
	pb := domain.NewPlayback()
	pb.Volume = defaultVolume
	return &PlayerService{
		player:    player,
		extractor: extractor,
		queue:     queue,
		playback:  pb,
	}
}

// Subscribe registers a callback for player events.
func (s *PlayerService) Subscribe(fn func(PlayerStatus)) {
	s.subMu.Lock()
	defer s.subMu.Unlock()
	s.subscribers = append(s.subscribers, fn)
}

func (s *PlayerService) emit(evt PlayerEvent, err error) {
	s.subMu.Lock()
	subs := make([]func(PlayerStatus), len(s.subscribers))
	copy(subs, s.subscribers)
	s.subMu.Unlock()

	status := PlayerStatus{
		Event:    evt,
		Track:    s.current,
		Playback: s.playback,
		Err:      err,
	}
	for _, fn := range subs {
		fn(status)
	}
}

// PlayCurrent extracts the stream URL for the current queue track and starts playback.
func (s *PlayerService) PlayCurrent(ctx context.Context) error {
	track := s.queue.CurrentTrack()
	if track == nil {
		return nil
	}

	streamURL, err := s.extractor.Extract(ctx, track.ID)
	if err != nil {
		s.emit(EventError, err)
		return err
	}
	track.StreamURL = streamURL

	if err := s.player.Play(ctx, streamURL); err != nil {
		s.emit(EventError, err)
		return err
	}

	if err := s.player.SetVolume(ctx, s.playback.Volume); err != nil {
		// Non-fatal: continue playback
	}

	s.mu.Lock()
	s.current = track
	s.playback.Play()
	s.mu.Unlock()

	s.startProgressPoll(ctx)
	s.emit(EventTrackStarted, nil)
	return nil
}

// PlayIndex sets the queue to the given index and plays that track.
func (s *PlayerService) PlayIndex(ctx context.Context, index int) error {
	if !s.queue.SetCurrent(index) {
		return nil
	}
	return s.PlayCurrent(ctx)
}

// TogglePause pauses if playing, resumes if paused.
func (s *PlayerService) TogglePause(ctx context.Context) error {
	s.mu.Lock()
	state := s.playback.State
	s.mu.Unlock()

	switch state {
	case domain.StatePlaying:
		return s.Pause(ctx)
	case domain.StatePaused:
		return s.Resume(ctx)
	}
	return nil
}

// Pause pauses playback.
func (s *PlayerService) Pause(ctx context.Context) error {
	if err := s.player.Pause(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	s.playback.Pause()
	s.mu.Unlock()
	s.emit(EventPaused, nil)
	return nil
}

// Resume resumes playback.
func (s *PlayerService) Resume(ctx context.Context) error {
	if err := s.player.Resume(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	s.playback.Play()
	s.mu.Unlock()
	s.emit(EventResumed, nil)
	return nil
}

// Stop stops playback completely.
func (s *PlayerService) Stop(ctx context.Context) error {
	s.stopProgressPoll()
	if err := s.player.Stop(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	s.playback.Stop()
	s.current = nil
	s.mu.Unlock()
	s.emit(EventStopped, nil)
	return nil
}

// Next advances to the next track and plays it.
func (s *PlayerService) Next(ctx context.Context) error {
	s.stopProgressPoll()

	s.mu.Lock()
	repeat := s.playback.Repeat
	s.mu.Unlock()

	switch repeat {
	case domain.RepeatOne:
		// Replay the same track.
		return s.PlayCurrent(ctx)
	case domain.RepeatAll:
		if !s.queue.Next() {
			// Wrap around to the beginning.
			s.queue.SetCurrent(0)
		}
		return s.PlayCurrent(ctx)
	default:
		if !s.queue.Next() {
			// End of queue — stop.
			return s.Stop(ctx)
		}
		return s.PlayCurrent(ctx)
	}
}

// Previous goes to the previous track and plays it.
func (s *PlayerService) Previous(ctx context.Context) error {
	s.stopProgressPoll()

	// If we're more than 3 seconds in, restart the current track instead.
	s.mu.Lock()
	pos := s.playback.Position
	s.mu.Unlock()
	if pos > 3*time.Second {
		return s.PlayCurrent(ctx)
	}

	if !s.queue.Previous() {
		// Already at the first track — just restart.
		return s.PlayCurrent(ctx)
	}
	return s.PlayCurrent(ctx)
}

// SeekRelative seeks forward or backward.
func (s *PlayerService) SeekRelative(ctx context.Context, delta time.Duration) error {
	return s.player.SeekRelative(ctx, delta)
}

// SetVolume sets the volume and syncs to the player.
func (s *PlayerService) SetVolume(ctx context.Context, volume int) error {
	if err := s.player.SetVolume(ctx, volume); err != nil {
		return err
	}
	s.mu.Lock()
	s.playback.Volume = volume
	s.mu.Unlock()
	s.emit(EventVolumeChanged, nil)
	return nil
}

// VolumeUp increases volume by 5.
func (s *PlayerService) VolumeUp(ctx context.Context) error {
	s.mu.Lock()
	s.playback.VolumeUp()
	vol := s.playback.Volume
	s.mu.Unlock()
	return s.player.SetVolume(ctx, vol)
}

// VolumeDown decreases volume by 5.
func (s *PlayerService) VolumeDown(ctx context.Context) error {
	s.mu.Lock()
	s.playback.VolumeDown()
	vol := s.playback.Volume
	s.mu.Unlock()
	return s.player.SetVolume(ctx, vol)
}

// ToggleMute toggles mute on/off.
func (s *PlayerService) ToggleMute(ctx context.Context) error {
	s.mu.Lock()
	s.playback.ToggleMute()
	muted := s.playback.Muted
	vol := s.playback.Volume
	s.mu.Unlock()

	if muted {
		return s.player.SetVolume(ctx, 0)
	}
	return s.player.SetVolume(ctx, vol)
}

// ToggleShuffle toggles shuffle mode.
func (s *PlayerService) ToggleShuffle() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.playback.ToggleShuffle()
}

// CycleRepeat cycles through repeat modes.
func (s *PlayerService) CycleRepeat() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.playback.CycleRepeat()
}

// Status returns a snapshot of the current player state.
func (s *PlayerService) Status() PlayerStatus {
	s.mu.Lock()
	defer s.mu.Unlock()
	return PlayerStatus{
		Track:    s.current,
		Playback: s.playback,
	}
}

// IsPlaying returns true if currently playing.
func (s *PlayerService) IsPlaying() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.playback.IsPlaying()
}

// Close shuts down the underlying player.
func (s *PlayerService) Close() error {
	s.stopProgressPoll()
	return s.player.Close()
}

// startProgressPoll begins a background goroutine that polls the player
// for the current position every second and emits progress events.
func (s *PlayerService) startProgressPoll(ctx context.Context) {
	s.stopProgressPoll()

	stop := make(chan struct{})
	s.mu.Lock()
	s.stopPoll = stop
	s.mu.Unlock()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				pos, err := s.player.GetPosition(ctx)
				if err != nil {
					continue
				}
				dur, _ := s.player.GetDuration(ctx)

				s.mu.Lock()
				s.playback.Position = pos
				if s.current != nil && dur > 0 {
					s.current.Duration = dur
				}
				s.mu.Unlock()

				s.emit(EventPositionUpdate, nil)

				// Detect track completion: position reached or exceeded duration.
				if dur > 0 && pos >= dur-time.Second {
					s.Next(ctx)
					return
				}
			case <-stop:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *PlayerService) stopProgressPoll() {
	s.mu.Lock()
	if s.stopPoll != nil {
		close(s.stopPoll)
		s.stopPoll = nil
	}
	s.mu.Unlock()
}
