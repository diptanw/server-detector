package main

import (
	"context"
	"net/http"
	"os"

	"github.com/diptanw/server-detector/cmd/commander/transport"
	"github.com/diptanw/server-detector/internal/platform/events"
	httplib "github.com/diptanw/server-detector/internal/platform/http"
	"github.com/diptanw/server-detector/internal/platform/logger"
	"github.com/diptanw/server-detector/internal/platform/retry"
	"github.com/diptanw/server-detector/internal/platform/server"
	"github.com/diptanw/server-detector/internal/platform/storage"
	"github.com/diptanw/server-detector/internal/platform/worker"
	"github.com/diptanw/server-detector/internal/processor"
)

func main() {
	// Read command line arguments
	config := ReadConfig()
	// Setup logger output and level
	log := logger.New(os.Stdout, config.LogLevel)

	// Setup persistence layer
	db := storage.NewInMemory()

	// Setup messenger
	messenger, err := events.Connect(config.NATSServer)
	if err != nil {
		panic(err)
	}

	defer messenger.Close()

	// Setup workers pool
	ctx := context.Background()
	pool := worker.NewPool(config.WorkersNum)

	pool.Run(ctx)
	defer pool.Close()

	// Setup services
	repository := processor.NewRepository(db)
	detect := processor.NewDetector(retry.NewHTTPClient(retry.DefaultPolicy()))
	service := processor.NewService(repository, detect, pool, log, messenger, config.NATSChannel)

	// Register all endpoints
	mux := httplib.Mux{}
	controller := transport.NewController(log, service)
	controller.RegisterHandlers(&mux)

	srv := server.New(&http.Server{
		Addr:    config.HTTPAddr,
		Handler: &mux,
	}, log)

	if err := srv.Serve(ctx); err != nil {
		panic(err)
	}
}
