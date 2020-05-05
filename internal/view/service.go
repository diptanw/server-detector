package view

import "errors"

// ErrBadRequestID is an error when request ID is malformed
var ErrBadRequestID = errors.New("request ID cannot be empty")

// Service is a struct that provides main view operations
type Service struct {
	repo Repository
}

// NewService returns a new Service instance
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetByRequestID returns a view for a given request ID
func (s *Service) GetByRequestID(id string) (DetectView, error) {
	if id == "" {
		return DetectView{}, ErrBadRequestID
	}

	return s.repo.Get(id)
}

// Get returns all created views
func (s *Service) Get() ([]DetectView, error) {
	return s.repo.GetAll()
}
