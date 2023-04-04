package models

type User struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Post     *Post  `json:"post"`
}