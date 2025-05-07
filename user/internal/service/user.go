package service

import "user/internal/repository"

type UserService interface {
	Register(userID int, UserName string) (err error)
}

type UserServiceImpl struct {
	repository repository.UserRepository
}

func NewUserServiceImpl(repository repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{repository: repository}
}

func (u *UserServiceImpl) Register(userID int, userName string) (err error) {
	return u.repository.CreateUser(userID, userName)
}
