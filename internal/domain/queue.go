package domain

// Queue represents the playback queue
type Queue struct {
	Tracks  []Track `json:"tracks" yaml:"tracks"`
	Current int     `json:"current" yaml:"current"`
}

// NewQueue creates an empty queue
func NewQueue() *Queue {
	return &Queue{
		Tracks:  make([]Track, 0),
		Current: -1,
	}
}

// Add appends a track to the queue
func (q *Queue) Add(track Track) {
	q.Tracks = append(q.Tracks, track)
	if q.Current == -1 {
		q.Current = 0
	}
}

// Remove removes a track at the given index
func (q *Queue) Remove(index int) bool {
	if index < 0 || index >= len(q.Tracks) {
		return false
	}
	q.Tracks = append(q.Tracks[:index], q.Tracks[index+1:]...)
	if q.Current >= len(q.Tracks) {
		q.Current = len(q.Tracks) - 1
	}
	return true
}

// MoveUp moves the track at index one position up
func (q *Queue) MoveUp(index int) bool {
	if index <= 0 || index >= len(q.Tracks) {
		return false
	}
	q.Tracks[index], q.Tracks[index-1] = q.Tracks[index-1], q.Tracks[index]
	if q.Current == index {
		q.Current--
	} else if q.Current == index-1 {
		q.Current++
	}
	return true
}

// MoveDown moves the track at index one position down
func (q *Queue) MoveDown(index int) bool {
	if index < 0 || index >= len(q.Tracks)-1 {
		return false
	}
	q.Tracks[index], q.Tracks[index+1] = q.Tracks[index+1], q.Tracks[index]
	if q.Current == index {
		q.Current++
	} else if q.Current == index+1 {
		q.Current--
	}
	return true
}

// Clear removes all tracks from the queue
func (q *Queue) Clear() {
	q.Tracks = make([]Track, 0)
	q.Current = -1
}

// CurrentTrack returns the currently selected track
func (q *Queue) CurrentTrack() *Track {
	if q.Current < 0 || q.Current >= len(q.Tracks) {
		return nil
	}
	return &q.Tracks[q.Current]
}

// Next advances to the next track, returns false if at end
func (q *Queue) Next() bool {
	if q.Current < len(q.Tracks)-1 {
		q.Current++
		return true
	}
	return false
}

// Previous goes to the previous track, returns false if at start
func (q *Queue) Previous() bool {
	if q.Current > 0 {
		q.Current--
		return true
	}
	return false
}

// IsEmpty returns true if the queue has no tracks
func (q *Queue) IsEmpty() bool {
	return len(q.Tracks) == 0
}

// Len returns the number of tracks in the queue
func (q *Queue) Len() int {
	return len(q.Tracks)
}
