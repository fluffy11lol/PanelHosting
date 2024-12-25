package main

import (
	"context"
	"fmt"
	"storage-microservice/internal/config"
	"storage-microservice/internal/transport/rest"
	"storage-microservice/pkg/logger"

	"log"

	"os"
	"os/signal"
	"syscall"
)

const (
	serviceName = "storage-microservice"
)

func main() {
	ctx := context.Background()
	mainLogger := logger.New(serviceName)
	ctx = context.WithValue(ctx, logger.LoggerKey, mainLogger)
	cfg := config.New()
	if cfg == nil {
		panic("failed to load config")
	}

	client, err := rest.NewS3Client(cfg.EndPoint, cfg.AccessKeyID, cfg.SecretAccessKey)
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	err = rest.EnsureBucketExists(client, cfg.BucketName)
	if err != nil {
		log.Fatalf("Failed to ensure bucket exists: %v", err)
	}
	fmt.Printf("Bucket '%s' is ready\n", cfg.BucketName)
	restServer, err := rest.New(ctx, cfg.RestServerPort, client, cfg.BucketName)
	if err != nil {
		mainLogger.Error(ctx, err.Error())
		return
	}

	graceCh := make(chan os.Signal, 1)
	signal.Notify(graceCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := restServer.Start(ctx); err != nil {
			mainLogger.Error(ctx, err.Error())
		}
	}()

	<-graceCh

	if err := restServer.Stop(ctx); err != nil {
		mainLogger.Error(ctx, err.Error())
	}
	mainLogger.Info(ctx, "Server stopped")
}
