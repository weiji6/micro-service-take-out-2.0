package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"net"
	"pay/config"
	"pay/discovery"
	"pay/internal/handler"
	"pay/internal/repository"
	"pay/internal/service"
)

func main() {
	config.InitConfig()

	userServiceAddr := viper.GetString("service.userAddress")
	itemServiceAddr := viper.GetString("service.itemAddress")

	connUser, err := grpc.Dial(userServiceAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	connItem, err := grpc.Dial(itemServiceAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer connUser.Close()
	defer connItem.Close()

	userService := service.NewUserServiceClient(connUser)
	itemService := service.NewItemServiceClient(connItem)

	payService := service.NewPayService(userService, itemService)

	payHandler := handler.NewPayHandler(payService)

	etcdAddress := []string{viper.GetString("etcd.address")}

	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	grpcAddress := viper.GetString("service.grpcAddress")
	userNode := discovery.Server{
		Name:    viper.GetString("service.domain"),
		Address: grpcAddress,
	}

	server := grpc.NewServer()
	defer server.Stop()

	service.RegisterPayServiceServer(server, payHandler)

	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}

	if _, err := etcdRegister.Register(userNode, 10); err != nil {
		panic(err)
	}

	if err = server.Serve(lis); err != nil {
		panic(err)
	}

	repository.InitDB()
}
