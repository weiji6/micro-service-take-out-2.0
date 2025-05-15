package main

import (
	"gateway/config"
	"gateway/handler"
	"gateway/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	middleware.InitLimiter()

	r := gin.Default()

	userHandler := handler.NewUserHandler()
	userHandler.Register(r)

	payHandler := handler.NewPayHandler()
	payHandler.Register(r)

	itemHandler := handler.NewItemHandler()
	itemHandler.Register(r)

	r.Run(":8080")
}
