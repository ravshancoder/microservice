package services

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"template/api-gateway/config"
	u "template/api-gateway/genproto/user_service"
	p "template/api-gateway/genproto/post_service"
)

type IServiceManager interface {
	UserService() u.UserServiceClient
	PostService() p.PostServiceClient
}

type serviceManager struct {
	userService u.UserServiceClient
	postService p.PostServiceClient
}

func (s *serviceManager) PostService() p.PostServiceClient {
	return s.postService
}

func (s *serviceManager) UserService() u.UserServiceClient {
	return s.userService
}

func NewServiceManager(conf *config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns")

	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.UserServiceHost, conf.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	connPost, err := grpc.Dial(
		fmt.Sprintf("%s:%d", conf.PostServiceHost, conf.PostServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	serviceManager := &serviceManager{
		userService: u.NewUserServiceClient(connUser),
		postService: p.NewPostServiceClient(connPost),
	}

	return serviceManager, nil
}
