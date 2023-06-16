package usecase

import (
	"github.com/labstack/echo/v4"
	"dbms/internal/models"
)

type Usecase interface {
	GetFullPost(ctx echo.Context, id uint64) (models.FullPost, error)
	EditPost(ctx echo.Context, id uint64, thread models.Thread) (models.Post, error)
}
