package handler

import (
	"context"
	"gateway/client"
	"gateway/internal/service"
	"gateway/middleware"
	"gateway/pkg/request"
	"gateway/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userClient service.UserServiceClient
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userClient: client.NewUserClient(),
	}
}

func (u *UserHandler) Register(r *gin.Engine) {
	r.POST("/register", u.CreateUser)
}

func (u *UserHandler) CreateUser(c *gin.Context) {
	if !middleware.Allow() {
		c.JSON(http.StatusTooManyRequests, response.Response{Code: 429, Message: "请求过于频繁，请稍后再试"})
		return
	}

	var user request.RegisterRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "无效的请求参数"})
		return
	}

	req := &service.RegisterRequest{
		UserId:   int32(user.UserID),
		UserName: user.UserName,
	}

	resp, err := u.userClient.Register(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "注册失败" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: resp.Message})
}
