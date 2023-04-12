package v1

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/microservice/api_gateway/api/docs"
	"github.com/microservice/api_gateway/api/handlers/models"
	"github.com/microservice/api_gateway/api/handlers/token"
	pu "github.com/microservice/api_gateway/genproto/user"
	l "github.com/microservice/api_gateway/pkg/logger"
)

// Login user
// @Summary login user api
// @Description this api login user
// @Tags Register
// @Accept json
// @Produce json
// @Param        email  	path string true "email"
// @Param        password   path string true "password"
// @Succes        200		{object}	models.LoginUser
// Failure       500        {object}  models.Error
// Failure       400        {object}  models.Error
// Failure       404        {object}  models.Error
// @Router /v1/login/{email}/{password} [get]
func (h *handlerV1) Login(c *gin.Context) {
	var (
		password = c.Param("password")
		email    = c.Param("email")
	)
	res, err := h.serviceManager.UserService().Login(
		context.Background(), &pu.LoginRequest{
			Email:    email,
			Password: password,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error login token", l.Error(err))
		return
	}

	h.jwtHandler = token.JWTHandler{
		SigninKEY: h.cfg.SigninKey,
		Sub:       res.Id,
		Iss:       "user",
		Role:      "authorized",
		Aud: []string{
			"ucook_frontend",
		},
		Log: h.log,
	}

	accesToken, refreshToken, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error generating token", l.Error(err))
		return
	}
	// keys for update tokens
	ucReq := &pu.RequestForTokens{
		Id:           res.Id,
		AccesToken:   accesToken,
		RefreshToken: refreshToken,
	}

	// Update tokens
	res, err = h.serviceManager.UserService().UpdateToken(context.Background(), &pu.RequestForTokens{Id: ucReq.Id, RefreshToken: ucReq.RefreshToken, AccesToken: ucReq.AccesToken})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while updating client token", l.Error(err))
		return
	}

	// Just forresponse
	response := &models.LoginUser{
		Id:        res.Id,
		Email:     res.Email,
		FirstName: res.FirstName,
		LastName:  res.LastName,
	}
	response.Refreshtoken = refreshToken
	response.Accessestoken = accesToken
	c.JSON(http.StatusOK, response)
}
