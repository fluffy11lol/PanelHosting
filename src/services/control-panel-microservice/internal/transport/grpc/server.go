package transport

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"control-panel/internal/interceptors"
	models "control-panel/internal/models"
	"control-panel/internal/transport/gateway"
	"control-panel/pkg/api/logger"
	panel "control-panel/pkg/api/panel"
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

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8082") // Укажите здесь ваш фронтенд-домен
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func New(ctx context.Context, port, restPort string, allowedMethods []string) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", port))
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to listen", zap.Error(err))
		return nil, err
	}
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptors.ContextWithLogger(logger.GetLoggerFromCtx(ctx)),
			interceptors.AuthInterceptor(allowedMethods, logger.GetLoggerFromCtx(ctx)),
		),
		grpc.MaxRecvMsgSize(50 * 1024 * 1024),
		grpc.MaxSendMsgSize(50 * 1024 * 1024),
	}

	grpcServer := grpc.NewServer(opts...)
	server, err := NewControlPanel(
		ctx,
		models.EnvsVars.Host_psql, models.EnvsVars.User_psql, models.EnvsVars.Password_psql, models.EnvsVars.Dbname_psql, models.EnvsVars.Port_psql,
		models.EnvsVars.Host_mysql, models.EnvsVars.User_mysql, models.EnvsVars.Password_mysql, models.EnvsVars.Port_mysql, models.EnvsVars.Dbname_mysql,
	)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to start", zap.Error(err))
	}
	panel.RegisterPanelServiceServer(grpcServer, server)

	gwmux := runtime.NewServeMux()
	err = panel.RegisterPanelServiceHandlerServer(context.Background(), gwmux, server)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to register gateway", zap.Error(err))
	}

	mux := http.NewServeMux()

	// Добавляем маршрут для статики
	mux.HandleFunc("/dashboard", gateway.Dashboardhandler)
	mux.HandleFunc("/server/{id}", gateway.ServerHandler)

	// Обрабатываем всё остальное через gRPC Gateway
	mux.Handle("/", enableCORS(gwmux))
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", restPort),
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
