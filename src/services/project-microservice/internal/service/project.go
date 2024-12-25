package service

import (
	"context"

	"project-microservice/internal/models"
)

type ProjectRepo interface {
	CreateProject(ctx context.Context, project models.Project) (*models.Project, error)
	GetProject(ctx context.Context, project models.Project) (*models.Project, error)
	UpdateProject(ctx context.Context, project models.Project) (*models.Project, error)
	DeleteProject(ctx context.Context, project models.Project) error
	ListProjects(ctx context.Context, project models.Project) ([]*models.Project, error)
}

type ProjectService struct {
	Repo ProjectRepo
}

func NewProjectService(repo ProjectRepo) *ProjectService {
	return &ProjectService{repo}
}

func (s *ProjectService) CreateProject(ctx context.Context, project models.Project) (*models.Project, error) {
	return s.Repo.CreateProject(ctx, project)
}

func (s *ProjectService) GetProject(ctx context.Context, project models.Project) (*models.Project, error) {
	return s.Repo.GetProject(ctx, project)
}

func (s *ProjectService) UpdateProject(ctx context.Context, project models.Project) (*models.Project, error) {
	return s.Repo.UpdateProject(ctx, project)
}

func (s *ProjectService) DeleteProject(ctx context.Context, project models.Project) error {
	return s.Repo.DeleteProject(ctx, project)
}

func (s *ProjectService) ListProjects(ctx context.Context, project models.Project) ([]*models.Project, error) {
	return s.Repo.ListProjects(ctx, project)
}
