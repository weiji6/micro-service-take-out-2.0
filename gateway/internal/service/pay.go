package service

import (
	"context"
	"errors"
)

type PayService struct {
	UserClient UserServiceClient
	ItemClient ItemServiceClient
}

func NewPayService(userClient UserServiceClient, ItemClient ItemServiceClient) *PayService {
	return &PayService{UserClient: userClient, ItemClient: ItemClient}
}

func (p *PayService) ProcessPayment(userID int32, itemID int32, amount float32) error {
	userResp, err := p.UserClient.CheckBalance(context.Background(), &CheckBalanceRequest{
		UserId: userID,
	})
	if err != nil {
		return errors.New("无法查询用户余额")
	}

	if userResp.Balance < amount {
		return errors.New("余额不足")
	}

	_, err = p.ItemClient.DecreaseStock(context.Background(), &DecreaseStockRequest{
		ItemId: itemID,
	})
	if err != nil {
		return errors.New("库存不足")
	}

	_, err = p.UserClient.DecreaseBalance(context.Background(), &DecreaseBalanceRequest{
		UserId: userID,
		Amount: amount,
	})
	if err != nil {
		return errors.New("扣款失败")
	}

	return nil
}
