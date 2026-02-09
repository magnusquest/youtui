# UI Components

This directory contains the Bubbletea TUI components and views for YouTui.

## Architecture

The UI follows the Bubbletea pattern with:
- **Components**: Reusable, stateless UI elements
- **Views**: Stateful view models that compose components
- **Main App**: Root model that manages view switching

## Components (`internal/ui/components/`)

### Searchbar
Text input component for search queries with focus states.

**Features:**
- Focus/blur states with visual feedback
- Placeholder text
- Character limit (100 chars)
- Width adjustment

**Usage:**
```go
searchbar := components.NewSearchbar()
searchbar.Focus()
searchbar.SetWidth(50)
```

### Tracklist
Scrollable list component for displaying tracks with selection.

**Features:**
- Keyboard navigation (j/k, up/down)
- Visual selection indicator
- Scroll indicators
- Configurable dimensions
- Thumbnail support

**Usage:**
```go
tracklist := components.NewTracklist()
tracklist.SetTracks(tracks)
tracklist.Focus()
```

### ProgressBar
Shows track playback position with visual progress indication.

**Features:**
- Percentage display
- Time labels (current / total)
- Customizable width
- Smooth progress rendering

**Usage:**
```go
progress := components.NewProgressBar()
progress.SetDuration(track.Duration)
progress.SetPosition(currentPosition)
```

### Timer
Displays Pomodoro countdown timer with phase information.

**Features:**
- Phase indicators (Work/Break)
- Running/paused states
- Profile name display
- Session counter
- Visual blink for active timer

**Usage:**
```go
timer := components.NewTimer()
timer.SetSession(pomodoroSession)
```

### NowPlaying
Composite component showing currently playing track with album art.

**Features:**
- Mosaic album art rendering
- Track metadata display
- Progress bar integration
- Playback state indicators
- Volume/shuffle/repeat indicators

**Usage:**
```go
nowPlaying := components.NewNowPlaying()
nowPlaying.SetTrack(track)
nowPlaying.SetThumbnail(mosaicArt)
nowPlaying.SetPlayback(playbackState)
```

## Views (`internal/ui/`)

### SearchViewModel
Search interface with query input and results display.

**Key bindings:**
- `/`: Focus search bar
- `Enter`: Perform search or play selected track
- `a`: Add track to queue
- `j/k`: Navigate results

### QueueViewModel
Queue management interface.

**Key bindings:**
- `Enter`: Play selected track
- `d`: Remove from queue
- `s`: Toggle shuffle
- `r`: Cycle repeat mode

### NowPlayingViewModel
Currently playing track display with controls.

**Key bindings:**
- `Space`: Play/pause
- `n/p`: Next/previous track
- `</>`: Seek backward/forward
- `+/-`: Volume up/down
- `m`: Toggle mute
- `s`: Toggle shuffle
- `r`: Cycle repeat

### PomodoroViewModel
Pomodoro timer interface.

**Key bindings:**
- `Enter`: Start/stop timer
- `[/]`: Switch profile
- `S`: Skip phase
- `R`: Reset timer

## Styling (`internal/ui/styles.go`)

Consistent Lipgloss theme with:
- **Primary**: #FF6B6B (coral red)
- **Secondary**: #4ECDC4 (turquoise)
- **Accent**: #FFE66D (yellow)
- **Success**: #51CF66 (green)

All components use these predefined styles for consistency.

## Integration

Views are composed in the root `Model` and switched via number keys (1-4).

The main app handles:
- Global key bindings (help, quit, view switching)
- Window size management
- Message routing to active view

## Service Integration (Phase 5)

The view models accept service dependencies through their constructors:
- `SearchViewModel` → `SearchService`
- `QueueViewModel` → `QueueService`
- `NowPlayingViewModel` → `PlayerService`
- `PomodoroViewModel` → `PomodoroService`

Currently, views work with `nil` services for testing. Full integration will be completed in Phase 5.
