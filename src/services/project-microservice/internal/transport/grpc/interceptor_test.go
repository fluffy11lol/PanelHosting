package grpc

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func TestContextWithLogger_LogsRequestStart(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Info", mock.Anything, "request started", mock.Anything).Return()

	interceptor := ContextWithLogger(mockLogger)
	ctx := context.Background()
	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test/method"}, func(c context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestAuthInterceptor_ValidToken(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Info", mock.Anything, "user authenticated", mock.Anything).Return()

	token := generateValidJWT(t, "test_user")
	md := metadata.New(map[string]string{"authorization": "Bearer " + token})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	interceptor := AuthInterceptor(mockLogger)
	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test/method"}, func(c context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.NoError(t, err)
	mockLogger.AssertExpectations(t)
}

func TestAuthInterceptor_EmptyAuthorizationHeader(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Error", mock.Anything, "missing token", mock.Anything).Return()

	md := metadata.New(map[string]string{"authorization": ""})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	interceptor := AuthInterceptor(mockLogger)
	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test/method"}, func(c context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing token")
	mockLogger.AssertExpectations(t)
}

func TestAuthInterceptor_MissingMetadata(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Error", mock.Anything, "missing metadata", mock.Anything).Return()

	ctx := context.Background()
	interceptor := AuthInterceptor(mockLogger)
	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test/method"}, func(c context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing metadata")
	mockLogger.AssertExpectations(t)
}

func TestAuthInterceptor_InvalidAuthorizationHeader(t *testing.T) {
	mockLogger := new(MockLogger)
	mockLogger.On("Error", mock.Anything, "invalid token", mock.Anything).Return()

	md := metadata.New(map[string]string{"authorization": "Bearer invalid_token"})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	interceptor := AuthInterceptor(mockLogger)
	_, err := interceptor(ctx, nil, &grpc.UnaryServerInfo{FullMethod: "/test/method"}, func(c context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
	mockLogger.AssertExpectations(t)
}

func generateValidJWT(t *testing.T, username string) string {
	claims := JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	assert.NoError(t, err)
	return tokenString
}
