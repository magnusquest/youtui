package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mganuesquest/youtui/internal/domain"
)

// NowPlaying displays the currently playing track with mosaic album art.
type NowPlaying struct {
	track       *domain.Track
	thumbnail   string // Mosaic-rendered thumbnail art
	playback    *domain.Playback
	progressBar ProgressBar
	width       int
	height      int
}

// NewNowPlaying creates a new now playing component.
func NewNowPlaying() NowPlaying {
	return NowPlaying{
		track:       nil,
		thumbnail:   "",
		playback:    nil,
		progressBar: NewProgressBar(),
		width:       80,
		height:      24,
	}
}

// SetTrack sets the currently playing track.
func (n *NowPlaying) SetTrack(track *domain.Track) {
	n.track = track
	if track != nil {
		n.progressBar.SetDuration(track.Duration)
	}
}

// SetThumbnail sets the mosaic-rendered thumbnail art.
func (n *NowPlaying) SetThumbnail(thumbnail string) {
	n.thumbnail = thumbnail
}

// SetPlayback sets the current playback state.
func (n *NowPlaying) SetPlayback(playback *domain.Playback) {
	n.playback = playback
	if playback != nil {
		n.progressBar.SetPosition(playback.Position)
	}
}

// SetSize sets the dimensions of the component.
func (n *NowPlaying) SetSize(width, height int) {
	n.width = width
	n.height = height
	n.progressBar.SetWidth(width - 4)
}

// View renders the now playing display.
func (n NowPlaying) View() string {
	if n.track == nil {
		return n.renderEmpty()
	}

	var b strings.Builder

	// Title section
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Bold(true).
		Width(n.width).
		Align(lipgloss.Center)
	
	b.WriteString(titleStyle.Render("♪ NOW PLAYING ♪"))
	b.WriteString("\n\n")

	// If we have a thumbnail, display it
	if n.thumbnail != "" {
		// Center the thumbnail
		thumbnailStyle := lipgloss.NewStyle().
			Align(lipgloss.Center).
			Width(n.width)
		
		b.WriteString(thumbnailStyle.Render(n.thumbnail))
		b.WriteString("\n\n")
	}

	// Track info
	trackTitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFEEEE")).
		Bold(true).
		Width(n.width).
		Align(lipgloss.Center)
	
	b.WriteString(trackTitleStyle.Render(n.track.Title))
	b.WriteString("\n")

	artistStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4ECDC4")).
		Width(n.width).
		Align(lipgloss.Center)
	
	b.WriteString(artistStyle.Render(n.track.Artist))
	b.WriteString("\n")

	if n.track.Album != "" {
		albumStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true).
			Width(n.width).
			Align(lipgloss.Center)
		
		b.WriteString(albumStyle.Render("from " + n.track.Album))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Progress bar
	progressStyle := lipgloss.NewStyle().
		Width(n.width).
		Align(lipgloss.Center)
	
	b.WriteString(progressStyle.Render(n.progressBar.View()))
	b.WriteString("\n\n")

	// Playback controls hint
	if n.playback != nil {
		controlsHint := n.renderControlsHint()
		b.WriteString(controlsHint)
	}

	// Add border
	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF6B6B")).
		Padding(1, 2).
		Width(n.width)

	return borderStyle.Render(b.String())
}

func (n NowPlaying) renderEmpty() string {
	emptyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Italic(true).
		Width(n.width).
		Height(n.height).
		Align(lipgloss.Center, lipgloss.Center)
	
	return emptyStyle.Render("♪ No track playing ♪\n\nPress Enter on a track to start playback")
}

func (n NowPlaying) renderControlsHint() string {
	if n.playback == nil {
		return ""
	}

	var controls []string

	// Playback state
	var stateIcon string
	var stateText string
	var stateColor lipgloss.Color

	switch n.playback.State {
	case domain.StatePlaying:
		stateIcon = "▶"
		stateText = "Playing"
		stateColor = lipgloss.Color("#51CF66")
	case domain.StatePaused:
		stateIcon = "⏸"
		stateText = "Paused"
		stateColor = lipgloss.Color("#FFD93D")
	case domain.StateStopped:
		stateIcon = "⏹"
		stateText = "Stopped"
		stateColor = lipgloss.Color("#888888")
	}

	stateStyle := lipgloss.NewStyle().
		Foreground(stateColor).
		Bold(true)
	
	controls = append(controls, stateStyle.Render(stateIcon+" "+stateText))

	// Volume
	var volumeIcon string
	if n.playback.Muted {
		volumeIcon = "🔇"
	} else if n.playback.Volume > 66 {
		volumeIcon = "🔊"
	} else if n.playback.Volume > 33 {
		volumeIcon = "🔉"
	} else {
		volumeIcon = "🔈"
	}

	volumeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4ECDC4"))
	
	controls = append(controls, volumeStyle.Render(fmt.Sprintf("%s %d%%", volumeIcon, n.playback.Volume)))

	// Shuffle
	if n.playback.Shuffle {
		shuffleStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFE66D"))
		controls = append(controls, shuffleStyle.Render("🔀 Shuffle"))
	}

	// Repeat
	var repeatText string
	switch n.playback.Repeat {
	case domain.RepeatOne:
		repeatText = "🔂 Repeat One"
	case domain.RepeatAll:
		repeatText = "🔁 Repeat All"
	}
	if repeatText != "" {
		repeatStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFE66D"))
		controls = append(controls, repeatStyle.Render(repeatText))
	}

	// Join controls
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#444444")).
		Render("  │  ")

	controlsStr := strings.Join(controls, separator)

	// Center the controls
	controlsStyle := lipgloss.NewStyle().
		Width(n.width).
		Align(lipgloss.Center)
	
	return controlsStyle.Render(controlsStr)
}
