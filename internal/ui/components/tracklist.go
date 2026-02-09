package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mganuesquest/youtui/internal/domain"
)

// Tracklist displays a scrollable list of tracks with selection.
type Tracklist struct {
	tracks   []domain.Track
	cursor   int
	offset   int
	height   int
	width    int
	focused  bool
	showThumbnail bool
	thumbnails map[string]string // track ID -> rendered thumbnail
}

// NewTracklist creates a new tracklist component.
func NewTracklist() Tracklist {
	return Tracklist{
		tracks:     []domain.Track{},
		cursor:     0,
		offset:     0,
		height:     20,
		width:      80,
		thumbnails: make(map[string]string),
	}
}

// SetTracks sets the list of tracks to display.
func (t *Tracklist) SetTracks(tracks []domain.Track) {
	t.tracks = tracks
	if t.cursor >= len(tracks) {
		t.cursor = len(tracks) - 1
	}
	if t.cursor < 0 {
		t.cursor = 0
	}
}

// SetThumbnail sets the rendered thumbnail for a track.
func (t *Tracklist) SetThumbnail(trackID, thumbnail string) {
	t.thumbnails[trackID] = thumbnail
}

// SetSize sets the dimensions of the tracklist.
func (t *Tracklist) SetSize(width, height int) {
	t.width = width
	t.height = height
}

// Focus sets the tracklist as focused.
func (t *Tracklist) Focus() {
	t.focused = true
}

// Blur removes focus from the tracklist.
func (t *Tracklist) Blur() {
	t.focused = false
}

// SelectedTrack returns the currently selected track, or nil if none.
func (t Tracklist) SelectedTrack() *domain.Track {
	if t.cursor >= 0 && t.cursor < len(t.tracks) {
		return &t.tracks[t.cursor]
	}
	return nil
}

// SelectedIndex returns the index of the selected track.
func (t Tracklist) SelectedIndex() int {
	return t.cursor
}

// Update handles messages for the tracklist.
func (t Tracklist) Update(msg tea.Msg) (Tracklist, tea.Cmd) {
	if !t.focused || len(t.tracks) == 0 {
		return t, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if t.cursor < len(t.tracks)-1 {
				t.cursor++
				// Scroll down if cursor goes below visible area
				if t.cursor >= t.offset+t.height {
					t.offset++
				}
			}
		case "k", "up":
			if t.cursor > 0 {
				t.cursor--
				// Scroll up if cursor goes above visible area
				if t.cursor < t.offset {
					t.offset--
				}
			}
		case "g":
			// Go to top
			t.cursor = 0
			t.offset = 0
		case "G":
			// Go to bottom
			t.cursor = len(t.tracks) - 1
			if t.cursor >= t.height {
				t.offset = t.cursor - t.height + 1
			}
		}
	}

	return t, nil
}

// View renders the tracklist.
func (t Tracklist) View() string {
	if len(t.tracks) == 0 {
		style := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true).
			Padding(2)
		return style.Render("No tracks to display")
	}

	var b strings.Builder

	// Calculate visible range
	start := t.offset
	end := t.offset + t.height
	if end > len(t.tracks) {
		end = len(t.tracks)
	}

	for i := start; i < end; i++ {
		track := t.tracks[i]
		isSelected := i == t.cursor

		var line string
		if isSelected && t.focused {
			cursor := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF6B6B")).
				Bold(true).
				Render("▶ ")

			title := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFEEEE")).
				Bold(true).
				Render(track.Title)

			artist := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#4ECDC4")).
				Render(track.Artist)

			duration := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				Render(track.FormatDuration())

			line = cursor + title + " - " + artist + "  " + duration
			
			line = lipgloss.NewStyle().
				Background(lipgloss.Color("#6B4EFF")).
				Width(t.width).
				Render(line)
		} else if isSelected {
			cursor := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				Render("▶ ")

			line = cursor + track.Title + " - " + track.Artist + "  " + track.FormatDuration()
		} else {
			line = "  " + track.Title + " - " + track.Artist + "  " + track.FormatDuration()
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	// Show scroll indicator
	if len(t.tracks) > t.height {
		scrollInfo := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Italic(true).
			Render(fmt.Sprintf("  [%d/%d]", t.cursor+1, len(t.tracks)))
		b.WriteString("\n")
		b.WriteString(scrollInfo)
	}

	return b.String()
}
