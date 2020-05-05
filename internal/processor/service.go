package processor

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/diptanw/server-detector/internal/platform/events"
	"github.com/diptanw/server-detector/internal/platform/logger"
	"github.com/diptanw/server-detector/internal/platform/worker"
)

var (
	// ErrBadCommandID is an error when ID is empty
	ErrBadCommandID = errors.New("command ID cannot be empty")
	// ErrBadRequestID is an error when request ID is empty
	ErrBadRequestID = errors.New("request ID cannot be empty")
	// ErrInvalidHost is an error when host is empty of malformed
	ErrInvalidHost = errors.New("invalid hostname")
)

// Service is a type that provides main detects operations
type Service struct {
	repository Repository
	detector   Detector
	pool       worker.Pool
	logger     logger.Logger
	messenger  *events.Messenger
	subject    string
}

// NewService returns a new detect service instance
func NewService(r Repository, d Detector, p worker.Pool, l logger.Logger, m *events.Messenger, s string) *Service {
	return &Service{
		repository: r,
		detector:   d,
		pool:       p,
		logger:     l,
		messenger:  m,
		subject:    s,
	}
}

// Submit submits new detect command and validates it's state
func (s *Service) Submit(cmd DetectCommand) (DetectCommand, error) {
	if cmd.RequestID == "" {
		return cmd, ErrBadRequestID
	}

	if cmd.Host == "" {
		return cmd, ErrInvalidHost
	}

	cmd, err := s.repository.Create(cmd)

	s.pool.Enqueue(func(ctx context.Context) {
		peer, err := s.detector.Detect(ctx, cmd.Host)
		if err != nil {
			s.logger.Errorf("unable to detect peer for host %s: %s", cmd.Host, err)
			return
		}

		if peer.Server == "" {
			s.logger.Warnf("no software information detected")
		}

		data, err := json.Marshal(struct {
			EventID   string `json:"eventId"`
			RequestID string `json:"requestId"`
			Host      Host   `json:"host"`
		}{
			EventID:   string(cmd.ID),
			RequestID: cmd.RequestID,
			Host:      peer,
		})

		if err != nil {
			s.logger.Errorf("unable to marshal event data: %s", err)
			return
		}

		err = s.messenger.Publish(s.subject, events.Message{
			Data: data,
		})

		if err != nil {
			s.logger.Errorf("unable to publish event: %s", err)
		}
	})

	return cmd, err
}

// Fetch returns a detect command for the given ID
func (s *Service) Fetch(id string) (DetectCommand, error) {
	if id == "" {
		return DetectCommand{}, ErrBadCommandID
	}

	return s.repository.Get(id)
}
