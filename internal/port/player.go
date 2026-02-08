package port

import (
	"context"
	"time"
)

// AudioPlayer defines operations for audio playback control
type AudioPlayer interface {
	// Play starts playing the given stream URL
	Play(ctx context.Context, streamURL string) error

	// Pause pauses playback
	Pause(ctx context.Context) error

	// Resume resumes playback
	Resume(ctx context.Context) error

	// Stop stops playback completely
	Stop(ctx context.Context) error

	// Seek seeks to the given position
	Seek(ctx context.Context, position time.Duration) error

	// SeekRelative seeks forward or backward by the given duration
	SeekRelative(ctx context.Context, delta time.Duration) error

	// SetVolume sets the volume (0-100)
	SetVolume(ctx context.Context, volume int) error

	// GetPosition returns the current playback position
	GetPosition(ctx context.Context) (time.Duration, error)

	// GetDuration returns the total duration
	GetDuration(ctx context.Context) (time.Duration, error)

	// IsPlaying returns true if currently playing
	IsPlaying(ctx context.Context) (bool, error)

	// Close shuts down the player
	Close() error
}
