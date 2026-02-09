package components

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Searchbar is a text input component for search queries.
type Searchbar struct {
	input   textinput.Model
	focused bool
	width   int
}

// NewSearchbar creates a new searchbar component.
func NewSearchbar() Searchbar {
	ti := textinput.New()
	ti.Placeholder = "Search for music..."
	ti.CharLimit = 100
	ti.Width = 50

	return Searchbar{
		input:   ti,
		focused: false,
	}
}

// Focus sets the searchbar as focused and ready for input.
func (s *Searchbar) Focus() tea.Cmd {
	s.focused = true
	return s.input.Focus()
}

// Blur removes focus from the searchbar.
func (s *Searchbar) Blur() {
	s.focused = false
	s.input.Blur()
}

// SetWidth sets the width of the searchbar.
func (s *Searchbar) SetWidth(width int) {
	s.width = width
	if width > 10 {
		s.input.Width = width - 10 // Account for prompt and padding
	}
}

// Value returns the current search query.
func (s Searchbar) Value() string {
	return s.input.Value()
}

// SetValue sets the search query.
func (s *Searchbar) SetValue(value string) {
	s.input.SetValue(value)
}

// Reset clears the search query.
func (s *Searchbar) Reset() {
	s.input.SetValue("")
}

// Update handles messages for the searchbar.
func (s Searchbar) Update(msg tea.Msg) (Searchbar, tea.Cmd) {
	var cmd tea.Cmd
	s.input, cmd = s.input.Update(msg)
	return s, cmd
}

// View renders the searchbar.
func (s Searchbar) View() string {
	var style lipgloss.Style
	
	if s.focused {
		style = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6B6B")).
			Padding(0, 1)
	} else {
		style = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#444444")).
			Padding(0, 1)
	}

	prompt := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4ECDC4")).
		Bold(true).
		Render("🔍 ")

	return style.Render(prompt + s.input.View())
}
