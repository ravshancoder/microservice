package v1

import (
	"context"
	"fmt"
	"net/http"

	"github.com/microservice/api_gateway/api/handlers/models"
	pu "github.com/microservice/api_gateway/genproto/post"
	l "github.com/microservice/api_gateway/pkg/logger"
	"github.com/microservice/api_gateway/pkg/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Summary create post
// @Description This api creates a post
// @Tags Post
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body models.PostRequest true "Create Post"
// @Success 200 {object} models.Post
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/post [post]
func (h *handlerV1) CreatePost(c *gin.Context) {

	fmt.Println(c.GetHeader("Authorization"))
	var (
		body        models.PostRequest
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

	response, err := h.serviceManager.PostService().CreatePost(context.Background(), &pu.PostRequest{
		Title:       body.Title,
		Description: body.Description,
		UserId:      body.UserId,
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

// @Summary update post
// @Description This api updates a post
// @Tags Post
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body models.UpdatePostRequest true "Update Post"
// @Success 200 {object} models.Post
// @Failure 400 {object} models.StandartErrorModel
// @Failure 404 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/posts/{id} [put]
func (h *handlerV1) UpdatePost(c *gin.Context) {
	var (
		body        models.UpdatePostRequest
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

	response, err := h.serviceManager.PostService().UpdatePost(context.Background(), &pu.UpdatePostRequest{
		Id:          body.Id,
		Title:       body.Title,
		Description: body.Description,
	})

	if err != nil {
		if status.Code(err) == codes.NotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Post not found",
			})
			h.log.Error("Post not found", l.Error(err))
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update Post", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, &models.Post{
		Id:          response.Id,
		Title:       response.Title,
		Description: response.Description,
	})
}

// @Summary get post by id
// @Description This api gets a post by id
// @Tags Post
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 200 {object} models.Post
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/post/{id} [get]
func (h *handlerV1) GetPostById(c *gin.Context) {
	id := c.Param("id")

	response, err := h.serviceManager.PostService().GetPostById(context.Background(), &pu.IdRequest{Id: id})
	if err != nil {
		statusCode := http.StatusInternalServerError
		if status.Code(err) == codes.NotFound {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get post by id: ", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}


// @Summary delete post
// @Description This api deletes a post
// @Tags Post
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Id"
// @Success 200 {object} models.Post
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/post/{id} [delete]
func (h *handlerV1) DeletePost(c *gin.Context) {
	id := c.Param("id")
	
	response, err := h.serviceManager.PostService().DeletePost(context.Background(), &pu.IdRequest{
		Id: id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete post", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary search users by name
// @Description This api searches for users by first name
// @Tags Post
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param search query string true "search"
// @Success 200 {object}  models.Users
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/post/search [get]
func (h *handlerV1) SearchPost(c *gin.Context) {

	queryParams := c.Request.URL.Query()
	params, strerr := utils.ParseQueryParams(queryParams)
	
	if strerr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": strerr[0],
		})
		h.log.Error("Failed to get all users: " + strerr[0])
		return
	}

	response, err := h.serviceManager.PostService().SearchByTitle(context.Background(), &pu.Search{
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

// // @Summary get all posts
// // @Description This api gets all users
// // @Tags User
// // @Accept json
// // @Produce json
// // @Param limit query int true "Limit"
// // @Param page query int true "Page"
// // @Success 200 {object} []models.User
// // @Failure 400 {object} models.StandartErrorModel
// // @Failure 500 {object} models.StandartErrorModel
// // @Router /users [get]
// func (h *handlerV1) GetAllUsers(c *gin.Context) {
// 	queryParams := c.Request.URL.Query()
// 	params, errstr := utils.ParseQueryParams(queryParams)
// 	if errstr != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": errstr[0],
// 		})
// 		h.log.Error("Failed to get all users: " + errstr[0])
// 		return
// 	}
// 	response, err := h.serviceManager.UserService().GetAllUsers(context.Background(), &pu.AllUsersRequest{
// 		Limit: params.Limit,
// 		Page:  params.Page,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("Failed to get all users: ", l.Error(err))
// 		return
// 	}
// 	c.JSON(http.StatusOK, response)
// }