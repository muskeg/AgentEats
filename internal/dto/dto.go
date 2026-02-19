package dto

// --- Request DTOs ---

// RestaurantIn is the payload for creating/updating a restaurant.
type RestaurantIn struct {
	Name        string           `json:"name"`
	Description string           `json:"description,omitempty"`
	Cuisines    []string         `json:"cuisines"`
	PriceRange  string           `json:"price_range"`
	Address     string           `json:"address"`
	City        string           `json:"city"`
	State       string           `json:"state,omitempty"`
	ZipCode     string           `json:"zip_code,omitempty"`
	Country     string           `json:"country"`
	Latitude    *float64         `json:"latitude,omitempty"`
	Longitude   *float64         `json:"longitude,omitempty"`
	Phone       string           `json:"phone,omitempty"`
	Email       string           `json:"email,omitempty"`
	Website     string           `json:"website,omitempty"`
	Features    []string         `json:"features"`
	TotalSeats  int              `json:"total_seats"`
	Hours       []OperatingHoursIn `json:"hours"`
}

type OperatingHoursIn struct {
	Day       string `json:"day"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
	IsClosed  bool   `json:"is_closed"`
}

type MenuItemIn struct {
	Category      string   `json:"category"`
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	Price         float64  `json:"price"`
	Currency      string   `json:"currency"`
	DietaryLabels []string `json:"dietary_labels"`
	IsAvailable   bool     `json:"is_available"`
	IsPopular     bool     `json:"is_popular"`
	ImageURL      string   `json:"image_url,omitempty"`
	Calories      *int     `json:"calories,omitempty"`
}

type ReservationIn struct {
	CustomerName    string `json:"customer_name"`
	CustomerEmail   string `json:"customer_email,omitempty"`
	CustomerPhone   string `json:"customer_phone,omitempty"`
	PartySize       int    `json:"party_size"`
	Date            string `json:"date"`
	Time            string `json:"time"`
	SpecialRequests string `json:"special_requests,omitempty"`
}

// --- Response DTOs ---

type RestaurantSummary struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Cuisines    []string `json:"cuisines"`
	PriceRange  string   `json:"price_range"`
	City        string   `json:"city"`
	Rating      *float64 `json:"rating"`
	ReviewCount int      `json:"review_count"`
	Address     string   `json:"address"`
	Features    []string `json:"features"`
}

type OperatingHoursOut struct {
	Day       string `json:"day"`
	OpenTime  string `json:"open_time"`
	CloseTime string `json:"close_time"`
	IsClosed  bool   `json:"is_closed"`
}

type RestaurantDetail struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description,omitempty"`
	Cuisines    []string            `json:"cuisines"`
	PriceRange  string              `json:"price_range"`
	Address     string              `json:"address"`
	City        string              `json:"city"`
	State       string              `json:"state,omitempty"`
	ZipCode     string              `json:"zip_code,omitempty"`
	Country     string              `json:"country"`
	Latitude    *float64            `json:"latitude,omitempty"`
	Longitude   *float64            `json:"longitude,omitempty"`
	Phone       string              `json:"phone,omitempty"`
	Email       string              `json:"email,omitempty"`
	Website     string              `json:"website,omitempty"`
	Features    []string            `json:"features"`
	TotalSeats  int                 `json:"total_seats"`
	Rating      *float64            `json:"rating"`
	ReviewCount int                 `json:"review_count"`
	IsActive    bool                `json:"is_active"`
	Hours       []OperatingHoursOut `json:"hours"`
}

type MenuItemOut struct {
	ID            string   `json:"id"`
	Category      string   `json:"category"`
	Name          string   `json:"name"`
	Description   string   `json:"description,omitempty"`
	Price         float64  `json:"price"`
	Currency      string   `json:"currency"`
	DietaryLabels []string `json:"dietary_labels"`
	IsAvailable   bool     `json:"is_available"`
	IsPopular     bool     `json:"is_popular"`
	ImageURL      string   `json:"image_url,omitempty"`
	Calories      *int     `json:"calories,omitempty"`
}

type MenuOut struct {
	RestaurantID   string                  `json:"restaurant_id"`
	RestaurantName string                  `json:"restaurant_name"`
	Currency       string                  `json:"currency"`
	Categories     map[string][]MenuItemOut `json:"categories"`
}

type ReservationOut struct {
	ID              string `json:"id"`
	RestaurantID    string `json:"restaurant_id"`
	RestaurantName  string `json:"restaurant_name,omitempty"`
	CustomerName    string `json:"customer_name"`
	CustomerEmail   string `json:"customer_email,omitempty"`
	CustomerPhone   string `json:"customer_phone,omitempty"`
	PartySize       int    `json:"party_size"`
	Date            string `json:"date"`
	Time            string `json:"time"`
	Status          string `json:"status"`
	SpecialRequests string `json:"special_requests,omitempty"`
	CreatedAt       string `json:"created_at"`
}

type AvailabilityOut struct {
	RestaurantID   string   `json:"restaurant_id"`
	RestaurantName string   `json:"restaurant_name"`
	Date           string   `json:"date"`
	AvailableTimes []string `json:"available_times"`
	MaxPartySize   int      `json:"max_party_size"`
}

type RecommendationOut struct {
	Restaurant     RestaurantSummary `json:"restaurant"`
	MatchReasons   []string          `json:"match_reasons"`
	RelevanceScore float64           `json:"relevance_score"`
}

type HealthOut struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Service string `json:"service"`
}

// ErrorOut is a standard error response.
type ErrorOut struct {
	Error string `json:"error"`
}

// --- Owner DTOs ---

// RegisterOwnerIn is the payload for owner registration.
type RegisterOwnerIn struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// RegisterOwnerOut is the response after registration.
// The API key is only returned once — store it securely.
type RegisterOwnerOut struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	APIKey string `json:"api_key"` // only returned on creation
}

// OwnerOut is a safe owner representation (no API key).
type OwnerOut struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// --- Bulk Import DTOs ---

// BulkMenuImportIn is the payload for bulk menu import.
type BulkMenuImportIn struct {
	// Strategy: "replace" deletes all existing items first; "merge" adds/updates.
	Strategy string       `json:"strategy"` // "replace" (default) or "merge"
	Items    []MenuItemIn `json:"items"`
}

// BulkMenuImportOut is the response after a bulk import.
type BulkMenuImportOut struct {
	RestaurantID string `json:"restaurant_id"`
	Imported     int    `json:"imported"`
	Strategy     string `json:"strategy"`
}
