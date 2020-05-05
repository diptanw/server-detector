package main

import (
	"flag"

	"github.com/diptanw/server-detector/internal/platform/logger"
)

// Config contains server configuration
type Config struct {
	HTTPAddr    string
	LogLevel    int
	WorkersNum  int
	NATSServer  string
	NATSChannel string
}

// ReadConfig reads config values from command args
func ReadConfig() Config {
	config := Config{}

	flag.StringVar(&config.HTTPAddr, "addr", ":8080", "an address for http server")
	flag.StringVar(&config.NATSServer, "nats-server", "nats://0.0.0.0:4222", "an address to NATS server")
	flag.StringVar(&config.NATSChannel, "nats-channel", "events/detect", "a channel where to publish events")
	flag.IntVar(&config.LogLevel, "level", logger.Info, "a logging level")
	flag.IntVar(&config.WorkersNum, "workers", 3, "a concurrent workers number")
	flag.Parse()

	return config
}
