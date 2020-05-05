// Package transport implements various handlers for http server
package transport

import (
	"net/http"

	httplib "github.com/diptanw/server-detector/internal/platform/http"
	"github.com/diptanw/server-detector/internal/platform/jsonapi"
	"github.com/diptanw/server-detector/internal/platform/logger"
	"github.com/diptanw/server-detector/internal/view"
)

// Controller is a type that provides view handlers
type Controller struct {
	log logger.Logger
	svc *view.Service
}

// NewController creates a new instance of the Controller type
func NewController(log logger.Logger, service *view.Service) *Controller {
	return &Controller{
		log: log,
		svc: service,
	}
}

// RegisterHandlers registers all endpoints for the views handler
func (c Controller) RegisterHandlers(mux *httplib.Mux) {
	mux.AddRoute("GET", "v1/views", c.getHandler(mux))
	mux.AddRoute("GET", "v1/views/:reqId", c.getByRequestIDHandler(mux))
	mux.AddRoute("GET", "health", c.healthHandler)
}

func (c *Controller) getHandler(_ *httplib.Mux) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		res, err := c.svc.Get()
		if err != nil {
			c.handleError(wr, err)
			return
		}

		jsonapi.Write(wr, toPageResponse(res))
	}
}

func (c *Controller) getByRequestIDHandler(mux *httplib.Mux) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		params := mux.GetParams(req.Method, req.URL.Path)

		res, err := c.svc.GetByRequestID(params["reqId"])
		if err != nil {
			c.handleError(wr, err)
			return
		}

		jsonapi.Write(wr, toSingleResponse(res))
	}
}

func (c *Controller) healthHandler(wr http.ResponseWriter, req *http.Request) {
	wr.WriteHeader(http.StatusOK)
}

func (c *Controller) handleError(w http.ResponseWriter, err error) {
	if status := toStatus(err); status != http.StatusOK {
		c.log.Warnf("request has failed with error: %s", err)

		w.WriteHeader(status)
		jsonapi.Write(w, jsonapi.ErrorResponse{Errors: []string{err.Error()}})
	}
}
