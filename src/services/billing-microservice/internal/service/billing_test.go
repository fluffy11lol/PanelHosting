package service

import (
	"billing-microservice/internal/models"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBillingRepo struct {
	mock.Mock
}

func (m *MockBillingRepo) GetTariffs(ctx context.Context) (*[]models.Tariff, error) {
	args := m.Called(ctx)
	if args.Get(0) != nil {
		return args.Get(0).(*[]models.Tariff), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestBillingService_GetTariffs_Success(t *testing.T) {

	mockRepo := new(MockBillingRepo)
	service := NewBillingService(mockRepo)
	ctx := context.Background()

	expectedTariffs := &[]models.Tariff{
		{ID: "1", Name: "Basic", SSD: 256, CPU: 2, RAM: 4, Price: 100},
		{ID: "2", Name: "Premium", SSD: 512, CPU: 4, RAM: 8, Price: 200},
	}

	mockRepo.On("GetTariffs", ctx).Return(expectedTariffs, nil)

	result, err := service.GetTariffs(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(*expectedTariffs), len(*result))
	assert.Equal(t, (*expectedTariffs)[0].Name, (*result)[0].Name)

	mockRepo.AssertExpectations(t)
}

func TestBillingService_GetTariffs_Error(t *testing.T) {

	mockRepo := new(MockBillingRepo)
	service := NewBillingService(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetTariffs", ctx).Return(nil, errors.New("repo error"))

	result, err := service.GetTariffs(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "repo error")

	mockRepo.AssertExpectations(t)
}
