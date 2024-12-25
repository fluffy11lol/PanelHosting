package grpc

import (
	"context"
	"errors"
	"project-microservice/internal/middleware"
	"project-microservice/internal/models"
	client "project-microservice/pkg/api/project/protobuf"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Service interface {
	CreateProject(ctx context.Context, project models.Project) (*models.Project, error)
	GetProject(ctx context.Context, project models.Project) (*models.Project, error)
	UpdateProject(ctx context.Context, project models.Project) (*models.Project, error)
	DeleteProject(ctx context.Context, project models.Project) error
	ListProjects(ctx context.Context, project models.Project) ([]*models.Project, error)
}

type ProjectService struct {
	client.UnimplementedProjectServiceServer
	service Service
}

func NewProjectService(srv Service) *ProjectService {
	return &ProjectService{service: srv}
}

func (s *ProjectService) CreateProject(ctx context.Context, req *client.CreateProjectRequest) (*client.CreateProjectResponse, error) {

	userID, ok := ctx.Value(middleware.UserIDKey).(string)

	if !ok || userID == "" {
		return nil, errors.New("userID not found in context")
	}

	// Создаем проект с привязкой к userID.
	_, err := s.service.CreateProject(ctx, models.Project{
		ID:             uuid.NewString(),
		Name:           req.GetName(),
		TariffID:       0, // Указать, если есть значение по умолчанию.
		TariffStatus:   0,
		ExpirationTime: time.Now().Add(time.Hour),
		UserID:         userID,
	})
	if err != nil {
		return &client.CreateProjectResponse{Status: false}, err
	}
	return &client.CreateProjectResponse{Status: true}, nil
}

func (s *ProjectService) GetProject(ctx context.Context, req *emptypb.Empty) (*client.GetProjectResponse, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, errors.New("userID not found in context")
	}

	// Передаем userID в бизнес-логику
	project, err := s.service.GetProject(ctx, models.Project{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}

	projectGRPC := client.Server{
		Id:     project.ID,
		Name:   project.Name,
		UserId: project.UserID,
		Networking: &client.Networking{
			Address: project.Networking.Address,
		},
		Status: int32(project.Status),
		TariffInfo: &client.TariffInfo{
			TariffId:     project.TariffID,
			TariffStatus: int32(project.TariffStatus),
		},
		CreatedAt: project.CreatedAt.String(),
	}

	// Формируем ответ
	return &client.GetProjectResponse{Server: &projectGRPC}, nil
}

func (s *ProjectService) UpdateProject(ctx context.Context, req *client.UpdateProjectRequest) (*client.UpdateProjectResponse, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, errors.New("userID not found in context")
	}

	reqProject := models.Project{
		ID:           req.GetID(),
		Name:         req.GetName(),
		TariffID:     req.GetTariffInfo().GetTariffId(),
		TariffStatus: models.TariffStatus(req.GetTariffInfo().TariffStatus),
		Status:       models.ProjectStatus(req.GetStatus()),
		UserID:       userID,
	}

	project, err := s.service.UpdateProject(ctx, reqProject)
	if err != nil {
		return nil, err
	}

	projectGRPC := client.Server{
		Id:   project.ID,
		Name: project.Name,
		Networking: &client.Networking{
			Address: project.Networking.Address,
		},
		Status: int32(project.Status),
		TariffInfo: &client.TariffInfo{
			TariffId:     project.TariffID,
			TariffStatus: int32(project.TariffStatus),
		},
		CreatedAt: project.CreatedAt.String(),
	}

	// Формируем ответ
	return &client.UpdateProjectResponse{Server: &projectGRPC}, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, req *client.DeleteProjectRequest) (*client.DeleteProjectResponse, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, errors.New("userID not found in context")
	}

	reqProject := models.Project{
		ID:     req.GetID(),
		UserID: userID,
	}

	err := s.service.DeleteProject(ctx, reqProject)
	if err != nil {
		return &client.DeleteProjectResponse{Status: false}, err
	}

	return &client.DeleteProjectResponse{Status: true}, nil
}

func (s *ProjectService) ListProjects(ctx context.Context, req *emptypb.Empty) (*client.ListProjectsResponse, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		return nil, errors.New("userID not found in context")
	}

	projects, err := s.service.ListProjects(ctx, models.Project{UserID: userID})
	if err != nil {
		return nil, err
	}

	var projectGRPs []*client.Server
	for _, project := range projects {
		projectGRPs = append(projectGRPs, &client.Server{
			Id:   project.ID,
			Name: project.Name,
			Networking: &client.Networking{
				Address: project.Networking.Address,
				Ports:   project.Networking.Ports,
			},
			Status: int32(project.Status),
			TariffInfo: &client.TariffInfo{
				TariffId:     project.TariffID,
				TariffStatus: int32(project.TariffStatus),
			},
			CreatedAt: project.CreatedAt.String(),
		})
	}

	return &client.ListProjectsResponse{Server: projectGRPs}, nil
}
