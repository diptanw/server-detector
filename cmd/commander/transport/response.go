package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/diptanw/server-detector/internal/platform/jsonapi"
	"github.com/diptanw/server-detector/internal/platform/storage"
	"github.com/diptanw/server-detector/internal/processor"
)

func getLocation(ID storage.ID) string {
	return fmt.Sprintf("v1/processor/%s", ID)
}

func toResponse(c processor.DetectCommand) jsonapi.DataResponse {
	return jsonapi.DataResponse{
		Data: c,
		Links: jsonapi.ResourceLinks{
			Self: jsonapi.NewHref(getLocation(c.ID)),
		},
	}
}

func toStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if isError(err, processor.ErrInvalidHost, processor.ErrBadRequestID, processor.ErrBadCommandID, &json.SyntaxError{}) {
		return http.StatusBadRequest
	}

	if isError(err, storage.ErrMissingID, storage.ErrNotFound) {
		return http.StatusNotFound
	}

	return http.StatusInternalServerError
}

func isError(err error, targets ...error) bool {
	for _, t := range targets {
		if errors.Is(err, t) {
			return true
		}
	}

	return false
}
