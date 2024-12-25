package service

import (
	"authentication-microservice/internal/models"
	"context"
)

type UserRepo interface {
	CreateUser(ctx context.Context, User models.User) (*models.User, error)
	LoginUser(ctx context.Context, User models.User) (*models.User, error)
}

type UserService struct {
	Repo UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{repo}
}

func (s *UserService) RegisterUser(ctx context.Context, User models.User) (*models.User, error) {
	return s.Repo.CreateUser(ctx, User)
}
func (s *UserService) LoginUser(ctx context.Context, User models.User) (*models.User, error) {
	return s.Repo.LoginUser(ctx, User)
}
