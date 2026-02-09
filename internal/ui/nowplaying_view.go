package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mganuesquest/youtui/internal/domain"
	"github.com/mganuesquest/youtui/internal/service"
	"github.com/mganuesquest/youtui/internal/ui/components"
)

// NowPlayingViewModel manages the now playing view state and components.
type NowPlayingViewModel struct {
	nowPlaying components.NowPlaying
	
	playerService *service.PlayerService
	
	currentTrack *domain.Track
	playback     *domain.Playback
	thumbnail    string
	
	width  int
	height int
}

// NewNowPlayingViewModel creates a new now playing view model.
func NewNowPlayingViewModel(playerService *service.PlayerService) NowPlayingViewModel {
	vm := NowPlayingViewModel{
		nowPlaying:    components.NewNowPlaying(),
		playerService: playerService,
	}
	
	// Subscribe to player events if service is available
	if playerService != nil {
		playerService.Subscribe(vm.onPlayerEvent)
	}
	
	return vm
}

// Init initializes the now playing view.
func (n NowPlayingViewModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the now playing view.
func (n NowPlayingViewModel) Update(msg tea.Msg) (NowPlayingViewModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			// Toggle play/pause
			return n, n.togglePlayPause()
			
		case "n":
			// Next track
			return n, n.nextTrack()
			
		case "p":
			// Previous track
			return n, n.previousTrack()
			
		case "<":
			// Seek backward
			return n, n.seekBackward()
			
		case ">":
			// Seek forward
			return n, n.seekForward()
			
		case "+":
			// Volume up
			return n, n.volumeUp()
			
		case "-":
			// Volume down
			return n, n.volumeDown()
			
		case "m":
			// Toggle mute
			return n, n.toggleMute()
			
		case "s":
			// Toggle shuffle
			return n, n.toggleShuffle()
			
		case "r":
			// Cycle repeat mode
			return n, n.cycleRepeat()
		}
	
	case playerStatusMsg:
		n.updateFromPlayerStatus(msg.status)
		return n, nil
	}

	return n, nil
}

// SetSize sets the view dimensions.
func (n *NowPlayingViewModel) SetSize(width, height int) {
	n.width = width
	n.height = height
	n.nowPlaying.SetSize(width, height)
}

// View renders the now playing view.
func (n NowPlayingViewModel) View() string {
	var b strings.Builder

	// Title
	title := TitleStyle.Render("♪ Now Playing")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Now playing component
	b.WriteString(n.nowPlaying.View())
	b.WriteString("\n\n")

	// Controls hint
	controlsHint := DimStyle.Render("Space: Play/Pause | n/p: Next/Prev | </>: Seek | +/-: Volume | m: Mute | s: Shuffle | r: Repeat")
	b.WriteString(controlsHint)

	return b.String()
}

// Message types
type playerStatusMsg struct {
	status service.PlayerStatus
}

func (n *NowPlayingViewModel) onPlayerEvent(status service.PlayerStatus) {
	// This is called from the player service's goroutine
	// We'll send it as a message to the Bubbletea program
	// TODO: Wire this up when integrating with the main app
}

func (n *NowPlayingViewModel) updateFromPlayerStatus(status service.PlayerStatus) {
	if status.Track != nil {
		n.currentTrack = status.Track
		n.nowPlaying.SetTrack(status.Track)
	}
	
	if status.Playback != nil {
		n.playback = status.Playback
		n.nowPlaying.SetPlayback(status.Playback)
	}
}

func (n NowPlayingViewModel) togglePlayPause() tea.Cmd {
	// TODO: Implement when player service is wired up
	return nil
}

func (n NowPlayingViewModel) nextTrack() tea.Cmd {
	// TODO: Implement when player service is wired up
	return nil
}

func (n NowPlayingViewModel) previousTrack() tea.Cmd {
	// TODO: Implement when player service is wired up
	return nil
}

func (n NowPlayingViewModel) seekBackward() tea.Cmd {
	// TODO: Implement when player service is wired up (seek -10s)
	return nil
}

func (n NowPlayingViewModel) seekForward() tea.Cmd {
	// TODO: Implement when player service is wired up (seek +10s)
	return nil
}

func (n NowPlayingViewModel) volumeUp() tea.Cmd {
	// TODO: Implement when player service is wired up
	return nil
}

func (n NowPlayingViewModel) volumeDown() tea.Cmd {
	// TODO: Implement when player service is wired up
	return nil
}

func (n NowPlayingViewModel) toggleMute() tea.Cmd {
	// TODO: Implement when player service is wired up
	return nil
}

func (n NowPlayingViewModel) toggleShuffle() tea.Cmd {
	// TODO: Implement when player service is wired up
	return nil
}

func (n NowPlayingViewModel) cycleRepeat() tea.Cmd {
	// TODO: Implement when player service is wired up
	return nil
}
