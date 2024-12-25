package grpc

import (
	"billing-microservice/internal/models"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"

	client "billing-microservice/pkg/api/billing/protobuf"
)

type Service interface {
	GetTariffs(ctx context.Context) (*[]models.Tariff, error)
}

type BillingService struct {
	client.UnimplementedUserServiceServer
	service Service
}

func NewBillingService(srv Service) *BillingService {
	return &BillingService{service: srv}
}
func (s *BillingService) GetTariffs(ctx context.Context, req *emptypb.Empty) (*client.GetTariffsResponse, error) {

	tariffs, err := s.service.GetTariffs(ctx)
	if err != nil {
		return nil, err
	}

	var grpcTariffs []*client.Tariff
	for _, t := range *tariffs {
		grpcTariffs = append(grpcTariffs, &client.Tariff{
			ID:    t.ID,
			SSD:   int64(t.SSD),
			CPU:   int64(t.CPU),
			RAM:   int64(t.RAM),
			Price: int64(t.Price),
		})
	}

	return &client.GetTariffsResponse{Tariffs: grpcTariffs}, nil
}
