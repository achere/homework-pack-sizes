package server

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log/slog"

	"github.com/achere/homework-pack-sizes/internal/pack"
	"github.com/joeshaw/envdecode"
)

//go:embed templates
var content embed.FS

const (
	defaultPort = 8080
)

type App struct {
	Config   *Config
	logger   *slog.Logger
	SizeRepo pack.PackSizeRepo
	template *template.Template
}

type Config struct {
	Port  int    `env:"PORT"`
	Order string `env:"ORDER"`
	DbUrl string `env:"DB_URL"`
}

// NewApp creates a new App, initialising the config from environment variables.
func NewApp(ctx context.Context, logger *slog.Logger) (*App, error) {
	app := &App{}

	app.logger = logger
	app.Config = &Config{}

	err := envdecode.Decode(app.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to load from env: %w", err)
	}

	if app.Config.DbUrl == "" {
		return nil, fmt.Errorf("DB_URL was not set")
	}
	if app.Config.Port == 0 {
		app.Config.Port = defaultPort
	}

	app.template, err = template.ParseFS(content, "templates/index.html")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return app, nil
}
