package processor

import (
	"fmt"
	"time"

	"github.com/diptanw/server-detector/internal/platform/storage"
)

// Store is an interface type that defines storage abstraction
type Store interface {
	Find(m storage.Matcher) ([]interface{}, error)
	Insert(r interface{}) error
}

// Repository is a type that manipulates detect entity data
type Repository struct {
	db Store
}

// DetectCommand is a type that represents detect comand model
type DetectCommand struct {
	ID        storage.ID `json:"id"`        // ID is an event identifier
	RequestID string     `json:"requestID"` // RequestID is the a group identifier of submitted request
	Host      string     `json:"host"`      // Host is a host name to detect
	CreatedAt time.Time  `json:"createdAt"` // CreatedAt is when event has been recorded
}

// NewRepository creates anew instance of Repository
func NewRepository(db Store) Repository {
	return Repository{
		db: db,
	}
}

// Get returns a detect command entity for the specified ID
func (r Repository) Get(id string) (DetectCommand, error) {
	res, err := r.db.Find(func(value interface{}) bool {
		if c, ok := value.(DetectCommand); ok {
			return string(c.ID) == id
		}

		return false
	})

	if err != nil {
		return DetectCommand{}, fmt.Errorf("get command: %w", err)
	}

	if len(res) == 0 {
		return DetectCommand{}, storage.ErrNotFound
	}

	return res[0].(DetectCommand), nil
}

// Create creates a new detect command entity
func (r Repository) Create(d DetectCommand) (DetectCommand, error) {
	d.ID = storage.NewID()
	d.CreatedAt = time.Now().UTC()

	if err := r.db.Insert(d); err != nil {
		return DetectCommand{}, fmt.Errorf("save command: %w", err)
	}

	return d, nil
}
