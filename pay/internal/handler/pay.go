package handler

import (
	"context"
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
	err := p.payService.ProcessPayment(req.UserId, req.ItemId, req.Amount)
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
