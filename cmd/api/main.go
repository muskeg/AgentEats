package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/agenteats/agenteats/internal/config"
	"github.com/agenteats/agenteats/internal/database"
	"github.com/agenteats/agenteats/internal/handlers"
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

	// Restaurants
	r.Get("/restaurants", handlers.SearchRestaurants)
	r.Post("/restaurants", handlers.CreateRestaurant)
	r.Get("/restaurants/{restaurantID}", handlers.GetRestaurant)
	r.Put("/restaurants/{restaurantID}", handlers.UpdateRestaurant)

	// Menu
	r.Get("/restaurants/{restaurantID}/menu", handlers.GetMenu)
	r.Post("/restaurants/{restaurantID}/menu/items", handlers.AddMenuItem)

	// Reservations
	r.Get("/restaurants/{restaurantID}/availability", handlers.CheckAvailability)
	r.Post("/restaurants/{restaurantID}/reservations", handlers.MakeReservation)
	r.Get("/restaurants/{restaurantID}/reservations", handlers.ListReservations)
	r.Delete("/reservations/{reservationID}", handlers.CancelReservation)

	// Recommendations
	r.Get("/recommendations", handlers.GetRecommendations)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("🚀 AgentEats API server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
