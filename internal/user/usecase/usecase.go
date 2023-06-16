package usecase

import (
	"context"
	"dbms/internal/models"
	user "dbms/internal/user/repository"
)

type Usecase interface {
	CreateUser(ctx context.Context, createUser models.User) (models.User, error)
	UpdateUserByNickname(ctx context.Context, user models.User) (models.User, error)
	GetUserByNickname(ctx context.Context, nickname string) (models.User, error)
}

type usecase struct {
	userRepository user.Repository
}

func (u usecase) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	err := u.userRepository.CheckExistUserByNickname(ctx, user.Nickname)
	if err != nil {
		return models.User{}, err
	}

	err = u.userRepository.CreateUser(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u usecase) UpdateUserByNickname(ctx context.Context, user models.User) (models.User, error) {
	updatedUser, err := u.userRepository.UpdateUserByNickname(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	return updatedUser, nil
}

func (u usecase) GetUserByNickname(ctx context.Context, nickname string) (models.User, error) {
	user, err := u.userRepository.GetUserByNickname(ctx, nickname)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func NewUsecase(userRepository user.Repository) Usecase {
	return &usecase{userRepository: userRepository}
}
