package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/agenteats/agenteats/internal/config"
	"github.com/agenteats/agenteats/internal/models"
)

// DB is the global database connection.
var DB *gorm.DB

// Init opens the database and runs migrations.
func Init(cfg *config.Config) {
	logLevel := logger.Silent
	if cfg.Debug {
		logLevel = logger.Info
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Auto-migrate all models
	if err := DB.AutoMigrate(
		&models.Restaurant{},
		&models.OperatingHours{},
		&models.MenuItem{},
		&models.Reservation{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	log.Println("Database initialized")
}
