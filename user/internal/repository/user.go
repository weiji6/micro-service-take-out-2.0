package repository

import (
	"errors"
	"user/internal/service"

	"gorm.io/gorm"
)

type User struct {
	UserID   int     `gorm:"primaryKey" json:"id"`
	UserName string  `json:"name"`
	Balance  float64 `json:"balance"`
}

// CheckUserExist 检查用户是否存在
func (u *User) CheckUserExist(req *service.RegisterRequest) bool {
	if err := DB.Where("user_name=?", req.UserName).First(&u).Error; err == gorm.ErrRecordNotFound {
		return false
	}

	return true
}

func (u *User) CreateUser(req *service.RegisterRequest) (user User, err error) {
	var count int64
	DB.Where("user_name=?", req.UserName).Count(&count)
	if count != 0 {
		return User{}, errors.New("用户存在")
	}

	user = User{
		UserName: req.UserName,
		Balance:  100.0,
	}

	err = DB.Create(&user).Error
	return user, err
}
