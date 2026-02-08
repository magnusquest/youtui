# 🎵 YouTui - Implementation Plan

## 📝 Project Overview

**YouTui** is a terminal-based music player that streams YouTube videos and playlists with queue management and an integrated Pomodoro timer. Built as a single Go binary with a beautiful TUI using the Charm ecosystem.

---

## 🎯 Core Features

1. **YouTube Search** — Search videos via YouTube Data API v3, display results with thumbnails
2. **Queue Management** — Full CRUD (add, remove, reorder, clear) with keyboard shortcuts
3. **Playback Control** — Play/pause/skip/seek/volume via `mpv` + `yt-dlp`
4. **Pomodoro Mode** — Preset profiles (Classic, Deep Work, Sprint) + custom profiles
   - Music plays during work sessions
   - Music pauses during breaks
   - Audio notification on state transitions

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    YouTui (Single Go Binary)            │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌─────────────────┐  ┌─────────────────────────────┐  │
│  │  Bubbletea TUI  │  │  Pomodoro Timer             │  │
│  │  - Search view  │  │  - Work/Break states        │  │
│  │  - Queue view   │  │  - Profile management       │  │
│  │  - Now Playing  │  │  - Playback pause/resume    │  │
│  │  - Pomodoro     │  │                             │  │
│  └────────┬────────┘  └──────────────┬──────────────┘  │
│           │                          │                  │
│           ▼                          ▼                  │
│  ┌──────────────────────────────────────────────────┐  │
│  │              Playback Controller                  │  │
│  │  - mpv IPC control (play/pause/seek/volume)      │  │
│  │  - Queue state management                        │  │
│  │  - Track progress tracking                       │  │
│  └────────────────────────┬─────────────────────────┘  │
│                           │                             │
│  ┌────────────────────────┴─────────────────────────┐  │
│  │              YouTube Service                      │  │
│  │  - Google YouTube Data API v3 (search/metadata)  │  │
│  │  - yt-dlp (stream URL extraction only)           │  │
│  │  - Mosaic thumbnail rendering                    │  │
│  └──────────────────────────────────────────────────┘  │
│                                                         │
├─────────────────────────────────────────────────────────┤
│  External: mpv (playback) + yt-dlp (streams) + API v3  │
└─────────────────────────────────────────────────────────┘
```

---

## 📦 Project Structure (SOLID-based)

```
youtui/
├── cmd/
│   └── youtui/
│       └── main.go              # Entry point, wires dependencies
│
├── internal/
│   ├── app/
│   │   └── app.go               # Application orchestrator
│   │
│   ├── domain/                  # Core business entities (no dependencies)
│   │   ├── track.go             # Track entity
│   │   ├── queue.go             # Queue entity
│   │   ├── pomodoro.go          # Pomodoro session/profile entities
│   │   └── playback.go          # Playback state entity
│   │
│   ├── service/                 # Business logic (depends on domain + ports)
│   │   ├── player/
│   │   │   └── player.go        # Playback orchestration
│   │   ├── queue/
│   │   │   └── queue.go         # Queue CRUD operations
│   │   ├── pomodoro/
│   │   │   └── pomodoro.go      # Timer logic, state transitions
│   │   └── search/
│   │       └── search.go        # Search orchestration
│   │
│   ├── port/                    # Interfaces (Dependency Inversion)
│   │   ├── youtube.go           # YouTubeClient interface
│   │   ├── player.go            # AudioPlayer interface
│   │   ├── storage.go           # Storage interface (queue persistence)
│   │   └── thumbnail.go         # ThumbnailRenderer interface
│   │
│   ├── adapter/                 # External implementations (implements ports)
│   │   ├── youtube/
│   │   │   └── apiv3.go         # YouTube Data API v3 client
│   │   ├── mpv/
│   │   │   └── mpv.go           # mpv IPC controller
│   │   ├── ytdlp/
│   │   │   └── ytdlp.go         # yt-dlp stream URL extractor
│   │   ├── mosaic/
│   │   │   └── renderer.go      # Mosaic thumbnail renderer
│   │   └── storage/
│   │       └── json.go          # JSON file storage
│   │
│   └── ui/                      # Bubbletea TUI (presentation layer)
│       ├── app.go               # Root TUI model
│       ├── styles/
│       │   └── styles.go        # Lipgloss styles
│       ├── components/          # Reusable UI components
│       │   ├── searchbar.go
│       │   ├── tracklist.go
│       │   ├── nowplaying.go
│       │   ├── progress.go
│       │   ├── timer.go
│       │   └── help.go          # Help modal overlay
│       ├── views/               # Full-screen views
│       │   ├── search.go
│       │   ├── queue.go
│       │   ├── playing.go
│       │   └── pomodoro.go
│       └── keymap/
│           └── keymap.go        # Key bindings
│
├── config/
│   └── config.go                # Config loading (os.Getenv + YAML)
│
├── config.yaml                  # Default configuration
├── go.mod
├── go.sum
└── README.md
```

### SOLID Principles Applied

| Principle | Implementation |
|-----------|----------------|
| **S**ingle Responsibility | Each package has one reason to change (e.g., `adapter/mpv` only changes if mpv IPC changes) |
| **O**pen/Closed | New adapters (e.g., VLC player) can be added without modifying services |
| **L**iskov Substitution | Any `AudioPlayer` implementation can replace `mpv` adapter |
| **I**nterface Segregation | Small, focused interfaces in `port/` (not one giant interface) |
| **D**ependency Inversion | Services depend on `port/` interfaces, not concrete adapters |

---

## 🛠️ Technical Dependencies

### Go Dependencies

```go
// go.mod
module github.com/user/youtui

go 1.22

require (
    // TUI Framework (Charm ecosystem)
    github.com/charmbracelet/bubbletea v1.2.4
    github.com/charmbracelet/bubbles v0.20.0
    github.com/charmbracelet/lipgloss v1.0.0
    github.com/charmbracelet/x/exp/term v0.0.0-20240814160751-e2dc55e1de44

    // YouTube Data API v3
    google.golang.org/api v0.214.0

    // Configuration: stdlib os.Getenv (no external dependency needed)

    // Logging
    github.com/charmbracelet/log v0.4.0
)
```

### System Dependencies

| Dependency | Version | Purpose |
|------------|---------|---------|
| **Go** | 1.22+ | Build toolchain |
| **mpv** | 0.35+ | Audio playback with IPC socket |
| **yt-dlp** | 2024.01+ | YouTube stream URL extraction |

### External APIs

| API | Purpose | Quota |
|-----|---------|-------|
| **YouTube Data API v3** | Search, video metadata, thumbnails | 10,000 units/day |

---

## 🎨 User Interface Design

### TUI Layout

```
╔══════════════════════════════════════════════════════════════╗
║  🎵 YouTui                              🍅 25:00 [Work]      ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║  ┌────────────┐  Now Playing                                ║
║  │  ▓▓▓▓▓▓▓▓  │  Tame Impala - The Less I Know The Better  ║
║  │  ▓██████▓  │  Album: Currents (2015)                    ║
║  │  ▓██████▓  │                                             ║
║  │  ▓██████▓  │   advancement━━━━━━━━━━━━━━●──────  02:45 / 03:38  ║
║  │  ▓▓▓▓▓▓▓▓  │                                             ║
║  └────────────┘   advancement⏸ Pause  ⏭ Next  🔀 Shuffle  🔁 Repeat   ║
║   (Mosaic art)                                              ║
╠══════════════════════════════════════════════════════════════╣
║  Search: [_____________________] 🔍                          ║
╠══════════════════════════════════════════════════════════════╣
║  📋 Queue (3 tracks)                                         ║
║  ┌──────────────────────────────────────────────────────┐   ║
║  │ ▶ 1. Tame Impala - The Less I Know...        [03:38] │   ║
║  │   2. MGMT - Electric Feel                    [03:49] │   ║
║  │   3. Daft Punk - Get Lucky                   [06:09] │   ║
║  └──────────────────────────────────────────────────────┘   ║
╠══════════════════════════════════════════════════════════════╣
║  [1]Search [2]Queue [3]Playing [4]Pomodoro    [?]Help       ║
╚══════════════════════════════════════════════════════════════╝
```

### Key Bindings

#### Global
| Key | Action |
|-----|--------|
| `1` | Switch to Search view |
| `2` | Switch to Queue view |
| `3` | Switch to Now Playing view |
| `4` | Switch to Pomodoro view |
| `?` | Show help overlay |
| `q` / `Ctrl+C` | Quit (with confirmation) |

#### Playback
| Key | Action |
|-----|--------|
| `Space` | Play/Pause toggle |
| `n` | Next track |
| `p` | Previous track |
| `>` | Seek forward 10s |
| `<` | Seek backward 10s |
| `+` / `=` | Volume up |
| `-` | Volume down |
| `m` | Mute toggle |
| `s` | Shuffle toggle |
| `r` | Repeat mode (none → one → all) |

#### Search View
| Key | Action |
|-----|--------|
| `/` | Focus search input |
| `Enter` | Execute search / Play selected |
| `a` | Add selected to queue |
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `Esc` | Clear search / Unfocus |

#### Queue View
| Key | Action |
|-----|--------|
| `Enter` | Play selected track |
| `d` / `Delete` | Remove selected from queue |
| `c` | Clear entire queue (with confirmation) |
| `J` | Move selected track down |
| `K` | Move selected track up |
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |

#### Pomodoro View
| Key | Action |
|-----|--------|
| `Enter` | Start/Stop timer |
| `r` | Reset current session |
| `]` | Next profile |
| `[` | Previous profile |
| `S` | Skip to next phase (work→break→work) |

---

## ❓ Help Modal

A simple popup overlay triggered by `?` from any view. Displays all keyboard shortcuts in a scrollable list.

```
╔════════════════════════════════════════════╗
║  ❓ Keyboard Shortcuts           [Esc] ✕   ║
╠════════════════════════════════════════════╣
║                                            ║
║  GLOBAL                                    ║
║  1/2/3/4    Switch views                   ║
║  ?          Toggle this help               ║
║  q          Quit                           ║
║                                            ║
║  PLAYBACK                                  ║
║  Space      Play/Pause                     ║
║  n/p        Next/Previous track            ║
║  </> 		  Seek ±10s                       ║
║  +/-        Volume up/down                 ║
║  m          Mute                           ║
║  s          Shuffle                        ║
║  r          Repeat mode                    ║
║                                            ║
║  NAVIGATION                                ║
║  j/k        Move cursor down/up            ║
║  Enter      Select/Play                    ║
║  a          Add to queue                   ║
║  d          Remove from queue              ║
║  /          Focus search                   ║
║  Esc        Back/Unfocus                   ║
║                                            ║
║  POMODORO                                  ║
║  Enter      Start/Stop timer               ║
║  [/]        Switch profile                 ║
║  S          Skip phase                     ║
║                                            ║
╚════════════════════════════════════════════╝
```

**Behavior:**
- `?` toggles the modal on/off from any view
- `Esc` closes the modal
- Modal renders as overlay on top of current view
- Simple vertical scroll if content exceeds height

---

## 🍅 Pomodoro Profiles

### Default Profiles

| Profile | Work | Short Break | Long Break | Sessions before Long |
|---------|------|-------------|------------|----------------------|
| **Classic** | 25 min | 5 min | 15 min | 4 |
| **Deep Work** | 50 min | 10 min | 30 min | 2 |
| **Sprint** | 15 min | 3 min | 10 min | 4 |

### Behavior

- **Work Phase**: Music plays normally
- **Break Phase**: Music pauses automatically
- **Transition**: Audio notification (configurable sound)
- **Display**: Timer visible in header bar across all views

---

## 📝 Configuration

```yaml
# config.yaml

# YouTube Data API v3
youtube:
  api_key: "${YOUTUBE_API_KEY}"  # Set via environment variable
  max_results: 25
  region_code: "US"

# mpv player settings
player:
  socket_path: "/tmp/youtui-mpv.sock"
  default_volume: 80

# Pomodoro settings
pomodoro:
  default_profile: "classic"
  notification_sound: true
  profiles:
    classic:
      work_minutes: 25
      short_break_minutes: 5
      long_break_minutes: 15
      sessions_before_long_break: 4
    deep_work:
      work_minutes: 50
      short_break_minutes: 10
      long_break_minutes: 30
      sessions_before_long_break: 2
    sprint:
      work_minutes: 15
      short_break_minutes: 3
      long_break_minutes: 10
      sessions_before_long_break: 4

# Queue persistence
storage:
  queue_file: "~/.config/youtui/queue.json"

# UI settings
ui:
  show_thumbnails: true
  thumbnail_width: 12
  thumbnail_height: 6
```

---

## 🎯 Implementation Phases

### Phase 1: Project Foundation
- [ ] Initialize Go module with dependencies
- [ ] Set up project structure (folders, packages)
- [ ] Create domain entities (Track, Queue, Playback, Pomodoro)
- [ ] Define port interfaces
- [ ] Set up configuration loading (os.Getenv + YAML parsing)

### Phase 2: External Adapters
- [ ] Implement YouTube Data API v3 client (search, video details)
- [ ] Implement yt-dlp adapter (stream URL extraction)
- [ ] Implement mpv IPC adapter (playback control)
- [ ] Implement Mosaic thumbnail renderer
- [ ] Implement JSON file storage adapter

### Phase 3: Core Services
- [ ] Implement search service
- [ ] Implement queue service (CRUD operations)
- [ ] Implement player service (orchestrates mpv + queue)
- [ ] Implement pomodoro service (timer, state machine)

### Phase 4: TUI Components
- [ ] Create Lipgloss styles (theme)
- [ ] Build reusable components (searchbar, tracklist, progress, timer)
- [ ] Build Now Playing component with Mosaic album art

### Phase 5: TUI Views
- [ ] Implement Search view
- [ ] Implement Queue view
- [ ] Implement Now Playing view
- [ ] Implement Pomodoro view
- [ ] Implement Help overlay

### Phase 6: Integration & Polish
- [ ] Wire all components in main.go
- [ ] Add keyboard shortcut handling
- [ ] Add error handling and user feedback
- [ ] Test all features end-to-end
- [ ] Write README with installation instructions

---

## 🔑 YouTube API Setup

### Step 1: Create Google Cloud Project
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create new project: "YouTui"
3. Note the Project ID

### Step 2: Enable YouTube Data API v3
1. Navigate to "APIs & Services" → "Library"
2. Search for "YouTube Data API v3"
3. Click "Enable"

### Step 3: Create API Credentials
1. Go to "APIs & Services" → "Credentials"
2. Click "Create Credentials" → "API Key"
3. Copy the API key
4. (Recommended) Restrict key to YouTube Data API v3 only

### Step 4: Configure YouTui
```bash
export YOUTUBE_API_KEY="your-api-key-here"
```

---

## 📊 API Quota Management

**Daily Quota:** 10,000 units

| Operation | Cost | Notes |
|-----------|------|-------|
| Search | 100 units | ~100 searches/day |
| Video details | 1 unit | Cheap, use freely |

**Optimization:**
- Cache search results in memory during session
- Batch video detail requests when possible
- Show quota usage in Pomodoro view settings

---

## 🚀 Installation (Future)

```bash
# Install system dependencies
brew install mpv yt-dlp  # macOS
# or
sudo apt install mpv yt-dlp  # Debian/Ubuntu

# Install YouTui
go install github.com/user/youtui/cmd/youtui@latest

# Set API key
export YOUTUBE_API_KEY="your-key"

# Run
youtui
```

---

## ✅ Success Criteria

- [ ] Can search YouTube and see results with thumbnails
- [ ] Can add tracks to queue and manage queue (add/remove/reorder/clear)
- [ ] Can play/pause/skip/seek tracks
- [ ] Pomodoro timer works with preset profiles
- [ ] Music pauses on break, resumes on work
- [ ] All interactions via keyboard shortcuts
- [ ] Clean exit saves queue state
- [ ] Runs smoothly on macOS and Linux
