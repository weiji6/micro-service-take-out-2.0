package repository

import "gorm.io/gorm"

type Pay struct {
	PayID  int `gorm:"primaryKey" json:"id"`
	UserID int `json:"userId"`
	ItemID int `json:"itemId"`
	Amount int `json:"amount"`
}

type PayRepository interface {
	CreatePay(userID int, itemID int, amount float64) error
}

type PayRepositoryImpl struct {
	db *gorm.DB
}

func NewPayRepositoryImpl(db *gorm.DB) *PayRepositoryImpl {
	return &PayRepositoryImpl{db: db}
}

func (p *PayRepositoryImpl) CreatePay(userID int, itemID int, amount int) error {
	pay := &Pay{
		UserID: userID,
		ItemID: itemID,
		Amount: amount,
	}

	return p.db.Create(&pay).Error
}
