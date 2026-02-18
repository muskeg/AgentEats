package models

import (
	"time"

	"github.com/google/uuid"
)

// --- Enums as string types ---

type PriceRange string

const (
	PriceBudget     PriceRange = "$"
	PriceModerate   PriceRange = "$$"
	PriceUpscale    PriceRange = "$$$"
	PriceFineDining PriceRange = "$$$$"
)

type ReservationStatus string

const (
	StatusConfirmed ReservationStatus = "confirmed"
	StatusCancelled ReservationStatus = "cancelled"
	StatusCompleted ReservationStatus = "completed"
	StatusNoShow    ReservationStatus = "no_show"
)

// --- Models ---

// Restaurant represents a restaurant listing.
type Restaurant struct {
	ID          string     `gorm:"primaryKey;size:36" json:"id"`
	Name        string     `gorm:"size:200;not null;index" json:"name"`
	Description string     `gorm:"type:text" json:"description,omitempty"`
	Cuisines    string     `gorm:"size:500" json:"cuisines"` // comma-separated
	PriceRange  PriceRange `gorm:"size:10;not null;default:'$$'" json:"price_range"`
	Address     string     `gorm:"size:500;not null" json:"address"`
	City        string     `gorm:"size:100;not null;index" json:"city"`
	State       string     `gorm:"size:100" json:"state,omitempty"`
	ZipCode     string     `gorm:"size:20" json:"zip_code,omitempty"`
	Country     string     `gorm:"size:100;not null;default:'US'" json:"country"`
	Latitude    *float64   `json:"latitude,omitempty"`
	Longitude   *float64   `json:"longitude,omitempty"`
	Phone       string     `gorm:"size:30" json:"phone,omitempty"`
	Email       string     `gorm:"size:200" json:"email,omitempty"`
	Website     string     `gorm:"size:500" json:"website,omitempty"`
	Features    string     `gorm:"size:500" json:"features"` // comma-separated
	TotalSeats  int        `gorm:"not null;default:50" json:"total_seats"`
	Rating      *float64   `json:"rating,omitempty"`
	ReviewCount int        `gorm:"not null;default:0" json:"review_count"`
	IsActive    bool       `gorm:"not null;default:true" json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	Hours      []OperatingHours `gorm:"foreignKey:RestaurantID;constraint:OnDelete:CASCADE" json:"hours,omitempty"`
	MenuItems  []MenuItem       `gorm:"foreignKey:RestaurantID;constraint:OnDelete:CASCADE" json:"menu_items,omitempty"`
	Reservations []Reservation  `gorm:"foreignKey:RestaurantID;constraint:OnDelete:CASCADE" json:"reservations,omitempty"`
}

// OperatingHours represents the hours for one day of the week.
type OperatingHours struct {
	ID           uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	RestaurantID string `gorm:"size:36;not null;index" json:"restaurant_id"`
	Day          string `gorm:"size:10;not null" json:"day"`
	OpenTime     string `gorm:"size:5;not null" json:"open_time"`
	CloseTime    string `gorm:"size:5;not null" json:"close_time"`
	IsClosed     bool   `gorm:"not null;default:false" json:"is_closed"`
}

// MenuItem represents a single dish on a restaurant's menu.
type MenuItem struct {
	ID            string   `gorm:"primaryKey;size:36" json:"id"`
	RestaurantID  string   `gorm:"size:36;not null;index" json:"restaurant_id"`
	Category      string   `gorm:"size:100;not null;default:'Main'" json:"category"`
	Name          string   `gorm:"size:200;not null" json:"name"`
	Description   string   `gorm:"type:text" json:"description,omitempty"`
	Price         float64  `gorm:"not null" json:"price"`
	Currency      string   `gorm:"size:3;not null;default:'USD'" json:"currency"`
	DietaryLabels string   `gorm:"size:300" json:"dietary_labels"` // comma-separated
	IsAvailable   bool     `gorm:"not null;default:true" json:"is_available"`
	IsPopular     bool     `gorm:"not null;default:false" json:"is_popular"`
	ImageURL      string   `gorm:"size:500" json:"image_url,omitempty"`
	Calories      *int     `json:"calories,omitempty"`
}

// Reservation represents a table reservation.
type Reservation struct {
	ID              string            `gorm:"primaryKey;size:36" json:"id"`
	RestaurantID    string            `gorm:"size:36;not null;index" json:"restaurant_id"`
	CustomerName    string            `gorm:"size:200;not null" json:"customer_name"`
	CustomerEmail   string            `gorm:"size:200" json:"customer_email,omitempty"`
	CustomerPhone   string            `gorm:"size:30" json:"customer_phone,omitempty"`
	PartySize       int               `gorm:"not null" json:"party_size"`
	Date            string            `gorm:"size:10;not null" json:"date"`  // YYYY-MM-DD
	Time            string            `gorm:"size:5;not null" json:"time"`   // HH:MM
	Status          ReservationStatus `gorm:"size:20;not null;default:'confirmed'" json:"status"`
	SpecialRequests string            `gorm:"type:text" json:"special_requests,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
}

// NewID generates a new UUID string.
func NewID() string {
	return uuid.New().String()
}
