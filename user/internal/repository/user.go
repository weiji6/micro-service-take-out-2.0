package repository

import (
	"errors"
	"gorm.io/gorm"
	"user/internal/service"
)

type User struct {
	UserID   int     `gorm:"primaryKey" json:"id"`
	UserName string  `json:"name"`
	Balance  float64 `json:"balance"`
}

//type UserRepository interface {
//	CreateUser(userID int, userName string) error
//	DecreaseBalance(userID int, amount float64) error
//}
//
//type UserRepositoryImpl struct {
//	db *gorm.DB
//}
//
//func NewUserRepositoryImpl(db *gorm.DB) *UserRepositoryImpl {
//	return &UserRepositoryImpl{db: db}
//}

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

//func (u *User) DecreaseBalance(userID int, amount float64) error {
//	return u.db.Model(&User{}).Where("id = ? AND balance >= ?", userID, amount).Update("balance", gorm.Expr("balance - ?", amount)).Error
//}
