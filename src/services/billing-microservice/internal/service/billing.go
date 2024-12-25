package service

import (
	"billing-microservice/internal/models"
	"context"
)

type BillingRepo interface {
	GetTariffs(ctx context.Context) (*[]models.Tariff, error)
}

type BillingService struct {
	Repo BillingRepo
}

func NewBillingService(repo BillingRepo) *BillingService {
	return &BillingService{repo}
}

func (s *BillingService) GetTariffs(ctx context.Context) (*[]models.Tariff, error) {
	return s.Repo.GetTariffs(ctx)
}
