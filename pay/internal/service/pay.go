package service

import (
	"context"
	"errors"
	"fmt"
	"pay/internal/repository"
	"pay/pkg/lock"
)

type PayService struct {
	UserClient UserServiceClient
	ItemClient ItemServiceClient
	PayRepo    repository.PayRepository
}

func NewPayService(userClient UserServiceClient, ItemClient ItemServiceClient, payRepo repository.PayRepository) *PayService {
	return &PayService{UserClient: userClient, ItemClient: ItemClient, PayRepo: payRepo}
}

func (p *PayService) ProcessPayment(userID int32, itemID int32, quantity int32, amount float32) error {
	ctx := context.Background()

	userResp, err := p.UserClient.CheckBalance(ctx, &CheckBalanceRequest{
		UserId: userID,
	})
	if err != nil {
		return errors.New("无法查询用户余额")
	}

	if userResp.Balance < amount {
		return errors.New("余额不足")
	}

	// 引入分布式锁
	mutex := lock.RedSync.NewMutex(fmt.Sprintf("lock:item:%d", itemID))

	// 加锁
	if err := mutex.Lock(); err != nil {
		return errors.New("获取库存锁失败")
	}
	defer func() {
		_, _ = mutex.Unlock()
	}()

	_, err = p.ItemClient.DecreaseStock(ctx, &DecreaseStockRequest{
		ItemId:   itemID,
		Quantity: quantity,
	})
	if err != nil {
		return errors.New("库存不足")
	}

	_, err = p.UserClient.DecreaseBalance(ctx, &DecreaseBalanceRequest{
		UserId: userID,
		Amount: amount,
	})
	if err != nil {
		return errors.New("扣款失败")
	}

	return nil
}

func (p *PayService) CreateOrder(userID int32, itemID int32, quantity int32, amount float32) error {
	err := p.ProcessPayment(userID, itemID, quantity, amount)
	if err != nil {
		return err
	}

	return p.PayRepo.CreatePay(int(userID), int(itemID), amount)
}
