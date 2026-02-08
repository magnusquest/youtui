package domain

import "time"

// PlaybackState represents the current playback status
type PlaybackState int

const (
	StateStopped PlaybackState = iota
	StatePlaying
	StatePaused
)

// RepeatMode represents the repeat behavior
type RepeatMode int

const (
	RepeatNone RepeatMode = iota
	RepeatOne
	RepeatAll
)

// Playback represents the current playback state
type Playback struct {
	State    PlaybackState `json:"state"`
	Position time.Duration `json:"position"`
	Volume   int           `json:"volume"` // 0-100
	Muted    bool          `json:"muted"`
	Shuffle  bool          `json:"shuffle"`
	Repeat   RepeatMode    `json:"repeat"`
}

// NewPlayback creates a new playback state with defaults
func NewPlayback() *Playback {
	return &Playback{
		State:   StateStopped,
		Volume:  80,
		Muted:   false,
		Shuffle: false,
		Repeat:  RepeatNone,
	}
}

// TogglePause toggles between playing and paused
func (p *Playback) TogglePause() {
	if p.State == StatePlaying {
		p.State = StatePaused
	} else if p.State == StatePaused {
		p.State = StatePlaying
	}
}

// Play sets state to playing
func (p *Playback) Play() {
	p.State = StatePlaying
}

// Pause sets state to paused
func (p *Playback) Pause() {
	p.State = StatePaused
}

// Stop sets state to stopped and resets position
func (p *Playback) Stop() {
	p.State = StateStopped
	p.Position = 0
}

// ToggleMute toggles mute state
func (p *Playback) ToggleMute() {
	p.Muted = !p.Muted
}

// ToggleShuffle toggles shuffle mode
func (p *Playback) ToggleShuffle() {
	p.Shuffle = !p.Shuffle
}

// CycleRepeat cycles through repeat modes
func (p *Playback) CycleRepeat() {
	p.Repeat = (p.Repeat + 1) % 3
}

// VolumeUp increases volume by 5, max 100
func (p *Playback) VolumeUp() {
	p.Volume += 5
	if p.Volume > 100 {
		p.Volume = 100
	}
}

// VolumeDown decreases volume by 5, min 0
func (p *Playback) VolumeDown() {
	p.Volume -= 5
	if p.Volume < 0 {
		p.Volume = 0
	}
}

// IsPlaying returns true if currently playing
func (p *Playback) IsPlaying() bool {
	return p.State == StatePlaying
}
