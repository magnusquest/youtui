package service

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/mganuesquest/youtui/internal/domain"
)

func setupPomodoroService() (*PomodoroService, *stubPlayer) {
	sp := &stubPlayer{duration: 3 * time.Minute}
	ext := &stubExtractor{url: "http://stream/audio"}
	q := NewQueueService(nil)
	q.Add(makeTrack("t1", "Track 1"))

	playerSvc := NewPlayerService(sp, ext, q, 80)
	pomSvc := NewPomodoroService(playerSvc, "Classic")
	return pomSvc, sp
}

func TestPomodoroService_InitialState(t *testing.T) {
	svc, _ := setupPomodoroService()
	session := svc.Session()

	if session.Phase != domain.PhaseWork {
		t.Errorf("expected PhaseWork, got %d", session.Phase)
	}
	if session.Running {
		t.Error("expected not running initially")
	}
	if session.TimeRemaining != 25*time.Minute {
		t.Errorf("expected 25m, got %v", session.TimeRemaining)
	}

	profile := svc.CurrentProfile()
	if profile.Name != "Classic" {
		t.Errorf("expected Classic profile, got %s", profile.Name)
	}
}

func TestPomodoroService_Toggle(t *testing.T) {
	svc, _ := setupPomodoroService()
	ctx := context.Background()

	svc.Toggle(ctx)
	session := svc.Session()
	if !session.Running {
		t.Error("expected running after toggle")
	}

	svc.Toggle(ctx)
	session = svc.Session()
	if session.Running {
		t.Error("expected stopped after second toggle")
	}
}

func TestPomodoroService_Reset(t *testing.T) {
	svc, _ := setupPomodoroService()
	ctx := context.Background()

	svc.Start(ctx)
	// Simulate some time passing via direct manipulation.
	svc.mu.Lock()
	svc.session.TimeRemaining = 20 * time.Minute
	svc.mu.Unlock()

	svc.Reset()
	session := svc.Session()

	if session.Running {
		t.Error("expected not running after reset")
	}
	if session.TimeRemaining != 25*time.Minute {
		t.Errorf("expected 25m after reset, got %v", session.TimeRemaining)
	}
}

func TestPomodoroService_Skip(t *testing.T) {
	svc, _ := setupPomodoroService()
	ctx := context.Background()

	// Skip from Work → Short Break.
	svc.Skip(ctx)
	session := svc.Session()
	if session.Phase != domain.PhaseShortBreak {
		t.Errorf("expected ShortBreak after skip, got %d", session.Phase)
	}
	if session.TimeRemaining != 5*time.Minute {
		t.Errorf("expected 5m short break, got %v", session.TimeRemaining)
	}

	// Skip from Short Break → Work.
	svc.Skip(ctx)
	session = svc.Session()
	if session.Phase != domain.PhaseWork {
		t.Errorf("expected Work after second skip, got %d", session.Phase)
	}
}

func TestPomodoroService_SkipPausesPlayerDuringBreak(t *testing.T) {
	svc, sp := setupPomodoroService()
	ctx := context.Background()

	// Start playback first.
	svc.player.PlayCurrent(ctx)

	// Skip to break — should pause the player.
	svc.Skip(ctx)

	sp.mu.Lock()
	paused := sp.paused
	sp.mu.Unlock()
	if !paused {
		t.Error("expected player paused during break")
	}

	// Skip back to work — should resume.
	svc.Skip(ctx)
	sp.mu.Lock()
	playing := sp.playing
	sp.mu.Unlock()
	if !playing {
		t.Error("expected player resumed during work")
	}
}

func TestPomodoroService_NextPreviousProfile(t *testing.T) {
	svc, _ := setupPomodoroService()

	// Classic → Deep Work
	svc.NextProfile()
	if svc.CurrentProfile().Name != "Deep Work" {
		t.Errorf("expected Deep Work, got %s", svc.CurrentProfile().Name)
	}

	// Deep Work → Sprint
	svc.NextProfile()
	if svc.CurrentProfile().Name != "Sprint" {
		t.Errorf("expected Sprint, got %s", svc.CurrentProfile().Name)
	}

	// Sprint → Classic (wrap)
	svc.NextProfile()
	if svc.CurrentProfile().Name != "Classic" {
		t.Errorf("expected Classic (wrap), got %s", svc.CurrentProfile().Name)
	}

	// Classic → Sprint (previous wrap)
	svc.PreviousProfile()
	if svc.CurrentProfile().Name != "Sprint" {
		t.Errorf("expected Sprint (prev wrap), got %s", svc.CurrentProfile().Name)
	}
}

func TestPomodoroService_ProfileSwitchResets(t *testing.T) {
	svc, _ := setupPomodoroService()
	ctx := context.Background()

	svc.Start(ctx)
	svc.NextProfile()

	session := svc.Session()
	if session.Running {
		t.Error("expected timer stopped after profile switch")
	}
	if session.TimeRemaining != 50*time.Minute {
		t.Errorf("expected Deep Work 50m, got %v", session.TimeRemaining)
	}
}

func TestPomodoroService_Subscribe(t *testing.T) {
	svc, _ := setupPomodoroService()
	ctx := context.Background()

	var events []PomodoroEvent
	var mu sync.Mutex
	svc.Subscribe(func(s PomodoroStatus) {
		mu.Lock()
		events = append(events, s.Event)
		mu.Unlock()
	})

	svc.Start(ctx)
	svc.Pause()
	svc.Reset()

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(events) < 3 {
		t.Fatalf("expected at least 3 events, got %d", len(events))
	}
	if events[0] != PomEventStarted {
		t.Errorf("expected Started, got %d", events[0])
	}
	if events[1] != PomEventPaused {
		t.Errorf("expected Paused, got %d", events[1])
	}
	if events[2] != PomEventReset {
		t.Errorf("expected Reset, got %d", events[2])
	}
}

func TestPomodoroService_NilPlayer(t *testing.T) {
	// Should not panic when player is nil.
	svc := NewPomodoroService(nil, "Classic")
	ctx := context.Background()

	svc.Start(ctx)
	svc.Skip(ctx)
	svc.Pause()
	svc.Reset()
}

func TestPomodoroService_LongBreakAfterSessions(t *testing.T) {
	svc, _ := setupPomodoroService()
	ctx := context.Background()

	// Classic: 4 sessions before long break.
	// Skip through 4 work+break cycles: W→SB, SB→W, W→SB, SB→W, W→SB, SB→W, W→LB.
	for i := 0; i < 3; i++ {
		svc.Skip(ctx) // work → short break
		svc.Skip(ctx) // short break → work
	}
	svc.Skip(ctx) // 4th work → long break

	session := svc.Session()
	if session.Phase != domain.PhaseLongBreak {
		t.Errorf("expected LongBreak after 4 sessions, got phase %d", session.Phase)
	}
	if session.TimeRemaining != 15*time.Minute {
		t.Errorf("expected 15m long break, got %v", session.TimeRemaining)
	}
}
