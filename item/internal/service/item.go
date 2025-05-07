package service

import "item/internal/repository"

type ItemService interface {
}

type ItemServiceImpl struct {
	repository repository.ItemRepository
}

func NewItemServiceImpl(repository repository.ItemRepository) *ItemServiceImpl {
	return &ItemServiceImpl{repository: repository}
}
