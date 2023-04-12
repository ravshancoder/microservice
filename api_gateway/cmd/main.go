package main

import (
	"github.com/gomodule/redigo/redis"
	"github.com/microservice/api_gateway/api"
	"github.com/microservice/api_gateway/config"
	"github.com/microservice/api_gateway/pkg/logger"
	"github.com/microservice/api_gateway/services"
	r "github.com/microservice/api_gateway/storage/redis"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api_gateway")

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error: ", logger.Error(err))
	}

	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cfg.RedisHost+":"+cfg.RedisPort)
		},
	}

	server := api.New(api.Option{
		Conf:           cfg,
		ServiceManager: serviceManager,
		RedisRepo:      r.NewRedisRepo(pool),
		Logger:         log,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run HTTP server: ", logger.Error(err))
		panic(err)
	}
}
