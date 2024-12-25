package authentication

import (
	"bytes"
	"context"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

// MockUserServiceClient is a mock implementation of UserServiceClient
type MockUserServiceClient struct {
	mock.Mock
}

func (m *MockUserServiceClient) RegisterUser(ctx context.Context, in *RegisterUserRequest, opts ...grpc.CallOption) (*RegisterUserResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RegisterUserResponse), args.Error(1)
}

func (m *MockUserServiceClient) LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*LoginUserResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LoginUserResponse), args.Error(1)
}

func (m *MockUserServiceClient) LogoutUser(ctx context.Context, in *LogoutUserRequest, opts ...grpc.CallOption) (*LogoutUserResponse, error) {
	args := m.Called(ctx, in, opts)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LogoutUserResponse), args.Error(1)
}

type MockUserServiceServer struct {
	mock.Mock
	UnimplementedUserServiceServer // Add this line to embed the UnimplementedUserServiceServer
}

// Остальные методы MockUserServiceServer остаются как прежде
func (m *MockUserServiceServer) RegisterUser(ctx context.Context, in *RegisterUserRequest) (*RegisterUserResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*RegisterUserResponse), args.Error(1)
}

func (m *MockUserServiceServer) LoginUser(ctx context.Context, in *LoginUserRequest) (*LoginUserResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LoginUserResponse), args.Error(1)
}
func (m *MockUserServiceServer) LogoutUser(ctx context.Context, in *LogoutUserRequest) (*LogoutUserResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*LogoutUserResponse), args.Error(1)
}

func Test_request_UserService_RegisterUser_0(t *testing.T) {
	tests := []struct {
		name        string
		reqBody     string
		setupMock   func(*MockUserServiceClient)
		expectCode  codes.Code
		expectError string
	}{
		{
			name:    "success",
			reqBody: `{"username": "testuser", "email": "test@example.com", "password": "password"}`,
			setupMock: func(m *MockUserServiceClient) {
				m.On("RegisterUser", mock.Anything, mock.Anything, mock.Anything).Return(&RegisterUserResponse{Status: true}, nil)
			},
			expectCode: codes.OK,
		},
		{
			name:        "invalid request body",
			reqBody:     `invalid json`,
			expectCode:  codes.InvalidArgument,
			expectError: "invalid character 'i' looking for beginning of value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockUserServiceClient)
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}
			req := httptest.NewRequest(http.MethodPost, "/v1/user/register", bytes.NewBufferString(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _, err := request_UserService_RegisterUser_0(context.Background(), &runtime.JSONPb{}, mockClient, req, nil)

			if tt.expectError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectError)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectCode, st.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			mockClient.AssertExpectations(t)

		})
	}
}
func Test_local_request_UserService_RegisterUser_0(t *testing.T) {
	tests := []struct {
		name        string
		reqBody     string
		setupMock   func(*MockUserServiceServer)
		expectCode  codes.Code
		expectError string
	}{
		{
			name:    "success",
			reqBody: `{"username": "testuser", "email": "test@example.com", "password": "password"}`,
			setupMock: func(m *MockUserServiceServer) {
				m.On("RegisterUser", mock.Anything, mock.Anything).Return(&RegisterUserResponse{Status: true}, nil)
			},
			expectCode: codes.OK,
		},
		{
			name:        "invalid request body",
			reqBody:     `invalid json`,
			expectCode:  codes.InvalidArgument,
			expectError: "invalid character 'i' looking for beginning of value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := new(MockUserServiceServer)
			if tt.setupMock != nil {
				tt.setupMock(mockServer)
			}
			req := httptest.NewRequest(http.MethodPost, "/v1/user/register", bytes.NewBufferString(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp, _, err := local_request_UserService_RegisterUser_0(context.Background(), &runtime.JSONPb{}, mockServer, req, nil)

			if tt.expectError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectError)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectCode, st.Code())

			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			mockServer.AssertExpectations(t)
		})
	}
}

func Test_request_UserService_LoginUser_0(t *testing.T) {
	tests := []struct {
		name        string
		reqBody     string
		setupMock   func(*MockUserServiceClient)
		expectCode  codes.Code
		expectError string
	}{
		{
			name:    "success",
			reqBody: `{"username": "testuser", "password": "password"}`,
			setupMock: func(m *MockUserServiceClient) {
				m.On("LoginUser", mock.Anything, mock.Anything, mock.Anything).Return(&LoginUserResponse{Token: "test-token"}, nil)
			},
			expectCode: codes.OK,
		},
		{
			name:        "invalid request body",
			reqBody:     `invalid json`,
			expectCode:  codes.InvalidArgument,
			expectError: "invalid character 'i' looking for beginning of value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockUserServiceClient)
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}
			req := httptest.NewRequest(http.MethodPost, "/v1/user/login", bytes.NewBufferString(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp, _, err := request_UserService_LoginUser_0(context.Background(), &runtime.JSONPb{}, mockClient, req, nil)

			if tt.expectError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectError)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectCode, st.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			mockClient.AssertExpectations(t)

		})
	}
}

func Test_local_request_UserService_LoginUser_0(t *testing.T) {
	tests := []struct {
		name        string
		reqBody     string
		setupMock   func(*MockUserServiceServer)
		expectCode  codes.Code
		expectError string
	}{
		{
			name:    "success",
			reqBody: `{"username": "testuser", "password": "password"}`,
			setupMock: func(m *MockUserServiceServer) {
				m.On("LoginUser", mock.Anything, mock.Anything).Return(&LoginUserResponse{Token: "test-token"}, nil)
			},
			expectCode: codes.OK,
		},
		{
			name:        "invalid request body",
			reqBody:     `invalid json`,
			expectCode:  codes.InvalidArgument,
			expectError: "invalid character 'i' looking for beginning of value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := new(MockUserServiceServer)
			if tt.setupMock != nil {
				tt.setupMock(mockServer)
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/user/login", bytes.NewBufferString(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")
			resp, _, err := local_request_UserService_LoginUser_0(context.Background(), &runtime.JSONPb{}, mockServer, req, nil)

			if tt.expectError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectError)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectCode, st.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			mockServer.AssertExpectations(t)

		})
	}
}

func Test_request_UserService_LogoutUser_0(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockUserServiceClient)
		expectCode  codes.Code
		expectError string
	}{
		{
			name: "success",
			setupMock: func(m *MockUserServiceClient) {
				m.On("LogoutUser", mock.Anything, mock.Anything, mock.Anything).Return(&LogoutUserResponse{Status: true}, nil)
			},
			expectCode: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := new(MockUserServiceClient)
			if tt.setupMock != nil {
				tt.setupMock(mockClient)
			}
			req := httptest.NewRequest(http.MethodPost, "/v1/user/logout", nil)

			resp, _, err := request_UserService_LogoutUser_0(context.Background(), &runtime.JSONPb{}, mockClient, req, nil)

			if tt.expectError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectError)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectCode, st.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}

			mockClient.AssertExpectations(t)
		})
	}
}

func Test_local_request_UserService_LogoutUser_0(t *testing.T) {
	tests := []struct {
		name        string
		setupMock   func(*MockUserServiceServer)
		expectCode  codes.Code
		expectError string
	}{
		{
			name: "success",
			setupMock: func(m *MockUserServiceServer) {
				m.On("LogoutUser", mock.Anything, mock.Anything).Return(&LogoutUserResponse{Status: true}, nil)
			},
			expectCode: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := new(MockUserServiceServer)
			if tt.setupMock != nil {
				tt.setupMock(mockServer)
			}
			req := httptest.NewRequest(http.MethodPost, "/v1/user/logout", nil)
			resp, _, err := local_request_UserService_LogoutUser_0(context.Background(), &runtime.JSONPb{}, mockServer, req, nil)

			if tt.expectError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectError)
				st, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectCode, st.Code())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
			mockServer.AssertExpectations(t)
		})
	}
}
