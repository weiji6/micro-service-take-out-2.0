package main

import (
	"gateway/handler"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	userHandler := handler.NewUserHandler()
	userHandler.Register(r)

	payHandler := handler.NewPayHandler()
	payHandler.Register(r)

	r.Run(":8080")
}
