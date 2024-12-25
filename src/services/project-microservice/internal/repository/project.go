package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"

	"project-microservice/internal/models"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewProjectRepository(db *gorm.DB, redis *redis.Client) *ProjectRepository {
	return &ProjectRepository{db: db, redis: redis}
}
func (r *ProjectRepository) CreateProject(ctx context.Context, project models.Project) (*models.Project, error) {
	var newProject *models.Project

	operation := func() error {
		if project.Name == "" {
			return errors.New("project name is required")
		}
		if err := r.db.WithContext(ctx).Where("name = ?", project.Name).First(&models.Project{}).Error; err == nil {
			return errors.New("project already exists")
		}

		newProject = &models.Project{
			ID:             uuid.New().String(),
			Name:           project.Name,
			UserID:         project.UserID,
			Status:         models.ProjectStopped,
			TariffID:       project.TariffID,
			TariffStatus:   project.TariffStatus,
			Networking:     models.NetworkingInfo{Address: "user1.ghst.tech", Ports: models.Int32Slice{25565, 25566}},
			ExpirationTime: project.ExpirationTime,
			CreatedAt:      time.Now(),
		}

		if err := r.db.WithContext(ctx).Create(newProject).Error; err != nil {
			return fmt.Errorf("failed to create project: %w", err)
		}
		return nil
	}
	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = 10 * time.Second
	err := backoff.Retry(operation, bo)
	if err != nil {
		return nil, fmt.Errorf("failed to create project after retries: %w", err)
	}

	cacheKeyProject := fmt.Sprintf("project:%s", project.UserID)
	if err := r.redis.Del(ctx, cacheKeyProject).Err(); err != nil {
		fmt.Println("error invalidating cache in Redis:", err)
	}

	return newProject, nil
}

func (r *ProjectRepository) GetProject(ctx context.Context, proj models.Project) (*models.Project, error) {
	cacheKey := fmt.Sprintf("project:%s", proj.UserID)
	cachedProject, err := r.redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var project models.Project
		if err := json.Unmarshal([]byte(cachedProject), &project); err == nil {
			return &project, nil
		}
	} else if !errors.Is(err, redis.Nil) {
		fmt.Println("error getting data from Redis:", err)
	}

	var project models.Project
	if err := r.db.WithContext(ctx).Where("user_id = ?", proj.UserID).First(&project).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	projectBytes, err := json.Marshal(project)
	if err == nil {
		err = r.redis.Set(ctx, cacheKey, projectBytes, time.Hour).Err()
		if err != nil {
			fmt.Println("error saving data to Redis:", err)
		}
	}

	return &project, nil
}

func (r *ProjectRepository) UpdateProject(ctx context.Context, proj models.Project) (*models.Project, error) {
	var existingProject models.Project
	if err := r.db.WithContext(ctx).First(&existingProject, "id = ?", proj.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("project with ID %s not found", proj.ID)
		}
		return nil, fmt.Errorf("failed to find project: %v", err)
	}

	if err := r.db.WithContext(ctx).Model(&existingProject).Updates(proj).Error; err != nil {
		return nil, fmt.Errorf("failed to update project: %v", err)
	}

	return &existingProject, nil
}

func (r *ProjectRepository) DeleteProject(ctx context.Context, proj models.Project) error {
	var project models.Project
	if err := r.db.WithContext(ctx).First(&project, "user_id = ?", proj.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("project with user_ID %s not found", proj.UserID)
		}
		return fmt.Errorf("failed to find project: %v", err)
	}

	if err := r.db.WithContext(ctx).Delete(&project).Error; err != nil {
		return fmt.Errorf("failed to delete project: %v", err)
	}

	return nil
}

func (r *ProjectRepository) ListProjects(ctx context.Context, proj models.Project) ([]*models.Project, error) {
	var projects []*models.Project
	// Запрос к базе для извлечения проектов, связанных с пользователем
	if err := r.db.WithContext(ctx).Where("user_id =?", proj.UserID).Find(&projects).Error; err != nil {
		return nil, fmt.Errorf("failed to list projects: %v", err)
	}

	return projects, nil
}
