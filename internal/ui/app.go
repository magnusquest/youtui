package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mganuesquest/youtui/config"
)

// View represents the current active view
type View int

const (
	ViewSearch View = iota
	ViewQueue
	ViewPlaying
	ViewPomodoro
)

// Model is the root TUI model
type Model struct {
	config      *config.Config
	currentView View
	showHelp    bool
	width       int
	height      int
}

// New creates the root TUI model
func New(cfg *config.Config) Model {
	return Model{
		config:      cfg,
		currentView: ViewSearch,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			m.showHelp = !m.showHelp
		case "1":
			m.currentView = ViewSearch
		case "2":
			m.currentView = ViewQueue
		case "3":
			m.currentView = ViewPlaying
		case "4":
			m.currentView = ViewPomodoro
		case "esc":
			m.showHelp = false
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	if m.showHelp {
		return m.renderHelp()
	}

	switch m.currentView {
	case ViewSearch:
		return "Search View (TODO)\n\nPress ? for help, q to quit"
	case ViewQueue:
		return "Queue View (TODO)\n\nPress ? for help, q to quit"
	case ViewPlaying:
		return "Now Playing View (TODO)\n\nPress ? for help, q to quit"
	case ViewPomodoro:
		return "Pomodoro View (TODO)\n\nPress ? for help, q to quit"
	default:
		return "Unknown view"
	}
}

func (m Model) renderHelp() string {
	return `
Keyboard Shortcuts
==================

GLOBAL
  1/2/3/4    Switch views
  ?          Toggle this help
  q          Quit

PLAYBACK
  Space      Play/Pause
  n/p        Next/Previous track
  </>        Seek +/-10s
  +/-        Volume up/down
  m          Mute
  s          Shuffle
  r          Repeat mode

NAVIGATION
  j/k        Move cursor down/up
  Enter      Select/Play
  a          Add to queue
  d          Remove from queue
  /          Focus search
  Esc        Back/Unfocus

POMODORO
  Enter      Start/Stop timer
  [/]        Switch profile
  S          Skip phase

Press Esc or ? to close
`
}
