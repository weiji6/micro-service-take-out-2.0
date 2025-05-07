package handler

import "pay/internal/service"

type PayHandler struct {
	service service.PayService
}

func NewPayHandler(service service.PayService) *PayHandler {
	return &PayHandler{service: service}
}
