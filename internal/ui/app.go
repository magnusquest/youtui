package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mganuesquest/youtui/config"
	"github.com/mganuesquest/youtui/internal/service"
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
	
	// View models
	searchView      SearchViewModel
	queueView       QueueViewModel
	nowPlayingView  NowPlayingViewModel
	pomodoroView    PomodoroViewModel
}

// New creates the root TUI model
func New(cfg *config.Config) Model {
	// Services are nil for now - they'll be properly initialized in Phase 5
	return Model{
		config:         cfg,
		currentView:    ViewSearch,
		searchView:     NewSearchViewModel(nil),
		queueView:      NewQueueViewModel(nil),
		nowPlayingView: NewNowPlayingViewModel(nil),
		pomodoroView:   NewPomodoroViewModel(nil),
	}
}

// NewWithServices creates the root TUI model with service dependencies
func NewWithServices(cfg *config.Config, 
	searchService *service.SearchService,
	queueService *service.QueueService,
	playerService *service.PlayerService,
	pomodoroService *service.PomodoroService) Model {
	
	return Model{
		config:         cfg,
		currentView:    ViewSearch,
		searchView:     NewSearchViewModel(searchService),
		queueView:      NewQueueViewModel(queueService),
		nowPlayingView: NewNowPlayingViewModel(playerService),
		pomodoroView:   NewPomodoroViewModel(pomodoroService),
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle global keys first
		if m.showHelp {
			switch msg.String() {
			case "?", "esc":
				m.showHelp = false
			}
			return m, nil
		}
		
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			m.showHelp = true
			return m, nil
		case "1":
			m.currentView = ViewSearch
			return m, nil
		case "2":
			m.currentView = ViewQueue
			return m, nil
		case "3":
			m.currentView = ViewPlaying
			return m, nil
		case "4":
			m.currentView = ViewPomodoro
			return m, nil
		}
		
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Update all view sizes
		m.searchView.SetSize(msg.Width, msg.Height)
		m.queueView.SetSize(msg.Width, msg.Height)
		m.nowPlayingView.SetSize(msg.Width, msg.Height)
		m.pomodoroView.SetSize(msg.Width, msg.Height)
		
		return m, nil
	}
	
	// Route messages to the current view
	switch m.currentView {
	case ViewSearch:
		m.searchView, cmd = m.searchView.Update(msg)
		cmds = append(cmds, cmd)
	case ViewQueue:
		m.queueView, cmd = m.queueView.Update(msg)
		cmds = append(cmds, cmd)
	case ViewPlaying:
		m.nowPlayingView, cmd = m.nowPlayingView.Update(msg)
		cmds = append(cmds, cmd)
	case ViewPomodoro:
		m.pomodoroView, cmd = m.pomodoroView.Update(msg)
		cmds = append(cmds, cmd)
	}
	
	return m, tea.Batch(cmds...)
}

// View implements tea.Model
func (m Model) View() string {
	if m.showHelp {
		return m.renderHelp()
	}
	
	// Add view tabs/navigation bar
	viewTabs := m.renderViewTabs()
	
	// Get the current view content
	var viewContent string
	switch m.currentView {
	case ViewSearch:
		viewContent = m.searchView.View()
	case ViewQueue:
		viewContent = m.queueView.View()
	case ViewPlaying:
		viewContent = m.nowPlayingView.View()
	case ViewPomodoro:
		viewContent = m.pomodoroView.View()
	default:
		viewContent = "Unknown view"
	}
	
	return viewTabs + "\n" + viewContent
}

func (m Model) renderViewTabs() string {
	views := []struct {
		index View
		icon  string
		name  string
	}{
		{ViewSearch, "🔍", "Search"},
		{ViewQueue, "📑", "Queue"},
		{ViewPlaying, "♪", "Now Playing"},
		{ViewPomodoro, "⏱", "Pomodoro"},
	}
	
	var tabs []string
	for _, v := range views {
		var tab string
		if v.index == m.currentView {
			tab = SelectedStyle.Render(" " + v.icon + " " + v.name + " ")
		} else {
			tab = DimStyle.Render(" " + v.icon + " " + v.name + " ")
		}
		tabs = append(tabs, tab)
	}
	
	separator := DimStyle.Render(" │ ")
	tabBar := ""
	for i, tab := range tabs {
		if i > 0 {
			tabBar += separator
		}
		tabBar += tab
	}
	
	helpHint := DimStyle.Render("  [Press ? for help, q to quit]")
	
	return tabBar + helpHint
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
