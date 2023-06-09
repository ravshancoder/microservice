package v1

import (
	"context"
	"fmt"
	// "fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/microservice/api_gateway/api/docs"
	"github.com/microservice/api_gateway/api/handlers/models"
	pu "github.com/microservice/api_gateway/genproto/user"
	"github.com/microservice/api_gateway/pkg/etc"
	l "github.com/microservice/api_gateway/pkg/logger"
)

// Login user
// @Summary login user api
// @Description this api login user
// @Tags Register
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param        email  	path string true "email"
// @Param        password   path string true "password"
// @Succes       200		{object}	models.LoginUser
// Failure       500        {object}  models.Error
// Failure       400        {object}  models.Error
// Failure       404        {object}  models.Error
// @Router /v1/login/{email}/{password} [get]
func (h *handlerV1) Login(c *gin.Context) {
	var (
		email    = c.Param("email")
		password = c.Param("password")
	)
	// fmt.Println("Password: ", password, "	Email: ", email)

	res, err := h.serviceManager.UserService().GetByEmail(
		context.Background(), &pu.EmailReq{
			Email: email,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error GetByEmail login func", l.Error(err))
		return
	}
	

	if !etc.CheckPasswordHash(password, res.Password) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Passwrod or Email error",
		})
		return
	}

	h.jwtHandler.Iss = "user"
	h.jwtHandler.Sub = res.Id
	fmt.Println(res.UserType)
	if res.UserType == "owner" {
		h.jwtHandler.Role = "owner"
	} else if res.UserType == "admin" {
		h.jwtHandler.Role = "admin"
	} else if res.UserType == "user"{
		h.jwtHandler.Role = "user"
	} else {
		h.jwtHandler.Role = "unauthorized"
	}
	h.jwtHandler.Aud = []string{"microservice"}
	h.jwtHandler.SiginKey = h.cfg.SiginKey
	h.jwtHandler.Log = h.log
	tokens, err := h.jwtHandler.GenerateAuthJWT()
	accessToken := tokens[0]
	refreshToken := tokens[1]
	
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
		RefreshToken: refreshToken,
	}

	// Update tokens
	res, err = h.serviceManager.UserService().UpdateToken(context.Background(), &pu.RequestForTokens{Id: ucReq.Id, RefreshToken: ucReq.RefreshToken})
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
	response.AccesToken = accessToken
	response.Refreshtoken = refreshToken
	c.JSON(http.StatusOK, response)
}
