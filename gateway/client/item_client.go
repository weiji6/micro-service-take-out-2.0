package client

import (
	"gateway/internal/service"
	"google.golang.org/grpc"
	"log"
)

func NewItemClient() service.ItemServiceClient {
	conn, err := grpc.Dial("0.0.0.0:30003", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("无法连接到Item服务:%v", err)
	}

	return service.NewItemServiceClient(conn)
}
