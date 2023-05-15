package repo

import (
	p "github.com/microservice/post_service/genproto/post"
)

type PostStorageI interface {
	CreatePost(*p.PostRequest) (*p.PostResponse, error)
	GetPostById(*p.IdRequest) (*p.PostResponse, error)
	GetPostByUserId(*p.IdRequest) (*p.Posts, error)
	GetPostForUser(*p.IdRequest) (*p.Posts, error)
	GetPostForComment(*p.IdRequest) (*p.PostResponse, error)
	SearchByTitle(*p.Search) (*p.Posts, error)
	LikePost(*p.LikeRequest) (*p.PostResponse, error)
	UpdatePost(*p.UpdatePostRequest) error
	DeletePost(*p.IdRequest) (*p.PostResponse, error)
}
