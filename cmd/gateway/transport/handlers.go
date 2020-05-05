// Package transport implements various handlers for http server
package transport

import (
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"

	httplib "github.com/diptanw/server-detector/internal/platform/http"
	"github.com/diptanw/server-detector/internal/platform/jsonapi"
	"github.com/diptanw/server-detector/internal/platform/logger"
	"github.com/diptanw/server-detector/internal/relay"
)

// Controller is a type that provides view handlers
type Controller struct {
	log logger.Logger
	svc *relay.Relay
}

// NewController creates a new instance of the Controller type
func NewController(log logger.Logger, service *relay.Relay) *Controller {
	return &Controller{
		log: log,
		svc: service,
	}
}

// RegisterHandlers registers all endpoints for the views handler
func (c *Controller) RegisterHandlers(mux *httplib.Mux) {
	mux.AddRoute("GET", "v1/detects/:reqID", c.getByRequestIDHandler(mux))
	mux.AddRoute("GET", "v1/detects", c.getHandler(mux))
	mux.AddRoute("POST", "v1/detects", c.postHandler(mux))
	mux.AddRoute("GET", "health", c.healthHandler)
}

// RegisterOpenAPIHandler registers an Open API Specification endpoint
func (c *Controller) RegisterOpenAPIHandler(mux *httplib.Mux, path string) {
	mux.AddRoute("GET", "v1/openapi", c.getOpenAPIHandler(path))
}

func (c *Controller) postHandler(_ *httplib.Mux) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		var hosts struct {
			Data []string
		}

		if err := jsonapi.Read(req, &hosts); err != nil {
			c.handleError(wr, err)
			return
		}

		reqID, err := c.svc.PostCommand(hosts.Data)
		if err != nil {
			c.handleError(wr, err)
			return
		}

		wr.Header().Set("Location", getLocation(reqID))
		wr.WriteHeader(http.StatusAccepted)
	}
}

func (c *Controller) getByRequestIDHandler(mux *httplib.Mux) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		params := mux.GetParams(req.Method, req.URL.Path)

		res, err := c.svc.QueryView(params["reqID"])
		if err != nil {
			c.handleError(wr, err)
			return
		}

		c.write(wr, res)
	}
}

func (c *Controller) getHandler(_ *httplib.Mux) http.HandlerFunc {
	return func(wr http.ResponseWriter, req *http.Request) {
		res, err := c.svc.QueryAllViews()
		if err != nil {
			c.handleError(wr, err)
			return
		}

		c.write(wr, res)
	}
}

func (c *Controller) write(wr http.ResponseWriter, b []byte) {
	wr.Header().Set("Content-Length", strconv.Itoa(len(b)))
	wr.Header().Set("Content-Type", "application/json; charset=utf-8")

	if _, err := wr.Write(b); err != nil {
		c.handleError(wr, err)
		return
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

func (c *Controller) getOpenAPIHandler(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadFile(filepath.Clean(path))
		if err != nil {
			c.handleError(w, err)
			return
		}

		w.Header().Set("Content-Length", strconv.Itoa(len(data)))
		w.Header().Set("Content-Type", "application/x-yaml; charset=utf-8")

		if _, err := w.Write(data); err != nil {
			c.handleError(w, err)
			return
		}
	}
}
