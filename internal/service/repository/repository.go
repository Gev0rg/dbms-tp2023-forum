package repository

import (
	"context"
	"dbms/internal/models"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	Clear(ctx context.Context) error
	GetStatus(ctx context.Context) (models.Status, error)
}

type repository struct {
	db *sqlx.DB
}

func (r repository) Clear(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `TRUNCATE TABLE users, forums, threads, posts, users_votes, users_forums CASCADE;`)
	return err
}

func (r repository) GetStatus(ctx context.Context) (models.Status, error) {
	var status models.Status

	err := r.db.GetContext(ctx, &status, `
		SELECT 
			(SELECT COUNT(*) FROM forums) AS forums,
			(SELECT COUNT(*) FROM posts) AS posts,
            (SELECT COUNT(*) FROM threads) AS threads,
            (SELECT COUNT(*) FROM users) AS users
	`)
	if err!= nil {
        return models.Status{}, err
    }

	return status, nil
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}
