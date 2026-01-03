package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/vondr/identity-go/internal/infrastructure/database/models"
)

var DB *gorm.DB

func InitDB(databaseURL string) error {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB = db
	log.Println("Database connection established")

	return nil
}

func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.Organization{},
		&models.OrganizationMember{},
		&models.App{},
		&models.Session{},
		&models.UserGroup{},
		&models.UserGroupMember{},
		&models.AppAllowedCountry{},
		&models.Invitation{},
	)
}

func GetDB() *gorm.DB {
	return DB
}
