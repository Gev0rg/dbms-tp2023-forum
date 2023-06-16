package repository

import (
	"context"
	"dbms/internal/models"
)

type Repository interface {
	CheckExistPostBySlug(ctx context.Context, slug string) error
	CreatePost(ctx context.Context, forum models.Forum) error
	GetPostById(ctx context.Context, slug string) (models.Forum, error)
}
