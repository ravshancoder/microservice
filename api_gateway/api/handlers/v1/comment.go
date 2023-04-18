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
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body models.CommentRequest true "Write Comment"
// @Success 200 {object} models.Comment
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/comment [post]
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

// @Summary get comments for post
// @Description This api gets a comment for post
// @Tags Comment
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 200 {object} models.Comment
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/comment/{id} [get]
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

// @Summary delete comment
// @Description This api deletes a comment
// @Tags Comment
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id path int true "Id"
// @Success 200 {object} models.Comment
// @Failure 400 {object} models.StandartErrorModel
// @Failure 500 {object} models.StandartErrorModel
// @Router /v1/comment/{id} [delete]
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
