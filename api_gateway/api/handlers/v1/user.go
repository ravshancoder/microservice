package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/microservice/api_gateway/api/handlers/models"
	pu "github.com/microservice/api_gateway/genproto/user"
	l "github.com/microservice/api_gateway/pkg/logger"
	"github.com/microservice/api_gateway/pkg/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Summary create user
// @Description This api creates a user
// @Tags User
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body models.UserRequest true "Create User"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/users [post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		body        models.UserRequest
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to bind json: ", l.Error(err))
		return
	}

	response, err := h.serviceManager.UserService().CreateUser(context.Background(), &pu.UserRequest{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary update user
// @Description This api updates a user
// @Tags User
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body models.UpdateUserRequest true "Update User"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandartErrorModel
// @Failure 404 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/user/{id} [put]
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var (
		body        models.UpdateUserRequest
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to bind json: ", l.Error(err))
		return
	}

	response, err := h.serviceManager.UserService().UpdateUser(context.Background(), &pu.UpdateUserRequest{
		Id:        body.Id,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
	})

	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			h.log.Error("user not found", l.Error(err))
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.User{
		Id:        response.Id,
		FirstName: response.FirstName,
		LastName:  response.LastName,
		Email:     response.Email,
	})
}

// @Summary get user by id
// @Description This api gets a user by id
// @Tags User
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/user/{id} [get]
func (h *handlerV1) GetUserById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to get user by id: ", l.Error(err))
		return
	}

	

	response, err := h.serviceManager.UserService().GetUserById(context.Background(), &pu.IdRequest{Id: int64(id)})
	if err != nil {
		statusCode := http.StatusInternalServerError
		if status.Code(err) == codes.NotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user by id: ", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary get all users
// @Description This api gets all users
// @Tags User
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param limit query int true "Limit"
// @Param page query int true "Page"
// @Success 200 {object} []models.User
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/users [get]
func (h *handlerV1) GetAllUsers(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	params, errstr := utils.ParseQueryParams(queryParams)
	if errstr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errstr[0],
		})
		h.log.Error("Failed to get all users: " + errstr[0])
		return
	}

	response, err := h.serviceManager.UserService().GetAllUsers(context.Background(), &pu.AllUsersRequest{
		Limit: params.Limit,
		Page:  params.Page,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to get all users: ", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary delete user
// @Description This api deletes a user
// @Tags User
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Id"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/user/{id} [delete]
func (h *handlerV1) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		h.log.Error("Failed to parse user ID: ", l.Error(err))
		return
	}

	response, err := h.serviceManager.UserService().DeleteUser(context.Background(), &pu.IdRequest{
		Id: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.User{
		Id:        response.Id,
		FirstName: response.FirstName,
		LastName:  response.LastName,
		Email:     response.Email,
		CreatedAt: response.CreatedAt,
		UpdatedAt: response.UpdatedAt,
	})
}

// @Summary search users by name
// @Description This api searches for users by first name
// @Tags User
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param first_name query string true "FirstName"
// @Success 200 {object}  models.Users
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/users/{search} [get]
func (h *handlerV1) SearchUsers(c *gin.Context) {

	queryParams := c.Request.URL.Query()
	params, strerr := utils.ParseQueryParams(queryParams)

	if strerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": strerr[0],
		})
		h.log.Error("Failed to get all users: " + strerr[0])
		return
	}

	response, err := h.serviceManager.UserService().SearchUsersByName(context.Background(), &pu.SearchUsers{
		Search: params.Search,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to search users: ", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
