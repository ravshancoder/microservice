package repo

import (
	"context"
	u "github.com/microservice/user_service/genproto/user"
)

type UserStoreI interface {
	CreateUser(context.Context, *u.UserRequest) (*u.UserResponse, error)
	GetUserById(context.Context, *u.IdRequest) (*u.UserResponse, error)
	GetUserForClient(*u.IdRequest) (*u.UserResponse, error)
	GetAllUsers(*u.AllUsersRequest) (*u.Users, error)
	SearchUsersByName(*u.SearchUsers) (*u.Users, error)
	UpdateUser(*u.UpdateUserRequest) error
	DeleteUser(*u.IdRequest) (*u.UserResponse, error)
	CheckFiedld(*u.CheckFieldReq) (*u.CheckFieldRes, error)
	GetByEmail(*u.EmailReq) (*u.UserResponse, error)
	UpdateToken(*u.RequestForTokens) (*u.UserResponse, error)
	GetUserIdByToken(*u.Token) (*u.IdResp, error)
}
