package main

import (
	"gateway/config"
	"gateway/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()

	r := gin.Default()

	userHandler := handler.NewUserHandler()
	userHandler.Register(r)

	payHandler := handler.NewPayHandler()
	payHandler.Register(r)

	r.Run(":8080")
}
