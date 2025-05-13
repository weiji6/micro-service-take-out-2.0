package client

import (
	"gateway/internal/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
)

func NewPayClient() service.PayServiceClient {
	payServiceAddress := viper.GetString("service.payServiceAddress")

	conn, err := grpc.Dial(payServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("无法连接到Pay服务:%v", err)
	}

	return service.NewPayServiceClient(conn)
}
