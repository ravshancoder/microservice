package main

import (
	"fmt"
	"net"

	"github.com/microservice/post_service/config"
	p "github.com/microservice/post_service/genproto/post"
	"github.com/microservice/post_service/pkg/db"
	"github.com/microservice/post_service/pkg/logger"
	"github.com/microservice/post_service/service"
	grpcclient "github.com/microservice/post_service/service/grpc_client"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"

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

	grpcClient, err := grpcclient.New(cfg)
	if err != nil {
		fmt.Println("failed while grpc client", err.Error())
	}

	postService := service.NewPostService(connDb, log, grpcClient)

	lis, err := net.Listen("tcp", cfg.PostServicePort)
	if err != nil {
		log.Fatal("failed while listening port: %v", logger.Error(err))
	}

	s := grpc.NewServer()
	reflection.Register(s)
	p.RegisterPostServiceServer(s, postService)

	log.Info("main: server running",
		logger.String("port", cfg.PostServicePort))
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed while listening: %v", logger.Error(err))
	}

}
