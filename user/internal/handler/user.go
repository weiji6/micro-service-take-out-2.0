package handler

import (
	"context"
	"user/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// Register 注册账号
// @Summary 注册账号
// @Description 用于用户注册账号
// @Tags kill
// @Accept json
// @Produce json
// @Param request body request.UserRequest true "请求参数"
// @Success 200 {object} response.Response "注册成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /takeout/register [post]
func (u *UserHandler) Register(ctx *context.Context, req *service.RegisterRequest) (resp *service.RegisterResponse, err error) {
	if err := u.service.Register(int(req.UserId), req.UserName); err != nil {
		return &service.RegisterResponse{Code: 500, Message: "注册失败"}, err
	}

	return &service.RegisterResponse{Code: 200, Message: "注册成功"}, nil
}
