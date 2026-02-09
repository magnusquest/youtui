package components

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/mganuesquest/youtui/internal/domain"
)

// Timer displays a Pomodoro countdown in the header area.
type Timer struct {
	session *domain.PomodoroSession
	width   int
}

// NewTimer creates a new timer component.
func NewTimer() Timer {
	return Timer{
		session: nil,
		width:   40,
	}
}

// SetSession sets the pomodoro session to display.
func (t *Timer) SetSession(session *domain.PomodoroSession) {
	t.session = session
}

// SetWidth sets the width of the timer display.
func (t *Timer) SetWidth(width int) {
	t.width = width
}

// View renders the timer.
func (t Timer) View() string {
	if t.session == nil {
		return t.renderEmpty()
	}

	// Phase indicator
	var phaseStyle lipgloss.Style
	var phaseIcon string

	switch t.session.Phase {
	case domain.PhaseWork:
		phaseStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#51CF66")).
			Bold(true)
		phaseIcon = "🎯"
	case domain.PhaseShortBreak, domain.PhaseLongBreak:
		phaseStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4ECDC4")).
			Bold(true)
		phaseIcon = "☕"
	}

	phaseName := phaseStyle.Render(phaseIcon + " " + t.session.PhaseName())

	// Timer display
	var timeStyle lipgloss.Style
	if t.session.Running {
		timeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFE66D")).
			Bold(true).
			Blink(true)
	} else {
		timeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
	}

	timeStr := timeStyle.Render(formatPomodoroTime(t.session.TimeRemaining))

	// Profile name
	profileStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)
	
	profileName := profileStyle.Render(fmt.Sprintf("[%s]", t.session.Profile.Name))

	// Status indicator
	var statusStyle lipgloss.Style
	var statusIcon string
	if t.session.Running {
		statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#51CF66")).
			Bold(true)
		statusIcon = "▶"
	} else {
		statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))
		statusIcon = "⏸"
	}
	status := statusStyle.Render(statusIcon)

	// Session counter
	sessionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888"))
	
	sessionCount := sessionStyle.Render(fmt.Sprintf("Session %d", t.session.CompletedSessions+1))

	// Combine elements
	return fmt.Sprintf("%s  %s %s  %s  %s", 
		phaseName, 
		status, 
		timeStr, 
		profileName,
		sessionCount,
	)
}

func (t Timer) renderEmpty() string {
	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true)
	
	return emptyStyle.Render("⏱  Pomodoro timer not started")
}

func formatPomodoroTime(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
