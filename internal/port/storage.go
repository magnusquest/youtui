package port

import "github.com/mganuesquest/youtui/internal/domain"

// Storage defines operations for persisting queue state
type Storage interface {
	// SaveQueue persists the queue to storage
	SaveQueue(queue *domain.Queue) error

	// LoadQueue loads the queue from storage
	LoadQueue() (*domain.Queue, error)
}
