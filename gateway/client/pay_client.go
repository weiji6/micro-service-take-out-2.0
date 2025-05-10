package client

import (
	"gateway/internal/service"
	"google.golang.org/grpc"
	"log"
)

func NewPayClient() service.PayServiceClient {
	conn, err := grpc.Dial("0.0.0.0:30002", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("无法连接到Pay服务:%v", err)
	}

	return service.NewPayServiceClient(conn)
}
