package ui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mganuesquest/youtui/internal/domain"
	"github.com/mganuesquest/youtui/internal/service"
	"github.com/mganuesquest/youtui/internal/ui/components"
)

// PomodoroViewModel manages the pomodoro view state and components.
type PomodoroViewModel struct {
	timer components.Timer
	
	pomodoroService *service.PomodoroService
	
	session  *domain.PomodoroSession
	profiles []domain.PomodoroProfile
	
	width  int
	height int
}

// NewPomodoroViewModel creates a new pomodoro view model.
func NewPomodoroViewModel(pomodoroService *service.PomodoroService) PomodoroViewModel {
	vm := PomodoroViewModel{
		timer:           components.NewTimer(),
		pomodoroService: pomodoroService,
		profiles:        domain.DefaultProfiles(),
	}
	
	// Subscribe to pomodoro events if service is available
	if pomodoroService != nil {
		pomodoroService.Subscribe(vm.onPomodoroEvent)
	}
	
	return vm
}

// Init initializes the pomodoro view.
func (p PomodoroViewModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the pomodoro view.
func (p PomodoroViewModel) Update(msg tea.Msg) (PomodoroViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Start/stop timer
			return p, p.toggleTimer()
			
		case "[":
			// Previous profile
			return p, p.previousProfile()
			
		case "]":
			// Next profile
			return p, p.nextProfile()
			
		case "S":
			// Skip phase
			return p, p.skipPhase()
			
		case "R":
			// Reset timer
			return p, p.resetTimer()
		}
	
	case pomodoroStatusMsg:
		p.updateFromPomodoroStatus(msg.status)
		return p, nil
	}

	return p, nil
}

// SetSize sets the view dimensions.
func (p *PomodoroViewModel) SetSize(width, height int) {
	p.width = width
	p.height = height
	p.timer.SetWidth(width)
}

// View renders the pomodoro view.
func (p PomodoroViewModel) View() string {
	var b strings.Builder

	// Title
	title := TitleStyle.Render("⏱  Pomodoro Timer")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Timer display (centered)
	timerStyle := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center)
	
	b.WriteString(timerStyle.Render(p.timer.View()))
	b.WriteString("\n\n")

	// Session info
	if p.session != nil {
		b.WriteString(p.renderSessionInfo())
		b.WriteString("\n\n")
	}

	// Profile selector
	b.WriteString(p.renderProfileSelector())
	b.WriteString("\n\n")

	// Pomodoro explanation
	b.WriteString(p.renderExplanation())
	b.WriteString("\n\n")

	// Controls hint
	controlsHint := DimStyle.Render("Enter: Start/Stop | [/]: Switch Profile | S: Skip Phase | R: Reset")
	b.WriteString(controlsHint)

	return b.String()
}

func (p PomodoroViewModel) renderSessionInfo() string {
	if p.session == nil {
		return ""
	}

	var b strings.Builder

	// Statistics
	statsStyle := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Foreground(colorTextDim)

	completedSessions := fmt.Sprintf("Completed sessions: %d", p.session.CompletedSessions)
	b.WriteString(statsStyle.Render(completedSessions))
	b.WriteString("\n")

	// Next phase preview
	var nextPhase string
	switch p.session.Phase {
	case domain.PhaseWork:
		if (p.session.CompletedSessions+1)%p.session.Profile.SessionsBeforeLongBreak == 0 {
			nextPhase = fmt.Sprintf("Next: Long Break (%s)", formatDurationShort(p.session.Profile.LongBreakDuration))
		} else {
			nextPhase = fmt.Sprintf("Next: Short Break (%s)", formatDurationShort(p.session.Profile.ShortBreakDuration))
		}
	case domain.PhaseShortBreak, domain.PhaseLongBreak:
		nextPhase = fmt.Sprintf("Next: Work Session (%s)", formatDurationShort(p.session.Profile.WorkDuration))
	}

	nextStyle := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center).
		Foreground(colorSecondary).
		Italic(true)
	
	b.WriteString(nextStyle.Render(nextPhase))

	return b.String()
}

func (p PomodoroViewModel) renderProfileSelector() string {
	var b strings.Builder

	profileStyle := lipgloss.NewStyle().
		Width(p.width).
		Align(lipgloss.Center)

	label := DimStyle.Render("Available Profiles:")
	b.WriteString(profileStyle.Render(label))
	b.WriteString("\n")

	// Show all profiles
	var profiles []string
	currentProfile := ""
	if p.session != nil {
		currentProfile = p.session.Profile.Name
	}

	for _, prof := range p.profiles {
		var profStr string
		if prof.Name == currentProfile {
			selectedStyle := lipgloss.NewStyle().
				Foreground(colorPrimary).
				Bold(true)
			profStr = selectedStyle.Render("▶ " + prof.Name)
		} else {
			profStr = "  " + prof.Name
		}
		
		details := DimStyle.Render(fmt.Sprintf(" (%s work / %s break)", 
			formatDurationShort(prof.WorkDuration),
			formatDurationShort(prof.ShortBreakDuration)))
		
		profiles = append(profiles, profStr+details)
	}

	profilesStr := strings.Join(profiles, "\n")
	b.WriteString(profileStyle.Render(profilesStr))

	return b.String()
}

func (p PomodoroViewModel) renderExplanation() string {
	explanationStyle := lipgloss.NewStyle().
		Width(p.width).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorBorder).
		Foreground(colorTextDim)

	explanation := `The Pomodoro Technique helps maintain focus and productivity:

• Work Phase: Focus on your task without distractions
• Break Phase: Rest and recharge
• After several sessions, take a longer break

Music will pause during breaks and resume during work phases.`

	return explanationStyle.Render(explanation)
}

// Message types
type pomodoroStatusMsg struct {
	status service.PomodoroStatus
}

func (p *PomodoroViewModel) onPomodoroEvent(status service.PomodoroStatus) {
	// This is called from the pomodoro service's goroutine
	// We'll send it as a message to the Bubbletea program
	// TODO: Wire this up when integrating with the main app
}

func (p *PomodoroViewModel) updateFromPomodoroStatus(status service.PomodoroStatus) {
	if status.Session != nil {
		p.session = status.Session
		p.timer.SetSession(status.Session)
	}
}

func (p PomodoroViewModel) toggleTimer() tea.Cmd {
	// TODO: Implement when pomodoro service is wired up
	return nil
}

func (p PomodoroViewModel) previousProfile() tea.Cmd {
	// TODO: Implement when pomodoro service is wired up
	return nil
}

func (p PomodoroViewModel) nextProfile() tea.Cmd {
	// TODO: Implement when pomodoro service is wired up
	return nil
}

func (p PomodoroViewModel) skipPhase() tea.Cmd {
	// TODO: Implement when pomodoro service is wired up
	return nil
}

func (p PomodoroViewModel) resetTimer() tea.Cmd {
	// TODO: Implement when pomodoro service is wired up
	return nil
}

func formatDurationShort(d time.Duration) string {
	minutes := int(d.Minutes())
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}
	hours := minutes / 60
	mins := minutes % 60
	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh%dm", hours, mins)
}
