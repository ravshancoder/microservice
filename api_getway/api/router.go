package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "template/api-gateway/api/docs" // swag
	v1 "template/api-gateway/api/handlers/v1"
	"template/api-gateway/config"
	"template/api-gateway/pkg/logger"
	"template/api-gateway/services"
)

// Option ...
type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
}

// New ...
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
	})

	api := router.Group("/v1")
	// users
	api.POST("/users", handlerV1.CreateUser)
	api.GET("/users/:id", handlerV1.GetUser)
	api.GET("/users", handlerV1.ListUsers)
	api.PUT("/users/:id", handlerV1.UpdateUser)
	api.DELETE("/users/:id", handlerV1.DeleteUser)

	//post 
	api.POST("/post", handlerV1.CreatePost)
	api.GET("/post/get/:id", handlerV1.GetPost)
	api.PUT("/post/update", handlerV1.UpdatePost)
	api.DELETE("/post/delete/:id", handlerV1.DeletePost)
	api.GET("/post/list/:page/:limit/:search", handlerV)

	return router
}
