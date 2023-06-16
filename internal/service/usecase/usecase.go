package usecase

import (
	"context"
	"dbms/internal/models"
)

type Usecase interface {
	Clear(ctx context.Context) error
	GetStatus(ctx context.Context) (models.Status, error)
}
