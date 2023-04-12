package v1

import (
	"github.com/microservice/api_gateway/api/handlers/token"
	"github.com/microservice/api_gateway/config"
	"github.com/microservice/api_gateway/pkg/logger"
	"github.com/microservice/api_gateway/services"
	"github.com/microservice/api_gateway/storage/repo"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	redis          repo.RedisRepo
	jwtHandler     token.JWTHandler
}

type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Redis          repo.RedisRepo
	JwtHandler     token.JWTHandler
}

func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		redis:          c.Redis,
		jwtHandler:     c.JwtHandler,
	}
}
