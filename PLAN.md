# YouTui

A terminal-based music player that streams YouTube videos and playlists with queue management and an integrated Pomodoro timer. Built as a single Go binary with a beautiful TUI using the Charm ecosystem.

## Features

- **YouTube Search** — Search videos via YouTube Data API v3, display results with thumbnails
- **Queue Management** — Full CRUD (add, remove, reorder, clear) with keyboard shortcuts
- **Playback Control** — Play/pause/skip/seek/volume via `mpv` + `yt-dlp`
- **Pomodoro Mode** — Preset profiles (Classic, Deep Work, Sprint) + custom profiles
  - Music plays during work sessions
  - Music pauses during breaks
  - Audio notification on state transitions

## Requirements

| Dependency | Version | Purpose |
|------------|---------|---------|
| **Go** | 1.22+ | Build toolchain |
| **mpv** | 0.35+ | Audio playback with IPC socket |
| **yt-dlp** | 2024.01+ | YouTube stream URL extraction |

## Installation

```bash
# Install system dependencies
brew install mpv yt-dlp  # macOS
# or
sudo apt install mpv yt-dlp  # Debian/Ubuntu

# Install YouTui
go install github.com/magnusquest/youtui/cmd/youtui@latest

# Set API key (see YouTube API Setup below)
export YOUTUBE_API_KEY="your-key"

# Run
youtui
```

## YouTube API Setup

### Step 1: Create Google Cloud Project
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create new project: "YouTui"

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

## Keyboard Shortcuts

### Global
| Key | Action |
|-----|--------|
| `1` | Switch to Search view |
| `2` | Switch to Queue view |
| `3` | Switch to Now Playing view |
| `4` | Switch to Pomodoro view |
| `?` | Show help overlay |
| `q` / `Ctrl+C` | Quit |

### Playback
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

### Search View
| Key | Action |
|-----|--------|
| `/` | Focus search input |
| `Enter` | Execute search / Play selected |
| `a` | Add selected to queue |
| `j` / `↓` | Move cursor down |
| `k` / `↑` | Move cursor up |
| `Esc` | Clear search / Unfocus |

### Queue View
| Key | Action |
|-----|--------|
| `Enter` | Play selected track |
| `d` / `Delete` | Remove selected from queue |
| `c` | Clear entire queue |
| `J` | Move selected track down |
| `K` | Move selected track up |

### Pomodoro View
| Key | Action |
|-----|--------|
| `Enter` | Start/Stop timer |
| `r` | Reset current session |
| `]` | Next profile |
| `[` | Previous profile |
| `S` | Skip to next phase |

## Pomodoro Profiles

| Profile | Work | Short Break | Long Break | Sessions before Long |
|---------|------|-------------|------------|----------------------|
| **Classic** | 25 min | 5 min | 15 min | 4 |
| **Deep Work** | 50 min | 10 min | 30 min | 2 |
| **Sprint** | 15 min | 3 min | 10 min | 4 |

## Configuration

YouTui uses a YAML config file at `~/.config/youtui/config.yaml`:

```yaml
youtube:
  api_key: "${YOUTUBE_API_KEY}"
  max_results: 25
  region_code: "US"

player:
  socket_path: "/tmp/youtui-mpv.sock"
  default_volume: 80

pomodoro:
  default_profile: "classic"
  notification_sound: true

storage:
  queue_file: "~/.config/youtui/queue.json"

ui:
  show_thumbnails: true
  thumbnail_width: 12
  thumbnail_height: 6
```

## Architecture

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

## License

MIT
