package models

type PostRequest struct {
	Title string `json:"title"`
	Description  string `json:"description"`
	UserId     string `json:"user_id"`
}

type Post struct {
	Id        string  `json:"id"`
	Title string `json:"title"`
	Description  string `json:"description"`
	UserId     int64 `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type IdPostRequest struct {
	Id int64 `json:"id"`
}

type GetAllPostsRequest struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

type UpdatePostRequest struct {
	Id           string  `json:"id"`
	Title string `json:"title"`
	Description  string `json:"description"`
	UserId     string `json:"user_id"`
}

type Posts struct {
	Posts []Post `json:"users"`
}
