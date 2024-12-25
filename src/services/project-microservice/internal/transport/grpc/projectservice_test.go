package grpc

import (
	"context"
	"errors"
	"project-microservice/internal/middleware"
	"project-microservice/internal/models"
	client "project-microservice/pkg/api/project/protobuf"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Mock for Service interface
type MockService struct {
	mock.Mock
}

func (m *MockService) CreateProject(ctx context.Context, project models.Project) (*models.Project, error) {
	args := m.Called(ctx, project)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}
func (m *MockService) GetProject(ctx context.Context, project models.Project) (*models.Project, error) {
	args := m.Called(ctx, project)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *MockService) UpdateProject(ctx context.Context, project models.Project) (*models.Project, error) {
	args := m.Called(ctx, project)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *MockService) DeleteProject(ctx context.Context, project models.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockService) ListProjects(ctx context.Context, project models.Project) ([]*models.Project, error) {
	args := m.Called(ctx, project)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Project), args.Error(1)
}

func TestProjectService_CreateProject(t *testing.T) {
	testCases := []struct {
		name             string
		inputContext     context.Context
		inputRequest     *client.CreateProjectRequest
		mockServiceSetup func(mockService *MockService)
		expectedResponse *client.CreateProjectResponse
		expectedError    error
	}{
		{
			name:         "success",
			inputContext: context.WithValue(context.Background(), middleware.UserIDKey, "test_user"),
			inputRequest: &client.CreateProjectRequest{
				Name: "test_project",
			},
			mockServiceSetup: func(mockService *MockService) {
				mockService.On("CreateProject", mock.Anything, mock.AnythingOfType("models.Project")).Return(&models.Project{
					ID: "some_id",
				}, nil)
			},
			expectedResponse: &client.CreateProjectResponse{Status: true},
			expectedError:    nil,
		},
		{
			name:             "missing_user_id",
			inputContext:     context.Background(),
			inputRequest:     &client.CreateProjectRequest{Name: "test_project"},
			mockServiceSetup: func(mockService *MockService) {},
			expectedResponse: nil,
			expectedError:    errors.New("userID not found in context"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockService)
			tc.mockServiceSetup(mockService)

			svc := NewProjectService(mockService)
			resp, err := svc.CreateProject(tc.inputContext, tc.inputRequest)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, resp)
			}
			mockService.AssertExpectations(t)
		})
	}
}

func TestProjectService_GetProject(t *testing.T) {
	testCases := []struct {
		name             string
		inputContext     context.Context
		inputRequest     *emptypb.Empty
		mockServiceSetup func(mockService *MockService)
		expectedResponse *client.GetProjectResponse
		expectedError    error
	}{
		{
			name:         "success",
			inputContext: context.WithValue(context.Background(), middleware.UserIDKey, "test_user"),
			inputRequest: &emptypb.Empty{},
			mockServiceSetup: func(mockService *MockService) {
				mockService.On("GetProject", mock.Anything, models.Project{UserID: "test_user"}).Return(&models.Project{
					ID:     "test_id",
					Name:   "test_project",
					UserID: "test_user",
					Networking: models.NetworkingInfo{
						Address: "test_address",
					},
					Status:       models.ProjectRunning,
					TariffID:     1,
					TariffStatus: models.TariffActive,
					CreatedAt:    time.Now(),
				}, nil)
			},
			expectedResponse: &client.GetProjectResponse{
				Server: &client.Server{
					Id:     "test_id",
					Name:   "test_project",
					UserId: "test_user",
					Networking: &client.Networking{
						Address: "test_address",
					},
					Status: int32(models.ProjectRunning),
					TariffInfo: &client.TariffInfo{
						TariffId:     1,
						TariffStatus: int32(models.TariffActive),
					},
					// Игнорируем поле CreatedAt при сравнении
				},
			},
			expectedError: nil,
		},
		{
			name:             "missing_user_id",
			inputContext:     context.Background(),
			inputRequest:     &emptypb.Empty{},
			mockServiceSetup: func(mockService *MockService) {},
			expectedResponse: nil,
			expectedError:    errors.New("userID not found in context"),
		},
		{
			name:         "service_error",
			inputContext: context.WithValue(context.Background(), middleware.UserIDKey, "test_user"),
			inputRequest: &emptypb.Empty{},
			mockServiceSetup: func(mockService *MockService) {
				mockService.On("GetProject", mock.Anything, models.Project{UserID: "test_user"}).Return(nil, errors.New("service error"))
			},
			expectedResponse: nil,
			expectedError:    errors.New("service error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := new(MockService)
			tc.mockServiceSetup(mockService)

			svc := NewProjectService(mockService)
			resp, err := svc.GetProject(tc.inputContext, tc.inputRequest)

			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				// Игнорируем поле CreatedAt при сравнении
				assert.True(t, compareServersIgnoringCreatedAt(tc.expectedResponse.Server, resp.Server))
			}
			mockService.AssertExpectations(t)
		})
	}
}

// Функция для сравнения объектов Server, игнорируя поле CreatedAt
func compareServersIgnoringCreatedAt(expected, actual *client.Server) bool {
	// Сравниваем все поля, кроме CreatedAt
	return expected.Id == actual.Id &&
		expected.Name == actual.Name &&
		expected.UserId == actual.UserId &&
		expected.Status == actual.Status &&
		expected.TariffInfo.TariffId == actual.TariffInfo.TariffId &&
		expected.TariffInfo.TariffStatus == actual.TariffInfo.TariffStatus &&
		expected.Networking.Address == actual.Networking.Address
}
