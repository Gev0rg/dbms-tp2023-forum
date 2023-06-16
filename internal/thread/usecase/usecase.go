package usecase

import (
	"github.com/labstack/echo/v4"
	"dbms/internal/models"
)

type Usecase interface {
	CreatePostsById(ctx echo.Context, id uint64) ([]models.Post, error)
	CreatePostsBySlug(ctx echo.Context, slug string) ([]models.Post, error)

	EditThreadById(ctx echo.Context, thread models.Thread, id uint64) (models.Thread, error)
	EditThreadBySlug(ctx echo.Context, thread models.Thread, slug string) (models.Thread, error)

	VoteById(ctx echo.Context, id uint64) (models.Thread, error)
	VoteBySlug(ctx echo.Context, slug string) (models.Thread, error)

	GetThreadById(ctx echo.Context, id uint64) (models.Thread, error)
	GetThreadBySlug(ctx echo.Context, slug string) (models.Thread, error)

	GetThreadPostsById(ctx echo.Context, id uint64) ([]models.Post, error)
	GetThreadPostsBySlug(ctx echo.Context, slug string) ([]models.Post, error)
}