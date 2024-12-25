package grpc

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/proto"
	"net"
	"net/http"

	client "authentication-microservice/pkg/api/authentication/protobuf"
	"authentication-microservice/pkg/logger"
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
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8082")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
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
			AuthInterceptor(logger.GetLoggerFromCtx(ctx)),
			ContextWithLogger(logger.GetLoggerFromCtx(ctx)),
		),
		grpc.MaxRecvMsgSize(50 * 1024 * 1024),
		grpc.MaxSendMsgSize(50 * 1024 * 1024),
	}

	grpcServer := grpc.NewServer(opts...)
	client.RegisterUserServiceServer(grpcServer, NewUserService(service))

	restSrv := runtime.NewServeMux(
		runtime.WithForwardResponseOption(func(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
			if md, ok := runtime.ServerMetadataFromContext(ctx); ok {
				if cookies := md.HeaderMD.Get("Set-Cookie"); len(cookies) > 0 {
					for _, cookie := range cookies {
						w.Header().Add("Set-Cookie", cookie)
					}
				}
			}
			return nil
		}),
	)
	if err := client.RegisterUserServiceHandlerServer(context.Background(), restSrv, NewUserService(service)); err != nil {
		logger.GetLoggerFromCtx(ctx).Error(ctx, "failed to register rest handler", zap.Error(err))
		return nil, err
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/authentication", func(w http.ResponseWriter, r *http.Request) {

		http.ServeFile(w, r, "./static/login.html")
	})
	mux.Handle("/", enableCORS(restSrv))
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
	}
	return s.restServer.Shutdown(ctx)
}
