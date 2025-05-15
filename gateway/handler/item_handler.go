package handler

import (
	"context"
	"fmt"
	"gateway/client"
	"gateway/internal/service"
	"gateway/middleware"
	"gateway/pkg/request"
	"gateway/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ItemHandler struct {
	itemClient service.ItemServiceClient
}

func NewItemHandler() *ItemHandler {
	return &ItemHandler{
		itemClient: client.NewItemClient(),
	}
}

func (i *ItemHandler) Register(r *gin.Engine) {
	r.GET("/items", i.GetItemList)
	r.POST("/items", i.CreateItem)
}

func (i *ItemHandler) GetItemList(c *gin.Context) {
	if !middleware.Allow() {
		c.JSON(http.StatusTooManyRequests, response.Response{Code: 429, Message: "请求过于频繁，请稍后再试"})
		return
	}

	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")

	req := &service.GetItemListRequest{
		Page:     int32(parseInt(page)),
		PageSize: int32(parseInt(pageSize)),
	}

	resp, err := i.itemClient.GetItemList(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "获取商品列表失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    int(resp.Code),
		Message: resp.Message,
		Data:    resp,
	})
}

func (i *ItemHandler) CreateItem(c *gin.Context) {
	if !middleware.Allow() {
		c.JSON(http.StatusTooManyRequests, response.Response{Code: 429, Message: "请求过于频繁，请稍后再试"})
		return
	}

	var item request.ItemRequest
	if err := c.ShouldBindJSON(&item); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "无效的请求参数"})
		return
	}

	req := &service.CreateItemRequest{
		ItemId:   int32(item.ItemID),
		ItemName: item.ItemName,
		Price:    float32(item.Price),
		Stock:    int32(item.Stock),
	}

	resp, err := i.itemClient.CreateItem(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "创建商品失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Response{
		Code:    int(resp.Code),
		Message: resp.Message,
	})
}

// 辅助函数：将字符串转换为整数
func parseInt(s string) int {
	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return 1 // 默认返回 1
	}
	return result
}
