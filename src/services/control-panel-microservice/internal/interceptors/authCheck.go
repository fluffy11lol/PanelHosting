package interceptors

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"control-panel/pkg/api/logger"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func GenerateJWTtoken(id int) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":  string(rune(id)),
			"nbf": now.Unix(),
			"exp": now.Add(2 * time.Hour).Unix(),
			"iat": now.Unix(),
		})

	tokenString, err := token.SignedString([]byte(string(rune(id))))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AuthInterceptor(allowedMethods []string, l logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		l.Info(ctx, "request started", zap.String("method", info.FullMethod))

		methodName := info.FullMethod

		for _, allowed := range allowedMethods {
			if strings.HasSuffix(methodName, allowed) {
				return handler(ctx, req)
			}
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			l.Error(ctx, "missing metadata")
			return nil, status.Errorf(http.StatusUnauthorized, "Missing metadata")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			l.Error(ctx, "missing authorization header")
			return nil, status.Errorf(http.StatusUnauthorized, "Missing Authorization header")
		}

		tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")
		if tokenString == "" {
			l.Error(ctx, "invalid token format")
			return nil, status.Errorf(http.StatusUnauthorized, "Invalid token format")
		}

		var UserID string

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return nil, errors.New("invalid token format")
			}

			id, ok := claims["id"].(string)
			if !ok {
				return nil, fmt.Errorf("invalid field id.")
			}
			UserID = id
			return []byte(id), nil
		})

		if err != nil {
			l.Error(ctx, "invalid token")
			return nil, status.Errorf(http.StatusUnauthorized, "Invalid token: %v", err)
		}

		if !token.Valid {
			l.Error(ctx, "invalid token claims")
			return nil, status.Errorf(http.StatusUnauthorized, "Invalid token claims")
		}

		l.Info(ctx, "user authenticated", zap.String("user id", UserID))

		return handler(ctx, req)
	}
}
