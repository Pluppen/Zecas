// internal/database/database.go
package database

import (
	"backend/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a connection to the database
func Connect(dsn string) (*gorm.DB, error) {
	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	return gorm.Open(postgres.Open(dsn), config)
}

// Migrate performs database migrations
func Migrate(db *gorm.DB) error {

	err := enableUUIDExtension(db)
	if err != nil {
		return fmt.Errorf("failed to enable UUID extension: %w", err)
	}

	return db.AutoMigrate(
		&models.Project{},
		&models.Target{},
		&models.ScanConfig{},
		&models.Scan{},
		&models.Application{},
		&models.Finding{},
		&models.Report{},
		&models.ScanTask{},
		&models.Session{},
		&models.Account{},
		&models.VerificationToken{},
		&models.User{},
		&models.TargetRelation{},
		&models.Service{},
	)
}

func enableUUIDExtension(db *gorm.DB) error {
	// Check if extension is already enabled
	var count int64
	err := db.Raw(`
		SELECT COUNT(*) FROM pg_extension WHERE extname = 'uuid-ossp'
	`).Count(&count).Error

	if err != nil {
		return fmt.Errorf("failed to check if UUID extension is enabled: %w", err)
	}

	// Enable extension if not already enabled
	if count == 0 {
		log.Println("Enabling uuid-ossp extension...")
		err = db.Exec(`CREATE EXTENSION "uuid-ossp"`).Error
		if err != nil {
			return fmt.Errorf("failed to enable UUID extension: %w", err)
		}
	}

	return nil
}
