package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mganuesquest/youtui/internal/domain"
	"github.com/mganuesquest/youtui/internal/service"
	"github.com/mganuesquest/youtui/internal/ui/components"
)

// QueueViewModel manages the queue view state and components.
type QueueViewModel struct {
	tracklist components.Tracklist
	
	queueService *service.QueueService
	
	queue  *domain.Queue
	error  string
	
	width  int
	height int
}

// NewQueueViewModel creates a new queue view model.
func NewQueueViewModel(queueService *service.QueueService) QueueViewModel {
	vm := QueueViewModel{
		tracklist:    components.NewTracklist(),
		queueService: queueService,
	}
	
	vm.tracklist.Focus()
	vm.loadQueue()
	
	return vm
}

// Init initializes the queue view.
func (q QueueViewModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the queue view.
func (q QueueViewModel) Update(msg tea.Msg) (QueueViewModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Play selected track
			if q.tracklist.SelectedTrack() != nil {
				return q, q.playTrack()
			}
			
		case "d":
			// Remove from queue
			if q.tracklist.SelectedTrack() != nil {
				return q, q.removeFromQueue()
			}
			
		case "s":
			// Toggle shuffle
			return q, q.toggleShuffle()
			
		case "r":
			// Cycle repeat mode
			return q, q.cycleRepeat()
		}
	
	case queueUpdatedMsg:
		q.queue = msg.queue
		if q.queue != nil {
			q.tracklist.SetTracks(q.queue.Tracks)
		}
		return q, nil
	}

	// Update tracklist
	q.tracklist, cmd = q.tracklist.Update(msg)

	return q, cmd
}

// SetSize sets the view dimensions.
func (q *QueueViewModel) SetSize(width, height int) {
	q.width = width
	q.height = height
	q.tracklist.SetSize(width-4, height-8)
}

// View renders the queue view.
func (q QueueViewModel) View() string {
	var b strings.Builder

	// Title
	title := TitleStyle.Render("📑 Queue")
	b.WriteString(title)
	b.WriteString("\n")

	// Queue info
	if q.queue != nil {
		info := DimStyle.Render(fmt.Sprintf("%d tracks in queue", len(q.queue.Tracks)))
		b.WriteString(info)
		b.WriteString("\n\n")
		
		// Current track indicator
		if q.queue.Current >= 0 && q.queue.Current < len(q.queue.Tracks) {
			currentStyle := lipgloss.NewStyle().
				Foreground(colorSuccess).
				Bold(true)
			currentInfo := currentStyle.Render(fmt.Sprintf("▶ Currently playing: #%d", q.queue.Current+1))
			b.WriteString(currentInfo)
			b.WriteString("\n\n")
		}
	}

	// Error message
	if q.error != "" {
		errorMsg := ErrorStyle.Render("Error: " + q.error)
		b.WriteString(errorMsg)
		b.WriteString("\n\n")
	}

	// Tracklist
	if q.queue != nil && len(q.queue.Tracks) > 0 {
		b.WriteString(q.tracklist.View())
	} else {
		hint := HelpStyle.Render("Queue is empty\nPress 'a' on search results to add tracks")
		b.WriteString(hint)
	}

	b.WriteString("\n\n")

	// Controls hint
	controlsHint := DimStyle.Render("Enter: Play | d: Remove | s: Shuffle | r: Repeat")
	b.WriteString(controlsHint)

	return b.String()
}

// Message types
type queueUpdatedMsg struct {
	queue *domain.Queue
}

func (q *QueueViewModel) loadQueue() {
	// Load queue from service
	if q.queueService != nil {
		tracks := q.queueService.Tracks()
		currentIndex := q.queueService.CurrentIndex()
		
		// Reconstruct queue for display
		q.queue = &domain.Queue{
			Tracks:  tracks,
			Current: currentIndex,
		}
		
		if q.queue != nil {
			q.tracklist.SetTracks(q.queue.Tracks)
		}
	}
}

func (q QueueViewModel) playTrack() tea.Cmd {
	// TODO: Implement when player service is wired up
	return nil
}

func (q QueueViewModel) removeFromQueue() tea.Cmd {
	// TODO: Implement when queue service methods are available
	return nil
}

func (q QueueViewModel) toggleShuffle() tea.Cmd {
	// TODO: Implement when player service is wired up
	return nil
}

func (q QueueViewModel) cycleRepeat() tea.Cmd {
	// TODO: Implement when player service is wired up
	return nil
}
