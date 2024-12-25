package project

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/emptypb"
)

// MockProjectServiceServer -  мок для ProjectServiceServer
type MockProjectServiceServer struct {
	CreateProjectFn func(ctx context.Context, req *CreateProjectRequest) (*CreateProjectResponse, error)
	GetProjectFn    func(ctx context.Context, req *emptypb.Empty) (*GetProjectResponse, error)
	UpdateProjectFn func(ctx context.Context, req *UpdateProjectRequest) (*UpdateProjectResponse, error)
	DeleteProjectFn func(ctx context.Context, req *DeleteProjectRequest) (*DeleteProjectResponse, error)
	ListProjectsFn  func(ctx context.Context, req *emptypb.Empty) (*ListProjectsResponse, error)
}

func (m *MockProjectServiceServer) CreateProject(ctx context.Context, req *CreateProjectRequest) (*CreateProjectResponse, error) {
	if m.CreateProjectFn != nil {
		return m.CreateProjectFn(ctx, req)
	}
	return nil, errors.New("CreateProject function not implemented in mock")
}

func (m *MockProjectServiceServer) GetProject(ctx context.Context, req *emptypb.Empty) (*GetProjectResponse, error) {
	if m.GetProjectFn != nil {
		return m.GetProjectFn(ctx, req)
	}
	return nil, errors.New("GetProject function not implemented in mock")
}

func (m *MockProjectServiceServer) UpdateProject(ctx context.Context, req *UpdateProjectRequest) (*UpdateProjectResponse, error) {
	if m.UpdateProjectFn != nil {
		return m.UpdateProjectFn(ctx, req)
	}
	return nil, errors.New("UpdateProject function not implemented in mock")
}
func (m *MockProjectServiceServer) DeleteProject(ctx context.Context, req *DeleteProjectRequest) (*DeleteProjectResponse, error) {
	if m.DeleteProjectFn != nil {
		return m.DeleteProjectFn(ctx, req)
	}
	return nil, errors.New("DeleteProject function not implemented in mock")
}

func (m *MockProjectServiceServer) ListProjects(ctx context.Context, req *emptypb.Empty) (*ListProjectsResponse, error) {
	if m.ListProjectsFn != nil {
		return m.ListProjectsFn(ctx, req)
	}
	return nil, errors.New("ListProjects function not implemented in mock")
}
func (m *MockProjectServiceServer) mustEmbedUnimplementedProjectServiceServer() {}
func TestCreateProject(t *testing.T) {
	testCases := []struct {
		name               string
		requestBody        interface{}
		mockServerResponse *CreateProjectResponse
		mockServerError    error
		expectedStatus     int
		expectedBody       string
	}{
		{
			name: "successful project creation",
			requestBody: map[string]interface{}{
				"name": "test_project",
			},
			mockServerResponse: &CreateProjectResponse{Status: true},
			expectedStatus:     http.StatusOK,
			expectedBody:       `{"Status":true}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockServer := &MockProjectServiceServer{
				CreateProjectFn: func(ctx context.Context, req *CreateProjectRequest) (*CreateProjectResponse, error) {
					if tc.mockServerError != nil {
						return nil, tc.mockServerError
					}
					return tc.mockServerResponse, nil
				},
			}

			mux := runtime.NewServeMux()
			err := RegisterProjectServiceHandlerServer(context.Background(), mux, mockServer)
			assert.NoError(t, err)

			bodyBytes, _ := json.Marshal(tc.requestBody)
			req, err := http.NewRequest("POST", "/v1/project/create", bytes.NewReader(bodyBytes))
			assert.NoError(t, err)

			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.JSONEq(t, tc.expectedBody, recorder.Body.String())

		})
	}
}

func TestDeleteProject(t *testing.T) {
	testCases := []struct {
		name               string
		requestQuery       string
		mockServerResponse *DeleteProjectResponse
		mockServerError    error
		expectedStatus     int
		expectedBody       string
	}{
		{
			name:               "successful project deletion",
			requestQuery:       "ID=test-id",
			mockServerResponse: &DeleteProjectResponse{Status: true},
			expectedStatus:     http.StatusOK,
			expectedBody:       `{"Status":true}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockServer := &MockProjectServiceServer{
				DeleteProjectFn: func(ctx context.Context, req *DeleteProjectRequest) (*DeleteProjectResponse, error) {
					if tc.mockServerError != nil {
						return nil, tc.mockServerError
					}
					return tc.mockServerResponse, nil
				},
			}

			mux := runtime.NewServeMux()
			err := RegisterProjectServiceHandlerServer(context.Background(), mux, mockServer)
			assert.NoError(t, err)

			req, err := http.NewRequest("DELETE", "/v1/project/delete?"+tc.requestQuery, nil)
			assert.NoError(t, err)

			recorder := httptest.NewRecorder()
			mux.ServeHTTP(recorder, req)

			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.JSONEq(t, tc.expectedBody, recorder.Body.String())
		})
	}
}
