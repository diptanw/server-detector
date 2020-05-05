package main

import (
	"flag"

	"github.com/diptanw/server-detector/internal/platform/logger"
)

// Config contains server configuration
type Config struct {
	HTTPAddr    string
	QueryAddr   string
	CommandAddr string
	APIPath     string
	LogLevel    int
}

// ReadConfig reads config values from command args
func ReadConfig() Config {
	config := Config{}

	flag.StringVar(&config.HTTPAddr, "addr", ":8080", "an address for http server")
	flag.StringVar(&config.QueryAddr, "query-addr", "http://0.0.0.0:6080", "a n address to querier service")
	flag.StringVar(&config.CommandAddr, "command-addr", "http://0.0.0.0:6081", "a n address to commander service")
	flag.StringVar(&config.APIPath, "api", "./api/openapi.yml", "a path to the Open API Specification file")
	flag.IntVar(&config.LogLevel, "level", logger.Info, "a logging level")
	flag.Parse()

	return config
}
