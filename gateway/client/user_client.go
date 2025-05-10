package client

import (
	"gateway/internal/service"
	"google.golang.org/grpc"
	"log"
)

func NewUserClient() service.UserServiceClient {
	conn, err := grpc.Dial("0.0.0.0:30001", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("无法连接到User服务:%v", err)
	}

	return service.NewUserServiceClient(conn)
}
