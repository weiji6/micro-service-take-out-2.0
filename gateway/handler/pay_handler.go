package handler

import (
	"context"
	"gateway/client"
	"gateway/internal/service"
	"gateway/middleware"
	"gateway/pkg/request"
	"gateway/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
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
	if !middleware.Allow() {
		c.JSON(http.StatusTooManyRequests, response.Response{Code: 429, Message: "请求过于频繁，请稍后再试"})
		return
	}

	var pay request.PayRequest
	if err := c.ShouldBindJSON(&pay); err != nil {
		c.JSON(http.StatusBadRequest, response.Response{Code: 400, Message: "无效的请求参数"})
		return
	}

	//// 创建 DTM SAGA
	//dtmAddress := viper.GetString("dtm.address")
	//saga := dtmgrpc.NewSagaGrpc(dtmAddress, dtmgrpc.MustGenGid(dtmAddress)).
	//	Add("localhost:30002/item.ItemService/DecreaseStock", "localhost:30002/item.ItemService/DecreaseStockRevert",
	//		&service.DecreaseStockRequest{ItemId: int32(pay.ItemID), Quantity: int32(pay.Quantity)}).
	//	Add("localhost:30001/user.UserService/DecreaseBalance", "localhost:30001/user.UserService/DecreaseBalanceRevert",
	//		&service.DecreaseBalanceRequest{UserId: int32(pay.UserID), Amount: pay.Amount}).
	//	Add("localhost:30003/pay.PayService/Pay", "localhost:30003/pay.PayService/PayRevert",
	//		&service.PayRequest{UserId: int32(pay.UserID), ItemId: int32(pay.ItemID), Quantity: int32(pay.Quantity), Amount: pay.Amount})
	//
	//// 提交事务
	//if err := saga.Submit(); err != nil {
	//	c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "支付失败" + err.Error()})
	//	return
	//}

	req := &service.PayRequest{
		UserId:   int32(pay.UserID),
		ItemId:   int32(pay.ItemID),
		Quantity: int32(pay.Quantity),
		Amount:   pay.Amount,
	}

	_, err := p.payClient.CreateOrder(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Response{Code: 500, Message: "创建订单失败，已回滚支付" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, response.Response{Code: 200, Message: "支付成功并创建订单"})
}
