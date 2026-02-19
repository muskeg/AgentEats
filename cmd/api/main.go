package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	mcphttp "github.com/mark3labs/mcp-go/server"

	"github.com/agenteats/agenteats/internal/config"
	"github.com/agenteats/agenteats/internal/database"
	"github.com/agenteats/agenteats/internal/handlers"
	"github.com/agenteats/agenteats/internal/mcpserver"
	authmw "github.com/agenteats/agenteats/internal/middleware"
)

func main() {
	cfg := config.Load()
	database.Init(cfg)

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Get("/health", handlers.Health)

	// --- Public (read-only) ---
	r.Get("/restaurants", handlers.SearchRestaurants)
	r.Get("/restaurants/{restaurantID}", handlers.GetRestaurant)
	r.Get("/restaurants/{restaurantID}/menu", handlers.GetMenu)
	r.Get("/restaurants/{restaurantID}/availability", handlers.CheckAvailability)
	r.Post("/restaurants/{restaurantID}/reservations", handlers.MakeReservation)
	r.Get("/restaurants/{restaurantID}/reservations", handlers.ListReservations)
	r.Delete("/reservations/{reservationID}", handlers.CancelReservation)
	r.Get("/recommendations", handlers.GetRecommendations)

	// --- Owner registration (no auth) ---
	r.Post("/owners/register", handlers.RegisterOwner)

	// --- Authenticated owner routes ---
	r.Group(func(r chi.Router) {
		r.Use(authmw.RequireAPIKey)

		r.Post("/owners/rotate-key", handlers.RotateKey)

		// Restaurant management
		r.Post("/restaurants", handlers.CreateOwnedRestaurant)
		r.Put("/restaurants/{restaurantID}", handlers.UpdateOwnedRestaurant)

		// Menu management
		r.Post("/restaurants/{restaurantID}/menu/items", handlers.AddOwnedMenuItem)
		r.Post("/restaurants/{restaurantID}/menu/import", handlers.BulkImportMenu)
	})

	// --- Remote MCP (Streamable HTTP) ---
	mcpSrv := mcpserver.NewServer()
	httpMCP := mcphttp.NewStreamableHTTPServer(mcpSrv,
		mcphttp.WithStateLess(true), // fully stateless — scale-to-zero friendly
	)
	r.Mount("/mcp", httpMCP)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("🚀 AgentEats API server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
