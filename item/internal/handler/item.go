package handler

import (
	"context"
	"item/internal/repository"
	"item/internal/service"
)

type ItemHandler struct {
	service.UnimplementedItemServiceServer
	itemRepo repository.ItemRepository
}

// NewItemHandler 创建新的 ItemHandler 实例
func NewItemHandler(itemRepo repository.ItemRepository) *ItemHandler {
	return &ItemHandler{
		itemRepo: itemRepo,
	}
}

func (i *ItemHandler) CreateItem(ctx context.Context, req *service.CreateItemRequest) (*service.CreateItemResponse, error) {
	err := i.itemRepo.CreateItem(int(req.ItemId), req.ItemName, float64(req.Price), int(req.Stock))
	if err != nil {
		return &service.CreateItemResponse{
			Code:    500,
			Message: "创建商品失败",
		}, err
	}
	return &service.CreateItemResponse{
		Code:    200,
		Message: "商品创建成功",
	}, nil
}

func (i *ItemHandler) DecreaseStock(ctx context.Context, req *service.DecreaseStockRequest) (*service.DecreaseStockResponse, error) {
	err := i.itemRepo.DecreaseStock(int(req.ItemId))
	if err != nil {
		return &service.DecreaseStockResponse{
			Code:    500,
			Message: "减少库存失败",
		}, err
	}
	return &service.DecreaseStockResponse{
		Code:    200,
		Message: "库存减少成功",
	}, nil
}
