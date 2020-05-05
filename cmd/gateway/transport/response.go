package transport

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/diptanw/server-detector/internal/relay"
)

func getLocation(id string) string {
	return fmt.Sprintf("v1/detects/%s", id)
}

func toStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if errors.Is(err, relay.ErrInvalidHost) || errors.Is(err, relay.ErrBadRequestID) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}
