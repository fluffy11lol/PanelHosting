package internal

import (
	"authentication-microservice/internal/models"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		panic("Failed to migrate database")
	}
}
