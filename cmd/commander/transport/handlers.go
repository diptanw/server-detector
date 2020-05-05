// Package transport implements various handlers for http server
package transport

import (
	"net/http"

	httplib "github.com/diptanw/server-detector/internal/platform/http"
	"github.com/diptanw/server-detector/internal/platform/jsonapi"
	"github.com/diptanw/server-detector/internal/platform/logger"
	"github.com/diptanw/server-detector/internal/processor"
)

// Controller is a type that provides command handlers
type Controller struct {
	log logger.Logger
	srv *processor.Service
}

// NewController creates a new instance of the Controller type
func NewController(log logger.Logger, service *processor.Service) *Controller {
	return &Controller{
		log: log,
		srv: service,
	}
}

// RegisterHandlers registers all transport routes for the command handler
func (c *Controller) RegisterHandlers(mux *httplib.Mux) {
	mux.AddRoute("POST", "v1/processor", c.postHandler(mux))
	mux.AddRoute("GET", "v1/processor/:id", c.getByIDHandler(mux))
	mux.AddRoute("GET", "health", c.healthHandler)
}

func (c *Controller) postHandler(_ *httplib.Mux) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		var cmd processor.DetectCommand

		if err := jsonapi.Read(req, &cmd); err != nil {
			c.handleError(wr, err)
			return
		}

		res, err := c.srv.Submit(cmd)
		if err != nil {
			c.handleError(wr, err)
			return
		}

		wr.Header().Set("Location", getLocation(res.ID))
		wr.WriteHeader(http.StatusCreated)
	}
}

func (c *Controller) getByIDHandler(mux *httplib.Mux) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		params := mux.GetParams(req.Method, req.URL.Path)

		res, err := c.srv.Fetch(params["id"])
		if err != nil {
			c.handleError(wr, err)
			return
		}

		jsonapi.Write(wr, toResponse(res))
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
