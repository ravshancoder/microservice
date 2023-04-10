package v1

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/project/api_gateway/api/handlers/models"
	"github.com/project/api_gateway/pkg/etc"
	l "github.com/project/api_gateway/pkg/logger"
)

// Register register user
// @Summary register user api
// @Description this api registers
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.RegisterUserModel true "register user"
// @Success 200{object} models.User
// @Router /v1/register[post]
func (h *handlerV1) Register(c *gin.Context) {

	var body models.RegisterModel

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("eror while bind json", l.Error(err))
		return
	}
	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)
	body.Password, err = etc.HashPassword(body.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseError{
			Error: err,
		})
		h.log.Error("couldn't hash the password")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	existsEmail, err := h.serviceManager.UserService().CheckField(ctx, &user.CheckFieldReq{
		Field: "email",
		Value: body.Email,
	})
}
