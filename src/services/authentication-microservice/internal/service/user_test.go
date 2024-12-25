package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"authentication-microservice/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) CreateUser(ctx context.Context, user models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) LoginUser(ctx context.Context, user models.User) (*models.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestNewUserService(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := NewUserService(mockRepo)
	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.Repo)
}

func TestUserService_RegisterUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("successful registration", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Password: "password123",
			Email:    "test@example.com",
		}
		expectedUser := &models.User{
			ID:       "123",
			Username: "testuser",
			Password: "hashedpassword",
			Email:    "test@example.com",
		}

		mockRepo.On("CreateUser", ctx, user).Return(expectedUser, nil).Once()

		result, err := service.RegisterUser(ctx, user)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("registration with empty username", func(t *testing.T) {
		user := models.User{
			Username: "",
			Password: "password123",
			Email:    "test@example.com",
		}
		expectedError := errors.New("username is required")

		mockRepo.On("CreateUser", ctx, user).Return(nil, expectedError).Once()

		result, err := service.RegisterUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("registration with empty email", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Password: "password123",
			Email:    "",
		}
		expectedError := errors.New("email is required")

		mockRepo.On("CreateUser", ctx, user).Return(nil, expectedError).Once()

		result, err := service.RegisterUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("registration with invalid email format", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Password: "password123",
			Email:    "invalid-email",
		}
		expectedError := errors.New("invalid email format")

		mockRepo.On("CreateUser", ctx, user).Return(nil, expectedError).Once()

		result, err := service.RegisterUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("registration with duplicate username", func(t *testing.T) {
		user := models.User{
			Username: "existinguser",
			Password: "password123",
			Email:    "test@example.com",
		}
		expectedError := errors.New("username already exists")

		mockRepo.On("CreateUser", ctx, user).Return(nil, expectedError).Once()

		result, err := service.RegisterUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("registration with very long password", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Password: string(make([]byte, 73)),
			Email:    "test@example.com",
		}
		expectedError := errors.New("password too long")

		mockRepo.On("CreateUser", ctx, user).Return(nil, expectedError).Once()

		result, err := service.RegisterUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("registration with database error", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Password: "password123",
			Email:    "test@example.com",
		}
		expectedError := errors.New("database connection error")

		mockRepo.On("CreateUser", ctx, user).Return(nil, expectedError).Once()

		result, err := service.RegisterUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_LoginUser(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("successful login", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Password: "password123",
		}
		expectedUser := &models.User{
			ID:       "123",
			Username: "testuser",
			Password: "hashedpassword",
			Email:    "test@example.com",
			Token:    "jwt-token",
		}

		mockRepo.On("LoginUser", ctx, user).Return(expectedUser, nil).Once()

		result, err := service.LoginUser(ctx, user)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		assert.NotEmpty(t, result.Token)
		mockRepo.AssertExpectations(t)
	})

	t.Run("login with empty username", func(t *testing.T) {
		user := models.User{
			Username: "",
			Password: "password123",
		}
		expectedError := errors.New("username is required")

		mockRepo.On("LoginUser", ctx, user).Return(nil, expectedError).Once()

		result, err := service.LoginUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("login with empty password", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Password: "",
		}
		expectedError := errors.New("password is required")

		mockRepo.On("LoginUser", ctx, user).Return(nil, expectedError).Once()

		result, err := service.LoginUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("login with invalid credentials", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Password: "wrongpassword",
		}
		expectedError := errors.New("invalid credentials")

		mockRepo.On("LoginUser", ctx, user).Return(nil, expectedError).Once()

		result, err := service.LoginUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("login with non-existent user", func(t *testing.T) {
		user := models.User{
			Username: "nonexistent",
			Password: "password123",
		}
		expectedError := errors.New("user not found")

		mockRepo.On("LoginUser", ctx, user).Return(nil, expectedError).Once()

		result, err := service.LoginUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("login with database error", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
			Password: "password123",
		}
		expectedError := errors.New("database connection error")

		mockRepo.On("LoginUser", ctx, user).Return(nil, expectedError).Once()

		result, err := service.LoginUser(ctx, user)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("login with context timeout", func(t *testing.T) {
		timeoutCtx, cancel := context.WithTimeout(ctx, time.Millisecond)
		defer cancel()

		user := models.User{
			Username: "testuser",
			Password: "password123",
		}
		expectedError := context.DeadlineExceeded

		mockRepo.On("LoginUser", timeoutCtx, user).Return(nil, expectedError).Once()

		result, err := service.LoginUser(timeoutCtx, user)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("login with canceled context", func(t *testing.T) {
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel()

		user := models.User{
			Username: "testuser",
			Password: "password123",
		}
		expectedError := context.Canceled

		mockRepo.On("LoginUser", cancelCtx, user).Return(nil, expectedError).Once()

		result, err := service.LoginUser(cancelCtx, user)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
