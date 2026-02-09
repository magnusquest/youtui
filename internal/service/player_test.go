package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/mganuesquest/youtui/internal/domain"
)

// --- stubs ---

type stubPlayer struct {
	mu        sync.Mutex
	playing   bool
	paused    bool
	stopped   bool
	volume    int
	position  time.Duration
	duration  time.Duration
	streamURL string
	closed    bool
}

func (p *stubPlayer) Play(_ context.Context, url string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.streamURL = url
	p.playing = true
	p.paused = false
	p.stopped = false
	return nil
}

func (p *stubPlayer) Pause(_ context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused = true
	p.playing = false
	return nil
}

func (p *stubPlayer) Resume(_ context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.paused = false
	p.playing = true
	return nil
}

func (p *stubPlayer) Stop(_ context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.playing = false
	p.paused = false
	p.stopped = true
	return nil
}

func (p *stubPlayer) Seek(_ context.Context, _ time.Duration) error  { return nil }
func (p *stubPlayer) SeekRelative(_ context.Context, _ time.Duration) error { return nil }

func (p *stubPlayer) SetVolume(_ context.Context, v int) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.volume = v
	return nil
}

func (p *stubPlayer) GetPosition(_ context.Context) (time.Duration, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.position, nil
}

func (p *stubPlayer) GetDuration(_ context.Context) (time.Duration, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.duration, nil
}

func (p *stubPlayer) IsPlaying(_ context.Context) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.playing, nil
}

func (p *stubPlayer) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closed = true
	return nil
}

type stubExtractor struct {
	url string
	err error
}

func (e *stubExtractor) Extract(_ context.Context, _ string) (string, error) {
	return e.url, e.err
}

// --- helpers ---

func setupPlayerService() (*PlayerService, *stubPlayer) {
	sp := &stubPlayer{duration: 3 * time.Minute}
	ext := &stubExtractor{url: "http://stream/audio"}
	q := NewQueueService(nil)
	q.Add(makeTrack("t1", "Track 1"))
	q.Add(makeTrack("t2", "Track 2"))
	q.Add(makeTrack("t3", "Track 3"))

	svc := NewPlayerService(sp, ext, q, 80)
	return svc, sp
}

// --- tests ---

func TestPlayerService_PlayCurrent(t *testing.T) {
	svc, sp := setupPlayerService()
	defer svc.Close()

	ctx := context.Background()
	if err := svc.PlayCurrent(ctx); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sp.mu.Lock()
	defer sp.mu.Unlock()
	if sp.streamURL != "http://stream/audio" {
		t.Errorf("expected stream URL, got %q", sp.streamURL)
	}
	if !sp.playing {
		t.Error("expected player to be playing")
	}
	if !svc.IsPlaying() {
		t.Error("expected service to report playing")
	}
}

func TestPlayerService_TogglePause(t *testing.T) {
	svc, _ := setupPlayerService()
	defer svc.Close()

	ctx := context.Background()
	svc.PlayCurrent(ctx)

	// Pause.
	if err := svc.TogglePause(ctx); err != nil {
		t.Fatalf("toggle pause error: %v", err)
	}
	if svc.IsPlaying() {
		t.Error("expected paused")
	}

	// Resume.
	if err := svc.TogglePause(ctx); err != nil {
		t.Fatalf("toggle resume error: %v", err)
	}
	if !svc.IsPlaying() {
		t.Error("expected playing")
	}
}

func TestPlayerService_Stop(t *testing.T) {
	svc, sp := setupPlayerService()
	defer svc.Close()

	ctx := context.Background()
	svc.PlayCurrent(ctx)
	svc.Stop(ctx)

	sp.mu.Lock()
	defer sp.mu.Unlock()
	if !sp.stopped {
		t.Error("expected player stopped")
	}
}

func TestPlayerService_NextPrevious(t *testing.T) {
	svc, _ := setupPlayerService()
	defer svc.Close()

	ctx := context.Background()
	svc.PlayCurrent(ctx)

	// Next track.
	if err := svc.Next(ctx); err != nil {
		t.Fatalf("Next error: %v", err)
	}
	status := svc.Status()
	if status.Track == nil || status.Track.ID != "t2" {
		t.Errorf("expected track t2, got %v", status.Track)
	}

	// Previous (position is 0, so it should go to previous track).
	if err := svc.Previous(ctx); err != nil {
		t.Fatalf("Previous error: %v", err)
	}
	status = svc.Status()
	if status.Track == nil || status.Track.ID != "t1" {
		t.Errorf("expected track t1, got %v", status.Track)
	}
}

func TestPlayerService_VolumeUpDown(t *testing.T) {
	svc, sp := setupPlayerService()
	defer svc.Close()

	ctx := context.Background()
	svc.PlayCurrent(ctx)

	svc.VolumeUp(ctx)
	status := svc.Status()
	if status.Playback.Volume != 85 {
		t.Errorf("expected volume 85, got %d", status.Playback.Volume)
	}

	svc.VolumeDown(ctx)
	status = svc.Status()
	if status.Playback.Volume != 80 {
		t.Errorf("expected volume 80, got %d", status.Playback.Volume)
	}

	sp.mu.Lock()
	vol := sp.volume
	sp.mu.Unlock()
	if vol != 80 {
		t.Errorf("expected stub volume 80, got %d", vol)
	}
}

func TestPlayerService_ToggleMute(t *testing.T) {
	svc, sp := setupPlayerService()
	defer svc.Close()

	ctx := context.Background()
	svc.PlayCurrent(ctx)

	svc.ToggleMute(ctx)
	sp.mu.Lock()
	v := sp.volume
	sp.mu.Unlock()
	if v != 0 {
		t.Errorf("expected muted volume 0, got %d", v)
	}

	svc.ToggleMute(ctx)
	sp.mu.Lock()
	v = sp.volume
	sp.mu.Unlock()
	if v != 80 {
		t.Errorf("expected unmuted volume 80, got %d", v)
	}
}

func TestPlayerService_CycleRepeatShuffle(t *testing.T) {
	svc, _ := setupPlayerService()
	defer svc.Close()

	svc.CycleRepeat()
	status := svc.Status()
	if status.Playback.Repeat != domain.RepeatOne {
		t.Errorf("expected RepeatOne, got %d", status.Playback.Repeat)
	}

	svc.ToggleShuffle()
	status = svc.Status()
	if !status.Playback.Shuffle {
		t.Error("expected shuffle on")
	}
}

func TestPlayerService_Subscribe(t *testing.T) {
	svc, _ := setupPlayerService()
	defer svc.Close()

	var received []PlayerEvent
	var mu sync.Mutex
	svc.Subscribe(func(s PlayerStatus) {
		mu.Lock()
		received = append(received, s.Event)
		mu.Unlock()
	})

	ctx := context.Background()
	svc.PlayCurrent(ctx)
	svc.Pause(ctx)

	// Give a moment for events to propagate.
	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(received) < 2 {
		t.Fatalf("expected at least 2 events, got %d", len(received))
	}
	if received[0] != EventTrackStarted {
		t.Errorf("expected EventTrackStarted, got %d", received[0])
	}
	if received[1] != EventPaused {
		t.Errorf("expected EventPaused, got %d", received[1])
	}
}

func TestPlayerService_PlayIndex(t *testing.T) {
	svc, _ := setupPlayerService()
	defer svc.Close()

	ctx := context.Background()
	if err := svc.PlayIndex(ctx, 2); err != nil {
		t.Fatalf("PlayIndex error: %v", err)
	}
	status := svc.Status()
	if status.Track == nil || status.Track.ID != "t3" {
		t.Errorf("expected track t3, got %v", status.Track)
	}
}

func TestPlayerService_Close(t *testing.T) {
	svc, sp := setupPlayerService()
	svc.Close()

	sp.mu.Lock()
	defer sp.mu.Unlock()
	if !sp.closed {
		t.Error("expected player closed")
	}
}
