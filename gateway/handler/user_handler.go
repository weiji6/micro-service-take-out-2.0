package handler

import (
	"context"
	"gateway/client"
	"gateway/internal/service"
	"gateway/pkg/request"
	"gateway/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
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
