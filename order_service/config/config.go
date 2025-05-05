package config

import (
	"github.com/recktt77/Microservices-First-/order_service/pkg/mongo"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Mongo   mongo.Config
	Server  Server
	Version string `env:"VERSION"`
}

type Server struct {
	HTTPServer HTTPServer
}

type HTTPServer struct {
	Port           int           `env:"HTTP_PORT"`
	ReadTimeout    time.Duration `env:"HTTP_READ_TIMEOUT"`
	WriteTimeout   time.Duration `env:"HTTP_WRITE_TIMEOUT"`
	IdleTimeout    time.Duration `env:"HTTP_IDLE_TIMEOUT"`
	MaxHeaderBytes int           `env:"HTTP_MAX_HEADER_BYTES"`
	TrustedProxies []string      `env:"HTTP_TRUSTED_PROXIES" envSeparator:","`
	Mode           string        `env:"GIN_MODE"`
}

func New() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	// Set default values
	if cfg.Server.HTTPServer.Port == 0 {
		cfg.Server.HTTPServer.Port = 8080
	}
	if cfg.Server.HTTPServer.ReadTimeout == 0 {
		cfg.Server.HTTPServer.ReadTimeout = 30 * time.Second
	}
	if cfg.Server.HTTPServer.WriteTimeout == 0 {
		cfg.Server.HTTPServer.WriteTimeout = 30 * time.Second
	}
	if cfg.Server.HTTPServer.IdleTimeout == 0 {
		cfg.Server.HTTPServer.IdleTimeout = 60 * time.Second
	}
	if cfg.Server.HTTPServer.MaxHeaderBytes == 0 {
		cfg.Server.HTTPServer.MaxHeaderBytes = 1048576
	}
	if cfg.Server.HTTPServer.Mode == "" {
		cfg.Server.HTTPServer.Mode = "release"
	}

	return &cfg, nil
}
