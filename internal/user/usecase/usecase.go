package usecase

import (
	"context"
	"dbms/internal/models"
	user "dbms/internal/user/repository"
)

type Usecase interface {
	CreateUser(ctx context.Context, createUser model.User) (model.User, error)
	UpdateUserByNickname(ctx context.Context, user model.User) (model.User, error)
	GetUserByNickname(ctx context.Context, nickname string) (model.User, error)
}

type usecase struct {
	userRepository user.Repository
}

func NewUserUsecase(userRepository user.Repository) Usecase {
	return &usecase{userRepository: userRepository}
}

func (u usecase) CreateUser(ctx context.Context, user model.User) (model.User, error) {
	err := u.userRepository.CheckExistUserByNickname(ctx, user.Nickname)
	if err != nil {
		return model.User{}, err
	}

	err = u.userRepository.CreateUser(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (u usecase) UpdateUserByNickname(ctx context.Context, user model.User) (model.User, error) {
	updatedUser, err := u.userRepository.UpdateUserByNickname(ctx, user)
	if err != nil {
		return model.User{}, err
	}

	return updatedUser, nil
}

func (u usecase) GetUserByNickname(ctx context.Context, nickname string) (model.User, error) {
	user, err := u.userRepository.GetUserByNickname(ctx, nickname)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
