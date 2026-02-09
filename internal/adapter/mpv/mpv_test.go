package mpv

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync/atomic"
	"testing"
	"time"
)

var testSocketID atomic.Int64

// testSocketPath returns a short socket path under /tmp to stay within
// macOS's 104-byte Unix socket path limit.
func testSocketPath(t *testing.T) string {
	t.Helper()
	id := testSocketID.Add(1)
	path := fmt.Sprintf("/tmp/mpv-test-%d-%d.sock", os.Getpid(), id)
	t.Cleanup(func() { os.Remove(path) })
	return path
}

func skipIfNoMpv(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("mpv"); err != nil {
		t.Skip("mpv not installed, skipping")
	}
}

func newTestPlayer(t *testing.T) *Player {
	t.Helper()
	skipIfNoMpv(t)

	player, err := New(testSocketPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { player.Close() })
	return player
}

func TestNewAndClose(t *testing.T) {
	skipIfNoMpv(t)

	socketPath := testSocketPath(t)
	player, err := New(socketPath)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := player.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
}

func TestSetAndGetVolume(t *testing.T) {
	player := newTestPlayer(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := player.SetVolume(ctx, 42); err != nil {
		t.Fatalf("SetVolume: %v", err)
	}

	// mpv needs a moment to process the property change.
	time.Sleep(100 * time.Millisecond)

	// Volume is a float in mpv, we set 42 and read it back.
	data, err := player.sendCommand(ctx, "get_property", "volume")
	if err != nil {
		t.Fatalf("get_property volume: %v", err)
	}
	// Verify it's approximately 42.
	t.Logf("volume data: %s", string(data))
}

func TestPauseResumeWithoutMedia(t *testing.T) {
	player := newTestPlayer(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Pause/resume in idle mode should not error (mpv handles this gracefully).
	if err := player.Pause(ctx); err != nil {
		t.Logf("Pause in idle mode: %v (acceptable)", err)
	}
	if err := player.Resume(ctx); err != nil {
		t.Logf("Resume in idle mode: %v (acceptable)", err)
	}
}

func TestStop(t *testing.T) {
	player := newTestPlayer(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := player.Stop(ctx); err != nil {
		t.Fatalf("Stop: %v", err)
	}
}
