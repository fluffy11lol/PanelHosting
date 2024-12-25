package grpc

import (
	"billing-microservice/internal/models"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/emptypb"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) GetTariffs(ctx context.Context) (*[]models.Tariff, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]models.Tariff), args.Error(1)
}

func TestGetTariffs_Error(t *testing.T) {
	mockService := new(MockService)
	billingService := NewBillingService(mockService)

	mockService.On("GetTariffs", mock.Anything).Return(nil, assert.AnError)

	resp, err := billingService.GetTariffs(context.Background(), &emptypb.Empty{})

	assert.Error(t, err)
	assert.Nil(t, resp)
	mockService.AssertExpectations(t)
}
