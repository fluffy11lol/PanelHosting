package grpc

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"authentication-microservice/internal/models"
	client "authentication-microservice/pkg/api/authentication/protobuf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) RegisterUser(ctx context.Context, user models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockService) LoginUser(ctx context.Context, user models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestNewUserService(t *testing.T) {
	mockSrv := new(MockService)
	userService := NewUserService(mockSrv)
	assert.NotNil(t, userService)
	assert.Equal(t, mockSrv, userService.service)
}

func TestUserService_RegisterUser(t *testing.T) {
	tests := []struct {
		name          string
		request       *client.RegisterUserRequest
		setupMock     func(*MockService)
		expectedResp  *client.RegisterUserResponse
		expectedError error
	}{
		{
			name: "successful registration",
			request: &client.RegisterUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockService) {
				m.On("RegisterUser", mock.Anything, models.User{
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password123",
				}).Return(&models.User{
					Username: "testuser",
					Email:    "test@example.com",
				}, nil)
			},
			expectedResp:  &client.RegisterUserResponse{Status: true},
			expectedError: nil,
		},
		{
			name: "registration error",
			request: &client.RegisterUserRequest{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password123",
			},
			setupMock: func(m *MockService) {
				m.On("RegisterUser", mock.Anything, models.User{
					Username: "testuser",
					Email:    "test@example.com",
					Password: "password123",
				}).Return(nil, errors.New("registration failed"))
			},
			expectedResp:  &client.RegisterUserResponse{Status: false},
			expectedError: errors.New("registration failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSrv := new(MockService)
			tt.setupMock(mockSrv)
			userService := NewUserService(mockSrv)

			resp, err := userService.RegisterUser(context.Background(), tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedResp, resp)
			mockSrv.AssertExpectations(t)
		})
	}
}

func TestUserService_LoginUser(t *testing.T) {
	tests := []struct {
		name          string
		request       *client.LoginUserRequest
		setupMock     func(*MockService)
		expectedResp  *client.LoginUserResponse
		expectedError error
	}{
		{
			name: "login error",
			request: &client.LoginUserRequest{
				Username: "testuser",
				Password: "wrong-password",
			},
			setupMock: func(m *MockService) {
				m.On("LoginUser", mock.Anything, models.User{
					Username: "testuser",
					Password: "wrong-password",
				}).Return(nil, errors.New("invalid credentials"))
			},
			expectedResp:  nil,
			expectedError: errors.New("invalid credentials"),
		},
	}
	originalSendHeader := grpcSendHeader
	defer func() { grpcSendHeader = originalSendHeader }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSrv := new(MockService)
			tt.setupMock(mockSrv)
			userService := NewUserService(mockSrv)
			var md metadata.MD
			grpcSendHeader = func(ctx context.Context, m metadata.MD) error {
				md = m
				if tt.name == "send header error" {
					return errors.New("failed to send header")
				}
				return nil
			}
			resp, err := userService.LoginUser(context.Background(), tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if !errors.Is(tt.expectedError, errors.New("failed to send header")) {
					assert.Equal(t, tt.expectedError.Error(), err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, md.Get("Set-Cookie")[0], fmt.Sprintf("token=%s; HttpOnly; Path=/; Secure", "test-token"))
			}
			assert.Equal(t, tt.expectedResp, resp)
			mockSrv.AssertExpectations(t)

		})
	}
}

func TestUserService_LogoutUser(t *testing.T) {
	tests := []struct {
		name          string
		request       *client.LogoutUserRequest
		expectedResp  *client.LogoutUserResponse
		expectedError error
	}{
		{
			name:          "logout error",
			request:       &client.LogoutUserRequest{},
			expectedResp:  nil,
			expectedError: errors.New("failed to send header"),
		},
	}
	originalSendHeader := grpcSendHeader
	defer func() { grpcSendHeader = originalSendHeader }()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userService := NewUserService(new(MockService))
			var md metadata.MD

			grpcSendHeader = func(ctx context.Context, m metadata.MD) error {
				md = m
				if tt.name == "logout error" {
					return errors.New("failed to send header")
				}
				return nil
			}

			resp, err := userService.LogoutUser(context.Background(), tt.request)

			if tt.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, md.Get("Set-Cookie")[0], "token=; HttpOnly; Path=/; Secure; Max-Age=0")
			}

			assert.Equal(t, tt.expectedResp, resp)
		})
	}
}

var grpcSendHeader = func(ctx context.Context, m metadata.MD) error {
	return grpc.SendHeader(ctx, m)
}
