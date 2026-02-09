package storage

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/mganuesquest/youtui/internal/domain"
)

// JSONStorage implements port.Storage using JSON files.
type JSONStorage struct {
	filePath string
}

// New creates a JSONStorage that persists queue state to the given file path.
func New(filePath string) *JSONStorage {
	return &JSONStorage{filePath: filePath}
}

// SaveQueue marshals the queue to JSON and writes it atomically to disk.
func (s *JSONStorage) SaveQueue(queue *domain.Queue) error {
	data, err := json.MarshalIndent(queue, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(s.filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	// Atomic write: write to temp file, then rename.
	tmp := s.filePath + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.filePath)
}

// LoadQueue reads the queue from the JSON file.
// If the file does not exist, it returns a new empty queue.
func (s *JSONStorage) LoadQueue() (*domain.Queue, error) {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return domain.NewQueue(), nil
		}
		return nil, err
	}

	var queue domain.Queue
	if err := json.Unmarshal(data, &queue); err != nil {
		return nil, err
	}
	return &queue, nil
}
