package domain

import "time"

// PomodoroPhase represents the current phase
type PomodoroPhase int

const (
	PhaseWork PomodoroPhase = iota
	PhaseShortBreak
	PhaseLongBreak
)

// PomodoroProfile defines timer durations
type PomodoroProfile struct {
	Name                   string        `json:"name" yaml:"name"`
	WorkDuration           time.Duration `json:"work_duration" yaml:"work_minutes"`
	ShortBreakDuration     time.Duration `json:"short_break_duration" yaml:"short_break_minutes"`
	LongBreakDuration      time.Duration `json:"long_break_duration" yaml:"long_break_minutes"`
	SessionsBeforeLongBreak int           `json:"sessions_before_long_break" yaml:"sessions_before_long_break"`
}

// DefaultProfiles returns the built-in pomodoro profiles
func DefaultProfiles() []PomodoroProfile {
	return []PomodoroProfile{
		{
			Name:                   "Classic",
			WorkDuration:           25 * time.Minute,
			ShortBreakDuration:     5 * time.Minute,
			LongBreakDuration:      15 * time.Minute,
			SessionsBeforeLongBreak: 4,
		},
		{
			Name:                   "Deep Work",
			WorkDuration:           50 * time.Minute,
			ShortBreakDuration:     10 * time.Minute,
			LongBreakDuration:      30 * time.Minute,
			SessionsBeforeLongBreak: 2,
		},
		{
			Name:                   "Sprint",
			WorkDuration:           15 * time.Minute,
			ShortBreakDuration:     3 * time.Minute,
			LongBreakDuration:      10 * time.Minute,
			SessionsBeforeLongBreak: 4,
		},
	}
}

// PomodoroSession tracks the current timer state
type PomodoroSession struct {
	Profile          PomodoroProfile `json:"profile"`
	Phase            PomodoroPhase   `json:"phase"`
	TimeRemaining    time.Duration   `json:"time_remaining"`
	CompletedSessions int             `json:"completed_sessions"`
	Running          bool            `json:"running"`
}

// NewPomodoroSession creates a new session with the given profile
func NewPomodoroSession(profile PomodoroProfile) *PomodoroSession {
	return &PomodoroSession{
		Profile:       profile,
		Phase:         PhaseWork,
		TimeRemaining: profile.WorkDuration,
		Running:       false,
	}
}

// Start begins the timer
func (s *PomodoroSession) Start() {
	s.Running = true
}

// Pause stops the timer
func (s *PomodoroSession) Pause() {
	s.Running = false
}

// Toggle switches between running and paused
func (s *PomodoroSession) Toggle() {
	s.Running = !s.Running
}

// Reset restarts the current phase
func (s *PomodoroSession) Reset() {
	s.TimeRemaining = s.currentPhaseDuration()
	s.Running = false
}

// Skip advances to the next phase
func (s *PomodoroSession) Skip() {
	s.advancePhase()
}

// Tick decrements the timer by one second, returns true if phase completed
func (s *PomodoroSession) Tick() bool {
	if !s.Running {
		return false
	}
	s.TimeRemaining -= time.Second
	if s.TimeRemaining <= 0 {
		s.advancePhase()
		return true
	}
	return false
}

func (s *PomodoroSession) advancePhase() {
	switch s.Phase {
	case PhaseWork:
		s.CompletedSessions++
		if s.CompletedSessions%s.Profile.SessionsBeforeLongBreak == 0 {
			s.Phase = PhaseLongBreak
			s.TimeRemaining = s.Profile.LongBreakDuration
		} else {
			s.Phase = PhaseShortBreak
			s.TimeRemaining = s.Profile.ShortBreakDuration
		}
	case PhaseShortBreak, PhaseLongBreak:
		s.Phase = PhaseWork
		s.TimeRemaining = s.Profile.WorkDuration
	}
	s.Running = false
}

func (s *PomodoroSession) currentPhaseDuration() time.Duration {
	switch s.Phase {
	case PhaseWork:
		return s.Profile.WorkDuration
	case PhaseShortBreak:
		return s.Profile.ShortBreakDuration
	case PhaseLongBreak:
		return s.Profile.LongBreakDuration
	}
	return s.Profile.WorkDuration
}

// PhaseName returns a human-readable name for the current phase
func (s *PomodoroSession) PhaseName() string {
	switch s.Phase {
	case PhaseWork:
		return "Work"
	case PhaseShortBreak:
		return "Short Break"
	case PhaseLongBreak:
		return "Long Break"
	}
	return "Unknown"
}

// FormatTimeRemaining returns time as MM:SS
func (s *PomodoroSession) FormatTimeRemaining() string {
	minutes := int(s.TimeRemaining.Minutes())
	seconds := int(s.TimeRemaining.Seconds()) % 60
	return pad(minutes) + ":" + pad(seconds)
}

// IsBreak returns true if currently in a break phase
func (s *PomodoroSession) IsBreak() bool {
	return s.Phase == PhaseShortBreak || s.Phase == PhaseLongBreak
}
