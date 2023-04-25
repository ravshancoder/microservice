package models

type CommentRequest struct {
	PostId string `json:"post_id"`
	UserId string `json:"user_id"`
	Text   string `json:"text"`
}

type Comment struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type IdCommentRequest struct {
	Id string `json:"id"`
}

type GetAllCommentsRequest struct {
	Limit int64 `json:"limit"`
	Page  int64 `json:"page"`
}

type UpdateCommentRequest struct {
	Id        string  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type Comments struct {
	Comments []Comment `json:"Comments"`
}
