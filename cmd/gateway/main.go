package main

import (
	"context"
	"net/http"
	"os"

	"github.com/diptanw/server-detector/cmd/gateway/transport"
	httplib "github.com/diptanw/server-detector/internal/platform/http"
	"github.com/diptanw/server-detector/internal/platform/logger"
	"github.com/diptanw/server-detector/internal/platform/retry"
	"github.com/diptanw/server-detector/internal/platform/server"
	"github.com/diptanw/server-detector/internal/relay"
)

func main() {
	// Read command line arguments
	config := ReadConfig()
	// Setup logger output and level
	log := logger.New(os.Stdout, config.LogLevel)

	client := retry.NewHTTPClient(retry.DefaultPolicy())
	service := relay.New(config.CommandAddr, config.QueryAddr, client)

	// Register all endpoints
	mux := httplib.Mux{}
	// Allow CORS for swagger-ui requests
	mux.WithMiddleware(httplib.CORSMiddleware)

	controller := transport.NewController(log, service)
	controller.RegisterHandlers(&mux)
	controller.RegisterOpenAPIHandler(&mux, config.APIPath)

	srv := server.New(&http.Server{
		Addr:    config.HTTPAddr,
		Handler: &mux,
	}, log)

	if err := srv.Serve(context.Background()); err != nil {
		panic(err)
	}
}
