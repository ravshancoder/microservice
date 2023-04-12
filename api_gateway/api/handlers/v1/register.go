package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	r "github.com/gomodule/redigo/redis"
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

	var body models.UserRegister

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

	err = h.redis.SetWithTTL(body.Email, string(userBodyBte), 86400)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while redis user body", l.Error(err))
	}

	c.JSON(http.StatusAccepted, code)
}

// Verify verify user
// @Summary verify user api
// @Description this api verifies
// @Tags Register
// @Accept json
// @Produce json
// @Param email path string true "email"
// @Param code path string true "code"
// @Succes 200{object} models.RegisterModel
// @Router /verify/{email}/{code} [get]
func (h *handlerV1) Verify(c *gin.Context) {
	var (
		code  = c.Param("code")
		email = c.Param("email")
		body  = models.UserRegister{}
	)
	fmt.Println(email)
	userBody, err := h.redis.Get(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error get from redis user body", l.Error(err))
		return
	}

	byteData, err := r.String(userBody, err)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error get from redis user body", l.Error(err))
		return
	}

	err = json.Unmarshal([]byte(byteData), &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while unmarshalling user data", l.Error(err))
		return
	}

	if body.Code != code {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while checking code ", l.Error(err))
		return
	}

	hashePassword, err := etc.GeneratePasswordHash(body.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while generating hash for password", l.Error(err))
		return
	}
	checkEmail, err := h.serviceManager.UserService().CheckField(context.Background(), &pu.CheckFieldReq{
		Field: "email",
		Value: body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("Failed to check uniqueness of email", l.Error(err))
		return
	}

	if checkEmail.Exists {
		c.JSON(http.StatusConflict, models.StandartErrorModel{
			Error: models.Error{
				Message: "mail already exists",
			},
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	fmt.Println(body)
	user, err := h.serviceManager.UserService().CreateUser(ctx, &pu.UserRequest{
		Email:     body.Email,
		Password:  string(hashePassword),
		FirstName: body.FirstName,
		LastName:  body.LastName,
	})
	user.Email = body.Email
	user.FirstName = body.FirstName
	user.LastName = body.LastName
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while creating user to db", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, user)
}
