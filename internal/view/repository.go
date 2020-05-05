package view

import (
	"fmt"
	"time"

	"github.com/diptanw/server-detector/internal/platform/storage"
)

// Store is an interface type that defines storage abstraction
type Store interface {
	Find(m storage.Matcher) ([]interface{}, error)
	Insert(r interface{}) error
	Update(r interface{}) error
}

// Repository is a type that manipulates view entity data
type Repository struct {
	db Store
}

// DetectView is a type that represents view data model
type DetectView struct {
	ID        storage.ID `json:"id"`
	RequestID string     `json:"requestId"`
	Hosts     []Host     `json:"hosts"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

// Host represens the detected host info
type Host struct {
	Domain string   `json:"domain"`
	Server string   `json:"server"`
	IPs    []string `json:"ips"`
}

// NewRepository creates anew instance of Repository type
func NewRepository(db Store) Repository {
	return Repository{
		db: db,
	}
}

// Get returns a view for the specified request ID
func (r Repository) Get(requestID string) (DetectView, error) {
	res, err := r.db.Find(func(value interface{}) bool {
		if c, ok := value.(DetectView); ok {
			return c.RequestID == requestID
		}

		return false
	})

	if err != nil {
		return DetectView{}, fmt.Errorf("get view: %w", err)
	}

	if len(res) == 0 {
		return DetectView{}, storage.ErrNotFound
	}

	return res[0].(DetectView), nil
}

// GetAll returns all views from the data store
func (r Repository) GetAll() ([]DetectView, error) {
	res, err := r.db.Find(func(value interface{}) bool {
		_, ok := value.(DetectView)
		return ok
	})

	if err != nil {
		return []DetectView{}, fmt.Errorf("get all views: %w", err)
	}

	views := make([]DetectView, len(res))

	for i, withID := range res {
		views[i] = withID.(DetectView)
	}

	return views, nil
}

// Create creates a new view entity
func (r Repository) Create(c DetectView) (DetectView, error) {
	c.ID = storage.NewID()
	c.CreatedAt = time.Now().UTC()

	err := r.db.Insert(c)
	if err != nil {
		return DetectView{}, fmt.Errorf("create view: %w", err)
	}

	return c, nil
}

// Update updates an existing view entity
func (r Repository) Update(c DetectView) (DetectView, error) {
	now := time.Now().UTC()
	c.UpdatedAt = &now

	err := r.db.Update(c)
	if err != nil {
		return DetectView{}, fmt.Errorf("update view: %w", err)
	}

	return c, nil
}
