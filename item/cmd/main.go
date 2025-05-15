package main

import (
	"item/config"
	"item/discovery"
	"item/internal/handler"
	"item/internal/repository"
	"item/internal/service"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func main() {
	config.InitConfig()

	db, err := repository.InitDB()
	if err != nil {
		panic("数据库连接失败:" + err.Error())
	}

	// 创建 ItemRepository 实例
	itemRepository := repository.NewItemRepositoryImpl(db)

	// etcd 地址
	etcdAddress := []string{viper.GetString("etcd.address")}

	// 服务的注册
	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	grpcRegisterAddress := viper.GetString("service.grpcRegisterAddress")
	userNode := discovery.Server{
		Name:    viper.GetString("service.domain"),
		Address: grpcRegisterAddress,
	}

	server := grpc.NewServer()
	defer server.Stop()

	// 绑定服务
	service.RegisterItemServiceServer(server, handler.NewItemHandler(itemRepository))

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
