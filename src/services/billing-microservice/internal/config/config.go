package config

import (
	"billing-microservice/pkg/db/cache"
	"billing-microservice/pkg/db/postgres"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	postgres.Config
	cache.RedisConfig

	GRPCServerPort int `env:"GRPC_SERVER_PORT" env-default:"9090"`
	RestServerPort int `env:"REST_SERVER_PORT" env-default:"8080"`
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
