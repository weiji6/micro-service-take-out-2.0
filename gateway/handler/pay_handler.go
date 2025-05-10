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

type PayHandler struct {
	payClient service.PayServiceClient
}

func NewPayHandler() *PayHandler {
	return &PayHandler{
		payClient: client.NewPayClient(),
	}
}

func (p *PayHandler) Register(r *gin.Engine) {
	r.POST("/pay", p.ProcessPayment)
}

func (p *PayHandler) ProcessPayment(c *gin.Context) {
	var pay request.PayRequest
	if err := c.ShouldBindJSON(&pay); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "无效的请求参数"})
		return
	}

	req := &service.PayRequest{
		UserId: int32(pay.UserID),
		ItemId: int32(pay.ItemID),
		Amount: pay.Amount,
	}

	resp, err := p.payClient.Pay(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "支付失败" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: resp.Message})
}
