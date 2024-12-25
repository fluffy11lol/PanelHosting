package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"project-microservice/internal/models"
)

func setupTestDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Exec("DELETE FROM projects")

	if err := db.AutoMigrate(&models.Project{}); err != nil {
		return nil, err
	}
	return db, nil
}

func setupTestRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	return client
}

func TestProjectRepository_GetProject(t *testing.T) {

	db, err := setupTestDB()
	require.NoError(t, err)

	redisClient := setupTestRedis()
	repo := NewProjectRepository(db, redisClient)
	ctx := context.Background()
	defer redisClient.FlushAll(ctx)

	project := models.Project{Name: "Test Project", UserID: "user1"}
	createdProject, err := repo.CreateProject(ctx, project)
	require.NoError(t, err)

	t.Run("Success_CacheHit", func(t *testing.T) {
		cachedProject, err := repo.GetProject(ctx, models.Project{UserID: "user1"})
		require.NoError(t, err)
		assert.Equal(t, createdProject.ID, cachedProject.ID)
		assert.Equal(t, "Test Project", cachedProject.Name)
	})

	t.Run("Success_CacheMiss_DBHit", func(t *testing.T) {
		redisClient.Del(ctx, fmt.Sprintf("project:%s", "user1"))
		cachedProject, err := repo.GetProject(ctx, models.Project{UserID: "user1"})
		require.NoError(t, err)
		assert.Equal(t, createdProject.ID, cachedProject.ID)
		assert.Equal(t, "Test Project", cachedProject.Name)
	})
}

func TestProjectRepository_DeleteProject(t *testing.T) {

	db, err := setupTestDB()
	require.NoError(t, err)
	redisClient := setupTestRedis()

	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		repo := NewProjectRepository(tx, redisClient)
		project := models.Project{Name: "Test Project", UserID: "user1"}
		_, err = repo.CreateProject(ctx, project)
		require.NoError(t, err)

		err = repo.DeleteProject(ctx, models.Project{UserID: "user1"})
		require.NoError(t, err)

		_, err = repo.GetProject(ctx, models.Project{UserID: "user1"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "project not found")
	})

	t.Run("Failed_ProjectNotFound", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()
		repo := NewProjectRepository(tx, redisClient)
		err := repo.DeleteProject(ctx, models.Project{UserID: "nonexistentuser"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "project with user_ID")
	})
}
