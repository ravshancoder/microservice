package kafka

import (
	"context"
	"time"

	"github.com/microservice/user_service/config"
	kafka "github.com/segmentio/kafka-go"

	// "github.com/casbin/casbin/v2/config"
	"github.com/microservice/user_service/pkg/logger"
	"github.com/microservice/user_service/pkg/messagebroker"
)

type KafkaProduce struct {
	kafkaWriter *kafka.Writer
	log         logger.Logger
}

func NewKafkaProducer(conf config.Config, log logger.Logger, topic string) messagebroker.Producer {
	connString := "kafka:9092"

	return &KafkaProduce{
		kafkaWriter: &kafka.Writer{
			Addr:         kafka.TCP(connString),
			Topic:        topic,
			BatchTimeout: time.Millisecond * 10,
		},
		log: log,
	}

}

func (p *KafkaProduce) Start() error {
	return nil
}

func (p *KafkaProduce) Stop() error {
	err := p.kafkaWriter.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *KafkaProduce) Produce(key, body []byte, logBody string) error {
	message := kafka.Message{
		Key:   key,
		Value: body,
	}

	if err := p.kafkaWriter.WriteMessages(context.Background(), message); err != nil {
		return err
	}
	return nil

}
