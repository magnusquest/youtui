package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ProgressBar shows track position with seek capability.
type ProgressBar struct {
	position time.Duration
	duration time.Duration
	width    int
}

// NewProgressBar creates a new progress bar component.
func NewProgressBar() ProgressBar {
	return ProgressBar{
		position: 0,
		duration: 0,
		width:    40,
	}
}

// SetPosition sets the current playback position.
func (p *ProgressBar) SetPosition(position time.Duration) {
	p.position = position
	if p.position < 0 {
		p.position = 0
	}
	if p.duration > 0 && p.position > p.duration {
		p.position = p.duration
	}
}

// SetDuration sets the total track duration.
func (p *ProgressBar) SetDuration(duration time.Duration) {
	p.duration = duration
}

// SetWidth sets the width of the progress bar.
func (p *ProgressBar) SetWidth(width int) {
	if width > 20 {
		p.width = width - 20 // Account for time labels
	}
}

// Position returns the current position.
func (p ProgressBar) Position() time.Duration {
	return p.position
}

// Duration returns the total duration.
func (p ProgressBar) Duration() time.Duration {
	return p.duration
}

// Percentage returns the completion percentage (0-100).
func (p ProgressBar) Percentage() float64 {
	if p.duration == 0 {
		return 0
	}
	return float64(p.position) / float64(p.duration) * 100
}

// View renders the progress bar.
func (p ProgressBar) View() string {
	if p.duration == 0 {
		return p.renderEmpty()
	}

	// Format time strings
	currentTime := formatDuration(p.position)
	totalTime := formatDuration(p.duration)

	timeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))
	
	timeLabel := timeStyle.Render(fmt.Sprintf("%s / %s", currentTime, totalTime))

	// Calculate filled portion
	barWidth := p.width
	if barWidth < 10 {
		barWidth = 10
	}

	percentage := p.Percentage()
	filled := int(float64(barWidth) * percentage / 100)
	if filled > barWidth {
		filled = barWidth
	}

	// Build progress bar
	filledStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B"))
	
	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444"))

	filledBar := filledStyle.Render(strings.Repeat("━", filled))
	emptyBar := emptyStyle.Render(strings.Repeat("━", barWidth-filled))

	bar := filledBar + emptyBar

	// Show percentage
	percentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFE66D")).
		Bold(true)
	
	percentLabel := percentStyle.Render(fmt.Sprintf(" %3.0f%%", percentage))

	return bar + percentLabel + "  " + timeLabel
}

func (p ProgressBar) renderEmpty() string {
	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444"))
	
	bar := emptyStyle.Render(strings.Repeat("━", p.width))
	
	timeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))
	
	timeLabel := timeStyle.Render("--:-- / --:--")

	return bar + "    " + timeLabel
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
