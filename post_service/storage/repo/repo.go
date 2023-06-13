package repo

import (
	"context"

	p "github.com/microservice/post_service/genproto/post"
)

type PostStorageI interface {
	CreatePost(context.Context, *p.PostRequest) (*p.PostResponse, error)
	GetPostById(context.Context, *p.IdRequest) (*p.PostResponse, error)
	GetPostByUserId(context.Context, *p.IdRequest) (*p.Posts, error)
	GetPostForUser(context.Context, *p.IdRequest) (*p.Posts, error)
	GetPostForComment(context.Context, *p.IdRequest) (*p.PostResponse, error)
	SearchByTitle(context.Context, *p.Search) (*p.Posts, error)
	LikePost(context.Context, *p.LikeRequest) (*p.PostResponse, error)
	UpdatePost(context.Context, *p.UpdatePostRequest) error
	DeletePost(context.Context, *p.IdRequest) (*p.PostResponse, error)
}
