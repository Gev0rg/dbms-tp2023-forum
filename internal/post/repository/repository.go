package repository

import (
	"context"
	"dbms/internal/models"
	myErrors "dbms/internal/models/errors"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CheckExistPostById(ctx context.Context, id int64) error

	CreatePost(ctx context.Context, post models.Post) error
	EditPostById(ctx context.Context, id int64, post models.UpdatePost) (models.Post, error)

	GetPostById(ctx context.Context, id int64) (models.Post, error)
}

type repository struct {
	db *sqlx.DB
}

func (r *repository) CheckExistPostById(ctx context.Context, id int64) error {
	var exist bool

	err := r.db.GetContext(ctx, &exist, "SELECT EXISTS(SELECT 1 FROM posts WHERE post_id = $1)", id)
	if err != nil {
		return err
	}

	if !exist {
        return myErrors.ErrPostNotFound
    }

	return nil
}

func (r *repository) CreatePost(ctx context.Context, post models.Post) error {
	// TODO: ТЯ ЖЕ ЛО
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO posts VALUES (:parent, :author, :message, :is_edited, :forum, :thread_id)`, post)
	return err
}

func (r *repository) EditPostById(ctx context.Context, id int64, post models.UpdatePost) (models.Post, error) {
	var updatedPost models.Post

	err := r.db.GetContext(ctx, &updatedPost, `UPDATE posts SET message=$1 WHERE nickname=$2 RETURNING *`,
		post.Message, id)
	if err != nil {
		return models.Post{}, err
	}

	return updatedPost, nil
}

func (r *repository) GetPostById(ctx context.Context, id int64) (models.Post, error) {
	var post models.Post

	err := r.db.GetContext(ctx, &post, `SELECT * FROM posts WHERE post_id=$1`, id)
	if err!= nil {
		return models.Post{}, err
	}

	return post, nil
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}
