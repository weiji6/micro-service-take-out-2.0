package client

import (
	"gateway/internal/service"
	"log"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewItemClient() service.ItemServiceClient {
	itemServiceAddress := viper.GetString("service.itemServiceAddress")

	conn, err := grpc.Dial(itemServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("无法连接到Item服务:%v", err)
	}

	return service.NewItemServiceClient(conn)
}
