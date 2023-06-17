package usecase

import (
	"context"
	"dbms/internal/models"
)

type Usecase interface {
	CreatePostsById(ctx context.Context, id int64, posts []models.CreatePost) ([]models.Post, error)
	CreatePostsBySlug(ctx context.Context, slug string, posts []models.CreatePost) ([]models.Post, error)

	EditThreadById(ctx context.Context, id int64, thread models.UpdateThread) (models.Thread, error)
	EditThreadBySlug(ctx context.Context, slug string, thread models.UpdateThread) (models.Thread, error)

	VoteById(ctx context.Context, id int64, vote models.Vote) (models.Thread, error)
	VoteBySlug(ctx context.Context, slug string, vote models.Vote) (models.Thread, error)

	GetThreadById(ctx context.Context, id int64) (models.Thread, error)
	GetThreadBySlug(ctx context.Context, slug string) (models.Thread, error)

	GetThreadPostsById(ctx context.Context, id int64, threadPostsInfo models.GetThreadPostsById) ([]models.Post, error)
	GetThreadPostsBySlug(ctx context.Context, slug string, threadPostsInfo models.GetThreadPostsBySlug) ([]models.Post, error)
}