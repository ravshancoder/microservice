package v1

import (
	"context"
	"net/http"
	"strconv"

	pc "github.com/project/api_gateway/genproto/comment"
	l "github.com/project/api_gateway/pkg/logger"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// route /v1/comment [comment]
func (h *handlerV1) WriteComment(c *gin.Context) {
	var (
		body        pc.CommentRequest
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

	response, err := h.serviceManager.CommentService().WriteComment(context.Background(), &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to write comment", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// route /v1/comment/{id} [GET]
func (h *handlerV1) GetComments(c *gin.Context) {
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

	response, err := h.serviceManager.CommentService().GetComments(context.Background(), &pc.GetAllCommentsRequest{PostId: int64(intID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get comment by post id", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// route /v1/comment/{id} [DELETE]
func (h *handlerV1) DeleteComment(c *gin.Context) {
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

	response, err := h.serviceManager.CommentService().DeleteComment(context.Background(), &pc.IdRequest{Id: int64(newID)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to delete comment", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}
