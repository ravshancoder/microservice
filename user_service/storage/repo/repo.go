package repo

import (
	u "najottalim/6_part_microservice/service/user_service/genproto/user"
)

type UserStoreI interface {
	CreateUser(*u.UserRequest) (*u.UserResponse, error)
	GetUserById(*u.IdRequest) (*u.UserResponse, error)
	GetUserForClient(*u.IdRequest) (*u.UserResponse, error)
	GetAllUsers(*u.AllUsersRequest) (*u.Users, error)
	SearchUsersByName(*u.SearchUsers) (*u.Users, error)
	UpdateUser(*u.UpdateUserRequest) error
	DeleteUser(*u.IdRequest) (*u.UserResponse, error)
}
