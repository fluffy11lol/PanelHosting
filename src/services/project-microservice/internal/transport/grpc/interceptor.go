package grpc

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc/metadata"
	"project-microservice/pkg/logger"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func ContextWithLogger(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		l.Info(ctx, "request started", zap.String("method", info.FullMethod))
		return handler(ctx, req)
	}
}

var jwtSecret = []byte("My_key")

type JWTClaims struct {
	Username string `json:"username"`
	jwt2.RegisteredClaims
}

func AuthInterceptor(l logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		if strings.HasSuffix(info.FullMethod, "LoginUser") || strings.HasSuffix(info.FullMethod, "RegisterUser") {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			l.Error(ctx, "missing metadata")
			return nil, grpc.Errorf(401, "missing metadata")
		}

		tokenString := getTokenFromMetadata(md)
		if tokenString == "" {
			l.Error(ctx, "missing token")
			return nil, grpc.Errorf(401, "missing token")
		}

		claims, err := validateJWT(tokenString)
		if err != nil {
			l.Error(ctx, "invalid token", zap.Error(err))
			return nil, grpc.Errorf(401, "invalid token")
		}

		l.Info(ctx, "user authenticated", zap.String("username", claims.Username))

		return handler(ctx, ctx)
	}
}

func getTokenFromMetadata(md metadata.MD) string {
	tokens := md["authorization"]
	if len(tokens) == 0 {
		return ""
	}

	parts := strings.SplitN(tokens[0], " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}
	return parts[1]
}

func validateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
