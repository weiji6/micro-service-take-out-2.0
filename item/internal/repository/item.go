package repository

import (
	"errors"

	"gorm.io/gorm"
)

type Item struct {
	ItemID   int     `gorm:"primaryKey"  json:"id"`
	ItemName string  `json:"itemName"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
}

type ItemRepository interface {
	CreateItem(itemID int, itemName string, price float64, stock int) error
	DecreaseStock(itemID int, quantity int) error
	GetItemList(page, pageSize int) ([]Item, int64, error)
}

type ItemRepositoryImpl struct {
	db *gorm.DB
}

func NewItemRepositoryImpl(db *gorm.DB) *ItemRepositoryImpl {
	return &ItemRepositoryImpl{db: db}
}

func (i *ItemRepositoryImpl) CreateItem(itemID int, itemName string, price float64, stock int) error {
	item := &Item{
		ItemID:   itemID,
		ItemName: itemName,
		Price:    price,
		Stock:    stock,
	}

	return i.db.Create(&item).Error
}

func (i *ItemRepositoryImpl) DecreaseStock(itemID int, quantity int) error {
	var item Item
	if err := i.db.First(&item, "item_id = ?", itemID).Error; err != nil {
		return err
	}

	if item.Stock < quantity {
		return errors.New("库存不足")
	}

	item.Stock -= quantity
	return i.db.Save(&item).Error
}

func (i *ItemRepositoryImpl) GetItemList(page, pageSize int) ([]Item, int64, error) {
	var items []Item
	var total int64

	offset := (page - 1) * pageSize

	// 获取总数
	if err := i.db.Model(&Item{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := i.db.Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
