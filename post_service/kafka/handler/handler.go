package kafka

import (
	"fmt"

	"github.com/microservice/post_service/config"
	"github.com/microservice/post_service/pkg/logger"
	"github.com/microservice/post_service/storage"

	pb "github.com/microservice/post_service/genproto/post"
)

type KafkaHandler struct {
	config  config.Config
	storage storage.IStorage
	log     logger.Logger
}

func NewKafkaHandlerFunc(config config.Config, storage storage.IStorage, log logger.Logger) *KafkaHandler {
	return &KafkaHandler{
		config:  config,
		storage: storage,
		log:     log,
	}
}

func (h *KafkaHandler) Handle(value []byte) error {
	post := pb.PostRequest{}
	err := post.Unmarshal(value)
	if err != nil {
		return err
	}
	fmt.Println("aaaaaaaaa")
	fmt.Println(post)
	_, err = h.storage.Post().CreatePost(&pb.PostRequest{
		UserId:      post.UserId,
		Title: post.Title,
		Description: post.Description,
	})
	if err != nil {
		return err
	}
	return nil
}
