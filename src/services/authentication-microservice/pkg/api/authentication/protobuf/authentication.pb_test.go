package authentication

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

// TestUserServiceServer is a mock implementation of UserServiceServer for testing.
type TestUserServiceServer struct {
	mockRegisterUser func(ctx context.Context, in *RegisterUserRequest) (*RegisterUserResponse, error)
	mockLoginUser    func(ctx context.Context, in *LoginUserRequest) (*LoginUserResponse, error)
	mockLogoutUser   func(ctx context.Context, in *LogoutUserRequest) (*LogoutUserResponse, error)
}

func (m *TestUserServiceServer) RegisterUser(ctx context.Context, in *RegisterUserRequest) (*RegisterUserResponse, error) {
	if m.mockRegisterUser != nil {
		return m.mockRegisterUser(ctx, in)
	}
	return nil, nil
}

func (m *TestUserServiceServer) LoginUser(ctx context.Context, in *LoginUserRequest) (*LoginUserResponse, error) {
	if m.mockLoginUser != nil {
		return m.mockLoginUser(ctx, in)
	}
	return nil, nil
}

func (m *TestUserServiceServer) LogoutUser(ctx context.Context, in *LogoutUserRequest) (*LogoutUserResponse, error) {
	if m.mockLogoutUser != nil {
		return m.mockLogoutUser(ctx, in)
	}
	return nil, nil
}

// TestLogoutUserHandler tests the LogoutUser endpoint handler.
func TestLogoutUserHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := &TestUserServiceServer{
		mockLogoutUser: func(ctx context.Context, in *LogoutUserRequest) (*LogoutUserResponse, error) {
			return &LogoutUserResponse{Status: true}, nil
		},
	}

	mux := runtime.NewServeMux()
	RegisterUserServiceHandlerServer(context.Background(), mux, mockServer)

	server := httptest.NewServer(mux)
	defer server.Close()

	testCases := []struct {
		name       string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "valid request",
			wantStatus: http.StatusOK,
			wantBody:   `{"status":true}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", server.URL+"/v1/user/logout", nil)
			if err != nil {
				t.Fatalf("could not create request: %v", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("could not send request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.wantStatus {
				t.Errorf("expected status %d, got %d", tc.wantStatus, resp.StatusCode)
			}

			respBody := &LogoutUserResponse{}
			if err := protojson.Unmarshal(getBodyBytes(resp), respBody); err != nil {
				t.Fatalf("could not unmarshal response body: %v", err)
			}

			wantBody := &LogoutUserResponse{}
			if err := protojson.Unmarshal([]byte(tc.wantBody), wantBody); err != nil {
				t.Fatalf("could not unmarshal want body: %v", err)
			}

			if respBody.GetStatus() != wantBody.GetStatus() {
				t.Errorf("expected body %v, got %v", wantBody, respBody)
			}
		})
	}
}

// Helper function to extract the body bytes from an HTTP response.
func getBodyBytes(resp *http.Response) []byte {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		// Handle the error appropriately, perhaps logging or returning an empty slice
		// For this example, we'll just return an empty slice.
		return []byte{}
	}
	return body
}

// Implements the UserServiceServer interface
func (m *TestUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

func TestServiceRegistration(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockServer := &TestUserServiceServer{
		mockRegisterUser: func(ctx context.Context, in *RegisterUserRequest) (*RegisterUserResponse, error) {
			return &RegisterUserResponse{Status: true}, nil
		},
		mockLoginUser: func(ctx context.Context, in *LoginUserRequest) (*LoginUserResponse, error) {
			return &LoginUserResponse{Token: "testtoken"}, nil
		},
		mockLogoutUser: func(ctx context.Context, in *LogoutUserRequest) (*LogoutUserResponse, error) {
			return &LogoutUserResponse{Status: true}, nil
		},
	}

	server := grpc.NewServer()
	RegisterUserServiceServer(server, mockServer)

	serviceInfo := server.GetServiceInfo()

	if _, ok := serviceInfo["api.UserService"]; !ok {
		t.Fatalf("Service 'api.UserService' is not registered.")
	}

}
