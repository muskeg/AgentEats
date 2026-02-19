package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/agenteats/agenteats/internal/database"
	"github.com/agenteats/agenteats/internal/dto"
	authmw "github.com/agenteats/agenteats/internal/middleware"
	"github.com/agenteats/agenteats/internal/services"
)

// --- Helpers ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, dto.ErrorOut{Error: msg})
}

func parseCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			result = append(result, t)
		}
	}
	return result
}

// --- Health ---

func Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, dto.HealthOut{
		Status:  "ok",
		Version: "0.1.0",
		Service: "AgentEats",
	})
}

// --- Restaurants ---

func SearchRestaurants(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	city := r.URL.Query().Get("city")
	cuisine := r.URL.Query().Get("cuisine")
	priceRange := r.URL.Query().Get("price_range")
	features := parseCSV(r.URL.Query().Get("features"))

	limit := 20
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 100 {
		limit = l
	}
	offset := 0
	if o, err := strconv.Atoi(r.URL.Query().Get("offset")); err == nil && o >= 0 {
		offset = o
	}

	results := services.ListRestaurants(database.DB, q, city, cuisine, priceRange, features, limit, offset)
	writeJSON(w, http.StatusOK, results)
}

func GetRestaurant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantID")
	result, err := services.GetRestaurant(database.DB, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Restaurant not found")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func CreateRestaurant(w http.ResponseWriter, r *http.Request) {
	var in dto.RestaurantIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := services.CreateRestaurant(database.DB, in)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create restaurant")
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantID")
	var in dto.RestaurantIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := services.UpdateRestaurant(database.DB, id, in)
	if err != nil {
		writeError(w, http.StatusNotFound, "Restaurant not found")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// --- Menu ---

func GetMenu(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantID")
	result, err := services.GetMenu(database.DB, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Restaurant not found")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func AddMenuItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantID")
	var in dto.MenuItemIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := services.AddMenuItem(database.DB, id, in)
	if err != nil {
		writeError(w, http.StatusNotFound, "Restaurant not found")
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

// --- Reservations ---

func CheckAvailability(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantID")
	date := r.URL.Query().Get("date")
	if date == "" {
		writeError(w, http.StatusBadRequest, "date parameter is required")
		return
	}
	partySize := 2
	if ps, err := strconv.Atoi(r.URL.Query().Get("party_size")); err == nil && ps > 0 {
		partySize = ps
	}

	result, err := services.CheckAvailability(database.DB, id, date, partySize)
	if err != nil {
		writeError(w, http.StatusNotFound, "Restaurant not found")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func MakeReservation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantID")
	var in dto.ReservationIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := services.MakeReservation(database.DB, id, in)
	if err != nil {
		writeError(w, http.StatusNotFound, "Restaurant not found")
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func ListReservations(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "restaurantID")
	date := r.URL.Query().Get("date")
	results := services.ListReservations(database.DB, id, date)
	writeJSON(w, http.StatusOK, results)
}

func CancelReservation(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "reservationID")
	result, err := services.CancelReservation(database.DB, id)
	if err != nil {
		writeError(w, http.StatusNotFound, "Reservation not found")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// --- Recommendations ---

func GetRecommendations(w http.ResponseWriter, r *http.Request) {
	cuisine := r.URL.Query().Get("cuisine")
	city := r.URL.Query().Get("city")
	priceRange := r.URL.Query().Get("price_range")
	features := parseCSV(r.URL.Query().Get("features"))
	dietaryNeeds := parseCSV(r.URL.Query().Get("dietary_needs"))
	occasion := r.URL.Query().Get("occasion")

	limit := 5
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 && l <= 20 {
		limit = l
	}

	results := services.GetRecommendations(database.DB, cuisine, city, priceRange, features, dietaryNeeds, occasion, limit)
	writeJSON(w, http.StatusOK, results)
}

// --- Owner Registration ---

func RegisterOwner(w http.ResponseWriter, r *http.Request) {
	var in dto.RegisterOwnerIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if in.Name == "" || in.Email == "" {
		writeError(w, http.StatusBadRequest, "name and email are required")
		return
	}
	result, err := services.RegisterOwner(database.DB, in)
	if err != nil {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func RotateKey(w http.ResponseWriter, r *http.Request) {
	owner := authmw.OwnerFromContext(r.Context())
	if owner == nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	newKey, err := services.RotateAPIKey(database.DB, owner.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to rotate key")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"api_key": newKey})
}

// --- Authenticated Restaurant Management ---

func CreateOwnedRestaurant(w http.ResponseWriter, r *http.Request) {
	owner := authmw.OwnerFromContext(r.Context())
	if owner == nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	var in dto.RestaurantIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := services.CreateRestaurantForOwner(database.DB, owner.ID, in)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create restaurant")
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func UpdateOwnedRestaurant(w http.ResponseWriter, r *http.Request) {
	owner := authmw.OwnerFromContext(r.Context())
	if owner == nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	id := chi.URLParam(r, "restaurantID")
	if !services.RestaurantBelongsToOwner(database.DB, id, owner.ID) {
		writeError(w, http.StatusForbidden, "you do not own this restaurant")
		return
	}
	var in dto.RestaurantIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := services.UpdateRestaurant(database.DB, id, in)
	if err != nil {
		writeError(w, http.StatusNotFound, "Restaurant not found")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func AddOwnedMenuItem(w http.ResponseWriter, r *http.Request) {
	owner := authmw.OwnerFromContext(r.Context())
	if owner == nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	id := chi.URLParam(r, "restaurantID")
	if !services.RestaurantBelongsToOwner(database.DB, id, owner.ID) {
		writeError(w, http.StatusForbidden, "you do not own this restaurant")
		return
	}
	var in dto.MenuItemIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	result, err := services.AddMenuItem(database.DB, id, in)
	if err != nil {
		writeError(w, http.StatusNotFound, "Restaurant not found")
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func BulkImportMenu(w http.ResponseWriter, r *http.Request) {
	owner := authmw.OwnerFromContext(r.Context())
	if owner == nil {
		writeError(w, http.StatusUnauthorized, "not authenticated")
		return
	}
	id := chi.URLParam(r, "restaurantID")
	if !services.RestaurantBelongsToOwner(database.DB, id, owner.ID) {
		writeError(w, http.StatusForbidden, "you do not own this restaurant")
		return
	}
	var in dto.BulkMenuImportIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	if len(in.Items) == 0 {
		writeError(w, http.StatusBadRequest, "items array is required and must not be empty")
		return
	}
	result, err := services.BulkImportMenu(database.DB, id, in)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}
