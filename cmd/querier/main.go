package main

import (
	"context"
	"net/http"
	"os"

	"github.com/diptanw/server-detector/cmd/querier/transport"
	"github.com/diptanw/server-detector/internal/platform/events"
	httplib "github.com/diptanw/server-detector/internal/platform/http"
	"github.com/diptanw/server-detector/internal/platform/logger"
	"github.com/diptanw/server-detector/internal/platform/server"
	"github.com/diptanw/server-detector/internal/platform/storage"
	"github.com/diptanw/server-detector/internal/view"
)

func main() {
	// Read command line arguments
	config := ReadConfig()
	// Setup logger output and level
	log := logger.New(os.Stdout, config.LogLevel)

	// Setup persistence layer
	db := storage.NewInMemory()

	// Setup services
	repository := view.NewRepository(db)
	service := view.NewService(repository)
	handler := view.NewEventHandler(repository, log)

	// Setup messenger
	messenger, err := events.Connect(config.NATSServer)
	if err != nil {
		panic(err)
	}

	defer messenger.Close()

	if err := messenger.Subscribe(config.NATSChannel, handler.Handle); err != nil {
		panic(err)
	}

	// Register all endpoints
	mux := httplib.Mux{}
	controller := transport.NewController(log, service)
	controller.RegisterHandlers(&mux)

	srv := server.New(&http.Server{
		Addr:    config.HTTPAddr,
		Handler: &mux,
	}, log)

	if err := srv.Serve(context.Background()); err != nil {
		panic(err)
	}
}
