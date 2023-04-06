// package v1

// import (
// 	"context"
// 	"net/http"
// 	"strconv"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/project/api_gateway/api/handlers/models"
// 	pu "github.com/project/api_gateway/genproto/user"
// 	l "github.com/project/api_gateway/pkg/logger"
// 	"github.com/project/api_gateway/pkg/utils"
// 	"google.golang.org/protobuf/encoding/protojson"
// )

// // @Summary Create a user
// // @Description Create a user
// // @Tags User
// // @Accept json
// // @Produce json
// // @Param body body models.UserRequest true "Create User"
// // @Success 200 {object} models.User
// // @Failure 400 {object} models.StandartErrorModel
// // @Failure 500 {object} models.StandartErrorModel
// // @Router /v1/users [post]
// func (h *handlerV1) CreateUser(c *gin.Context) {
// 	var (
// 		body        models.UserRequest
// 		jspbMarshal protojson.MarshalOptions
// 	)
// 	jspbMarshal.UseProtoNames = true

// 	err := c.ShouldBindJSON(&body)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to bind json", l.Error(err))
// 		return
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	response, err := h.serviceManager.UserService().CreateUser(ctx, &pu.UserRequest{
// 		FirstName: body.FirstName,
// 		LastName:  body.LastName,
// 		Email:     body.Email,
// 	})

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to create user", l.Error(err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, &models.User{
// 		Id:        response.Id,
// 		FirstName: response.FirstName,
// 		LastName:  response.LastName,
// 		Email:     response.Email,
// 	})
// }

// // @Summary get user by id
// // @Description This api gets a user by id
// // @Tags User
// // @Accept json
// // @Produce json
// // @Param id path int true "User ID"
// // @Success 200 {object} models.User
// // @Failure 400 {object} models.StandartErrorModel
// // @Failure 500 {object} models.StandartErrorModel
// // @Router /v1/users/{id} [get]
// func (h *handlerV1) GetUser(c *gin.Context) {
// 	var jspbMarshal protojson.MarshalOptions
// 	jspbMarshal.UseProtoNames = true

// 	id := c.Param("id")
// 	intId, err := strconv.Atoi(id)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to convert id to int", l.Error(err))
// 		return
// 	}

// 	response, err := h.serviceManager.UserService().GetUserById(context.Background(), &pu.IdRequest{Id: int64(intId)})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to get user by id", l.Error(err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// // route /v1/users [GET]
// func (h *handlerV1) GetAllUsers(c *gin.Context) {
// 	queryParams := c.Request.URL.Query()

// 	params, errStr := utils.ParseQueryParams(queryParams)
// 	if errStr != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": errStr[0],
// 		})
// 		h.log.Error("failed to parse query params to json: " + errStr[0])
// 		return
// 	}

// 	var jspbMarshal protojson.MarshalOptions
// 	jspbMarshal.UseProtoNames = true

// 	response, err := h.serviceManager.UserService().GetAllUsers(context.Background(), &pu.AllUsersRequest{Limit: params.Limit, Page: params.Page})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to get all users", l.Error(err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// // route /v1/users/{id} [PUT]
// func (h *handlerV1) UpdateUser(c *gin.Context) {
// 	var (
// 		body        pu.UpdateUserRequest
// 		jspbMarshal protojson.MarshalOptions
// 	)
// 	jspbMarshal.UseProtoNames = true

// 	err := c.ShouldBindJSON(&body)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to bind JSON", l.Error(err))
// 		return
// 	}

// 	newID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to convert id to int", l.Error(err))
// 	}
// 	body.Id = int64(newID)

// 	response, err := h.serviceManager.UserService().UpdateUser(context.Background(), &body)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to update user", l.Error(err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

// // route /v1/users/{id} [DELETE]
// func (h *handlerV1) DeleteUser(c *gin.Context) {
// 	jspbMarshal := protojson.MarshalOptions{}
// 	jspbMarshal.UseProtoNames = true

// 	newID, err := strconv.Atoi(c.Param("id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to convert id to int: ", l.Error(err))
// 		return
// 	}

// 	response, err := h.serviceManager.UserService().DeleteUser(context.Background(), &pu.IdRequest{Id: int64(newID)})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to delete user", l.Error(err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }

package v1

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/project/api_gateway/api/handlers/models"
	pu "github.com/project/api_gateway/genproto/user"
	l "github.com/project/api_gateway/pkg/logger"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Summary create user
// @Description This api creates a user
// @Tags User
// @Accept json
// @Produce json
// @Param body body models.UserRequest true "Create User"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /users [post]
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

	c.JSON(http.StatusOK, &models.User{
		Id:        response.Id,
		FirstName: response.FirstName,
		LastName:  response.LastName,
		Email:     response.Email,
	})
}

// @Summary update user
// @Description This api updates a user
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body models.UpdateUserRequest true "Update User"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandartErrorModel
// @Failure 404 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/users/{id} [put]
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var (
		body        models.UpdateUserRequest
		jspbMarshal protojson.MarshalOptions
	)

	jspbMarshal.UseProtoNames = true

	// Get user ID from URL path
	userId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user ID",
		})
		h.log.Error("Failed to parse user ID: ", l.Error(err))
		return
	}

	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to bind json: ", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().UpdateUser(ctx, &pu.UpdateUserRequest{
		Id:        userId,
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
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/users/{id} [get]
func (h *handlerV1) GetUserById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		h.log.Error("Invalid user ID: ", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().GetUserById(ctx, &pu.IdRequest{Id: id})
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

	c.JSON(http.StatusOK, &models.User{
		Id:        response.Id,
		FirstName: response.FirstName,
		LastName:  response.LastName,
		Email:     response.Email,
	})
}

// @Summary get all users
// @Description This api gets all users
// @Tags User
// @Accept json
// @Produce json
// @Param body body models.GetAllUsersRequest true "Get All User"
// @Success 200 {object} []models.User
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/users [get]
func (h *handlerV1) GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	var userreq models.GetAllUsersRequest
	response, err := h.serviceManager.UserService().GetAllUsers(ctx, &pu.AllUsersRequest{
		Page:  userreq.Page,
		Limit: userreq.Limit,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to get all users: ", l.Error(err))
		return
	}

	var users []models.User
	for _, user := range response.Users {
		users = append(users, models.User{
			Id:        user.Id,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
		})
	}

	c.JSON(http.StatusOK, users)
}

// @Summary delete user
// @Description This api deletes a user
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/users/{id} [delete]
func (h *handlerV1) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		h.log.Error("Failed to parse user ID: ", l.Error(err))
		return
	}

	
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	response, err := h.serviceManager.UserService().DeleteUser(ctx, &pu.IdRequest{
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
	})
}

// @Summary search users by name
// @Description This api searches for users by first name
// @Tags User
// @Accept json
// @Produce json
// @Param first_name query string true "First name"
// @Success 200 {object}  models.Users
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/users/search [get]
func (h *handlerV1) SearchUsers(c *gin.Context) {
	var (
		searchQuery pu.SearchUsers
	)

	if searchQuery.FirstName = c.Query("first_name"); searchQuery.FirstName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing required parameter 'first_name'",
		})
		h.log.Error("Missing required parameter 'first_name'")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().SearchUsersByName(ctx, &searchQuery)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to search users: ", l.Error(err))
		return
	}
	var users []models.User
	for _, response := range response.Users {
		user := models.User{
			Id:        response.Id,
			FirstName: response.FirstName,
			LastName:  response.LastName,
			Email:     response.Email,
		}
		users = append(users, user)
	}
	c.JSON(http.StatusOK, &models.Users{
		Users: users,
	})
}
