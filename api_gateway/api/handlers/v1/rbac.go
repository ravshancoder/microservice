package v1


// @Summary		Add Policy User
// @Description Add Policy for user
// @Tags 		Sudo
// @Security 	BareAuth
// @Accept 		json
// @Produce 	json
// @Param 		policy body models.Policy true "Policy"
// @Succes 		200 {object} user.Empty
// @Router 		/v1/admin/add/policy [POST]