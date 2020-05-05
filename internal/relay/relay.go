package relay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/diptanw/server-detector/internal/platform/storage"
)

var (
	// ErrInvalidHost is an error hosts are invalid or not specified
	ErrInvalidHost = errors.New("invalid hosts")
	// ErrBadRequestID is an error when request ID is malformed
	ErrBadRequestID = errors.New("request ID cannot be empty")
)

// Relay is a struct that redirects the requests to corresponding parties
type Relay struct {
	client     *http.Client
	commandURL string
	queryURL   string
}

// New creates an ew instance of Relay service
func New(commandURL, queryURL string, client *http.Client) *Relay {
	return &Relay{
		client:     client,
		commandURL: commandURL,
		queryURL:   queryURL,
	}
}

// PostCommand sends individual detect commands for all requested hosts
func (r *Relay) PostCommand(hosts []string) (string, error) {
	reqID := string(storage.NewID())

	for _, h := range hosts {
		b, err := json.Marshal(struct {
			RequestID string `json:"requestID"`
			Host      string `json:"host"`
		}{
			RequestID: reqID,
			Host:      h,
		})

		if err != nil {
			return "", fmt.Errorf("unable to marshal command body: %w", err)
		}

		resp, err := r.client.Post(r.commandURL, "application/json", bytes.NewBuffer(b))
		if err != nil {
			return "", fmt.Errorf("command request failed: %w", err)
		}

		resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return "", fmt.Errorf("unexpected status code: %d ", resp.StatusCode)
		}
	}

	return reqID, nil
}

// QueryView returns the requested view data from query service
func (r *Relay) QueryView(reqID string) ([]byte, error) {
	if reqID == "" {
		return nil, ErrBadRequestID
	}

	return r.get(fmt.Sprintf("%s/%s", r.queryURL, reqID))
}

// QueryAllViews returns the all views from query service
func (r *Relay) QueryAllViews() ([]byte, error) {
	return r.get(r.queryURL)
}

func (r *Relay) get(url string) ([]byte, error) {
	resp, err := r.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("query request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d ", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read response body: %w", err)
	}

	return b, nil
}
