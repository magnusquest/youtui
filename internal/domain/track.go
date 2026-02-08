package domain

import "time"

// Track represents a YouTube video/audio track
type Track struct {
	ID           string        `json:"id" yaml:"id"`
	Title        string        `json:"title" yaml:"title"`
	Artist       string        `json:"artist" yaml:"artist"`
	Album        string        `json:"album,omitempty" yaml:"album,omitempty"`
	Duration     time.Duration `json:"duration" yaml:"duration"`
	ThumbnailURL string        `json:"thumbnail_url" yaml:"thumbnail_url"`
	StreamURL    string        `json:"-" yaml:"-"` // Not persisted, fetched on demand
}

// FormatDuration returns duration as MM:SS
func (t Track) FormatDuration() string {
	minutes := int(t.Duration.Minutes())
	seconds := int(t.Duration.Seconds()) % 60
	return formatTime(minutes, seconds)
}

func formatTime(m, s int) string {
	return pad(m) + ":" + pad(s)
}

func pad(n int) string {
	if n < 10 {
		return "0" + itoa(n)
	}
	return itoa(n)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	digits := ""
	for n > 0 {
		digits = string(rune('0'+n%10)) + digits
		n /= 10
	}
	return digits
}
