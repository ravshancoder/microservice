package repo

import (
	"context"

	c "github.com/microservice/comment_service/genproto/comment"
)

type CommentStorageI interface {
	WriteComment(context.Context, *c.CommentRequest) (*c.CommentResponse, error)
	GetComments(context.Context, *c.GetAllCommentsRequest) (*c.Comments, error)
	GetCommentsForPost(context.Context, *c.GetAllCommentsRequest) (*c.Comments, error)
	DeleteComment(context.Context, *c.IdRequest) (*c.CommentResponse, error)
}
