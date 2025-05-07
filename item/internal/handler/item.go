package handler

import "item/internal/service"

type ItemHandler struct {
	service service.UnimplementedItemServiceServer
}

func NewItemHandler() *ItemHandler {
	return &ItemHandler{}
}
