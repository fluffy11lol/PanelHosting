package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	models "control-panel/internal/models"
	transport "control-panel/internal/transport/grpc"
	"control-panel/pkg/api/logger"
)

const (
	serviceName = "control-panel"
)

func main() {
	ctx := context.Background()
	mainLogger := logger.New(serviceName)
	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)
	models.LoadEnvs()

	// Создаем gRPC сервер
	allowedMethods := []string{"/Login"}
	grpcserver, err := transport.New(ctx, models.EnvsVars.Grpc_port, models.EnvsVars.Rest_port, allowedMethods)
	if err != nil {
		mainLogger.Error(ctx, err.Error())
		return
	}

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := grpcserver.Start(ctx); err != nil {
			mainLogger.Error(ctx, err.Error())
		}
	}()

	<-graceCh

	if err := grpcserver.Stop(ctx); err != nil {
		mainLogger.Error(ctx, err.Error())
	}
	mainLogger.Info(ctx, "Server stopped")
}
