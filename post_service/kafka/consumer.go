package kafka

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/microservice/post_service/config"
	handler "github.com/microservice/post_service/kafka/handler"
	"github.com/microservice/post_service/pkg/logger"
	"github.com/microservice/post_service/pkg/messagebroker"
	"github.com/microservice/post_service/storage"

	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	KafkaConsumer *kafka.Reader
	KafkaHandler  handler.KafkaHandler
	log           logger.Logger
}

func NewKafkaConsumer(db *sqlx.DB, conf *config.Config, log logger.Logger, topic string) messagebroker.Consumer {
	connString := "post_service:8030"
	return &KafkaConsumer{
		KafkaConsumer: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{connString},
			Topic:    topic,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
		KafkaHandler: *handler.NewKafkaHandlerFunc(*conf, storage.NewStoragePg(db), log),
		log:          log,
	}
}

func (k KafkaConsumer) Start() {
	fmt.Println("Consumer started--> ")
	for {
		m, err := k.KafkaConsumer.ReadMessage(context.Background())
		fmt.Println("master of copy paste here", err, m)
		if err != nil {
			k.log.Error("Error on consuming a message:", logger.Error(err))
			break
		}
		err = k.KafkaHandler.Handle(m.Value)
		if err != nil {
			k.log.Error("failed to handle consumed topic: ",
				logger.String("on topic", m.Topic), logger.Error(err))
		} else {
			fmt.Println()
			k.log.Info("Succesfully consumed message",
				logger.String("on topic", m.Topic),
				logger.String("message ", "succes"))
			fmt.Println()
		}
	}
	err := k.KafkaConsumer.Close()
	if err != nil {
		k.log.Error("Error on closing consumer: ", logger.Error(err))
	}
}
