package main

import (
	"fmt"
	"net"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go"

	"github.com/microservice/user_service/config"
	u "github.com/microservice/user_service/genproto/user"
	"github.com/microservice/user_service/kafka"
	"github.com/microservice/user_service/pkg/db"
	"github.com/microservice/user_service/pkg/logger"
	"github.com/microservice/user_service/pkg/messagebroker"
	"github.com/microservice/user_service/service"

	grpcclient "github.com/microservice/user_service/service/grpc_client"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 10,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "jaeger:6831",
		},
	}

	closer, err := conf.InitGlobalTracer(
		"user-service",
	)
	if err != nil {
		fmt.Println(err)
	}
	defer closer.Close()

	cfg := config.Load()
	log := logger.New(cfg.LogLevel, "golang")
	defer logger.Cleanup(log)

	connDb, err := db.ConnectToDB(cfg)
	if err != nil {
		fmt.Println("failed connect database", err)
	}

	produceMap := make(map[string]messagebroker.Producer)
	topic := "topic1"
	userTopicProduce := kafka.NewKafkaProducer(cfg, log, topic)
	defer func() {
		err := userTopicProduce.Stop()
		if err != nil {
			log.Fatal("Failed to stopping Kafka", logger.Error(err))
		}
	}()
	produceMap["user"] = userTopicProduce

	grpcClient, err := grpcclient.New(cfg)
	if err != nil {
		fmt.Println("failed while grpc client", err.Error())
	}

	userService := service.NewUserService(connDb, log, grpcClient)

	lis, err := net.Listen("tcp", cfg.UserServicePort)
	if err != nil {
		log.Fatal("failed while listening port: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)
	u.RegisterUserServiceServer(s, userService)

	log.Info("main: server running",
		logger.String("port", cfg.UserServicePort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed while listening: %v", logger.Error(err))
	}
}
