package repository

import (
	"context"
	"dbms/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateThread(ctx context.Context, thread models.Thread) (models.Thread, error)

	GetThreadById(ctx context.Context, id int) (models.Thread, error)
	GetThreadBySlug(ctx context.Context, slug string) (models.Thread, error)
}

type repository struct {
	db *sqlx.DB
}

func (r repository) CreateThread(ctx context.Context, thread models.Thread) (models.Thread, error) {

}

func (r repository) GetThreadById(ctx context.Context, id int) (models.Thread, error) {

}

func (r repository) GetThreadBySlug(ctx context.Context, slug string) (models.Thread, error) {

}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}
