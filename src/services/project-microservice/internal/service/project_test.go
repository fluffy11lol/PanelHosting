package service

import (
	"context"
	"errors"
	"project-microservice/internal/models"
	"reflect"
	"strconv"
	"testing"
)

type MockProjectRepo struct {
	CreateProjectFunc func(ctx context.Context, project models.Project) (*models.Project, error)
	GetProjectFunc    func(ctx context.Context, project models.Project) (*models.Project, error)
	UpdateProjectFunc func(ctx context.Context, project models.Project) (*models.Project, error)
	DeleteProjectFunc func(ctx context.Context, project models.Project) error
	ListProjectsFunc  func(ctx context.Context, project models.Project) ([]*models.Project, error)
}

func (m *MockProjectRepo) CreateProject(ctx context.Context, project models.Project) (*models.Project, error) {
	if m.CreateProjectFunc != nil {
		return m.CreateProjectFunc(ctx, project)
	}
	return nil, errors.New("CreateProjectFunc not implemented")
}

func (m *MockProjectRepo) GetProject(ctx context.Context, project models.Project) (*models.Project, error) {
	if m.GetProjectFunc != nil {
		return m.GetProjectFunc(ctx, project)
	}
	return nil, errors.New("GetProjectFunc not implemented")
}

func (m *MockProjectRepo) UpdateProject(ctx context.Context, project models.Project) (*models.Project, error) {
	if m.UpdateProjectFunc != nil {
		return m.UpdateProjectFunc(ctx, project)
	}
	return nil, errors.New("UpdateProjectFunc not implemented")
}

func (m *MockProjectRepo) DeleteProject(ctx context.Context, project models.Project) error {
	if m.DeleteProjectFunc != nil {
		return m.DeleteProjectFunc(ctx, project)
	}
	return errors.New("DeleteProjectFunc not implemented")
}

func (m *MockProjectRepo) ListProjects(ctx context.Context, project models.Project) ([]*models.Project, error) {
	if m.ListProjectsFunc != nil {
		return m.ListProjectsFunc(ctx, project)
	}
	return nil, errors.New("ListProjectsFunc not implemented")
}

func TestNewProjectService(t *testing.T) {
	repo := &MockProjectRepo{}
	service := NewProjectService(repo)
	if service == nil {
		t.Error("NewProjectService should return a valid service")
	}

	if service.Repo != repo {
		t.Error("NewProjectService should set the repository")
	}
}

func TestProjectService_CreateProject(t *testing.T) {
	testProject := models.Project{ID: strconv.Itoa(1), Name: "Test Project"}
	testError := errors.New("create error")

	tests := []struct {
		name            string
		mockRepo        *MockProjectRepo
		expectedProject *models.Project
		expectedError   error
	}{
		{
			name: "Successful create",
			mockRepo: &MockProjectRepo{
				CreateProjectFunc: func(ctx context.Context, project models.Project) (*models.Project, error) {
					return &testProject, nil
				},
			},
			expectedProject: &testProject,
			expectedError:   nil,
		},
		{
			name: "Failed create",
			mockRepo: &MockProjectRepo{
				CreateProjectFunc: func(ctx context.Context, project models.Project) (*models.Project, error) {
					return nil, testError
				},
			},
			expectedProject: nil,
			expectedError:   testError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewProjectService(tt.mockRepo)
			project, err := service.CreateProject(context.Background(), testProject)
			if !reflect.DeepEqual(project, tt.expectedProject) {
				t.Errorf("CreateProject() = %v, want %v", project, tt.expectedProject)
			}
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("CreateProject() error = %v, want %v", err, tt.expectedError)
			}
		})
	}
}

func TestProjectService_GetProject(t *testing.T) {
	testProject := models.Project{ID: strconv.Itoa(1), Name: "Test Project"}
	testError := errors.New("get error")

	tests := []struct {
		name            string
		mockRepo        *MockProjectRepo
		expectedProject *models.Project
		expectedError   error
	}{
		{
			name: "Successful get",
			mockRepo: &MockProjectRepo{
				GetProjectFunc: func(ctx context.Context, project models.Project) (*models.Project, error) {
					return &testProject, nil
				},
			},
			expectedProject: &testProject,
			expectedError:   nil,
		},
		{
			name: "Failed get",
			mockRepo: &MockProjectRepo{
				GetProjectFunc: func(ctx context.Context, project models.Project) (*models.Project, error) {
					return nil, testError
				},
			},
			expectedProject: nil,
			expectedError:   testError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewProjectService(tt.mockRepo)
			project, err := service.GetProject(context.Background(), testProject)
			if !reflect.DeepEqual(project, tt.expectedProject) {
				t.Errorf("GetProject() = %v, want %v", project, tt.expectedProject)
			}
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("GetProject() error = %v, want %v", err, tt.expectedError)
			}
		})
	}
}

func TestProjectService_UpdateProject(t *testing.T) {
	testProject := models.Project{ID: strconv.Itoa(1), Name: "Test Project"}
	testError := errors.New("update error")

	tests := []struct {
		name            string
		mockRepo        *MockProjectRepo
		expectedProject *models.Project
		expectedError   error
	}{
		{
			name: "Successful update",
			mockRepo: &MockProjectRepo{
				UpdateProjectFunc: func(ctx context.Context, project models.Project) (*models.Project, error) {
					return &testProject, nil
				},
			},
			expectedProject: &testProject,
			expectedError:   nil,
		},
		{
			name: "Failed update",
			mockRepo: &MockProjectRepo{
				UpdateProjectFunc: func(ctx context.Context, project models.Project) (*models.Project, error) {
					return nil, testError
				},
			},
			expectedProject: nil,
			expectedError:   testError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewProjectService(tt.mockRepo)
			project, err := service.UpdateProject(context.Background(), testProject)
			if !reflect.DeepEqual(project, tt.expectedProject) {
				t.Errorf("UpdateProject() = %v, want %v", project, tt.expectedProject)
			}
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("UpdateProject() error = %v, want %v", err, tt.expectedError)
			}
		})
	}
}

func TestProjectService_DeleteProject(t *testing.T) {
	testProject := models.Project{ID: strconv.Itoa(1), Name: "Test Project"}
	testError := errors.New("delete error")

	tests := []struct {
		name          string
		mockRepo      *MockProjectRepo
		expectedError error
	}{
		{
			name: "Successful delete",
			mockRepo: &MockProjectRepo{
				DeleteProjectFunc: func(ctx context.Context, project models.Project) error {
					return nil
				},
			},
			expectedError: nil,
		},
		{
			name: "Failed delete",
			mockRepo: &MockProjectRepo{
				DeleteProjectFunc: func(ctx context.Context, project models.Project) error {
					return testError
				},
			},
			expectedError: testError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewProjectService(tt.mockRepo)
			err := service.DeleteProject(context.Background(), testProject)

			if !errors.Is(err, tt.expectedError) {
				t.Errorf("DeleteProject() error = %v, want %v", err, tt.expectedError)
			}
		})
	}
}

func TestProjectService_ListProjects(t *testing.T) {
	testProject1 := models.Project{ID: strconv.Itoa(1), Name: "Test Project 1"}
	testProject2 := models.Project{ID: strconv.Itoa(2), Name: "Test Project 2"}
	testProjects := []*models.Project{&testProject1, &testProject2}

	testError := errors.New("list error")

	tests := []struct {
		name             string
		mockRepo         *MockProjectRepo
		expectedProjects []*models.Project
		expectedError    error
	}{
		{
			name: "Successful list",
			mockRepo: &MockProjectRepo{
				ListProjectsFunc: func(ctx context.Context, project models.Project) ([]*models.Project, error) {
					return testProjects, nil
				},
			},
			expectedProjects: testProjects,
			expectedError:    nil,
		},
		{
			name: "Failed list",
			mockRepo: &MockProjectRepo{
				ListProjectsFunc: func(ctx context.Context, project models.Project) ([]*models.Project, error) {
					return nil, testError
				},
			},
			expectedProjects: nil,
			expectedError:    testError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewProjectService(tt.mockRepo)
			projects, err := service.ListProjects(context.Background(), models.Project{})
			if !reflect.DeepEqual(projects, tt.expectedProjects) {
				t.Errorf("ListProjects() = %v, want %v", projects, tt.expectedProjects)
			}
			if !errors.Is(err, tt.expectedError) {
				t.Errorf("ListProjects() error = %v, want %v", err, tt.expectedError)
			}
		})
	}
}
