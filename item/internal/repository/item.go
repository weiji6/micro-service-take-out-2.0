package repository

import "gorm.io/gorm"

type Item struct {
	ItemID   int     `gorm:"primaryKey"  json:"id"`
	ItemName string  `json:"itemName"`
	Price    float64 `json:"price"`
	Stock    int     `json:"stock"`
}

type ItemRepository interface {
	CreateItem(itemID int, itemName string, price float64, stock int) error
	DecreaseStock(itemID int) error
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

func (i *ItemRepositoryImpl) DecreaseStock(itemID int) error {
	return i.db.Model(&Item{}).Where("id = ? AND stock > 0", itemID).Update("stock", gorm.Expr("stock - 1")).Error
}
