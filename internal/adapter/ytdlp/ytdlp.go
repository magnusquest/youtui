package ytdlp

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Extractor implements port.StreamURLExtractor using yt-dlp.
type Extractor struct{}

// New creates a new yt-dlp stream URL extractor.
func New() *Extractor {
	return &Extractor{}
}

// Extract returns a playable audio stream URL for the given YouTube video ID.
func (e *Extractor) Extract(ctx context.Context, videoID string) (string, error) {
	if _, err := exec.LookPath("yt-dlp"); err != nil {
		return "", fmt.Errorf("yt-dlp not found: install from https://github.com/yt-dlp/yt-dlp")
	}

	url := "https://www.youtube.com/watch?v=" + videoID

	var stdout, stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "yt-dlp",
		"-f", "bestaudio",
		"--get-url",
		"--no-playlist",
		url,
	)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errMsg := strings.TrimSpace(stderr.String())
		if errMsg == "" {
			errMsg = err.Error()
		}
		return "", fmt.Errorf("yt-dlp failed: %s", errMsg)
	}

	streamURL := strings.TrimSpace(stdout.String())
	if streamURL == "" {
		return "", fmt.Errorf("yt-dlp returned no stream URL for video %s", videoID)
	}

	return streamURL, nil
}
