// Package server provides a transport service for the application
package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diptanw/server-detector/internal/platform/logger"
)

// Server is a struct that wraps an http transport listener
type Server struct {
	srv HTTPServer
	log logger.Logger
}

// HTTPServer is an interface for http server listener
type HTTPServer interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

// New returns a new instance of Server
func New(srv HTTPServer, l logger.Logger) *Server {
	return &Server{srv: srv, log: l}
}

// Serve runs the bootstrapped http server listener
func (s *Server) Serve(ctx context.Context) error {
	errsCh := make(chan error)

	go func() {
		s.log.Infof("starting server...")
		errsCh <- s.srv.ListenAndServe()
	}()

	// Listen for interrupt and termination signals
	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGINT, syscall.SIGTERM)

	// Block until an error or a signal is received
	select {
	case err := <-errsCh:
		if err == http.ErrServerClosed {
			return err
		}
	case <-shutdownCh:
		s.log.Warnf("process terminated")
	}

	s.log.Infof("shutting down...")

	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.log.Errorf("shutdown error: %s", err)
	}

	// Cancel slow operations when exiting
	_, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	return err
}
