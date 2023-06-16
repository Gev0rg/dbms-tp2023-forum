package repository

import (
	"context"
	"dbms/internal/models"
)

type Repository interface {
	CreateThread(ctx context.Context, thread models.Thread) (models.Thread, error)

	GetThreadById(ctx context.Context, id int) (models.Thread, error)
	GetThreadBySlug(ctx context.Context, slug string) (models.Thread, error)
}
