package client

import (
	"gateway/internal/service"
	"log"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewUserClient() service.UserServiceClient {
	userServiceAddress := viper.GetString("service.userServiceAddress")

	conn, err := grpc.Dial(userServiceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("无法连接到User服务:%v", err)
	}

	return service.NewUserServiceClient(conn)
}
