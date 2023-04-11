package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/microservice/api_gateway/api/handlers/models"
	"github.com/microservice/api_gateway/genproto/comment"
	l "github.com/microservice/api_gateway/pkg/logger"
	
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

// @Summary Write comment
// @Description This api write comment
// @Tags Comment
// @Accept json
// @Produce json
// @Param body body models.CommentRequest true "Write Comment"
// @Success 200 {object} models.Comment
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /comment [post]
func (h *handlerV1) WriteComment(c *gin.Context) {
	var (
		body        models.CommentRequest
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

	response, err := h.serviceManager.CommentService().WriteComment(context.Background(), &comment.CommentRequest{
		PostId: body.PostId,
		UserId: body.UserId,
		Text:   body.Text,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create comment", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// // @Summary update Comment
// // @Description This api updates a comment
// // @Tags Comment
// // @Accept json
// // @Produce json
// // @Param body body models.UpdateCommentRequest true "Update Comment"
// // @Success 200 {object} models.Comment
// // @Failure 400 {object} models.StandartErrorModel
// // @Failure 404 {object} models.StandartErrorModel
// // @Failure 500 {object} models.StandartErrorModel
// // @Router /comment/{id} [put]
// func (h *handlerV1) UpdateComment(c *gin.Context) {
// 	var (
// 		body        models.UpdateCommentRequest
// 		jspbMarshal protojson.MarshalOptions
// 	)
// 	jspbMarshal.UseProtoNames = true
// 	err := c.ShouldBindJSON(&body)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("Failed to bind json: ", l.Error(err))
// 		return
// 	}
// 	response, err := h.serviceManager.CommentService().(context.Background(), &pu.UpdateUserRequest{
// 		Id:        body.Id,
// 		FirstName: body.FirstName,
// 		LastName:  body.LastName,
// 		Email:     body.Email,
// 	})
// 	if err != nil {
// 		if status.Code(err) == codes.NotFound {
// 			c.JSON(http.StatusNotFound, gin.H{
// 				"error": "user not found",
// 			})
// 			h.log.Error("user not found", l.Error(err))
// 			return
// 		}
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to update user", l.Error(err))
// 		return
// 	}
// 	c.JSON(http.StatusOK, &models.User{
// 		Id:        response.Id,
// 		FirstName: response.FirstName,
// 		LastName:  response.LastName,
// 		Email:     response.Email,
// 	})
// }

// @Summary get comment by id
// @Description This api gets a comment by id
// @Tags Comment
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 200 {object} models.Comment
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /comment/{id} [get]
func (h *handlerV1) GetCommentsForPost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to get comment by id: ", l.Error(err))
		return
	}

	response, err := h.serviceManager.CommentService().GetCommentsForPost(context.Background(), &comment.GetAllCommentsRequest{PostId: int64(id)})
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

// // @Summary get all users
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
// func (h *handlerV1) GetAllComments(c *gin.Context) {
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

// @Summary delete comment
// @Description This api deletes a comment
// @Tags Comment
// @Accept json
// @Produce json
// @Param id path int true "Id"
// @Success 200 {object} models.Comment
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /comment/{id} [delete]
func (h *handlerV1) DeleteComment(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid comment ID",
		})
		h.log.Error("Failed to parse comment ID: ", l.Error(err))
		return
	}

	response, err := h.serviceManager.CommentService().DeleteComment(context.Background(), &comment.IdRequest{Id: id})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete comment", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// // @Summary search users by name
// // @Description This api searches for users by first name
// // @Tags User
// // @Accept json
// // @Produce json
// // @Param first_name query string true "FirstName"
// // @Success 200 {object}  models.Users
// // @Failure 400 {object} models.StandartErrorModel
// // @Failure 500 {object} models.StandartErrorModel
// // @Router /users/search [get]
// func (h *handlerV1) SearchComment(c *gin.Context) {
// 	queryParams := c.Request.URL.Query()
// 	params, strerr := utils.ParseQueryParams(queryParams)
// 	if strerr != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": strerr[0],
// 		})
// 		h.log.Error("Failed to get all users: " + strerr[0])
// 		return
// 	}
// 	response, err := h.serviceManager.UserService().SearchUsersByName(context.Background(), &pu.SearchUsers{
// 		Search: params.Search,
// 	})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("Failed to search users: ", l.Error(err))
// 		return
// 	}
// 	c.JSON(http.StatusOK, response)
// }
