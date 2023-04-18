package api

import (
	"github.com/casbin/casbin/v2"
	v1 "github.com/microservice/api_gateway/api/handlers/v1"
	jwthandler "github.com/microservice/api_gateway/api/token"

	"github.com/microservice/api_gateway/api/middileware"
	"github.com/microservice/api_gateway/config"
	"github.com/microservice/api_gateway/pkg/logger"
	"github.com/microservice/api_gateway/services"
	"github.com/microservice/api_gateway/storage/repo"

	//"github.com/gin-contrib/cors"
	_ "github.com/microservice/api_gateway/api/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

// Option ...
type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	RedisRepo      repo.RedisRepo
	CasbinEnforcer *casbin.Enforcer
}

// New ...
// @title           			Swagger for user api
// @securityDefinitions.apikey 	ApiKeyAuth
// @in header
// @name Authorization
// @version        				1.0
// @description     			This is a user service api.
// @Host localhost:8080
func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
		Redis:          option.RedisRepo,
		CasbinEnforcer: option.CasbinEnforcer,
	})

	

	jwt := jwthandler.JWTHandler{
		SigninKEY: option.Conf.SignKey,
		Log:       option.Logger,
	}

	router.Use(middileware.NewAuth(option.CasbinEnforcer, jwt, config.Load()))
	
	api := router.Group("/v1")
	// users
	api.POST("/users", handlerV1.CreateUser)
	api.GET("/user/:id", handlerV1.GetUserById)
	api.GET("/users", handlerV1.GetAllUsers)
	api.GET("/users/:search", handlerV1.SearchUsers)
	api.PUT("/user/:id", handlerV1.UpdateUser)
	api.DELETE("/user/:id", handlerV1.DeleteUser)

	// register
	api.POST("/users/register", handlerV1.Register)
	api.GET("/verify/:email/:code", handlerV1.Verify)
	api.GET("/login/:email/:password", handlerV1.Login)

	// posts
	api.POST("/post", handlerV1.CreatePost)
	api.GET("/post/:id", handlerV1.GetPostById)
	// api.GET("/posts/:id", handlerV1.GetAllPosts)
	api.GET("/post/search", handlerV1.SearchPost)
	api.PUT("/post/:id", handlerV1.UpdatePost)
	api.DELETE("/post/:id", handlerV1.DeletePost)

	// comment
	api.POST("/comment", handlerV1.WriteComment)
	api.GET("/comments/:id", handlerV1.GetCommentsForPost)
	api.DELETE("/comment/:id", handlerV1.DeleteComment)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
