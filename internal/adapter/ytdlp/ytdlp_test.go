package ytdlp

import (
	"context"
	"os/exec"
	"strings"
	"testing"
	"time"
)

func skipIfNoYtdlp(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		t.Skip("yt-dlp not installed, skipping")
	}
}

func TestExtract(t *testing.T) {
	skipIfNoYtdlp(t)
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}

	ext := New()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Use a well-known, stable YouTube video.
	streamURL, err := ext.Extract(ctx, "dQw4w9WgXcQ")
	if err != nil {
		t.Fatalf("Extract: %v", err)
	}

	if !strings.HasPrefix(streamURL, "https://") {
		t.Errorf("expected https:// URL, got: %s", streamURL[:min(len(streamURL), 80)])
	}
}

func TestExtractInvalidVideoID(t *testing.T) {
	skipIfNoYtdlp(t)
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}

	ext := New()
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := ext.Extract(ctx, "INVALID_VIDEO_ID_XXXXXX")
	if err == nil {
		t.Fatal("expected error for invalid video ID")
	}
}

func TestExtractContextCancellation(t *testing.T) {
	skipIfNoYtdlp(t)

	ext := New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately.

	_, err := ext.Extract(ctx, "dQw4w9WgXcQ")
	if err == nil {
		t.Fatal("expected error for cancelled context")
	}
}
