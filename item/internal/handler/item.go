package handler

import (
	"context"
	"fmt"
	"item/internal/repository"
	"item/internal/service"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

type ItemHandler struct {
	service.UnimplementedItemServiceServer
	itemRepo repository.ItemRepository
	sf       singleflight.Group
	// 用于记录每个商品的等待请求数
	waitingRequests sync.Map
	// 用于记录每个商品的降级阈值
	degradeThreshold sync.Map
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
	itemID := int(req.ItemId)
	
	// 增加等待请求计数
	waitingCount, _ := i.waitingRequests.LoadOrStore(itemID, int32(0))
	currentCount := waitingCount.(int32)
	i.waitingRequests.Store(itemID, currentCount+1)
	defer func() {
		// 减少等待请求计数
		if count, ok := i.waitingRequests.Load(itemID); ok {
			i.waitingRequests.Store(itemID, count.(int32)-1)
		}
	}()

	// 获取降级阈值，如果没有设置则默认为 100
	threshold, _ := i.degradeThreshold.LoadOrStore(itemID, int32(100))
	
	// 如果等待请求数超过阈值，直接返回降级响应
	if currentCount > threshold.(int32) {
		return &service.DecreaseStockResponse{
			Code:    429,
			Message: "系统繁忙，请稍后重试",
		}, nil
	}

	// 设置超时上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// 使用 channel 来模拟分布式锁的等待
	lockCh := make(chan struct{})
	go func() {
		// 模拟获取分布式锁的过程
		time.Sleep(100 * time.Millisecond)
		close(lockCh)
	}()

	select {
	case <-timeoutCtx.Done():
		return &service.DecreaseStockResponse{
			Code:    408,
			Message: "请求超时",
		}, nil
	case <-lockCh:
		// 获取到锁后执行库存扣减
		err := i.itemRepo.DecreaseStock(itemID, int(req.Quantity))
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
}

func (i *ItemHandler) DecreaseStockRevert(ctx context.Context, req *service.DecreaseStockRequest) (*service.DecreaseStockResponse, error) {
	fmt.Printf("[Item] 回滚库存，商品ID: %d 数量: %d\n", req.ItemId, req.Quantity)

	err := repository.DB.Model(&repository.Item{}).
		Where("id = ?", req.ItemId).
		Update("stock", gorm.Expr("stock + ?", req.Quantity)).Error

	if err != nil {
		return &service.DecreaseStockResponse{Code: 500, Message: "回滚失败: " + err.Error()}, err
	}

	return &service.DecreaseStockResponse{Code: 200, Message: "库存回滚成功"}, nil
}

func (i *ItemHandler) GetItemList(ctx context.Context, req *service.GetItemListRequest) (*service.GetItemListResponse, error) {
	// 使用 singleflight 来防止缓存击穿
	key := fmt.Sprintf("item_list_%d_%d", req.Page, req.PageSize)
	result, err, _ := i.sf.Do(key, func() (interface{}, error) {
		items, total, err := i.itemRepo.GetItemList(int(req.Page), int(req.PageSize))
		if err != nil {
			return nil, err
		}

		itemInfos := make([]*service.ItemInfo, 0, len(items))
		for _, item := range items {
			itemInfos = append(itemInfos, &service.ItemInfo{
				ItemId:   int32(item.ItemID),
				ItemName: item.ItemName,
				Price:    float32(item.Price),
				Stock:    int32(item.Stock),
			})
		}

		return &service.GetItemListResponse{
			Code:    200,
			Message: "获取商品列表成功",
			Items:   itemInfos,
			Total:   int32(total),
		}, nil
	})

	if err != nil {
		return &service.GetItemListResponse{
			Code:    500,
			Message: "获取商品列表失败: " + err.Error(),
		}, err
	}

	return result.(*service.GetItemListResponse), nil
}
