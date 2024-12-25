package logger_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"storage-microservice/pkg/logger"
	"testing"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	m.Called(ctx, msg, fields)
}

func (m *MockLogger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	m.Called(ctx, msg, fields)
}

func (m *MockLogger) Sync() error {
	args := m.Called()
	return args.Error(0)
}

func TestLogger_WithLogger(t *testing.T) {
	mockLogger := new(MockLogger)

	ctx := context.Background()
	ctx = logger.WithLogger(ctx, mockLogger)

	retrievedLogger := logger.GetLoggerFromCtx(ctx)
	assert.Equal(t, mockLogger, retrievedLogger)

	mockLogger.On("Info", mock.Anything, "test message", mock.Anything).Return()
	retrievedLogger.Info(ctx, "test message", zap.String("key", "value"))
	mockLogger.AssertExpectations(t)

	mockLogger.On("Error", mock.Anything, "test error", mock.Anything).Return()
	retrievedLogger.Error(ctx, "test error", zap.String("key", "value"))
	mockLogger.AssertExpectations(t)
}

func TestLogger_New(t *testing.T) {
	loggerInstance := logger.New("test-service")

	assert.NotNil(t, loggerInstance)
	assert.Implements(t, (*logger.Logger)(nil), loggerInstance)
}

func TestLogger_Sync(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Sync").Return(nil)

	err := mockLogger.Sync()
	assert.NoError(t, err)

	mockLogger.AssertExpectations(t)
}
