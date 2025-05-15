package handler

import (
	"context"
	"fmt"
	"pay/internal/service"
)

type PayHandler struct {
	service.UnimplementedPayServiceServer
	payService *service.PayService
}

func NewPayHandler(PayService *service.PayService) *PayHandler {
	return &PayHandler{payService: PayService}
}

func (p *PayHandler) Pay(ctx context.Context, req *service.PayRequest) (*service.PayResponse, error) {
	fmt.Printf("[Pay] 支付，用户ID: %d 商品ID: %d 金额: %.2f\n", req.UserId, req.ItemId, req.Amount)

	err := p.payService.ProcessPayment(req.UserId, req.ItemId, req.Quantity, req.Amount)
	if err != nil {
		return &service.PayResponse{
			Code:    500,
			Message: "支付失败" + err.Error(),
		}, err
	}

	return &service.PayResponse{
		Code:    200,
		Message: "支付成功",
	}, nil
}

func (p *PayHandler) PayRevert(ctx context.Context, req *service.PayRequest) (*service.PayResponse, error) {
	fmt.Printf("[Pay] 回滚支付，用户ID: %d 商品ID: %d 金额: %.2f\n", req.UserId, req.ItemId, req.Amount)

	return &service.PayResponse{
		Code:    200,
		Message: "支付回滚成功",
	}, nil
}

func (p *PayHandler) CreateOrder(ctx context.Context, req *service.PayRequest) (*service.PayResponse, error) {
	amount := req.Amount * float32(req.Quantity)

	fmt.Printf("[CreateOrder] 创建订单，用户ID: %d 商品ID: %d 数量: %d 金额: %.2f\n", req.UserId, req.ItemId, req.Quantity, amount)

	err := p.payService.CreateOrder(req.UserId, req.ItemId, req.Quantity, amount)
	if err != nil {
		return &service.PayResponse{
			Code:    500,
			Message: "支付并创建订单失败: " + err.Error(),
		}, err
	}

	return &service.PayResponse{
		Code:    200,
		Message: "支付并创建订单成功",
	}, nil
}
