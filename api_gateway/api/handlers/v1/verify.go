package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	r "github.com/gomodule/redigo/redis"

	"github.com/gin-gonic/gin"
	
	"github.com/google/uuid"
	"github.com/microservice/api_gateway/api/handlers/models"
	pu "github.com/microservice/api_gateway/genproto/user"
	"github.com/microservice/api_gateway/pkg/etc"
	l "github.com/microservice/api_gateway/pkg/logger"
)

// Verify verify user
// @Summary verify user api
// @Description this api verifies
// @Tags Register
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param email path string true "email"
// @Param code path string true "code"
// @Succes 200{object} models.RegisterModel
// @Router /v1/verify/{email}/{code} [get]
func (h *handlerV1) Verify(c *gin.Context) {
	var (
		email = c.Param("email")
		code  = c.Param("code")
		body  = models.UserRegister{}
	)

	userBody, err := h.redis.Get(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error get from redis user body", l.Error(err))
		return
	}

	byteData, err := r.String(userBody, err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error get from redis user body", l.Error(err))
		return
	}

	err = json.Unmarshal([]byte(byteData), &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while unmarshalling user data", l.Error(err))
		return
	}

	if body.Code != code {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while checking code ", l.Error(err))
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while generating UUID", l.Error(err))
		return
	}
	h.jwtHandler.SiginKey = h.cfg.SiginKey
	h.jwtHandler.Sub = id.String()
	h.jwtHandler.Iss = "user"
	h.jwtHandler.Role = "authorized"
	h.jwtHandler.Aud = []string{"ucook-frontend"}
	h.jwtHandler.Log = h.log
	tokens, err := h.jwtHandler.GenerateAuthJWT()
	accessToken := tokens[0]
	refreshToken := tokens[1]
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while generating tokens", l.Error(err))
		return
	}

	hashedPassword, err := etc.HashPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while generating hash for password", l.Error(err))
		return
	}
	checkEmail, err := h.serviceManager.UserService().CheckField(context.Background(), &pu.CheckFieldReq{
		Field: "email",
		Value: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to check uniqueness of email", l.Error(err))
		return
	}

	if checkEmail.Exists {
		c.JSON(http.StatusConflict, models.StandartErrorModel{
			Error: models.Error{
				Message: "mail already exists",
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	fmt.Println(body)
	user, err := h.serviceManager.UserService().CreateUser(ctx, &pu.UserRequest{
		Email:     body.Email,
		Password:  hashedPassword,
		FirstName: body.FirstName,
		LastName:  body.LastName,

		RefreshToken: refreshToken,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while creating user to db", l.Error(err))
		return
	}
	h.jwtHandler.Sub = user.Id
	user.AccesToken = accessToken
	user.RefreshToken = refreshToken

	c.JSON(http.StatusCreated, user)
}
