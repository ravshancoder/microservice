package v1

import (
	"context"
	"net/http"
	"strconv"

	pp "github.com/project/api_gateway/genproto/post"
	l "github.com/project/api_gateway/pkg/logger"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// route /v1/post [POST]
func (h *handlerV1) CreatePost(c *gin.Context) {
	var (
		body        pp.PostRequest
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

	response, err := h.serviceManager.PostService().CreatePost(context.Background(), &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create post", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// route /v1/post/{id} [GET]
func (h *handlerV1) GetPost(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to convert id to int", l.Error(err))
		return
	}

	response, err := h.serviceManager.PostService().GetPostById(context.Background(), &pp.IdRequest{Id: int64(intID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get post by id", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// route /v1/posts/{id} [GET]
func (h *handlerV1) GetAllPosts(c *gin.Context) {
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

	response, err := h.serviceManager.PostService().GetPostByUserId(context.Background(), &pp.IdRequest{Id: int64(intId)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get posts by user id", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// route /v1/post/{id} [PUT]
func (h *handlerV1) UpdatePost(c *gin.Context) {
	var (
		body        pp.UpdatePostRequest
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

	response, err := h.serviceManager.PostService().UpdatePost(context.Background(), &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update post", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// route /v1/post/{id} [DELETE]
func (h *handlerV1) DeletePost(c *gin.Context) {
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

	response, err := h.serviceManager.PostService().DeletePost(context.Background(), &pp.IdRequest{Id: int64(newID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete post", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
