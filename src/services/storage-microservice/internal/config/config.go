package config

import (
	"fmt"

	"storage-microservice/pkg/db/cache"
	"storage-microservice/pkg/db/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.Config
	cache.RedisConfig
	EndPoint        string `env:"ENDPOINT" envDefault:"localhost:6379"`
	AccessKeyID     string `env:"ACCESS_KEY_ID" envDefault:""`
	SecretAccessKey string `env:"ACCESS_KEY" envDefault:""`
	BucketName      string `env:"BUCKET_NAME" envDefault:""`
	GRPCServerPort  int    `env:"GRPC_SERVER_PORT" env-default:"9090"`
	RestServerPort  int    `env:"REST_SERVER_PORT" env-default:"8085"`
	BasePath        string `env:"BASE_PATH" env-default:"./projects"`
}

func New() *Config {
	cfg := Config{}
	err := cleanenv.ReadConfig("./configs/local.env", &cfg)
	fmt.Println(err)
	if err != nil {
		return nil
	}
	return &cfg
}
