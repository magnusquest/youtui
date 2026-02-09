package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	// Primary colors
	colorPrimary   = lipgloss.Color("#FF6B6B")
	colorSecondary = lipgloss.Color("#4ECDC4")
	colorAccent    = lipgloss.Color("#FFE66D")
	
	// Text colors
	colorText      = lipgloss.Color("#FFFFFF")
	colorTextDim   = lipgloss.Color("#888888")
	colorTextBold  = lipgloss.Color("#FFEEEE")
	
	// UI colors
	colorBorder    = lipgloss.Color("#444444")
	colorSelection = lipgloss.Color("#6B4EFF")
	colorSuccess   = lipgloss.Color("#51CF66")
	colorWarning   = lipgloss.Color("#FFD93D")
	colorError     = lipgloss.Color("#FF6B6B")
)

// Base styles
var (
	// BaseStyle is the default style
	BaseStyle = lipgloss.NewStyle().
			Foreground(colorText)

	// TitleStyle for headers and titles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			MarginBottom(1)

	// SubtitleStyle for subtitles
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(colorSecondary)

	// BorderStyle for borders
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorBorder).
			Padding(1, 2)

	// FocusedBorderStyle for focused elements
	FocusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorPrimary).
				Padding(1, 2)

	// SelectedStyle for selected items
	SelectedStyle = lipgloss.NewStyle().
			Foreground(colorTextBold).
			Background(colorSelection).
			Bold(true)

	// DimStyle for less important text
	DimStyle = lipgloss.NewStyle().
			Foreground(colorTextDim)

	// BoldStyle for emphasis
	BoldStyle = lipgloss.NewStyle().
			Foreground(colorTextBold).
			Bold(true)

	// StatusBarStyle for status information
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(colorText).
			Background(lipgloss.Color("#1A1A1A")).
			Padding(0, 1)

	// HelpStyle for help text
	HelpStyle = lipgloss.NewStyle().
			Foreground(colorTextDim).
			Italic(true)

	// ErrorStyle for error messages
	ErrorStyle = lipgloss.NewStyle().
			Foreground(colorError).
			Bold(true)

	// SuccessStyle for success messages
	SuccessStyle = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Bold(true)

	// WarningStyle for warnings
	WarningStyle = lipgloss.NewStyle().
			Foreground(colorWarning).
			Bold(true)
)

// Component-specific styles

// Progress bar styles
var (
	ProgressBarFilledStyle = lipgloss.NewStyle().
				Foreground(colorPrimary)

	ProgressBarEmptyStyle = lipgloss.NewStyle().
				Foreground(colorBorder)

	ProgressBarTimeStyle = lipgloss.NewStyle().
				Foreground(colorTextDim)
)

// Timer styles
var (
	TimerStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Bold(true)

	TimerPhaseStyle = lipgloss.NewStyle().
			Foreground(colorSecondary)

	TimerRunningStyle = lipgloss.NewStyle().
				Foreground(colorSuccess).
				Bold(true)

	TimerPausedStyle = lipgloss.NewStyle().
				Foreground(colorTextDim)
)

// List/Tracklist styles
var (
	ListItemStyle = lipgloss.NewStyle().
			Padding(0, 2)

	ListSelectedItemStyle = lipgloss.NewStyle().
				Foreground(colorTextBold).
				Background(colorSelection).
				Bold(true).
				Padding(0, 2)

	ListCursorStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)
)

// Input styles
var (
	InputStyle = lipgloss.NewStyle().
			Foreground(colorText).
			Background(lipgloss.Color("#2A2A2A")).
			Padding(0, 1)

	InputFocusedStyle = lipgloss.NewStyle().
				Foreground(colorTextBold).
				Background(lipgloss.Color("#3A3A3A")).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(colorPrimary).
				Padding(0, 1)

	InputPromptStyle = lipgloss.NewStyle().
				Foreground(colorSecondary).
				Bold(true)
)
