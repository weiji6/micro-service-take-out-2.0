package handler

import (
	"context"
	"errors"
	"fmt"
	"user/internal/repository"
	"user/internal/service"

	"gorm.io/gorm"
)

type UserHandler struct {
	service.UnimplementedUserServiceServer
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// 实现 Register 方法
func (u *UserHandler) Register(ctx context.Context, req *service.RegisterRequest) (*service.RegisterResponse, error) {
	userRepo := repository.User{}
	if userRepo.CheckUserExist(req) {
		return &service.RegisterResponse{
			Code:    400,
			Message: "用户已存在",
		}, errors.New("用户已存在")
	}

	_, err := userRepo.CreateUser(req)
	if err != nil {
		return &service.RegisterResponse{
			Code:    500,
			Message: "服务器错误",
		}, err
	}

	return &service.RegisterResponse{
		Code:    200,
		Message: "注册成功",
	}, nil
}

// 实现 CheckBalance 方法
func (u *UserHandler) CheckBalance(ctx context.Context, req *service.CheckBalanceRequest) (*service.CheckBalanceResponse, error) {
	if repository.DB == nil {
		return &service.CheckBalanceResponse{
			Code:    500,
			Message: "数据库未初始化",
		}, errors.New("数据库未初始化")
	}

	var user repository.User
	result := repository.DB.First(&user, req.UserId)
	if result.Error != nil {
		return &service.CheckBalanceResponse{
			Code:    500,
			Message: "查询用户余额失败: " + result.Error.Error(),
		}, result.Error
	}

	return &service.CheckBalanceResponse{
		Code:    200,
		Message: "查询成功",
		Balance: float32(user.Balance),
	}, nil
}

// 实现 DecreaseBalance 方法
func (u *UserHandler) DecreaseBalance(ctx context.Context, req *service.DecreaseBalanceRequest) (*service.DecreaseBalanceResponse, error) {
	if repository.DB == nil {
		return &service.DecreaseBalanceResponse{
			Code:    500,
			Message: "数据库未初始化",
		}, errors.New("数据库未初始化")
	}

	var user repository.User
	result := repository.DB.First(&user, req.UserId)
	if result.Error != nil {
		return &service.DecreaseBalanceResponse{
			Code:    500,
			Message: "查询用户失败: " + result.Error.Error(),
		}, result.Error
	}

	if user.Balance < float64(req.Amount) {
		return &service.DecreaseBalanceResponse{
			Code:    400,
			Message: "余额不足",
		}, errors.New("余额不足")
	}

	// 使用事务来确保数据一致性
	err := repository.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&user).Update("balance", user.Balance-float64(req.Amount)).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return &service.DecreaseBalanceResponse{
			Code:    500,
			Message: "扣款失败: " + err.Error(),
		}, err
	}

	return &service.DecreaseBalanceResponse{
		Code:    200,
		Message: "扣款成功",
	}, nil
}

func (u *UserHandler) DecreaseBalanceRevert(ctx context.Context, req *service.DecreaseBalanceRequest) (*service.DecreaseBalanceResponse, error) {
	fmt.Printf("[User] 回滚扣款，用户ID: %d 金额: %.2f\n", req.UserId, req.Amount)

	err := repository.DB.Model(&repository.User{}).
		Where("id = ?", req.UserId).
		Update("balance", gorm.Expr("balance + ?", req.Amount)).Error

	if err != nil {
		return &service.DecreaseBalanceResponse{Code: 500, Message: "回滚失败: " + err.Error()}, err
	}

	return &service.DecreaseBalanceResponse{Code: 200, Message: "余额回滚成功"}, nil
}
