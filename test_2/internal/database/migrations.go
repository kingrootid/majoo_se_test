package database

import (
	"rootwritter/majoo_test_2_api/internal/models"

	"gorm.io/gorm"
)

// MigrateDB runs database migrations
func MigrateDB(db *gorm.DB) error {
	// Auto migrate the schema
	return db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
	)
}