package usecase

import (
	"context"
	"dbms/internal/models"
	service "dbms/internal/service/repository"
)

type Usecase interface {
	Clear(ctx context.Context) error
	GetStatus(ctx context.Context) (models.Status, error)
}

type usecase struct {
	serviceRepository service.Repository
}

func (u *usecase) Clear(ctx context.Context) error {
	err := u.serviceRepository.Clear(ctx)
	return err
}

func (u *usecase) GetStatus(ctx context.Context) (models.Status, error) {
	status, err := u.serviceRepository.GetStatus(ctx)
	if err != nil {
		return models.Status{}, err
	}

	return status, err
}

func NewUsecase(serviceRepository service.Repository) Usecase {
	return &usecase{serviceRepository: serviceRepository}
}

