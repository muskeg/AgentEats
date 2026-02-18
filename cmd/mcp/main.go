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
	case "sse":
		addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.MCPPort)
		sseServer := server.NewSSEServer(s, server.WithBaseURL(fmt.Sprintf("http://%s", addr)))
		log.Printf("🤖 AgentEats MCP server starting (SSE transport on %s)", addr)
		if err := sseServer.Start(addr); err != nil {
			log.Fatalf("MCP SSE server error: %v", err)
		}
	default:
		log.Println("🤖 AgentEats MCP server starting (stdio transport)")
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("MCP server error: %v", err)
		}
	}
}
