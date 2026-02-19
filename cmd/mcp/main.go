package main

import (
	"fmt"
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

	switch cfg.MCPTransport {
	case "http":
		addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.MCPPort)
		httpServer := server.NewStreamableHTTPServer(s, server.WithStateLess(true))
		log.Printf("🤖 AgentEats MCP server starting (Streamable HTTP on %s/mcp)", addr)
		if err := httpServer.Start(addr); err != nil {
			log.Fatalf("MCP HTTP server error: %v", err)
		}
	default:
		log.Println("🤖 AgentEats MCP server starting (stdio transport)")
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("MCP server error: %v", err)
		}
	}
}
