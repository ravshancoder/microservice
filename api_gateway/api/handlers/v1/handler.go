package v1

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/microservice/api_gateway/api/handlers/models"
	"github.com/microservice/api_gateway/api/token"
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
	casbinEnforcer *casbin.Enforcer
}

type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Redis          repo.RedisRepo
	JwtHandler     token.JWTHandler
	CasbinEnforcer *casbin.Enforcer
}

func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		redis:          c.Redis,
		jwtHandler:     c.JwtHandler,
		casbinEnforcer: c.CasbinEnforcer,
	}
}

func GetClaims(h *handlerV1, c *gin.Context) jwt.MapClaims {
	var (
		ErrUnauthorized = errors.New("unauthorized")
		authorization   models.GetProfileByJWTRequestModel
		claims          jwt.MapClaims
		err             error
	)

	fmt.Println(c.GetHeader("Authorization"))

	authorization.Token = c.GetHeader("Authorization")

	if c.Request.Header.Get("Authorization") == "" {
		c.JSON(http.StatusUnauthorized, models.Error{
			Message: "Unauthorized request",
		})
		h.log.Error("Unauthorized request ", logger.Error(ErrUnauthorized))
		return nil
	}

	authorization.Token = strings.TrimSpace(strings.Trim(authorization.Token, "Barer"))

	fmt.Println(authorization.Token)

	h.jwtHandler.Token = authorization.Token

	claims, err = h.jwtHandler.ExtractClaims()

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.Error{
			Message: "Unauthorized",
		})
		h.log.Error("Unauthorized request ", logger.Error(err))
		return nil
	}
	return claims
}

type BaseHandler struct{
	Config config.Config
}