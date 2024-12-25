package internal

import (
	"project-microservice/internal/models"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&models.NetworkingInfo{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
	err = db.AutoMigrate(&models.Project{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
}
