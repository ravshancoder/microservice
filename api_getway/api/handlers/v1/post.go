package v1

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"

	models "template/api-gateway/api/handlers/models"
	p "template/api-gateway/genproto/post_service"
	l "template/api-gateway/pkg/logger"
	"template/api-gateway/pkg/utils"
)

// CreatePost ...
// @Summary CreatePost
// @Description Api for creating a new post
// @Tags post
// @Accept json
// @Produce json
// @Param User body models.Post true "createPostModel"
// @Success 200 {object} models.Post
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/users/ [post]
func (h *handlerV1) CreatePost(c *gin.Context) {
	var (
		body        models.Post
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.PostService().Create(ctx, &p.Post{
		Id:        body.Id,
		OwnerId:   body.OwnerId,
		Name:      body.Name,
		CreatedAt: body.CreatedAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create post", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// GetPost gets post by id
// @Summary GetPost
// @Description Api for getting post by id
// @Tags post
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} models.Post
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/posts/{id} [get]
func (h *handlerV1) GetPost(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")
	intID, err := strconv.Atoi(id)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.PostService().Get(
		ctx, &p.GetPostRequest{
			Id: int64(intID),
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get post", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// ListUsers returns list of users
// route /v1/users/ [get]
func (h *handlerV1) GetPostsByUserId(c *gin.Context) {
	queryParams := c.Request.URL.Query()

	params, errStr := utils.ParseQueryParams(queryParams)
	if errStr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errStr[0],
		})
		h.log.Error("failed to parse query params json" + errStr[0])
		return
	}

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().GetAllUsers(
		ctx, &p.GetPostsByUserId{
			Limit: params.Limit,
			Page:  params.Page,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to list users", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUser updates user by id
// route /v1/users/{id} [put]
func (h *handlerV1) UpdatePost(c *gin.Context) {
	// var (
	// 	body        pb.User
	// 	jspbMarshal protojson.MarshalOptions
	// )
	// jspbMarshal.UseProtoNames = true

	// err := c.ShouldBindJSON(&body)
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	h.log.Error("failed to bind json", l.Error(err))
	// 	return
	// }
	// body.Id = c.Param("id")

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	// defer cancel()

	// response, err := h.serviceManager.UserService().Upda(ctx, &body)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	h.log.Error("failed to update user", l.Error(err))
	// 	return
	// }

	// c.JSON(http.StatusOK, response)
}

// DeleteUser deletes user by id
// route /v1/users/{id} [delete]
func (h *handlerV1) DeletePost(c *gin.Context) {
	// var jspbMarshal protojson.MarshalOptions
	// jspbMarshal.UseProtoNames = true

	// guid := c.Param("id")
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	// defer cancel()

	// response, err := h.serviceManager.UserService().Delete(
	// 	ctx, &pb.ByIdReq{
	// 		Id: guid,
	// 	})
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	h.log.Error("failed to delete user", l.Error(err))
	// 	return
	// }

	// c.JSON(http.StatusOK, response)
}
