package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"item/config"
	"item/discovery"
	"item/internal/handler"
	"item/internal/repository"
	"item/internal/service"
	"net"
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
	grpcAddress := viper.GetString("service.grpcAddress")
	userNode := discovery.Server{
		Name:    viper.GetString("service.domain"),
		Address: grpcAddress,
	}

	server := grpc.NewServer()
	defer server.Stop()

	// 绑定服务
	service.RegisterItemServiceServer(server, handler.NewItemHandler(itemRepository))
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
}
