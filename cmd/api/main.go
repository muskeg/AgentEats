package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
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

	// CORS — configurable via CORS_ORIGINS env var
	allowedOrigins := []string{"*"}
	if cfg.CORSOrigins != "*" {
		allowedOrigins = strings.Split(cfg.CORSOrigins, ",")
		for i := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
		}
	}
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Get("/health", handlers.Health)

	// --- Public (read-only) ---
	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(100, time.Minute))
		r.Get("/restaurants", handlers.SearchRestaurants)
		r.Get("/restaurants/{restaurantID}", handlers.GetRestaurant)
		r.Get("/restaurants/{restaurantID}/menu", handlers.GetMenu)
		r.Get("/restaurants/{restaurantID}/availability", handlers.CheckAvailability)
		r.Get("/restaurants/{restaurantID}/reservations", handlers.ListReservations)
		r.Get("/recommendations", handlers.GetRecommendations)
	})

	// --- Reservation endpoints (rate-limited) ---
	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(20, time.Minute))
		r.Post("/restaurants/{restaurantID}/reservations", handlers.MakeReservation)
		r.Delete("/reservations/{reservationID}", handlers.CancelReservation)
	})

	// --- Owner registration (strict rate limit) ---
	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(5, time.Minute))
		r.Post("/owners/register", handlers.RegisterOwner)
	})

	// --- Authenticated owner routes ---
	r.Group(func(r chi.Router) {
		r.Use(authmw.RequireAPIKey)

		r.Post("/owners/rotate-key", handlers.RotateKey)

		// Restaurant management
		r.Get("/owners/restaurants", handlers.ListOwnedRestaurants)
		r.Post("/restaurants", handlers.CreateOwnedRestaurant)
		r.Put("/restaurants/{restaurantID}", handlers.UpdateOwnedRestaurant)

		// Menu management
		r.Post("/restaurants/{restaurantID}/menu/items", handlers.AddOwnedMenuItem)
		r.Post("/restaurants/{restaurantID}/menu/import", handlers.BulkImportMenu)
	})

	// --- Remote MCP (Streamable HTTP, rate-limited) ---
	mcpSrv := mcpserver.NewServer()
	httpMCP := mcphttp.NewStreamableHTTPServer(mcpSrv,
		mcphttp.WithStateLess(true), // fully stateless — scale-to-zero friendly
	)
	r.Group(func(r chi.Router) {
		r.Use(httprate.LimitByIP(60, time.Minute))
		r.Mount("/mcp", httpMCP)
	})

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("🚀 AgentEats API server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
