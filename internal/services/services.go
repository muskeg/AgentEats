package services

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"gorm.io/gorm"

	"github.com/agenteats/agenteats/internal/dto"
	"github.com/agenteats/agenteats/internal/models"
)

// --- Helpers ---

func splitCSV(s string) []string {
	if s == "" {
		return []string{}
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

func joinCSV(items []string) string {
	return strings.Join(items, ",")
}

func toSummary(r *models.Restaurant) dto.RestaurantSummary {
	return dto.RestaurantSummary{
		ID:          r.ID,
		Name:        r.Name,
		Cuisines:    splitCSV(r.Cuisines),
		PriceRange:  string(r.PriceRange),
		City:        r.City,
		Rating:      r.Rating,
		ReviewCount: r.ReviewCount,
		Address:     r.Address,
		Features:    splitCSV(r.Features),
	}
}

func toDetail(r *models.Restaurant) dto.RestaurantDetail {
	hours := make([]dto.OperatingHoursOut, len(r.Hours))
	for i, h := range r.Hours {
		hours[i] = dto.OperatingHoursOut{
			Day:       h.Day,
			OpenTime:  h.OpenTime,
			CloseTime: h.CloseTime,
			IsClosed:  h.IsClosed,
		}
	}
	return dto.RestaurantDetail{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Cuisines:    splitCSV(r.Cuisines),
		PriceRange:  string(r.PriceRange),
		Address:     r.Address,
		City:        r.City,
		State:       r.State,
		ZipCode:     r.ZipCode,
		Country:     r.Country,
		Latitude:    r.Latitude,
		Longitude:   r.Longitude,
		Phone:       r.Phone,
		Email:       r.Email,
		Website:     r.Website,
		Features:    splitCSV(r.Features),
		TotalSeats:  r.TotalSeats,
		Rating:      r.Rating,
		ReviewCount: r.ReviewCount,
		IsActive:    r.IsActive,
		Hours:       hours,
	}
}

func toMenuItemOut(m *models.MenuItem) dto.MenuItemOut {
	return dto.MenuItemOut{
		ID:            m.ID,
		Category:      m.Category,
		Name:          m.Name,
		Description:   m.Description,
		Price:         m.Price,
		Currency:      m.Currency,
		DietaryLabels: splitCSV(m.DietaryLabels),
		IsAvailable:   m.IsAvailable,
		IsPopular:     m.IsPopular,
		ImageURL:      m.ImageURL,
		Calories:      m.Calories,
	}
}

func toReservationOut(r *models.Reservation, restaurantName string) dto.ReservationOut {
	return dto.ReservationOut{
		ID:              r.ID,
		RestaurantID:    r.RestaurantID,
		RestaurantName:  restaurantName,
		CustomerName:    r.CustomerName,
		CustomerEmail:   r.CustomerEmail,
		CustomerPhone:   r.CustomerPhone,
		PartySize:       r.PartySize,
		Date:            r.Date,
		Time:            r.Time,
		Status:          string(r.Status),
		SpecialRequests: r.SpecialRequests,
		CreatedAt:       r.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// --- Restaurant CRUD ---

// ListRestaurants searches and filters restaurants.
func ListRestaurants(db *gorm.DB, q, city, cuisine, priceRange string, features []string, limit, offset int) []dto.RestaurantSummary {
	query := db.Where("is_active = ?", true)

	if city != "" {
		query = query.Where("city LIKE ?", "%"+city+"%")
	}
	if cuisine != "" {
		query = query.Where("cuisines LIKE ?", "%"+cuisine+"%")
	}
	if priceRange != "" {
		query = query.Where("price_range = ?", priceRange)
	}
	for _, f := range features {
		query = query.Where("features LIKE ?", "%"+f+"%")
	}
	if q != "" {
		query = query.Where("name LIKE ? OR description LIKE ? OR cuisines LIKE ?",
			"%"+q+"%", "%"+q+"%", "%"+q+"%")
	}

	var restaurants []models.Restaurant
	query.Order("CASE WHEN rating IS NULL THEN 1 ELSE 0 END, rating DESC").
		Offset(offset).Limit(limit).Find(&restaurants)

	results := make([]dto.RestaurantSummary, len(restaurants))
	for i := range restaurants {
		results[i] = toSummary(&restaurants[i])
	}
	return results
}

// GetRestaurant returns full restaurant details.
func GetRestaurant(db *gorm.DB, id string) (*dto.RestaurantDetail, error) {
	var r models.Restaurant
	if err := db.Preload("Hours").First(&r, "id = ?", id).Error; err != nil {
		return nil, err
	}
	detail := toDetail(&r)
	return &detail, nil
}

// CreateRestaurant registers a new restaurant.
func CreateRestaurant(db *gorm.DB, in dto.RestaurantIn) (*dto.RestaurantDetail, error) {
	r := models.Restaurant{
		ID:          models.NewID(),
		Name:        in.Name,
		Description: in.Description,
		Cuisines:    joinCSV(in.Cuisines),
		PriceRange:  models.PriceRange(in.PriceRange),
		Address:     in.Address,
		City:        in.City,
		State:       in.State,
		ZipCode:     in.ZipCode,
		Country:     in.Country,
		Latitude:    in.Latitude,
		Longitude:   in.Longitude,
		Phone:       in.Phone,
		Email:       in.Email,
		Website:     in.Website,
		Features:    joinCSV(in.Features),
		TotalSeats:  in.TotalSeats,
		IsActive:    true,
	}

	if r.Country == "" {
		r.Country = "US"
	}
	if r.TotalSeats == 0 {
		r.TotalSeats = 50
	}
	if r.PriceRange == "" {
		r.PriceRange = models.PriceModerate
	}

	hours := make([]models.OperatingHours, len(in.Hours))
	for i, h := range in.Hours {
		hours[i] = models.OperatingHours{
			RestaurantID: r.ID,
			Day:          strings.ToLower(h.Day),
			OpenTime:     h.OpenTime,
			CloseTime:    h.CloseTime,
			IsClosed:     h.IsClosed,
		}
	}
	r.Hours = hours

	if err := db.Create(&r).Error; err != nil {
		return nil, err
	}

	detail := toDetail(&r)
	return &detail, nil
}

// UpdateRestaurant updates an existing restaurant.
func UpdateRestaurant(db *gorm.DB, id string, in dto.RestaurantIn) (*dto.RestaurantDetail, error) {
	var r models.Restaurant
	if err := db.First(&r, "id = ?", id).Error; err != nil {
		return nil, err
	}

	r.Name = in.Name
	r.Description = in.Description
	r.Cuisines = joinCSV(in.Cuisines)
	r.PriceRange = models.PriceRange(in.PriceRange)
	r.Address = in.Address
	r.City = in.City
	r.State = in.State
	r.ZipCode = in.ZipCode
	r.Country = in.Country
	r.Latitude = in.Latitude
	r.Longitude = in.Longitude
	r.Phone = in.Phone
	r.Email = in.Email
	r.Website = in.Website
	r.Features = joinCSV(in.Features)
	r.TotalSeats = in.TotalSeats

	// Replace hours
	db.Where("restaurant_id = ?", id).Delete(&models.OperatingHours{})
	hours := make([]models.OperatingHours, len(in.Hours))
	for i, h := range in.Hours {
		hours[i] = models.OperatingHours{
			RestaurantID: r.ID,
			Day:          strings.ToLower(h.Day),
			OpenTime:     h.OpenTime,
			CloseTime:    h.CloseTime,
			IsClosed:     h.IsClosed,
		}
	}
	r.Hours = hours

	if err := db.Save(&r).Error; err != nil {
		return nil, err
	}
	// Reload with hours
	return GetRestaurant(db, id)
}

// --- Menu ---

// GetMenu returns the full menu grouped by category.
func GetMenu(db *gorm.DB, restaurantID string) (*dto.MenuOut, error) {
	var r models.Restaurant
	if err := db.First(&r, "id = ?", restaurantID).Error; err != nil {
		return nil, err
	}

	var items []models.MenuItem
	db.Where("restaurant_id = ?", restaurantID).Order("category, name").Find(&items)

	categories := make(map[string][]dto.MenuItemOut)
	currency := "USD"
	for i := range items {
		out := toMenuItemOut(&items[i])
		currency = items[i].Currency
		categories[out.Category] = append(categories[out.Category], out)
	}

	return &dto.MenuOut{
		RestaurantID:   restaurantID,
		RestaurantName: r.Name,
		Currency:       currency,
		Categories:     categories,
	}, nil
}

// AddMenuItem adds a menu item to a restaurant.
func AddMenuItem(db *gorm.DB, restaurantID string, in dto.MenuItemIn) (*dto.MenuItemOut, error) {
	var r models.Restaurant
	if err := db.First(&r, "id = ?", restaurantID).Error; err != nil {
		return nil, err
	}

	item := models.MenuItem{
		ID:            models.NewID(),
		RestaurantID:  restaurantID,
		Category:      in.Category,
		Name:          in.Name,
		Description:   in.Description,
		Price:         in.Price,
		Currency:      in.Currency,
		DietaryLabels: joinCSV(in.DietaryLabels),
		IsAvailable:   in.IsAvailable,
		IsPopular:     in.IsPopular,
		ImageURL:      in.ImageURL,
		Calories:      in.Calories,
	}
	if item.Currency == "" {
		item.Currency = "USD"
	}
	if item.Category == "" {
		item.Category = "Main"
	}

	if err := db.Create(&item).Error; err != nil {
		return nil, err
	}

	out := toMenuItemOut(&item)
	return &out, nil
}

// --- Reservations ---

// CheckAvailability returns available time slots for a given date.
func CheckAvailability(db *gorm.DB, restaurantID, date string, partySize int) (*dto.AvailabilityOut, error) {
	var r models.Restaurant
	if err := db.First(&r, "id = ?", restaurantID).Error; err != nil {
		return nil, err
	}

	var existing []models.Reservation
	db.Where("restaurant_id = ? AND date = ? AND status = ?",
		restaurantID, date, models.StatusConfirmed).Find(&existing)

	bookedSeats := make(map[string]int)
	for _, res := range existing {
		bookedSeats[res.Time] += res.PartySize
	}

	var available []string
	for hour := 11; hour <= 21; hour++ {
		for _, minute := range []int{0, 30} {
			slot := fmt.Sprintf("%02d:%02d", hour, minute)
			if bookedSeats[slot]+partySize <= r.TotalSeats {
				available = append(available, slot)
			}
		}
	}

	return &dto.AvailabilityOut{
		RestaurantID:   restaurantID,
		RestaurantName: r.Name,
		Date:           date,
		AvailableTimes: available,
		MaxPartySize:   r.TotalSeats,
	}, nil
}

// MakeReservation creates a reservation.
func MakeReservation(db *gorm.DB, restaurantID string, in dto.ReservationIn) (*dto.ReservationOut, error) {
	var r models.Restaurant
	if err := db.First(&r, "id = ?", restaurantID).Error; err != nil {
		return nil, err
	}

	res := models.Reservation{
		ID:              models.NewID(),
		RestaurantID:    restaurantID,
		CustomerName:    in.CustomerName,
		CustomerEmail:   in.CustomerEmail,
		CustomerPhone:   in.CustomerPhone,
		PartySize:       in.PartySize,
		Date:            in.Date,
		Time:            in.Time,
		Status:          models.StatusConfirmed,
		SpecialRequests: in.SpecialRequests,
	}

	if err := db.Create(&res).Error; err != nil {
		return nil, err
	}

	out := toReservationOut(&res, r.Name)
	return &out, nil
}

// ListReservations returns reservations for a restaurant, optionally filtered by date.
func ListReservations(db *gorm.DB, restaurantID, date string) []dto.ReservationOut {
	var r models.Restaurant
	db.First(&r, "id = ?", restaurantID)

	query := db.Where("restaurant_id = ?", restaurantID)
	if date != "" {
		query = query.Where("date = ?", date)
	}

	var reservations []models.Reservation
	query.Order("date, time").Find(&reservations)

	results := make([]dto.ReservationOut, len(reservations))
	for i := range reservations {
		results[i] = toReservationOut(&reservations[i], r.Name)
	}
	return results
}

// CancelReservation cancels a reservation by ID.
func CancelReservation(db *gorm.DB, reservationID string) (*dto.ReservationOut, error) {
	var res models.Reservation
	if err := db.First(&res, "id = ?", reservationID).Error; err != nil {
		return nil, err
	}
	res.Status = models.StatusCancelled
	db.Save(&res)

	var r models.Restaurant
	db.First(&r, "id = ?", res.RestaurantID)

	out := toReservationOut(&res, r.Name)
	return &out, nil
}

// --- Recommendations ---

type scoredRestaurant struct {
	restaurant *models.Restaurant
	score      float64
	reasons    []string
}

// GetRecommendations generates ranked restaurant recommendations.
func GetRecommendations(db *gorm.DB, cuisine, city, priceRange string, features, dietaryNeeds []string, occasion string, limit int) []dto.RecommendationOut {
	query := db.Where("is_active = ?", true)
	if city != "" {
		query = query.Where("city LIKE ?", "%"+city+"%")
	}

	var candidates []models.Restaurant
	query.Find(&candidates)

	var scored []scoredRestaurant

	for i := range candidates {
		r := &candidates[i]
		score := 0.0
		var reasons []string

		rCuisines := splitCSV(r.Cuisines)
		rFeatures := splitCSV(r.Features)
		pr := string(r.PriceRange)

		// Cuisine match
		if cuisine != "" {
			for _, c := range rCuisines {
				if strings.Contains(strings.ToLower(c), strings.ToLower(cuisine)) {
					score += 0.3
					reasons = append(reasons, fmt.Sprintf("Serves %s cuisine", cuisine))
					break
				}
			}
		}

		// Price match
		if priceRange != "" && pr == priceRange {
			score += 0.15
			reasons = append(reasons, fmt.Sprintf("Matches %s budget", priceRange))
		}

		// Feature matches
		if len(features) > 0 {
			var matched []string
			for _, f := range features {
				for _, rf := range rFeatures {
					if strings.Contains(strings.ToLower(rf), strings.ToLower(f)) {
						matched = append(matched, f)
						break
					}
				}
			}
			if len(matched) > 0 {
				score += 0.15 * float64(len(matched)) / float64(len(features))
				reasons = append(reasons, fmt.Sprintf("Has: %s", strings.Join(matched, ", ")))
			}
		}

		// Dietary needs — check menu items
		if len(dietaryNeeds) > 0 {
			var menuItems []models.MenuItem
			db.Where("restaurant_id = ? AND is_available = ?", r.ID, true).Find(&menuItems)

			dietMatches := make(map[string]bool)
			for _, item := range menuItems {
				for _, label := range splitCSV(item.DietaryLabels) {
					for _, need := range dietaryNeeds {
						if strings.EqualFold(label, need) {
							dietMatches[label] = true
						}
					}
				}
			}
			if len(dietMatches) > 0 {
				score += 0.2 * float64(len(dietMatches)) / float64(len(dietaryNeeds))
				matched := make([]string, 0, len(dietMatches))
				for k := range dietMatches {
					matched = append(matched, k)
				}
				reasons = append(reasons, fmt.Sprintf("Menu has %s options", strings.Join(matched, ", ")))
			}
		}

		// Occasion heuristics
		if occasion != "" {
			occ := strings.ToLower(occasion)
			switch occ {
			case "date_night", "romantic":
				if pr == "$$$" || pr == "$$$$" {
					score += 0.1
					reasons = append(reasons, "Upscale ambiance for a date night")
				}
			case "business", "corporate":
				if pr == "$$$" || pr == "$$$$" {
					score += 0.1
					reasons = append(reasons, "Suitable for business dining")
				}
			case "family":
				if pr == "$" || pr == "$$" {
					score += 0.1
					reasons = append(reasons, "Family-friendly pricing")
				}
			case "casual":
				if pr == "$" || pr == "$$" {
					score += 0.1
					reasons = append(reasons, "Great for a casual meal")
				}
			}
		}

		// Rating boost
		if r.Rating != nil {
			score += 0.1 * (*r.Rating / 5.0)
			if *r.Rating >= 4.5 {
				reasons = append(reasons, fmt.Sprintf("Highly rated (%.1f★)", *r.Rating))
			} else if *r.Rating >= 4.0 {
				reasons = append(reasons, fmt.Sprintf("Well rated (%.1f★)", *r.Rating))
			}
		}

		noFilters := cuisine == "" && priceRange == "" && len(features) == 0 && len(dietaryNeeds) == 0 && occasion == ""
		if score > 0 || noFilters {
			if len(reasons) == 0 {
				reasons = append(reasons, "Popular in the area")
			}
			scored = append(scored, scoredRestaurant{
				restaurant: r,
				score:      math.Min(score, 1.0),
				reasons:    reasons,
			})
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	if limit > len(scored) {
		limit = len(scored)
	}

	results := make([]dto.RecommendationOut, limit)
	for i := 0; i < limit; i++ {
		results[i] = dto.RecommendationOut{
			Restaurant:     toSummary(scored[i].restaurant),
			MatchReasons:   scored[i].reasons,
			RelevanceScore: math.Round(scored[i].score*100) / 100,
		}
	}
	return results
}

// --- Owner Registration ---

// RegisterOwner creates a new owner account and returns the raw API key.
func RegisterOwner(db *gorm.DB, in dto.RegisterOwnerIn) (*dto.RegisterOwnerOut, error) {
	rawKey, keyHash := models.GenerateAPIKey()

	owner := models.Owner{
		ID:         models.NewID(),
		Name:       in.Name,
		Email:      in.Email,
		APIKeyHash: keyHash,
		IsActive:   true,
	}

	if err := db.Create(&owner).Error; err != nil {
		return nil, fmt.Errorf("failed to create owner (email may already exist): %w", err)
	}

	return &dto.RegisterOwnerOut{
		ID:     owner.ID,
		Name:   owner.Name,
		Email:  owner.Email,
		APIKey: rawKey,
	}, nil
}

// RotateAPIKey generates a new API key for an owner, invalidating the old one.
func RotateAPIKey(db *gorm.DB, ownerID string) (string, error) {
	rawKey, keyHash := models.GenerateAPIKey()
	if err := db.Model(&models.Owner{}).Where("id = ?", ownerID).Update("api_key_hash", keyHash).Error; err != nil {
		return "", err
	}
	return rawKey, nil
}

// --- Ownership Helpers ---

// RestaurantBelongsToOwner checks if the owner owns the restaurant.
func RestaurantBelongsToOwner(db *gorm.DB, restaurantID, ownerID string) bool {
	var count int64
	db.Model(&models.Restaurant{}).Where("id = ? AND owner_id = ?", restaurantID, ownerID).Count(&count)
	return count > 0
}

// CreateRestaurantForOwner creates a restaurant assigned to the authenticated owner.
func CreateRestaurantForOwner(db *gorm.DB, ownerID string, in dto.RestaurantIn) (*dto.RestaurantDetail, error) {
	r := models.Restaurant{
		ID:          models.NewID(),
		OwnerID:     ownerID,
		Name:        in.Name,
		Description: in.Description,
		Cuisines:    joinCSV(in.Cuisines),
		PriceRange:  models.PriceRange(in.PriceRange),
		Address:     in.Address,
		City:        in.City,
		State:       in.State,
		ZipCode:     in.ZipCode,
		Country:     in.Country,
		Latitude:    in.Latitude,
		Longitude:   in.Longitude,
		Phone:       in.Phone,
		Email:       in.Email,
		Website:     in.Website,
		Features:    joinCSV(in.Features),
		TotalSeats:  in.TotalSeats,
		IsActive:    true,
	}

	if r.Country == "" {
		r.Country = "US"
	}
	if r.TotalSeats == 0 {
		r.TotalSeats = 50
	}
	if r.PriceRange == "" {
		r.PriceRange = models.PriceModerate
	}

	hours := make([]models.OperatingHours, len(in.Hours))
	for i, h := range in.Hours {
		hours[i] = models.OperatingHours{
			RestaurantID: r.ID,
			Day:          strings.ToLower(h.Day),
			OpenTime:     h.OpenTime,
			CloseTime:    h.CloseTime,
			IsClosed:     h.IsClosed,
		}
	}
	r.Hours = hours

	if err := db.Create(&r).Error; err != nil {
		return nil, err
	}

	detail := toDetail(&r)
	return &detail, nil
}

// --- Bulk Menu Import ---

// BulkImportMenu imports menu items for a restaurant.
// Strategy "replace" deletes all existing items first. "merge" appends.
func BulkImportMenu(db *gorm.DB, restaurantID string, in dto.BulkMenuImportIn) (*dto.BulkMenuImportOut, error) {
	var r models.Restaurant
	if err := db.First(&r, "id = ?", restaurantID).Error; err != nil {
		return nil, fmt.Errorf("restaurant not found")
	}

	strategy := in.Strategy
	if strategy == "" {
		strategy = "replace"
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if strategy == "replace" {
			if err := tx.Where("restaurant_id = ?", restaurantID).Delete(&models.MenuItem{}).Error; err != nil {
				return fmt.Errorf("failed to clear existing menu: %w", err)
			}
		}

		for _, item := range in.Items {
			m := models.MenuItem{
				ID:            models.NewID(),
				RestaurantID:  restaurantID,
				Category:      item.Category,
				Name:          item.Name,
				Description:   item.Description,
				Price:         item.Price,
				Currency:      item.Currency,
				DietaryLabels: joinCSV(item.DietaryLabels),
				IsAvailable:   item.IsAvailable,
				IsPopular:     item.IsPopular,
				ImageURL:      item.ImageURL,
				Calories:      item.Calories,
			}
			if m.Currency == "" {
				m.Currency = "USD"
			}
			if m.Category == "" {
				m.Category = "Main"
			}
			if err := tx.Create(&m).Error; err != nil {
				return fmt.Errorf("failed to import item %q: %w", item.Name, err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &dto.BulkMenuImportOut{
		RestaurantID: restaurantID,
		Imported:     len(in.Items),
		Strategy:     strategy,
	}, nil
}
