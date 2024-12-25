package main

import (
	"billing-microservice/internal"
	"billing-microservice/internal/config"
	"billing-microservice/internal/repository"
	"billing-microservice/internal/service"
	"billing-microservice/internal/transport/grpc"
	"billing-microservice/pkg/db/postgres"
	"billing-microservice/pkg/logger"
	"context"

	"os"
	"os/signal"
	"syscall"
)

const (
	serviceName = "billing-microservice"
)

func main() {
	ctx := context.Background()
	mainLogger := logger.New(serviceName)
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
	internal.InsertData(ctx, db.Db)
	repo := repository.NewBillingRepository(db.Db)

	srv := service.NewBillingService(repo)

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
