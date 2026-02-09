package ui

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mganuesquest/youtui/internal/domain"
	"github.com/mganuesquest/youtui/internal/service"
	"github.com/mganuesquest/youtui/internal/ui/components"
)

// SearchViewModel manages the search view state and components.
type SearchViewModel struct {
	searchbar components.Searchbar
	tracklist components.Tracklist
	
	searchService *service.SearchService
	
	searching bool
	error     string
	results   []service.SearchResult
	
	width  int
	height int
}

// NewSearchViewModel creates a new search view model.
func NewSearchViewModel(searchService *service.SearchService) SearchViewModel {
	return SearchViewModel{
		searchbar:     components.NewSearchbar(),
		tracklist:     components.NewTracklist(),
		searchService: searchService,
		searching:     false,
		results:       []service.SearchResult{},
	}
}

// Init initializes the search view.
func (s SearchViewModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the search view.
func (s SearchViewModel) Update(msg tea.Msg) (SearchViewModel, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "/":
			// Focus search
			cmd = s.searchbar.Focus()
			s.tracklist.Blur()
			return s, cmd
			
		case "esc":
			// Unfocus search
			s.searchbar.Blur()
			s.tracklist.Focus()
			return s, nil
			
		case "enter":
			if s.searchbar.Value() != "" {
				// Perform search
				s.searching = true
				s.error = ""
				return s, s.performSearch()
			} else if s.tracklist.SelectedTrack() != nil {
				// Play selected track
				return s, s.playTrack()
			}
			
		case "a":
			// Add to queue
			if s.tracklist.SelectedTrack() != nil {
				return s, s.addToQueue()
			}
		}
	
	case searchResultMsg:
		s.searching = false
		if msg.err != nil {
			s.error = msg.err.Error()
		} else {
			s.results = msg.results
			tracks := make([]domain.Track, len(msg.results))
			for i, r := range msg.results {
				tracks[i] = r.Track
				if r.Thumbnail != "" {
					s.tracklist.SetThumbnail(r.Track.ID, r.Thumbnail)
				}
			}
			s.tracklist.SetTracks(tracks)
			s.tracklist.Focus()
		}
		return s, nil
	}

	// Update components
	s.searchbar, cmd = s.searchbar.Update(msg)
	cmds = append(cmds, cmd)
	
	s.tracklist, cmd = s.tracklist.Update(msg)
	cmds = append(cmds, cmd)

	return s, tea.Batch(cmds...)
}

// SetSize sets the view dimensions.
func (s *SearchViewModel) SetSize(width, height int) {
	s.width = width
	s.height = height
	s.searchbar.SetWidth(width - 4)
	s.tracklist.SetSize(width-4, height-10)
}

// View renders the search view.
func (s SearchViewModel) View() string {
	var b strings.Builder

	// Title
	title := TitleStyle.Render("🔍 Search")
	b.WriteString(title)
	b.WriteString("\n")

	// Searchbar
	b.WriteString(s.searchbar.View())
	b.WriteString("\n\n")

	// Error message
	if s.error != "" {
		errorMsg := ErrorStyle.Render("Error: " + s.error)
		b.WriteString(errorMsg)
		b.WriteString("\n\n")
	}

	// Searching indicator
	if s.searching {
		searchingMsg := lipgloss.NewStyle().
			Foreground(colorSecondary).
			Italic(true).
			Render("Searching...")
		b.WriteString(searchingMsg)
		b.WriteString("\n\n")
	}

	// Tracklist
	if len(s.results) > 0 {
		resultCount := DimStyle.Render(fmt.Sprintf("Found %d results:", len(s.results)))
		b.WriteString(resultCount)
		b.WriteString("\n")
		b.WriteString(s.tracklist.View())
	} else if !s.searching && s.searchbar.Value() == "" {
		hint := HelpStyle.Render("Press / to search for music")
		b.WriteString(hint)
	}

	return b.String()
}

// Message types
type searchResultMsg struct {
	results []service.SearchResult
	err     error
}

func (s SearchViewModel) performSearch() tea.Cmd {
	query := s.searchbar.Value()
	
	return func() tea.Msg {
		ctx := context.Background()
		results, err := s.searchService.Search(ctx, query, 20)
		return searchResultMsg{results: results, err: err}
	}
}

func (s SearchViewModel) playTrack() tea.Cmd {
	// TODO: Implement when player service is wired up
	return nil
}

func (s SearchViewModel) addToQueue() tea.Cmd {
	// TODO: Implement when queue service is wired up
	return nil
}
