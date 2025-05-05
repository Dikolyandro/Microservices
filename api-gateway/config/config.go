package config

import (
	"github.com/caarlos0/env/v10"
)

type Config struct {
	ProductServiceURL string `env:"PRODUCT_SERVICE_URL" envDefault:"http://localhost:8081"`
	OrderServiceURL   string `env:"ORDER_SERVICE_URL" envDefault:"http://localhost:8082"`
	Port              string `env:"PORT" envDefault:"8080"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
