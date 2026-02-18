package main

import (
	"log"

	"github.com/mark3labs/mcp-go/server"

	"github.com/agenteats/agenteats/internal/config"
	"github.com/agenteats/agenteats/internal/database"
	mcpserver "github.com/agenteats/agenteats/internal/mcpserver"
)

func main() {
	cfg := config.Load()
	database.Init(cfg)

	s := mcpserver.NewServer()

	log.Println("🤖 AgentEats MCP server starting (stdio transport)")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}
