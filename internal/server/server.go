package server

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/joeshaw/envdecode"
)

const (
	defaultSizes = "250,500,1000,2000,5000"
	defaultPort  = 8080
)

type App struct {
	Config *Config
	logger *slog.Logger
}

type Config struct {
	Port     int    `env:"PORT"`
	Order    string `env:"ORDER"`
	SizesStr string `env:"SIZES"`
	Sizes    []int
}

// NewApp creates a new App, initialising the config from environment variables.
func NewApp(logger *slog.Logger) (*App, error) {
	app := &App{}

	app.logger = logger
	app.Config = &Config{}

	err := envdecode.Decode(app.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to load from env: %w", err)
	}

	if app.Config.Port == 0 {
		app.Config.Port = defaultPort
	}
	if app.Config.SizesStr == "" {
		app.Config.SizesStr = defaultSizes
	}

	sizeStrings := strings.Split(app.Config.SizesStr, ",")
	sizes := make([]int, len(sizeStrings))
	for i, sizeStr := range sizeStrings {
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			return nil, fmt.Errorf("failed to convert size to int: %s", sizeStr)
		}
		sizes[i] = size
	}
	app.Config.Sizes = sizes

	return app, nil
}
