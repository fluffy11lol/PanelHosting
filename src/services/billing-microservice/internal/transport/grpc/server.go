package grpc

import (
	"billing-microservice/internal/middleware"
	"context"
	"fmt"
	"net"
	"net/http"

	client "billing-microservice/pkg/api/billing/protobuf"
	"billing-microservice/pkg/logger"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Server struct {
	grpcServer *grpc.Server
	restServer *http.Server
	listener   net.Listener
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
	client.RegisterUserServiceServer(grpcServer, NewBillingService(service))

	restSrv := runtime.NewServeMux()
	_ = middleware.Authorized(func(w http.ResponseWriter, r *http.Request) {
		restSrv.ServeHTTP(w, r)
	})
	if err := client.RegisterUserServiceHandlerServer(context.Background(), restSrv, NewBillingService(service)); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to register rest handler", zap.Error(err))
		return nil, err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/tariffs", func(w http.ResponseWriter, r *http.Request) {
		// Убедитесь, что путь к файлу верный
		http.ServeFile(w, r, "./static/index.html")
	})
	mux.Handle("/", restSrv)
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", restPort),
		Handler: mux,
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
