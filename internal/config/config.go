package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Config holds all application configuration, loaded from environment variables.
type Config struct {
	Host         string `envconfig:"HOST" default:"0.0.0.0"`
	Port         int    `envconfig:"PORT" default:"8000"`
	DatabaseURL  string `envconfig:"DATABASE_URL" default:"agenteats.db"`
	Debug        bool   `envconfig:"DEBUG" default:"false"`
	MCPTransport string `envconfig:"MCP_TRANSPORT" default:"stdio"` // "stdio" or "http"
	MCPPort      int    `envconfig:"MCP_PORT" default:"8001"`
	CORSOrigins  string `envconfig:"CORS_ORIGINS" default:"*"` // comma-separated allowed origins
}

// Load reads configuration from environment variables.
func Load() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	return &cfg
}
