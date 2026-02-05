# 🎵 YouTui - Implementation Plan

## 📝 Project Overview

**YouTui** is a self-hosted terminal music player that streams YouTube Music through your personal account, with dual interfaces: a beautiful TUI accessible via SSH and a mobile-friendly web interface. Designed to run on a Raspberry Pi on your home network.

**Inspiration:** Kreate Android music player + terminal.shop

---

## 🎯 Core Vision

- **Primary Access:** SSH-based TUI using Charm/Wish for command & control
- **Secondary Access:** Web GUI for mobile/tablet access
- **Deployment:** Self-hosted on Raspberry Pi (home network)
- **Architecture:** HTTP + WebSocket for real-time streaming and controls
- **Philosophy:** Beautiful, responsive interfaces (TUI & Web) with analogous functionality

---

## 🏗️ Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     Client Layer                             │
├──────────────────────────┬──────────────────────────────────┤
│   SSH Client (Wish TUI)  │    Web Browser (React/Svelte)    │
│   - Bubbletea Interface  │    - Mobile-first Design         │
│   - Real-time Updates    │    - WebSocket Audio Stream      │
│   - Config Management    │    - Touch-optimized Controls    │
└──────────────────────────┴──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    Go Backend (Raspberry Pi)                 │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐  │
│  │ Wish SSH     │  │ HTTP Server  │  │ WebSocket       │  │
│  │ Server       │  │ (Fiber/Echo) │  │ Handler         │  │
│  └──────────────┘  └──────────────┘  └─────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │         Core Music Engine (Go)                       │  │
│  │  - Playback Controller                               │  │
│  │  - Queue Management                                  │  │
│  │  - State Manager (playing, paused, etc.)             │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐  │
│  │ YouTube API  │  │ Redis Cache  │  │ Playlist        │  │
│  │ Client       │  │ Layer        │  │ Manager         │  │
│  └──────────────┘  └──────────────┘  └─────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                    External Services                         │
├─────────────────────────────────────────────────────────────┤
│  YouTube Data API v3 (Google) + Redis (Local Cache)         │
└─────────────────────────────────────────────────────────────┘
```

---

## 🛠️ Tech Stack

### Backend
- **Language:** Go 1.21+
- **SSH Server:** [Charm Wish](https://github.com/charmbracelet/wish) - SSH server framework
- **TUI Framework:** [Bubbletea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework (Model-View-Update pattern)
- **TUI Components:** [Bubbles](https://github.com/charmbracelet/bubbles) - UI components (list, spinner, progress)
- **TUI Styling:** [Lipgloss](https://github.com/charmbracelet/lipgloss) - Terminal layout and styling
- **Image Rendering:** [Mosaic](https://github.com/charmbracelet/x/tree/main/mosaic) - Terminal image rendering for album art
- **HTTP Framework:** [Fiber v3](https://github.com/gofiber/fiber) - Fast HTTP web framework
- **WebSocket:** Built-in Fiber WebSocket support
- **YouTube Integration:** [Google YouTube Data API v3](https://developers.google.com/youtube/v3) (Official API)
- **Caching:** [go-redis/v9](https://github.com/redis/go-redis) - Redis client for Go
- **Configuration:** [Viper](https://github.com/spf13/viper) - YAML config management
- **Google API Client:** [google-api-go-client](https://github.com/googleapis/google-api-go-client)

### Frontend (Web)
- **Framework:** SvelteKit 2 (compiler-based, reactive)
- **Styling:** TailwindCSS 3 (utility-first CSS)
- **UI Components:** DaisyUI (Tailwind component library)
- **WebSocket Client:** Native WebSocket API
- **Audio Player:** YouTube IFrame Player API (official, compliant)
- **Icons:** Lucide Svelte (consistent with TUI aesthetic)
- **Build:** Vite (fast, modern bundler)

### Infrastructure
- **Monorepo:** Nx Workspace
- **Containerization:** Docker + Docker Compose (optional)
- **Process Management:** systemd (for Pi deployment)
- **Reverse Proxy:** Caddy or Nginx (with auto HTTPS)

---

## 📦 Project Structure (Nx Monorepo)

```
youtui/
├── apps/
│   ├── backend/              # Go backend application
│   │   ├── cmd/
│   │   │   └── server/
│   │   │       └── main.go   # Entry point
│   │   ├── internal/
│   │   │   ├── api/          # HTTP API handlers
│   │   │   ├── ssh/          # Wish SSH server + TUI
│   │   │   ├── websocket/    # WebSocket handlers
│   │   │   ├── music/        # Music engine core
│   │   │   ├── youtube/      # YouTube Data API v3 integration
│   │   │   ├── cache/        # Redis caching layer
│   │   │   ├── playlist/     # Playlist management
│   │   │   └── config/       # Configuration
│   │   ├── pkg/              # Public packages
│   │   ├── Dockerfile
│   │   └── go.mod
│   │
│   └── web/                  # Svelte web application
│       ├── src/
│       │   ├── lib/
│       │   │   ├── components/   # Svelte components
│       │   │   ├── stores/       # State management
│       │   │   └── websocket/    # WS client
│       │   ├── routes/           # SvelteKit routes
│       │   └── app.html
│       ├── static/
│       ├── package.json
│       └── svelte.config.js
│
├── libs/
│   └── shared-types/         # Shared TypeScript types (API contracts)
│       └── src/
│           └── index.ts
│
├── tools/
│   ├── deploy/               # Deployment scripts for Raspberry Pi
│   │   ├── setup-pi.sh
│   │   └── systemd/
│   │       └── youtui.service
│   └── dev/
│       └── docker-compose.yml
│
├── docs/
│   ├── ARCHITECTURE.md
│   ├── API.md
│   └── DEPLOYMENT.md
│
├── nx.json
├── package.json
└── README.md
```

---

## 🎨 User Interface Design

### TUI (Bubbletea + Wish)

**Layout:**
```
╔══════════════════════════════════════════════════════════════╗
║  🎵 YouTui v1.0.0                    Connected: 2 clients   ║
╠══════════════════════════════════════════════════════════════╣
║                                                              ║
║  ┌────────────┐  Now Playing                                ║
║  │  ▓▓▓▓▓▓▓▓  │  Tame Impala - The Less I Know The Better  ║
║  │  ▓██████▓  │  Album: Currents (2015)                    ║
║  │  ▓██████▓  │                                             ║
║  │  ▓██████▓  │  ━━━━━━━━━━━━━━━━━━●──────  02:45 / 03:38  ║
║  │  ▓▓▓▓▓▓▓▓  │  [⏸️  Pause] [⏭️  Skip] [🔀 Shuffle] [🔁 Repeat] ║
║  └────────────┘                                             ║
║   Album Art                                                 ║
╠══════════════════════════════════════════════════════════════╣
║  Search: [_____________________] 🔍                          ║
╠══════════════════════════════════════════════════════════════╣
║  📋 Queue (3 tracks)                                         ║
║  ┌──────────────────────────────────────────────────────┐   ║
║  │ ▶ [▓] 1. Tame Impala - The Less I Know...   [03:38] │   ║
║  │   [▓] 2. MGMT - Electric Feel               [03:49] │   ║
║  │   [▓] 3. Daft Punk - Get Lucky               [06:09] │   ║
║  └──────────────────────────────────────────────────────┘   ║
║                                                              ║
║  💾 Playlists                                                ║
║  ┌──────────────────────────────────────────────────────┐   ║
║  │ • Chill Vibes (24 tracks)                            │   ║
║  │ • Workout Mix (18 tracks)                            │   ║
║  │ • Discover Weekly (50 tracks)                        │   ║
║  └──────────────────────────────────────────────────────┘   ║
║                                                              ║
╠══════════════════════════════════════════════════════════════╣
║  Tabs: Search | Queue | Playlists | Settings
╚══════════════════════════════════════════════════════════════╝
```

**Image Rendering:**
- Album art displayed using [Mosaic](https://github.com/charmbracelet/x/tree/main/mosaic) for terminal image rendering
- Thumbnails fetched from YouTube Data API v3 (`snippet.thumbnails`)
- Cached locally to minimize API calls and improve performance
- Adaptive rendering based on terminal capabilities (Sixel, Kitty, iTerm2, fallback to ASCII)
- Queue items show small thumbnail previews (▓ symbol represents tiny thumbnails)

**Key Bindings:**
- `j/k` or `↑/↓`: Navigate
- `Enter`: Play selected
- `Space`: Pause/Resume
- `n`: Next track
- `p`: Previous track
- `/`: Focus search within current view
- `Tab`: Switch panels (Search, Queue, Playlists, Settings)
- `Ctrl+C`: Quit (double tab)
### Web Interface (Svelte)

**Mobile-First Design:**
- Large album art display in now-playing view
- Bottom navigation bar (Search, Queue, Playlists, Settings)
- Now-playing bar on top (collapsed)
- Touch-optimized controls (large buttons)
- Pull-to-refresh for playlist updates
- Dark mode by default (matches TUI aesthetic)
- Thumbnail previews in queue and search results
- YouTube IFrame Player for audio playback (official API)

**Responsive Layout:**
- **Mobile (<768px):** Single column, bottom nav
- **Tablet (768-1024px):** Two column (queue + main)
- **Desktop (>1024px):** Three column (playlists, main, queue)

---

## 🔄 Real-time Communication

### WebSocket Protocol

**Client → Server Messages:**
```json
{
  "type": "play",
  "payload": { "trackId": "abc123" }
}

{
  "type": "pause"
}

{
  "type": "seek",
  "payload": { "position": 120 }
}

{
  "type": "queue_add",
  "payload": { "trackId": "def456" }
}
```

**Server → Client Messages:**
```json
{
  "type": "state_update",
  "payload": {
    "status": "playing",
    "currentTrack": {
      "id": "abc123",
      "title": "The Less I Know The Better",
      "artist": "Tame Impala",
      "duration": 218,
      "thumbnail": "https://..."
    },
    "position": 45,
    "volume": 80
  }
}

{
  "type": "queue_update",
  "payload": {
    "tracks": [...]
  }
}
```

---

## 💾 Caching Strategy

### Redis Cache Structure

**Track Metadata:**
```
Key: track:{video_id}
Value: JSON {
  "id": "abc123",
  "title": "Song Title",
  "artist": "Artist Name",
  "duration": 218,
  "thumbnail": "url",
  "audio_url": "cached_stream_url",
  "expires_at": "2024-02-05T12:00:00Z"
}
TTL: 24 hours
```

**Audio Files:**
```
Key: audio:{video_id}
Value: Binary audio data (MP3/Opus)
TTL: 7 days (for frequently played tracks)
Max Size: 10MB per track
```

**Search Results:**
```
Key: search:{query_hash}
Value: JSON array of track metadata
TTL: 1 hour
```

**Cache Eviction:**
- LRU policy for audio files
- Keep most recently/frequently played tracks
- Configurable max cache size (default: 5GB)

---

## 📡 API Endpoints

### REST API

**Core Endpoints:**
```
GET    /api/health               # Health check
GET    /api/status               # Current playback status
POST   /api/search               # Search YouTube Music
GET    /api/track/:id            # Get track details
POST   /api/play                 # Play track
POST   /api/pause                # Pause playback
POST   /api/resume               # Resume playback
POST   /api/next                 # Skip to next track
POST   /api/previous             # Previous track
POST   /api/seek                 # Seek to position

GET    /api/queue                # Get current queue
POST   /api/queue                # Add to queue
DELETE /api/queue/:index         # Remove from queue
PUT    /api/queue                # Reorder queue

GET    /api/playlists            # List playlists
POST   /api/playlists            # Create playlist
GET    /api/playlists/:id        # Get playlist
PUT    /api/playlists/:id        # Update playlist
DELETE /api/playlists/:id        # Delete playlist
POST   /api/playlists/:id/tracks # Add tracks to playlist

GET    /api/config               # Get server config
PUT    /api/config               # Update config (auth required)

GET    /stream/:track_id         # Audio stream endpoint
```

### WebSocket
```
WS /ws                           # WebSocket connection for real-time updates
```

---

## 🎯 MVP Development Phases

### Phase 1: Core Backend (2-3 weeks)
**Goal:** Basic music playback engine with YouTube Data API v3

**Tasks:**
- [ ] Set up Nx monorepo with Go backend
- [ ] Configure Google API key and YouTube Data API v3 client
- [ ] Implement YouTube API integration
  - [ ] Video search with proper filters (music category)
  - [ ] Video details retrieval
  - [ ] Stream URL extraction (within API quota limits)
  - [ ] Handle API rate limiting and quotas
- [ ] Create music playback controller (play, pause, skip)
- [ ] Build queue management system
- [ ] Set up Redis caching for track metadata and API responses
- [ ] Implement search functionality (with caching to minimize API calls)
- [ ] Create REST API endpoints (core playback + search)
- [ ] Add configuration system (Viper) with Google API key management

**Deliverable:** Working Go backend that can search and play music via YouTube Data API v3

---

### Phase 2: TUI Interface (2 weeks)
**Goal:** Beautiful terminal interface with Wish SSH server and image support

**Tasks:**
- [ ] Set up Wish SSH server with key-based auth
- [ ] Create Bubbletea TUI application
  - [ ] Now playing view with album art (using Mosaic)
  - [ ] Queue management with thumbnail previews
  - [ ] Search interface
  - [ ] Playlist view
- [ ] Implement image rendering
  - [ ] Integrate Mosaic for terminal images
  - [ ] Download and cache thumbnail images
  - [ ] Handle terminal capability detection
  - [ ] Fallback to ASCII art for unsupported terminals
- [ ] Implement key bindings
- [ ] Add real-time state updates via WebSocket
- [ ] Create configuration panel
- [ ] Test on local machine via SSH

**Deliverable:** Functional TUI accessible via `ssh user@raspberrypi -p 2222` with album art display

---

### Phase 3: Playlist Management (1 week)
**Goal:** Persistent playlist storage and management

**Tasks:**
- [ ] Design playlist data structure
- [ ] Implement playlist CRUD operations
- [ ] Add playlist storage (SQLite or JSON files)
- [ ] Create playlist UI in TUI
- [ ] Add "Add to playlist" functionality
- [ ] Implement playlist playback

**Deliverable:** Full playlist management in TUI and API

---

### Phase 4: Smart Caching (1 week)
**Goal:** Efficient audio caching to reduce bandwidth

**Tasks:**
- [ ] Implement audio file caching in Redis
- [ ] Add cache hit/miss tracking
- [ ] Create LRU eviction policy
- [ ] Add configurable cache size limits
- [ ] Implement pre-caching for next track in queue
- [ ] Add cache statistics endpoint
- [ ] Create cache management commands in TUI

**Deliverable:** Smart caching system that stores frequently played tracks

---

### Phase 5: WebSocket + Web Frontend (2-3 weeks)
**Goal:** Mobile-friendly web interface

**Tasks:**
- [ ] Implement WebSocket server in Go
- [ ] Create WebSocket protocol (message types)
- [ ] Set up Svelte/SvelteKit frontend in Nx
- [ ] Build web UI components
  - [ ] Now playing card
  - [ ] Playback controls
  - [ ] Queue view
  - [ ] Search interface
  - [ ] Playlist management
- [ ] Implement WebSocket client
- [ ] Add audio streaming via WebSocket or HTTP
- [ ] Style with TailwindCSS (dark theme)
- [ ] Test on mobile devices

**Deliverable:** Fully functional web interface with real-time sync

---

### Phase 6: Raspberry Pi Deployment (1 week)
**Goal:** Production deployment on Pi

**Tasks:**
- [ ] Create Docker images (optional)
- [ ] Write systemd service files
- [ ] Create deployment script for Pi
- [ ] Set up Caddy reverse proxy
- [ ] Configure Redis to start on boot
- [ ] Document Pi setup process
- [ ] Test on actual Raspberry Pi
- [ ] Optimize for Pi performance (arm64)

**Deliverable:** Running production system on Raspberry Pi accessible via home network

---

## 🚀 Post-MVP Features (Future Phases)

### Phase 7: Enhanced Features
- [ ] Lyrics display (via Genius API)
- [ ] Equalizer controls
- [ ] Crossfade between tracks
- [ ] Gapless playback
- [ ] Sleep timer
- [ ] Scrobbling to Last.fm
- [ ] Download tracks for offline playback
- [ ] Multi-room audio sync

### Phase 8: Anti-Ban Features
- [ ] Rotating proxy pool
- [ ] Request throttling
- [ ] User-agent rotation
- [ ] Cookie management
- [ ] Residential proxy support

### Phase 9: Social Features
- [ ] Shared queue (party mode)
- [ ] Collaborative playlists
- [ ] Listening history
- [ ] User accounts (multi-user support)
- [ ] "Now playing" notifications

---

## 🔒 Security Considerations

### SSH Access (Wish)
- Key-based authentication (no passwords)
- Rate limiting for SSH connections
- IP whitelisting (home network only)

### Web Interface
- HTTPS only (via Caddy)
- CORS configuration (restrict to home network)
- Optional basic auth for web interface
- CSRF protection

### API Security
- Rate limiting on endpoints
- Input validation and sanitization
- No public exposure (home network only)

---

## 📜 YouTube API Compliance & Fair Use

### Audio Streaming Implementation

**Important Note:** The YouTube Data API v3 provides metadata and search capabilities but doesn't directly expose audio stream URLs. Here are legitimate approaches:

#### Option 1: YouTube IFrame Player API (Recommended)
- Use the [IFrame Player API](https://developers.google.com/youtube/iframe_api_reference) for web client
- Player handles streaming within YouTube's terms
- Requires web browser or embedded player
- ✅ **Pros:** Fully compliant, official API, handles ads/monetization
- ❌ **Cons:** Requires browser, no TUI-only streaming

#### Option 2: Direct YouTube Embed Links
- Use `youtube.com/embed/{videoId}` for streaming
- Respects YouTube's terms and monetization
- Works with HTML5 audio/video players
- ✅ **Pros:** Simple, compliant
- ❌ **Cons:** Still needs web rendering

#### Option 3: YouTube Music API (If Available)
- Check if Google offers a YouTube Music API for developers
- May have different terms for music-focused apps
- Investigation needed

#### Recommended Architecture for Compliance:
```
TUI Client → Controls metadata and queue
Web Client → Handles actual audio playback via IFrame Player
Server → Manages state sync between clients
```

This approach:
- Uses Data API v3 for search and metadata (compliant)
- Uses IFrame Player for audio playback (compliant)
- TUI controls the experience, web client handles streaming
- Respects YouTube's monetization and terms

---

### API Quota Management
**Daily Quota:** 10,000 units per project (default)

**Cost per Operation:**
- Search query: 100 units
- Get video details: 1 unit
- Get video streams: Included in video details

**Optimization Strategies:**
1. **Aggressive Caching:**
   - Cache search results for 6 hours
   - Cache video metadata for 24 hours
   - Reduce redundant API calls by 90%

2. **Smart Search:**
   - Implement client-side filtering when possible
   - Use pagination efficiently
   - Deduplicate search queries

3. **Quota Monitoring:**
   - Track daily quota usage in Redis
   - Alert when approaching limits (80% threshold)
   - Implement graceful degradation (serve from cache only)

4. **Multi-User Quota Sharing:**
   - With 5 concurrent users, budget ~2,000 units per user per day
   - Approximately 20 searches per user per day (with caching)
   - Typical usage: 100-200 units/day with good caching

### Terms of Service Compliance
✅ **Allowed Use Cases:**
- Personal, non-commercial use
- Streaming content through official API
- Caching metadata (not video content itself)
- Building custom interfaces for YouTube content

❌ **Prohibited Actions:**
- Downloading and storing video files
- Circumventing ads or monetization
- Scraping without API
- Violating content creator rights

### Best Practices
1. Display YouTube attribution and links
2. Respect video availability and geographic restrictions
3. Honor age restrictions and content policies
4. Include links back to original YouTube videos
5. Do not remove or obscure YouTube branding from streams

---

## 📊 Performance Targets

### Raspberry Pi 4 (4GB RAM)
- **Startup Time:** < 5 seconds
- **Search Response:** < 500ms
- **Track Start Latency:** < 2 seconds (cache miss), < 500ms (cache hit)
- **Memory Usage:** < 500MB (with 5GB cache)
- **CPU Usage:** < 15% during playback
- **Concurrent Users:** 5-10 simultaneous streams

### Bandwidth
- **Average:** 128-320kbps per stream (user configurable)
- **Cache:** Reduce bandwidth by 70% after 1 week of usage

---

## 🧪 Testing Strategy

### Unit Tests
- Music engine logic
- Queue management
- Cache operations
- API handlers

### Integration Tests
- yt-dlp integration
- Redis caching
- WebSocket communication
- End-to-end playback flow

### Manual Testing
- TUI usability (keyboard navigation)
- Web UI responsiveness (mobile devices)
- Performance on Raspberry Pi
- Multi-client scenarios

---

## 📚 Dependencies & Installation

### Backend Dependencies
```go
// go.mod
module youtui

go 1.21

require (
    // TUI Framework
    github.com/charmbracelet/bubbletea v0.25.0
    github.com/charmbracelet/bubbles v0.18.0
    github.com/charmbracelet/lipgloss v0.9.1
    github.com/charmbracelet/wish v1.3.0
    github.com/charmbracelet/x/mosaic v0.0.0-20240125195854-abcdef123456  // Image rendering

    // HTTP & WebSocket
    github.com/gofiber/fiber/v3 v3.0.0-beta.3
    github.com/gofiber/contrib/websocket v1.3.2

    // YouTube API
    google.golang.org/api v0.156.0

    // Caching & Config
    github.com/redis/go-redis/v9 v9.4.0
    github.com/spf13/viper v1.18.2
)
```

**Note:** Mosaic version will be determined during implementation from the experimental repo.

### Frontend Dependencies
```json
{
  "name": "youtui-web",
  "type": "module",
  "dependencies": {
    "@sveltejs/kit": "^2.12.0",
    "svelte": "^5.0.0",
    "lucide-svelte": "^0.454.0"
  },
  "devDependencies": {
    "@sveltejs/adapter-static": "^3.0.5",
    "@sveltejs/vite-plugin-svelte": "^4.0.0",
    "tailwindcss": "^3.4.15",
    "daisyui": "^4.12.14",
    "vite": "^6.0.3"
  }
}
```

**Note:** No additional audio library needed - using YouTube IFrame Player API directly.

### System Requirements (Pi)
- Raspberry Pi 4 (4GB RAM recommended)
- Raspberry Pi OS (64-bit)
- Redis 7+
- Go 1.21+ (for building)
- **Google API Key** with YouTube Data API v3 enabled

---

## 🎓 Learning Resources

### Go + Bubbletea
- [Bubbletea Tutorial](https://github.com/charmbracelet/bubbletea/tree/master/tutorials)
- [Wish Examples](https://github.com/charmbracelet/wish/tree/main/examples)

### WebSockets in Go
- [Gorilla WebSocket Documentation](https://pkg.go.dev/github.com/gorilla/websocket)

### YouTube Data API v3
- [YouTube Data API v3 Overview](https://developers.google.com/youtube/v3/getting-started)
- [Go API Client Documentation](https://pkg.go.dev/google.golang.org/api/youtube/v3)
- [API Quota Management](https://developers.google.com/youtube/v3/determine_quota_cost)

### Svelte
- [Svelte Tutorial](https://svelte.dev/tutorial)
- [SvelteKit Documentation](https://kit.svelte.dev/docs)

---

## 📝 Configuration Example

```yaml
# config.yaml
server:
  ssh:
    host: "0.0.0.0"
    port: 2222
    key_path: "./ssh_host_key"

  http:
    host: "0.0.0.0"
    port: 8080

  websocket:
    path: "/ws"
    max_connections: 10

google:
  api_key: "${GOOGLE_API_KEY}"  # Set via environment variable
  youtube:
    max_results: 25              # Search results per query
    quota_daily_limit: 10000     # YouTube API daily quota
    region_code: "US"            # Region for content filtering
    safe_search: "moderate"      # none | moderate | strict

music:
  default_quality: "medium"  # low | medium | high
  audio_format: "opus"       # Preferred audio format
  buffer_size: 1024

  # API rate limiting to stay within quota
  rate_limit:
    requests_per_minute: 60
    burst: 10

cache:
  redis:
    host: "localhost"
    port: 6379
    db: 0

  # Cache API responses to minimize quota usage
  api_responses:
    enabled: true
    search_ttl_hours: 6        # Cache search results
    video_details_ttl_hours: 24
    stream_url_ttl_minutes: 30 # Stream URLs expire quickly

  metadata:
    ttl_hours: 24

playlists:
  storage_path: "./data/playlists"
  max_playlists: 100
  max_tracks_per_playlist: 1000
```

---

## 💻 Code Examples

### Bubbletea Model Structure

```go
package main

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type model struct {
    currentTrack  Track
    queue         []Track
    status        PlaybackStatus
    albumArtPath  string
    cursor        int
}

// Implement tea.Model interface
func (m model) Init() tea.Cmd {
    return tea.Batch(
        loadCurrentTrack(),
        fetchAlbumArt(),
        listenForWebSocketUpdates(),
    )
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        case "space":
            return m, togglePlayPause()
        case "n":
            return m, skipToNext()
        case "j", "down":
            if m.cursor < len(m.queue)-1 {
                m.cursor++
            }
        }

    case trackLoadedMsg:
        m.currentTrack = msg.track
        return m, fetchAlbumArt()

    case albumArtLoadedMsg:
        m.albumArtPath = msg.path
        return m, nil
    }

    return m, nil
}

func (m model) View() string {
    // Render UI with Lipgloss styling
    return renderNowPlaying(m) + renderQueue(m)
}
```

### Rendering Images with Mosaic

```go
package main

import (
    "image"
    _ "image/jpeg"
    _ "image/png"
    "os"

    "github.com/charmbracelet/x/mosaic"
)

func renderAlbumArt(imagePath string, width, height int) (string, error) {
    // Open the image file
    file, err := os.Open(imagePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    // Decode image
    img, _, err := image.Decode(file)
    if err != nil {
        return "", err
    }

    // Render to terminal using mosaic
    // Mosaic automatically detects terminal capabilities
    rendered := mosaic.Render(img, width, height)

    return rendered, nil
}

func downloadThumbnail(url string, cachePath string) error {
    // Download thumbnail from YouTube API response
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Save to cache
    out, err := os.Create(cachePath)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, resp.Body)
    return err
}
```

### YouTube Data API Integration

```go
package main

import (
    "context"
    "google.golang.org/api/option"
    "google.golang.org/api/youtube/v3"
)

func searchYouTube(ctx context.Context, query string, apiKey string) ([]*youtube.SearchResult, error) {
    service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
    if err != nil {
        return nil, err
    }

    call := service.Search.List([]string{"snippet"}).
        Q(query).
        Type("video").
        VideoCategoryId("10").  // Music category
        MaxResults(25).
        SafeSearch("moderate")

    response, err := call.Do()
    if err != nil {
        return nil, err
    }

    return response.Items, nil
}

func getVideoDetails(ctx context.Context, videoID string, apiKey string) (*youtube.Video, error) {
    service, err := youtube.NewService(ctx, option.WithAPIKey(apiKey))
    if err != nil {
        return nil, err
    }

    call := service.Videos.List([]string{"snippet", "contentDetails"}).
        Id(videoID)

    response, err := call.Do()
    if err != nil {
        return nil, err
    }

    if len(response.Items) == 0 {
        return nil, fmt.Errorf("video not found")
    }

    // Access thumbnail URLs
    video := response.Items[0]
    thumbnailURL := video.Snippet.Thumbnails.High.Url  // or .Medium, .Standard, .Maxres

    return video, nil
}
```

### Redis Caching

```go
package main

import (
    "context"
    "encoding/json"
    "time"

    "github.com/redis/go-redis/v9"
)

func setupRedis() *redis.Client {
    return redis.NewClient(&redis.Options{
        Addr:         "localhost:6379",
        Password:     "",
        DB:           0,
        DialTimeout:  10 * time.Second,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
        PoolSize:     10,
    })
}

func cacheTrackMetadata(ctx context.Context, rdb *redis.Client, track Track) error {
    data, err := json.Marshal(track)
    if err != nil {
        return err
    }

    key := fmt.Sprintf("track:%s", track.ID)
    return rdb.Set(ctx, key, data, 24*time.Hour).Err()
}

func getCachedTrack(ctx context.Context, rdb *redis.Client, trackID string) (*Track, error) {
    key := fmt.Sprintf("track:%s", trackID)

    data, err := rdb.Get(ctx, key).Result()
    if err == redis.Nil {
        return nil, nil  // Cache miss
    } else if err != nil {
        return nil, err
    }

    var track Track
    err = json.Unmarshal([]byte(data), &track)
    return &track, err
}
```

### Fiber HTTP Server with WebSocket

```go
package main

import (
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/contrib/websocket"
)

func setupServer() *fiber.App {
    app := fiber.New()

    // Serve static files for web UI
    app.Get("/*", static.New("./web/dist"))

    // API routes
    api := app.Group("/api")
    api.Get("/status", getStatus)
    api.Post("/search", searchMusic)
    api.Post("/play", playTrack)
    api.Get("/queue", getQueue)

    // WebSocket for real-time updates
    app.Use("/ws", func(c fiber.Ctx) error {
        if websocket.IsWebSocketUpgrade(c) {
            return c.Next()
        }
        return fiber.ErrUpgradeRequired
    })

    app.Get("/ws", websocket.New(func(c *websocket.Conn) {
        handleWebSocket(c)
    }))

    return app
}

func handleWebSocket(c *websocket.Conn) {
    defer c.Close()

    for {
        // Send state updates to client
        state := getCurrentState()
        if err := c.WriteJSON(state); err != nil {
            break
        }

        // Read commands from client
        var msg Command
        if err := c.ReadJSON(&msg); err != nil {
            break
        }

        handleCommand(msg)
    }
}
```

### Wish SSH Server Setup

```go
package main

import (
    "context"
    "net"
    "os"
    "os/signal"
    "syscall"
    "time"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/log"
    "github.com/charmbracelet/ssh"
    "github.com/charmbracelet/wish"
    "github.com/charmbracelet/wish/activeterm"
    "github.com/charmbracelet/wish/bubbletea"
    "github.com/charmbracelet/wish/logging"
)

func main() {
    srv, err := wish.NewServer(
        wish.WithAddress(net.JoinHostPort("0.0.0.0", "2222")),
        wish.WithHostKeyPath(".ssh/id_ed25519"),
        wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
            // Implement your auth logic
            return true  // Allow all for now
        }),
        wish.WithMiddleware(
            bubbletea.Middleware(teaHandler),
            activeterm.Middleware(),  // Ensure active PTY
            logging.Middleware(),
        ),
    )
    if err != nil {
        log.Error("Could not start server", "error", err)
        return
    }

    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
    log.Info("Starting SSH server", "port", "2222")

    go func() {
        if err = srv.ListenAndServe(); err != nil {
            log.Error("Could not start server", "error", err)
        }
    }()

    <-done
    log.Info("Stopping SSH server")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Error("Could not stop server", "error", err)
    }
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
    m := model{
        // Initialize your model
    }
    return m, []tea.ProgramOption{tea.WithAltScreen()}
}
```

### Svelte Component Example

```svelte
<!-- NowPlaying.svelte -->
<script>
  import { page } from '$app/state';

  let { currentTrack } = $props();
  let player;
  let isPlaying = $state(false);

  // YouTube IFrame Player API
  function loadYouTubePlayer() {
    player = new YT.Player('player', {
      videoId: currentTrack.videoId,
      events: {
        'onReady': onPlayerReady,
        'onStateChange': onPlayerStateChange
      }
    });
  }

  function togglePlay() {
    if (isPlaying) {
      player.pauseVideo();
    } else {
      player.playVideo();
    }
  }
</script>

<div class="now-playing">
  <img
    src={currentTrack.thumbnail}
    alt={currentTrack.title}
    class="album-art"
  />

  <div class="track-info">
    <h2>{currentTrack.title}</h2>
    <p>{currentTrack.artist}</p>
  </div>

  <div id="player"></div>

  <div class="controls">
    <button onclick={togglePlay}>
      {isPlaying ? '⏸️' : '▶️'}
    </button>
  </div>
</div>

<style>
  .album-art {
    width: 100%;
    max-width: 400px;
    border-radius: 8px;
  }
</style>
```

---

## 🚦 Success Metrics

### MVP Success Criteria
- ✅ Can search and play music from YouTube Music
- ✅ TUI accessible via SSH from any device on home network
- ✅ Web interface works on mobile (phone/tablet)
- ✅ Playlists can be created and managed
- ✅ Cache reduces bandwidth by 50%+ after 3 days of use
- ✅ Multiple clients can connect simultaneously
- ✅ Runs stably on Raspberry Pi for 7 days without restart

---

## 🔑 Google API Setup Guide

### Step 1: Create Google Cloud Project
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create new project: "YouTui Music Server"
3. Note the Project ID

### Step 2: Enable YouTube Data API v3
1. Navigate to "APIs & Services" → "Library"
2. Search for "YouTube Data API v3"
3. Click "Enable"

### Step 3: Create API Credentials
1. Go to "APIs & Services" → "Credentials"
2. Click "Create Credentials" → "API Key"
3. Copy the API key
4. (Optional) Restrict key:
   - Application restrictions: "IP addresses" (add your Pi's IP)
   - API restrictions: Only "YouTube Data API v3"

### Step 4: Configure YouTui
```bash
# Set environment variable
export GOOGLE_API_KEY="your-api-key-here"

# Or add to config.yaml
google:
  api_key: "your-api-key-here"
```

### Step 5: Monitor Quota
- Check usage: [Google Cloud Console → YouTube Data API v3 → Quotas](https://console.cloud.google.com/apis/api/youtube.googleapis.com/quotas)
- Default limit: 10,000 units/day
- Request quota increase if needed (for heavy usage)

---

## 🎉 Next Steps

1. **Review this plan** - Does it match your vision?
2. **Set up development environment** - Install Go, Node.js, Redis
3. **Create Nx monorepo** - Initialize project structure
4. **Start Phase 1** - Build core backend
5. **Iterate** - Build, test, refine

---

## 💭 Open Questions

1. **Audio format preference?** Opus (better quality/size) vs MP3 (wider compatibility)
2. **Authentication for web?** Open to home network or add password protection?
3. **Google API Key setup?** Will create new project in Google Cloud Console with YouTube Data API v3 enabled?
4. **API Quota monitoring?** Should the app show quota usage in settings/admin panel?
5. **Raspberry Pi model?** Confirming Pi 4 or considering Pi 5?
6. **Network exposure?** Home network only or considering ngrok/Cloudflare Tunnel for external access?
7. **Fallback behavior?** What should happen when API quota is exhausted? (cache-only mode, graceful error, etc.)

---

Let me know if you want to adjust anything in this plan! We can dive deeper into any specific phase or component.
