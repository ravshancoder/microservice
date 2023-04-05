package main

import (
	"github.com/project/api_gateway/api"
	"github.com/project/api_gateway/config"
	"github.com/project/api_gateway/pkg/logger"
	"github.com/project/api_gateway/services"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "api_gateway")

	serviceManager, err := services.NewServiceManager(&cfg)
	if err != nil {
		log.Error("gRPC dial error: ", logger.Error(err))
	}

	server := api.New(api.Option{
		Conf:           cfg,
		ServiceManager: serviceManager,
		Logger:         log,
	})

	if err := server.Run(cfg.HTTPPort); err != nil {
		log.Fatal("failed to run HTTP server: ", logger.Error(err))
		panic(err)
	}
}
