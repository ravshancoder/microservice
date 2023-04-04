package models

type Post struct {
	Id          int64  `json:"id"`
	OwnerId     int64  `json:"owner_id"`
	Name        string `json:"name"`
	CreatedAt   string `json:"created_at"`
	Description string `json:"description"`
}