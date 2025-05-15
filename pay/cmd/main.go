package main

import (
	"net"
	"pay/config"
	"pay/discovery"
	"pay/internal/handler"
	"pay/internal/repository"
	"pay/internal/service"
	"pay/pkg/lock"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	config.InitConfig()
	lock.InitRedisLock()
	repository.InitDB()

	userServiceAddr := viper.GetString("service.userRegisterAddress")
	itemServiceAddr := viper.GetString("service.itemRegisterAddress")

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

	payRepo := repository.NewPayRepositoryImpl(repository.DB)
	payService := service.NewPayService(userService, itemService, payRepo)

	payHandler := handler.NewPayHandler(payService)

	etcdAddress := []string{viper.GetString("etcd.address")}

	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	grpcRegisterAddress := viper.GetString("service.payRegisterAddress")
	userNode := discovery.Server{
		Name:    viper.GetString("service.domain"),
		Address: grpcRegisterAddress,
	}

	server := grpc.NewServer()
	defer server.Stop()

	service.RegisterPayServiceServer(server, payHandler)

	grpcListenAddress := viper.GetString("service.grpcListenAddress")
	lis, err := net.Listen("tcp", grpcListenAddress)
	if err != nil {
		panic(err)
	}

	if _, err := etcdRegister.Register(userNode, 10); err != nil {
		panic(err)
	}

	if err = server.Serve(lis); err != nil {
		panic(err)
	}
}
