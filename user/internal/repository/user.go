package repository

import "gorm.io/gorm"

type User struct {
	UserID   int `gorm:"primaryKey" json:"id"`
	UserName string`json:"name"`
	Balance  float64 `json:"balance"`
}

type UserRepository interface {
	CreateUser(userID int, userName string) error
	DecreaseBalance(userID int, amount float64) error
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (u *UserRepositoryImpl) CreateUser(userID int, userName string) error {
	user := &User{
		UserID:   userID,
		UserName: userName,
		Balance:  1000.0,
	}

	return u.db.Create(&user).Error
}

func (u *UserRepositoryImpl) DecreaseBalance(userID int, amount float64) error {
	return u.db.Model(&User{}).Where("id = ? AND balance >= ?", userID, amount).Update("balance", gorm.Expr("balance - ?", amount)).Error
}
