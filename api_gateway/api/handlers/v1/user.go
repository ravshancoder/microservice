package v1

import (
	"context"
	"net/http"
	"strconv"

	pu "github.com/project/api_gateway/genproto/user"
	l "github.com/project/api_gateway/pkg/logger"
	"github.com/project/api_gateway/pkg/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// route /v1/users [POST]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		body        pu.UserRequest
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	response, err := h.serviceManager.UserService().CreateUser(context.Background(), &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// route /v1/users/{id} [GET]
func (h *handlerV1) GetUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to convert id to int", l.Error(err))
		return
	}

	response, err := h.serviceManager.UserService().GetUserById(context.Background(), &pu.IdRequest{Id: int64(intId)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user by id", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// route /v1/users [GET]
func (h *handlerV1) GetAllUsers(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		h.log.Error("failed to parse query params to json: " + errStr[0])
		return
	}

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	response, err := h.serviceManager.UserService().GetAllUsers(context.Background(), &pu.AllUsersRequest{Limit: params.Limit, Page: params.Page})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get all users", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// route /v1/users/{id} [PUT]
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var (
		body        pu.UpdateUserRequest
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind JSON", l.Error(err))
		return
	}

	newID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to convert id to int", l.Error(err))
	}
	body.Id = int64(newID)

	response, err := h.serviceManager.UserService().UpdateUser(context.Background(), &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// route /v1/users/{id} [DELETE]
func (h *handlerV1) DeleteUser(c *gin.Context) {
	jspbMarshal := protojson.MarshalOptions{}
	jspbMarshal.UseProtoNames = true

	newID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to convert id to int: ", l.Error(err))
		return
	}

	response, err := h.serviceManager.UserService().DeleteUser(context.Background(), &pu.IdRequest{Id: int64(newID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
