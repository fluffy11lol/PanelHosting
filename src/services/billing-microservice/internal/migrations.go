package internal

import (
	"billing-microservice/internal/models"
	logger2 "billing-microservice/pkg/logger"
	"context"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	err := db.AutoMigrate(&models.Tariff{})
	if err != nil {
		panic("Failed to migrate database: " + err.Error())
	}
}

func InsertData(ctx context.Context, db *gorm.DB) {
	log := logger2.GetLoggerFromCtx(ctx)

	tariffs := []models.Tariff{
		{ID: "1", Name: "Host-0", SSD: 1, CPU: 1, RAM: 4, Price: 500},
		{ID: "2", Name: "Host-1", SSD: 5, CPU: 2, RAM: 8, Price: 1000},
		{ID: "3", Name: "Host-2", SSD: 10, CPU: 3, RAM: 16, Price: 1500},
	}

	for _, tariff := range tariffs {
		var existing models.Tariff

		err := db.First(&existing, "id = ?", tariff.ID).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				if createErr := db.Create(&tariff).Error; createErr != nil {
					log.Error(ctx, "Failed to insert tariff", zap.Error(createErr))
				} else {
					log.Info(ctx, "Tariff created", zap.String("tariffID", tariff.ID))
				}
			} else {
				log.Error(ctx, "Error querying tariff", zap.Error(err))
			}
		} else {
			log.Info(ctx, "Tariff already exists", zap.String("tariffID", tariff.ID))
		}
	}
}
