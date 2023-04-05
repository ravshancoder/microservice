package api

import (
	v1 "github.com/project/api_gateway/api/handlers/v1"
	"github.com/project/api_gateway/config"
	"github.com/project/api_gateway/pkg/logger"
	"github.com/project/api_gateway/services"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
}


// @title user
// @version v1
// @description user service
// @host localhost:8080
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
	api.POST("/user", handlerV1.CreateUser)
	api.GET("/user/:id", handlerV1.GetUser)
	api.GET("/users", handlerV1.GetAllUsers)
	api.PUT("/user/:id", handlerV1.UpdateUser)
	api.DELETE("/user/:id", handlerV1.DeleteUser)

	// posts
	api.POST("/post", handlerV1.CreatePost)
	api.GET("/post/:id", handlerV1.GetPost)
	api.GET("/posts/:id", handlerV1.GetAllPosts)
	api.PUT("/post/:id", handlerV1.UpdatePost)
	api.DELETE("/post/:id", handlerV1.DeletePost)

	// comment
	api.POST("/comment", handlerV1.WriteComment)
	api.GET("/comments/:id", handlerV1.GetComments)
	api.DELETE("/comment/:id", handlerV1.DeleteComment)

	// swagger
	url := ginSwagger.URL("swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
