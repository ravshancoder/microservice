package postgres

import (
	"fmt"
	"log"
	"time"

	u "github.com/project/user_service/genproto/user"
)

func (r *UserRepo) CreateUser(user *u.UserRequest) (*u.UserResponse, error) {
	var res u.UserResponse
	err := r.db.QueryRow(`
		insert into 
			users(first_name, last_name, email, password, acces_token, refresh_token) 
		values
			($1, $2, $3, $4, $5, $6) 
		returning 
			id, first_name, last_name, email, acces_token, refresh_token, created_at, updated_at`, user.FirstName, user.LastName, user.Email, user.Password, user.AccesToken, user.RefreshToken).
		Scan(
			&res.Id,
			&res.FirstName,
			&res.LastName,
			&res.Email,
			&res.AccesToken,
			&res.RefreshToken,
			&res.CreatedAt,
			&res.UpdatedAt,
		)

	if err != nil {
		log.Println("failed to create user")
		return &u.UserResponse{}, err
	}

	return &res, nil
}

func (r *UserRepo) GetUserById(user *u.IdRequest) (*u.UserResponse, error) {
	var res u.UserResponse
	err := r.db.QueryRow(`
		select 
			id, first_name, last_name, email, created_at, updated_at
		from 
			users 
		where id = $1 and deleted_at is null`, user.Id).
		Scan(
			&res.Id,
			&res.FirstName,
			&res.LastName,
			&res.Email,
			&res.CreatedAt,
			&res.UpdatedAt,
		)

	if err != nil {
		log.Println("failed to get user")
		return &u.UserResponse{}, err
	}

	return &res, nil
}

func (r *UserRepo) GetUserForClient(user_id *u.IdRequest) (*u.UserResponse, error) {
	var res u.UserResponse
	err := r.db.QueryRow(`
		select 
			id, first_name, last_name, email, created_at, updated_at 
		from 
			users 
		where id = $1`, user_id.Id).Scan(&res.Id, &res.FirstName, &res.LastName, &res.Email, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to get user for client")
		return &u.UserResponse{}, err
	}

	return &res, nil
}

func (r *UserRepo) GetAllUsers(req *u.AllUsersRequest) (*u.Users, error) {
	var res u.Users
	offset := (req.Page - 1) * req.Limit
	rows, err := r.db.Query(`
		select 
			id, first_name, last_name, email, refresh_token, created_at, updated_at 
		from 
			users 
		where 
			deleted_at is null 
		limit $1 offset $2`, req.Limit, offset)

	if err != nil {
		log.Println("failed to get all users")
		return &u.Users{}, err
	}

	for rows.Next() {
		temp := u.UserResponse{}

		err = rows.Scan(
			&temp.Id,
			&temp.FirstName,
			&temp.LastName,
			&temp.Email,
			&temp.RefreshToken,
			&temp.CreatedAt,
			&temp.UpdatedAt,
		)
		if err != nil {
			log.Println("failed to scanning user")
			return &u.Users{}, err
		}

		res.Users = append(res.Users, &temp)
	}

	return &res, nil
}

func (r *UserRepo) SearchUsersByName(req *u.SearchUsers) (*u.Users, error) {
	var res u.Users
	query := fmt.Sprint("select id, first_name, last_name, email, created_at, updated_at from users where first_name ilike '%" + req.Search + "%' and deleted_at is null")

	rows, err := r.db.Query(query)
	if err != nil {
		log.Println("failed to searching user")
		return &u.Users{}, err
	}

	for rows.Next() {
		temp := u.UserResponse{}

		err = rows.Scan(
			&temp.Id,
			&temp.FirstName,
			&temp.LastName,
			&temp.Email,
			&temp.CreatedAt,
			&temp.UpdatedAt,
		)
		if err != nil {
			log.Println("failed to searching user")
			return &u.Users{}, err
		}

		res.Users = append(res.Users, &temp)
	}

	return &res, nil
}

func (r *UserRepo) UpdateUser(user *u.UpdateUserRequest) error {
	res, err := r.db.Exec(`
		update
			users
		set
			first_name = $1, last_name = $2, email = $3, password = $4,  updated_at = $5
		where 
			id = $6`, user.FirstName, user.LastName, user.Email, user.Password, time.Now(), user.Id)

	if err != nil {
		log.Println("failed to update user")
		return err
	}

	fmt.Println(res.RowsAffected())
	return nil
}

func (r *UserRepo) DeleteUser(user *u.IdRequest) (*u.UserResponse, error) {
	temp := u.UserResponse{}
	err := r.db.QueryRow(`
		update 
			users
		set 
			deleted_at = $1 
		where 
			id = $2 
		returning 
			id, first_name, last_name, email, created_at, updated_at`, time.Now(), user.Id).Scan(&temp.Id, &temp.FirstName, &temp.LastName, &temp.Email, &temp.CreatedAt, &temp.UpdatedAt)

	if err != nil {
		log.Println("failed to delete user")
		return &u.UserResponse{}, err
	}

	return &temp, nil
}

func (r *UserRepo) CheckFiedld(req *u.CheckFieldReq) (*u.CheckFieldRes, error) {
	query := fmt.Sprintf("SELECT 1 FROM users WHERE %s=$1", req.Field)
	res := &u.CheckFieldRes{}
	temp := -1
	err := r.db.QueryRow(query, req.Value).Scan(&temp)
	if err != nil {
		res.Exists = false
		return res, nil
	}
	if temp == 0 {
		res.Exists = true
	} else {
		res.Exists = false
	}
	return res, nil
}

func (r *UserRepo) GetByEmail(req *u.EmailReq) (*u.LoginResponse, error) {
	res := u.LoginResponse{}
	err := r.db.QueryRow(`
	SELECT 
		id, 
		first_name,
		last_name,
		email, 
		password,
		refresh_token,
		created_at, 
		updated_at 
	FROM 
		users 
	WHERE 
		email=$1 AND deleted_at IS NULL`, req.Email).
		Scan(
			&res.Id,
			&res.FirstName,
			&res.LastName,
			&res.Email,
			&res.Password,
			&res.RefreshToken,
			&res.CreatedAt,
			&res.UpdatedAt,
		)
	if err != nil {
		fmt.Println("error while getting user login")
		return &u.LoginResponse{}, err
	}

	return &res, nil
}

func (r *UserRepo) UpdateToken(user *u.RequestForTokens) (*u.LoginResponse, error) {
	res := u.LoginResponse{}
	err := r.db.QueryRow(`
		update
			users
		set
			refresh_token = $1
		where 
			id = $2
		returning 
			id, first_name, last_name, email, refresh_token, updated_at`, user.RefreshToken, user.Id).
		Scan(
			&res.Id,
			&res.FirstName,
			&res.LastName,
			&res.Email,
			&res.RefreshToken,
			&res.UpdatedAt,
		)

	if err != nil {
		log.Println("failed to update user")
		return &u.LoginResponse{}, err
	}

	return &res, err
}
