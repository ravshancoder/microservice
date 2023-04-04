package repo

import (
	p "najottalim/6_part_microservice/service/post_service/genproto/post"
)

type PostStorageI interface {
	CreatePost(*p.PostRequest) (*p.GetPostResponse, error)
	GetPostById(*p.IdRequest) (*p.GetPostResponse, error)
	GetPostByUserId(*p.IdRequest) (*p.Posts, error)
	GetPostForUser(*p.IdRequest) (*p.Posts, error)
	GetPostForComment(*p.IdRequest) (*p.GetPostResponse, error)
	SearchByTitle(*p.Title) (*p.Posts, error)
	LikePost(*p.LikeRequest) (*p.GetPostResponse, error)
	UpdatePost(*p.UpdatePostRequest) error
	DeletePost(*p.IdRequest) (*p.GetPostResponse, error)
}
