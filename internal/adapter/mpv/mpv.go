package mpv

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"
)

// Player implements port.AudioPlayer using mpv's JSON IPC protocol.
type Player struct {
	cmd        *exec.Cmd
	conn       net.Conn
	socketPath string

	mu       sync.Mutex
	reqID    atomic.Int64
	pending  map[int64]chan response
	pendingM sync.Mutex

	done chan struct{}
}

type response struct {
	Data  json.RawMessage `json:"data"`
	Error string          `json:"error"`
	ReqID int64           `json:"request_id"`
}

type ipcCommand struct {
	Command   []interface{} `json:"command"`
	RequestID int64         `json:"request_id"`
}

// New launches mpv in idle mode and connects to its IPC socket.
func New(socketPath string) (*Player, error) {
	// Clean up any stale socket file.
	os.Remove(socketPath)

	cmd := exec.Command("mpv",
		"--idle",
		"--no-video",
		"--no-terminal",
		"--input-ipc-server="+socketPath,
	)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("mpv: failed to start: %w", err)
	}

	// Wait for the socket to become available.
	var conn net.Conn
	var err error
	for i := 0; i < 50; i++ {
		conn, err = net.Dial("unix", socketPath)
		if err == nil {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	if err != nil {
		cmd.Process.Kill()
		cmd.Wait()
		return nil, fmt.Errorf("mpv: failed to connect to IPC socket: %w", err)
	}

	p := &Player{
		cmd:        cmd,
		conn:       conn,
		socketPath: socketPath,
		pending:    make(map[int64]chan response),
		done:       make(chan struct{}),
	}

	go p.readLoop()
	return p, nil
}

// readLoop reads responses from the mpv socket and dispatches them.
func (p *Player) readLoop() {
	defer close(p.done)
	scanner := bufio.NewScanner(p.conn)
	for scanner.Scan() {
		var resp response
		if err := json.Unmarshal(scanner.Bytes(), &resp); err != nil {
			continue
		}
		// Only dispatch responses with a request_id (skip async events).
		if resp.ReqID == 0 {
			continue
		}
		p.pendingM.Lock()
		ch, ok := p.pending[resp.ReqID]
		if ok {
			delete(p.pending, resp.ReqID)
		}
		p.pendingM.Unlock()
		if ok {
			ch <- resp
		}
	}
}

// sendCommand sends a JSON IPC command and waits for the response.
func (p *Player) sendCommand(ctx context.Context, args ...interface{}) (json.RawMessage, error) {
	id := p.reqID.Add(1)

	ch := make(chan response, 1)
	p.pendingM.Lock()
	p.pending[id] = ch
	p.pendingM.Unlock()

	cmd := ipcCommand{Command: args, RequestID: id}
	data, err := json.Marshal(cmd)
	if err != nil {
		p.pendingM.Lock()
		delete(p.pending, id)
		p.pendingM.Unlock()
		return nil, err
	}
	data = append(data, '\n')

	p.mu.Lock()
	_, err = p.conn.Write(data)
	p.mu.Unlock()
	if err != nil {
		p.pendingM.Lock()
		delete(p.pending, id)
		p.pendingM.Unlock()
		return nil, fmt.Errorf("mpv: write failed: %w", err)
	}

	select {
	case resp := <-ch:
		if resp.Error != "" && resp.Error != "success" {
			return nil, fmt.Errorf("mpv: %s", resp.Error)
		}
		return resp.Data, nil
	case <-ctx.Done():
		p.pendingM.Lock()
		delete(p.pending, id)
		p.pendingM.Unlock()
		return nil, ctx.Err()
	case <-p.done:
		return nil, fmt.Errorf("mpv: connection closed")
	}
}

// Play loads and starts playing the given stream URL.
func (p *Player) Play(ctx context.Context, streamURL string) error {
	_, err := p.sendCommand(ctx, "loadfile", streamURL)
	return err
}

// Pause pauses playback.
func (p *Player) Pause(ctx context.Context) error {
	_, err := p.sendCommand(ctx, "set_property", "pause", true)
	return err
}

// Resume resumes playback.
func (p *Player) Resume(ctx context.Context) error {
	_, err := p.sendCommand(ctx, "set_property", "pause", false)
	return err
}

// Stop stops playback.
func (p *Player) Stop(ctx context.Context) error {
	_, err := p.sendCommand(ctx, "stop")
	return err
}

// Seek seeks to an absolute position.
func (p *Player) Seek(ctx context.Context, position time.Duration) error {
	_, err := p.sendCommand(ctx, "seek", position.Seconds(), "absolute")
	return err
}

// SeekRelative seeks forward or backward by the given duration.
func (p *Player) SeekRelative(ctx context.Context, delta time.Duration) error {
	_, err := p.sendCommand(ctx, "seek", delta.Seconds(), "relative")
	return err
}

// SetVolume sets the volume (0-100).
func (p *Player) SetVolume(ctx context.Context, volume int) error {
	_, err := p.sendCommand(ctx, "set_property", "volume", float64(volume))
	return err
}

// GetPosition returns the current playback position.
func (p *Player) GetPosition(ctx context.Context) (time.Duration, error) {
	data, err := p.sendCommand(ctx, "get_property", "time-pos")
	if err != nil {
		return 0, err
	}
	var seconds float64
	if err := json.Unmarshal(data, &seconds); err != nil {
		return 0, fmt.Errorf("mpv: failed to parse position: %w", err)
	}
	return time.Duration(seconds * float64(time.Second)), nil
}

// GetDuration returns the total duration of the current track.
func (p *Player) GetDuration(ctx context.Context) (time.Duration, error) {
	data, err := p.sendCommand(ctx, "get_property", "duration")
	if err != nil {
		return 0, err
	}
	var seconds float64
	if err := json.Unmarshal(data, &seconds); err != nil {
		return 0, fmt.Errorf("mpv: failed to parse duration: %w", err)
	}
	return time.Duration(seconds * float64(time.Second)), nil
}

// IsPlaying returns true if mpv is currently playing (not paused).
func (p *Player) IsPlaying(ctx context.Context) (bool, error) {
	data, err := p.sendCommand(ctx, "get_property", "pause")
	if err != nil {
		return false, err
	}
	var paused bool
	if err := json.Unmarshal(data, &paused); err != nil {
		return false, fmt.Errorf("mpv: failed to parse pause state: %w", err)
	}
	return !paused, nil
}

// Close shuts down the mpv process and cleans up.
func (p *Player) Close() error {
	// Try graceful quit first.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	p.sendCommand(ctx, "quit")

	p.conn.Close()

	// Wait for mpv to exit.
	done := make(chan error, 1)
	go func() { done <- p.cmd.Wait() }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		p.cmd.Process.Kill()
		<-done
	}

	os.Remove(p.socketPath)
	return nil
}
