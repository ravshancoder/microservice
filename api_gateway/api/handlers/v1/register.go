package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/microservice/api_gateway/api/docs"
	"github.com/microservice/api_gateway/api/handlers/models"
	"github.com/microservice/api_gateway/email"
	pu "github.com/microservice/api_gateway/genproto/user"
	"github.com/microservice/api_gateway/pkg/etc"
	l "github.com/microservice/api_gateway/pkg/logger"
)

// Register register user
// @Summary register user api
// @Description this api registers
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.RegisterModel true "register user"
// @Succes 200 {object}	models.StandartErrorModel
// @Failure	500 {object} models.StandartErrorModel
// @Router /users/register [post]
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

	existsEmail, err := h.serviceManager.UserService().CheckField(context.Background(), &pu.CheckFieldReq{
		Field: "email",
		Value: body.Email,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed check username uniques ", l.Error(err))
	}

	if existsEmail.Exists {
		c.JSON(http.StatusConflict, models.ResponseError{
			Error: models.ErrorMessage{
				Message: "mail already exists",
			},
		})
		h.log.Error("this email already exists ", l.Error(err))
	}

	code := etc.GenerateCode(6)

	msg := "Subject: User email verification\n Your verification code: " + code
	err = email.SendEmail([]string{body.Email}, []byte(msg))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"eroor": err.Error(),
		})
		return
	}

	body.Code = code

	userBodyBte, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while marshal user body", l.Error(err))
	}

	fmt.Println(userBodyBte)

	err = h.redis.SetWithTTL(body.Email, string(userBodyBte), 86400)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while redis user body", l.Error(err))
	}

	c.JSON(http.StatusAccepted, code)
}
