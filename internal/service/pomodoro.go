package service

import (
	"context"
	"sync"
	"time"

	"github.com/mganuesquest/youtui/internal/domain"
)

// PomodoroEvent is emitted when the pomodoro timer state changes.
type PomodoroEvent int

const (
	PomEventStarted PomodoroEvent = iota
	PomEventPaused
	PomEventReset
	PomEventSkipped
	PomEventPhaseCompleted // work → break or break → work
	PomEventTick
)

// PomodoroStatus holds a snapshot of the pomodoro state.
type PomodoroStatus struct {
	Event   PomodoroEvent
	Session *domain.PomodoroSession
	Profile domain.PomodoroProfile
}

// PomodoroService manages the pomodoro timer and coordinates with
// the PlayerService to pause music during breaks and resume during work.
type PomodoroService struct {
	mu       sync.Mutex
	session  *domain.PomodoroSession
	profiles []domain.PomodoroProfile
	profIdx  int

	player *PlayerService

	subMu       sync.Mutex
	subscribers []func(PomodoroStatus)

	stopTick chan struct{}
}

// NewPomodoroService creates a PomodoroService. If player is non-nil,
// playback will be paused during breaks and resumed during work phases.
func NewPomodoroService(player *PlayerService, defaultProfile string) *PomodoroService {
	profiles := domain.DefaultProfiles()
	idx := 0
	for i, p := range profiles {
		if p.Name == defaultProfile {
			idx = i
			break
		}
	}

	return &PomodoroService{
		session:  domain.NewPomodoroSession(profiles[idx]),
		profiles: profiles,
		profIdx:  idx,
		player:   player,
	}
}

// Subscribe registers a callback for pomodoro events.
func (s *PomodoroService) Subscribe(fn func(PomodoroStatus)) {
	s.subMu.Lock()
	defer s.subMu.Unlock()
	s.subscribers = append(s.subscribers, fn)
}

func (s *PomodoroService) emit(evt PomodoroEvent) {
	s.subMu.Lock()
	subs := make([]func(PomodoroStatus), len(s.subscribers))
	copy(subs, s.subscribers)
	s.subMu.Unlock()

	s.mu.Lock()
	status := PomodoroStatus{
		Event:   evt,
		Session: s.session,
		Profile: s.profiles[s.profIdx],
	}
	s.mu.Unlock()

	for _, fn := range subs {
		fn(status)
	}
}

// Toggle starts or pauses the timer.
func (s *PomodoroService) Toggle(ctx context.Context) {
	s.mu.Lock()
	running := s.session.Running
	s.mu.Unlock()

	if running {
		s.Pause()
	} else {
		s.Start(ctx)
	}
}

// Start begins the timer and starts the tick loop.
func (s *PomodoroService) Start(ctx context.Context) {
	s.mu.Lock()
	s.session.Start()
	isWork := s.session.Phase == domain.PhaseWork
	s.mu.Unlock()

	// Resume music during work, pause during breaks.
	if s.player != nil {
		if isWork && !s.player.IsPlaying() {
			s.player.Resume(ctx)
		} else if !isWork && s.player.IsPlaying() {
			s.player.Pause(ctx)
		}
	}

	s.startTickLoop(ctx)
	s.emit(PomEventStarted)
}

// Pause stops the timer.
func (s *PomodoroService) Pause() {
	s.stopTickLoop()
	s.mu.Lock()
	s.session.Pause()
	s.mu.Unlock()
	s.emit(PomEventPaused)
}

// Reset restarts the current phase.
func (s *PomodoroService) Reset() {
	s.stopTickLoop()
	s.mu.Lock()
	s.session.Reset()
	s.mu.Unlock()
	s.emit(PomEventReset)
}

// Skip advances to the next phase.
func (s *PomodoroService) Skip(ctx context.Context) {
	s.stopTickLoop()
	s.mu.Lock()
	s.session.Skip()
	isWork := s.session.Phase == domain.PhaseWork
	s.mu.Unlock()

	s.handlePhaseTransition(ctx, isWork)
	s.emit(PomEventSkipped)
}

// NextProfile switches to the next profile and resets the session.
func (s *PomodoroService) NextProfile() {
	s.stopTickLoop()
	s.mu.Lock()
	s.profIdx = (s.profIdx + 1) % len(s.profiles)
	s.session = domain.NewPomodoroSession(s.profiles[s.profIdx])
	s.mu.Unlock()
	s.emit(PomEventReset)
}

// PreviousProfile switches to the previous profile and resets the session.
func (s *PomodoroService) PreviousProfile() {
	s.stopTickLoop()
	s.mu.Lock()
	s.profIdx--
	if s.profIdx < 0 {
		s.profIdx = len(s.profiles) - 1
	}
	s.session = domain.NewPomodoroSession(s.profiles[s.profIdx])
	s.mu.Unlock()
	s.emit(PomEventReset)
}

// Session returns a copy of the current session.
func (s *PomodoroService) Session() domain.PomodoroSession {
	s.mu.Lock()
	defer s.mu.Unlock()
	return *s.session
}

// CurrentProfile returns the active profile.
func (s *PomodoroService) CurrentProfile() domain.PomodoroProfile {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.profiles[s.profIdx]
}

// handlePhaseTransition coordinates playback with phase changes.
func (s *PomodoroService) handlePhaseTransition(ctx context.Context, isWork bool) {
	if s.player == nil {
		return
	}
	if isWork {
		// Starting a work phase → resume music.
		s.player.Resume(ctx)
	} else {
		// Starting a break → pause music.
		s.player.Pause(ctx)
	}
}

// startTickLoop starts a goroutine that ticks the timer every second.
func (s *PomodoroService) startTickLoop(ctx context.Context) {
	s.stopTickLoop()

	stop := make(chan struct{})
	s.mu.Lock()
	s.stopTick = stop
	s.mu.Unlock()

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.mu.Lock()
				phaseCompleted := s.session.Tick()
				isWork := s.session.Phase == domain.PhaseWork
				s.mu.Unlock()

				if phaseCompleted {
					s.handlePhaseTransition(ctx, isWork)
					s.emit(PomEventPhaseCompleted)
				} else {
					s.emit(PomEventTick)
				}
			case <-stop:
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

func (s *PomodoroService) stopTickLoop() {
	s.mu.Lock()
	if s.stopTick != nil {
		close(s.stopTick)
		s.stopTick = nil
	}
	s.mu.Unlock()
}
