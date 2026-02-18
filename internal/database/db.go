package database

import (
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/agenteats/agenteats/internal/config"
	"github.com/agenteats/agenteats/internal/models"
)

// DB is the global database connection.
var DB *gorm.DB

// Init opens the database and runs migrations.
// It auto-detects the driver from DATABASE_URL:
//   - Starts with "postgres://" or "postgresql://" → Postgres
//   - Anything else → SQLite file path
func Init(cfg *config.Config) {
	logLevel := logger.Silent
	if cfg.Debug {
		logLevel = logger.Info
	}

	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	var err error
	dsn := cfg.DatabaseURL

	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		DB, err = gorm.Open(postgres.Open(dsn), gormCfg)
		log.Println("Using PostgreSQL database")
	} else {
		DB, err = gorm.Open(sqlite.Open(dsn), gormCfg)
		log.Println("Using SQLite database:", dsn)
	}
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
