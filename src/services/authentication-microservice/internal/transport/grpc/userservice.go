package grpc

import (
	"authentication-microservice/internal/models"
	client "authentication-microservice/pkg/api/authentication/protobuf"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Service interface {
	RegisterUser(ctx context.Context, User models.User) (*models.User, error)
	LoginUser(ctx context.Context, User models.User) (*models.User, error)
}

type UserService struct {
	client.UnimplementedUserServiceServer
	service Service
}

func NewUserService(srv Service) *UserService {
	return &UserService{service: srv}
}
func (s *UserService) RegisterUser(ctx context.Context, req *client.RegisterUserRequest) (*client.RegisterUserResponse, error) {
	_, err := s.service.RegisterUser(ctx, models.User{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Username: req.GetUsername(),
	})
	if err != nil {
		return &client.RegisterUserResponse{Status: false}, err
	}
	return &client.RegisterUserResponse{Status: true}, nil
}

func (s *UserService) LoginUser(ctx context.Context, req *client.LoginUserRequest) (*client.LoginUserResponse, error) {
	resp, err := s.service.LoginUser(ctx, models.User{
		Password: req.GetPassword(),
		Username: req.GetUsername(),
	})
	if err != nil {
		return nil, err
	}
	err = grpc.SendHeader(ctx, metadata.Pairs(
		"Set-Cookie",
		fmt.Sprintf("token=%s; HttpOnly; Path=/; Secure", resp.Token),
	))
	if err != nil {
		return nil, err
	}
	return &client.LoginUserResponse{Token: resp.Token}, nil
}

func (s *UserService) LogoutUser(ctx context.Context, req *client.LogoutUserRequest) (*client.LogoutUserResponse, error) {
	err := grpc.SendHeader(ctx, metadata.Pairs(
		"Set-Cookie",
		"token=; HttpOnly; Path=/; Secure; Max-Age=0",
	))
	if err != nil {
		return nil, err
	}
	return &client.LogoutUserResponse{Status: true}, nil
}
