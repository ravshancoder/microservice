package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/microservice/api_gateway/api/handlers/models"
	"github.com/microservice/api_gateway/genproto/user"
)

// @Summary		Add Policy User
// @Description Add Policy for user
// @Tags 		Sudo
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		policy body models.Policy true "Policy"
// @Succes 		200 {object} user.Empty
// @Router 		/v1/admin/add/policy [POST]
func (h *handlerV1) AddPolicyUser(c *gin.Context) {

	body := models.Policy{}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		fmt.Println(err)
	}

	ok, err := h.casbinEnforcer.AddPolicy(body.User, body.Domain, body.Action)
	if err != nil {
		fmt.Println(err)
	}
	h.casbinEnforcer.SavePolicy()
	fmt.Println(ok)
	c.JSON(http.StatusOK, user.Empty{})
}

// @Summary		Remove Policy User
// @Description Remove User Policy
// @Tags 		Sudo
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		policy body models.Policy true "Policy"
// @Succes 		200 {object} user.Empty
// @Router 		/v1/admin/remove/policy [POST]
func (h *handlerV1) RemovePolicyUser(c *gin.Context) {
	body := models.Policy{}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		fmt.Println(err)
	}

	ok, err := h.casbinEnforcer.RemovePolicy(body.User, body.Domain, body.Action)
	if err != nil {
		fmt.Println(err)
	}
	h.casbinEnforcer.SavePolicy()
	fmt.Println(ok)
	c.JSON(http.StatusOK, user.Empty{})
}


// @Summary		Add Role User
// @Description Add User Role
// @Tags 		Sudo
// @Security 	ApiKeyAuth
// @Accept 		json
// @Produce 	json
// @Param 		policy body models.Policy true "Policy"
// @Succes 		200 {object} user.Empty
// @Router 		/v1/admin/add/role [POST]
func (h *handlerV1) AddRoleUser(c *gin.Context) {

	body := models.Policy{}

	err := c.ShouldBindJSON(&body)
	if err != nil {
		fmt.Println(err)
	}

	ok, err := h.casbinEnforcer.AddRoleForUser(body.User, body.Domain, body.Action)
	if err != nil {
		fmt.Println(err)
	}
	h.casbinEnforcer.SavePolicy()
	fmt.Println(ok)
	c.JSON(http.StatusOK, user.Empty{})
}
