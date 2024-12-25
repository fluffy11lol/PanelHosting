package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"project-microservice/internal/middleware"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	client "project-microservice/pkg/api/project/protobuf"
	"project-microservice/pkg/logger"
)

type Server struct {
	grpcServer *grpc.Server
	restServer *http.Server
	listener   net.Listener
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Устанавливаем CORS-заголовки
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8082") // Укажите конкретный домен, если нужно
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Если это preflight-запрос (OPTIONS), отправляем пустой ответ
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Передаем запрос дальше
		next.ServeHTTP(w, r)
	})
}

func New(ctx context.Context, port, restPort int, service Service) (*Server, error) {
	if service == nil {
		return nil, fmt.Errorf("service cannot be nil")
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to listen", zap.Error(err))
		return nil, err
	}
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			ContextWithLogger(logger.GetLoggerFromCtx(ctx)),
			AuthInterceptor(logger.GetLoggerFromCtx(ctx)),
		),
		grpc.MaxRecvMsgSize(50 * 1024 * 1024),
		grpc.MaxSendMsgSize(50 * 1024 * 1024),
	}

	grpcServer := grpc.NewServer(opts...)
	client.RegisterProjectServiceServer(grpcServer, NewProjectService(service))

	restSrv := runtime.NewServeMux()
	handlerWithMiddleware := middleware.Authorized(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		restSrv.ServeHTTP(w, r)
	}))
	if err := client.RegisterProjectServiceHandlerServer(context.Background(), restSrv, NewProjectService(service)); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to register rest handler", zap.Error(err))
		return nil, err
	}
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", restPort),
		Handler: enableCORS(handlerWithMiddleware),
	}
	return &Server{grpcServer, httpServer, lis}, nil
}

func (s *Server) Start(ctx context.Context) error {
	eg := errgroup.Group{}

	eg.Go(func() error {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "starting gRPC server", zap.Int("port", s.listener.Addr().(*net.TCPAddr).Port))
		return s.grpcServer.Serve(s.listener)
	})

	eg.Go(func() error {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "starting Rest server", zap.String("port", s.restServer.Addr))
		return s.restServer.ListenAndServe()
	})

	return eg.Wait()
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	l := logger.GetLoggerFromCtx(ctx)
	if l != nil {
		l.Info(ctx, "gRPC server stopped")
		_ = l.Sync()
	}
	return s.restServer.Shutdown(ctx)
}
