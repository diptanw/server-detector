package view

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/diptanw/server-detector/internal/platform/events"
	"github.com/diptanw/server-detector/internal/platform/logger"
	"github.com/diptanw/server-detector/internal/platform/storage"
)

// EventHandler is a struct that materializes views
type EventHandler struct {
	repository Repository
	logger     logger.Logger
}

func NewEventHandler(r Repository, log logger.Logger) *EventHandler {
	return &EventHandler{
		repository: r,
		logger:     log,
	}
}

func (h *EventHandler) Handle(msg *events.Message) {
	event := struct {
		EventID   string `json:"eventId"`
		RequestID string `json:"requestId"`
		Host      Host   `json:"host"`
	}{}

	if err := json.Unmarshal(msg.Data, &event); err != nil {
		h.logger.Errorf("unable to unmarshal event data: %w", err)
		return
	}

	if event.RequestID == "" || event.EventID == "" {
		h.logger.Errorf("bad event data: %#v", event)
		return
	}

	view, err := h.repository.Get(event.RequestID)
	if errors.Is(err, storage.ErrNotFound) {
		if _, err := h.repository.Create(DetectView{
			ID:        storage.NewID(),
			RequestID: event.RequestID,
			Hosts:     []Host{event.Host},
			CreatedAt: time.Now().UTC(),
		}); err != nil {
			h.logger.Errorf("unable to create view: %w", err)
		}

		return
	}

	if err != nil {
		h.logger.Errorf("unable to fetch view: %w", err)
		return
	}

	view.Hosts = append(view.Hosts, event.Host)

	if _, err := h.repository.Update(view); err != nil {
		h.logger.Errorf("unable to update view: %w", err)
	}
}
