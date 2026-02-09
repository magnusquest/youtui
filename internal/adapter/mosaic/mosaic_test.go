package mosaic

import (
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}

	r := New()

	// Use a well-known YouTube thumbnail.
	output, err := r.Render("https://i.ytimg.com/vi/dQw4w9WgXcQ/default.jpg", 20, 10)
	if err != nil {
		t.Fatalf("Render: %v", err)
	}
	if output == "" {
		t.Fatal("expected non-empty output")
	}
	// The output should contain ANSI escape sequences for color.
	if !strings.Contains(output, "\x1b[") {
		t.Error("expected ANSI escape sequences in output")
	}
}

func TestRenderEmptyURL(t *testing.T) {
	r := New()

	output, err := r.Render("", 20, 10)
	if err != nil {
		t.Fatalf("Render with empty URL should not error: %v", err)
	}
	if output != "" {
		t.Errorf("expected empty output for empty URL, got %q", output)
	}
}

func TestRenderInvalidURL(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping network test in short mode")
	}

	r := New()

	_, err := r.Render("https://invalid.example.com/nonexistent.jpg", 20, 10)
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}
