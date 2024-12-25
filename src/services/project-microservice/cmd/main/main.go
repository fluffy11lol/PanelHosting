package main

import (
	"context"
	"project-microservice/internal"
	"project-microservice/internal/config"
	"project-microservice/internal/repository"
	"project-microservice/internal/service"
	"project-microservice/internal/transport/grpc"
	"project-microservice/pkg/db/cache"
	"project-microservice/pkg/db/postgres"
	"project-microservice/pkg/logger"

	"os"
	"os/signal"
	"syscall"
)

const (
	serviceName = "project-microservice"
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

	redis := cache.New(cfg.RedisConfig)
	repo := repository.NewProjectRepository(db.Db, redis)

	srv := service.NewProjectService(repo)

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
