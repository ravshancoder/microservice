package v1

import (
	"github.com/project/api_gateway/config"
	"github.com/project/api_gateway/pkg/logger"
	"github.com/project/api_gateway/services"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
}

type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
}

func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
	}
}

// func errorResponse(err error) *models.ErrorResponse {
// 	return &models.ErrorResponse{
// 		Error: err.Error(),
// 	}
// }
