package service

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	c "github.com/project/user_service/genproto/comment"
	p "github.com/project/user_service/genproto/post"
	u "github.com/project/user_service/genproto/user"
	"github.com/project/user_service/pkg/logger"
	grpcclient "github.com/project/user_service/service/grpc_client"
	"github.com/project/user_service/storage"
)

type UserService struct {
	storage storage.IStorage
	Logger  logger.Logger
	Client  grpcclient.Clients
}

func NewUserService(db *sqlx.DB, log logger.Logger, client grpcclient.Clients) *UserService {
	return &UserService{
		storage: storage.NewStoragePg(db),
		Logger:  log,
		Client:  client,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *u.UserRequest) (*u.UserResponse, error) {
	res, err := s.storage.User().CreateUser(req)
	if err != nil {
		log.Println("failed to creating user: ", err)
		return &u.UserResponse{}, err
	}

	return res, nil
}

func (s *UserService) GetUserById(ctx context.Context, req *u.IdRequest) (*u.UserResponse, error) {
	res, err := s.storage.User().GetUserById(req)
	if err != nil {
		log.Println("failed to getting user: ", err)
		return &u.UserResponse{}, err
	}

	post, err := s.Client.Post().GetPostForUser(ctx, &p.IdRequest{Id: req.Id})
	if err != nil {
		log.Println("failed to getting post for get user: ", err)
		return &u.UserResponse{}, err
	}

	for _, pt := range post.Posts {
		comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.GetAllCommentsRequest{PostId: pt.Id})
		if err != nil {
			log.Println("failed to get comments for post in user service: ", err)
			return &u.UserResponse{}, err
		}

		pst := u.Post{}

		for _, comment := range comments.Comments {
			comUser, err := s.storage.User().GetUserById(&u.IdRequest{Id: comment.UserId})
			if err != nil {
				log.Println("failed to get comment user: ", err)
				return &u.UserResponse{}, err
			}

			com := u.Comment{}
			com.UserId = comUser.Id
			com.UserName = comUser.FirstName + " " + comUser.LastName
			com.PostId = comment.PostId
			com.PostTitle = pt.Title
			com.PostUserName = res.FirstName + " " + res.LastName
			com.Text = comment.Text
			com.CreatedAt = comment.CreatedAt

			pst.Comments = append(pst.Comments, &com)
		}
		pst.Id = pt.Id
		pst.Title = pt.Title
		pst.Description = pt.Description
		pst.Likes = pt.Likes
		pst.CreatedAt = pt.CreatedAt
		pst.UpdatedAt = pt.UpdatedAt

		res.Posts = append(res.Posts, &pst)
	}

	return res, nil
}

func (s *UserService) GetUserForClient(ctx context.Context, req *u.IdRequest) (*u.UserResponse, error) {
	res, err := s.storage.User().GetUserById(req)
	if err != nil {
		log.Println("failed to getting user for clients: ", err)
		return &u.UserResponse{}, err
	}

	return res, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, req *u.AllUsersRequest) (*u.Users, error) {
	res, err := s.storage.User().GetAllUsers(req)
	if err != nil {
		log.Println("failed to getting all users: ", err)
		return &u.Users{}, err
	}

	for _, user := range res.Users {
		post, err := s.Client.Post().GetPostForUser(ctx, &p.IdRequest{Id: user.Id})
		if err != nil {
			log.Println("failed to getting post for get all users: ", err)
			return &u.Users{}, err
		}

		for _, pt := range post.Posts {
			comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.GetAllCommentsRequest{PostId: pt.Id})
			if err != nil {
				log.Println("failed to get comments for post in user service: ", err)
				return &u.Users{}, err
			}

			pst := u.Post{}

			for _, comment := range comments.Comments {
				comUser, err := s.storage.User().GetUserById(&u.IdRequest{Id: comment.UserId})
				if err != nil {
					log.Println("failed to get user comment: ", err)
					return &u.Users{}, err
				}

				com := u.Comment{}
				com.UserId = comment.UserId
				com.UserName = comUser.FirstName + " " + comUser.LastName
				com.PostId = comment.PostId
				com.PostTitle = pt.Title
				com.PostUserName = user.FirstName + " " + user.LastName
				com.Text = comment.Text
				com.CreatedAt = comment.CreatedAt

				pst.Comments = append(pst.Comments, &com)
			}
			pst.Id = pt.Id
			pst.Title = pt.Title
			pst.Description = pt.Description
			pst.Likes = pt.Likes
			pst.CreatedAt = pt.CreatedAt
			pst.UpdatedAt = pt.UpdatedAt

			user.Posts = append(user.Posts, &pst)
		}
	}

	return res, nil
}

func (s *UserService) SearchUsersByName(ctx context.Context, req *u.SearchUsers) (*u.Users, error) {
	res, err := s.storage.User().SearchUsersByName(req)
	if err != nil {
		log.Println("failed to searching user by name: ", err)
		return &u.Users{}, err
	}

	for _, user := range res.Users {
		post, err := s.Client.Post().GetPostForUser(ctx, &p.IdRequest{Id: user.Id})
		if err != nil {
			log.Println("failed to getting post for searching user by name: ", err)
			return &u.Users{}, err
		}

		for _, pt := range post.Posts {
			comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.GetAllCommentsRequest{PostId: pt.Id})
			if err != nil {
				log.Println("failed to get comments for post in user service: ", err)
				return &u.Users{}, err
			}

			pst := u.Post{}

			for _, comment := range comments.Comments {
				comUser, err := s.storage.User().GetUserById(&u.IdRequest{Id: comment.UserId})
				if err != nil {
					log.Println("failed to get user comment: ", err)
					return &u.Users{}, err
				}

				com := u.Comment{}
				com.UserId = comment.UserId
				com.UserName = comUser.FirstName + " " + comUser.LastName
				com.PostId = comment.PostId
				com.PostTitle = pt.Title
				com.PostUserName = user.FirstName + " " + user.LastName
				com.Text = comment.Text
				com.CreatedAt = comment.CreatedAt

				pst.Comments = append(pst.Comments, &com)
			}
			pst.Id = pt.Id
			pst.Title = pt.Title
			pst.Description = pt.Description
			pst.Likes = pt.Likes
			pst.CreatedAt = pt.CreatedAt
			pst.UpdatedAt = pt.UpdatedAt

			user.Posts = append(user.Posts, &pst)
		}
	}

	return res, nil
}

func (s *UserService) UpdateUser(ctx context.Context, req *u.UpdateUserRequest) (*u.UserResponse, error) {
	err := s.storage.User().UpdateUser(req)
	if err != nil {
		log.Println("failed to updating user: ", err)
		return &u.UserResponse{}, err
	}

	return &u.UserResponse{}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *u.IdRequest) (*u.UserResponse, error) {
	res, err := s.storage.User().DeleteUser(req)
	if err != nil {
		log.Println("failed to deleting user: ", err)
		return &u.UserResponse{}, err
	}

	post, err := s.Client.Post().GetPostForUser(ctx, &p.IdRequest{Id: req.Id})
	if err != nil {
		log.Println("failed to getting post for deleting user: ", err)
		return &u.UserResponse{}, err
	}

	for _, pt := range post.Posts {
		comments, err := s.Client.Comment().GetCommentsForPost(ctx, &c.GetAllCommentsRequest{PostId: pt.Id})
		if err != nil {
			log.Println("failed to get comments for post in user service: ", err)
			return &u.UserResponse{}, err
		}

		pst := u.Post{}

		for _, comment := range comments.Comments {
			comUser, err := s.storage.User().GetUserById(&u.IdRequest{Id: comment.UserId})
			if err != nil {
				log.Println("failed to get user comment: ", err)
				return &u.UserResponse{}, err
			}

			com := u.Comment{}
			com.UserId = comment.UserId
			com.UserName = comUser.FirstName + " " + comUser.LastName
			com.PostId = comment.PostId
			com.PostTitle = pt.Title
			com.UserName = res.FirstName + " " + res.LastName
			com.Text = comment.Text
			com.CreatedAt = comment.CreatedAt

			pst.Comments = append(pst.Comments, &com)
		}
		pst.Id = pt.Id
		pst.Title = pt.Title
		pst.Description = pt.Description
		pst.Likes = pt.Likes
		pst.CreatedAt = pt.CreatedAt
		pst.UpdatedAt = pt.UpdatedAt

		res.Posts = append(res.Posts, &pst)
	}

	return res, nil
}

func (s *UserService) CheckField(ctx context.Context, req *u.CheckFieldReq) (*u.CheckFieldRes, error) {
	res, err := s.storage.User().CheckFiedld(req)
	if err != nil {
		s.Logger.Error("error delete", logger.Any("Error delete users", err))
		return &u.CheckFieldRes{}, status.Error(codes.Internal, "something went wrong, please check user info")
	}
	return res, nil
}

func (s *UserService) GetByEmail(ctx context.Context, req *u.EmailReq) (*u.UserResponse, error) {
	customer, err := s.storage.User().GetByEmail(req)
	if err != nil {
		s.Logger.Error("Error while getting customer info by email", logger.Error(err))
		return nil, status.Error(codes.InvalidArgument, "Something went wrong")
	}
	return customer, nil
}

func (s *UserService) UpdateToken(ctx context.Context, req *u.RequestForTokens) (*u.UserResponse, error) {
	res, err := s.storage.User().UpdateToken(req)
	if err != nil {
		log.Println("failed to updating user: ", err)
		return &u.UserResponse{}, err
	}

	return res, err
}
