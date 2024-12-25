package main

import (
	"authentication-microservice/internal"
	"authentication-microservice/internal/config"
	"authentication-microservice/internal/repository"
	"authentication-microservice/internal/service"
	"authentication-microservice/internal/transport/grpc"
	"authentication-microservice/pkg/db/postgres"
	"authentication-microservice/pkg/logger"
	"context"

	"os"
	"os/signal"
	"syscall"
)

const (
	serviceName = "authentication-microservice"
)

func main() {
	ctx := context.Background()
	mainLogger := logger.New(serviceName)
	if mainLogger == nil {
		panic("failed to create logger")
	}
	defer mainLogger.Sync()
	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)
	cfg := config.New()
	if cfg == nil {
		panic("failed to load config")
	}

	db, err := postgres.New(cfg.Config)
	if err != nil {
		mainLogger.Error(ctx, err.Error())
		panic(err)
	}
	internal.RunMigrations(db.Db)
	repo := repository.NewUserRepository(db.Db)

	srv := service.NewUserService(repo)

	grpcserver, err := grpc.New(ctx, cfg.GRPCServerPort, cfg.RestServerPort, srv)
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
