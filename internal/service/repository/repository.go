package repository

import (
	"context"
	"dbms/internal/models"
)

type Repository interface {
	Clear(ctx context.Context) error
	GetStatus(ctx context.Context) (models.Status, error)
}
