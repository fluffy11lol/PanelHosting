package repository

import (
	"billing-microservice/internal/models"
	"context"
	"gorm.io/gorm"
)

type BillingRepository struct {
	db *gorm.DB
}

func NewBillingRepository(db *gorm.DB) *BillingRepository {
	return &BillingRepository{db: db}
}

func (r *BillingRepository) GetTariffs(ctx context.Context) (*[]models.Tariff, error) {
	var tariffs *[]models.Tariff

	if err := r.db.WithContext(ctx).Find(&tariffs).Error; err != nil {
		return nil, err
	}
	return tariffs, nil
}
