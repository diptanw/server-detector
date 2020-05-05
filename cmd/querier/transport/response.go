package transport

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/diptanw/server-detector/internal/platform/jsonapi"
	"github.com/diptanw/server-detector/internal/platform/storage"
	"github.com/diptanw/server-detector/internal/view"
)

func getLocation(id string) string {
	return fmt.Sprintf("v1/views/%s", id)
}

func toSingleResponse(v view.DetectView) jsonapi.DataResponse {
	return jsonapi.DataResponse{
		Data: &v,
		Links: jsonapi.ResourceLinks{
			Self: jsonapi.NewHref(getLocation(v.RequestID)),
		},
	}
}

func toPageResponse(vs []view.DetectView) jsonapi.PageResponse {
	data := make([]interface{}, len(vs))

	for i, v := range vs {
		data[i] = struct {
			view.DetectView
			jsonapi.LinksResponse
		}{
			LinksResponse: jsonapi.LinksResponse{
				Links: jsonapi.ResourceLinks{
					Self: jsonapi.NewHref(getLocation(v.RequestID)),
				},
			},
			DetectView: v,
		}
	}

	return jsonapi.PageResponse{
		Data: data,
		Links: jsonapi.PageLinks{
			Self: jsonapi.NewHref("v1/views"),
		},
	}
}

func toStatus(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if isError(err, view.ErrBadRequestID) {
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
