package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config holds all application configuration, loaded from environment variables.
type Config struct {
	Host        string `envconfig:"HOST" default:"0.0.0.0"`
	Port        int    `envconfig:"PORT" default:"8000"`
	DatabaseURL string `envconfig:"DATABASE_URL" default:"agenteats.db"`
	Debug       bool   `envconfig:"DEBUG" default:"false"`
}

// Load reads configuration from environment variables.
func Load() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return &cfg
}
