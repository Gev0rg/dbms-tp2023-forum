package http

import (
	"context"
	"dbms/internal/models"
	post "dbms/internal/post/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type postHandler struct {
	postUsecase post.Usecase
}

func (h *postHandler) GetFullPostHandler(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	var related []string
	// TODO: handle param related (string array)

	fullPost, err := h.postUsecase.GetFullPost(context.TODO(), id, related)
	if err!= nil {
        return err
    }

	return ctx.JSON(http.StatusOK, fullPost)
}

func (h *postHandler) EditPostHandler(ctx echo.Context) error {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err!= nil {
        return err
    }
	
	var updatePost models.UpdatePost
	err = ctx.Bind(&updatePost)
	if err != nil {
		return err
	}

	editPost, err := h.postUsecase.EditPost(context.TODO(), id, updatePost) 
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, editPost)
}

func NewHandler(e *echo.Echo, postUsecase post.Usecase) postHandler {
	handler := postHandler{postUsecase: postUsecase}

	e.POST("/post/:id/details", handler.EditPostHandler)
	e.GET("/post/:id/details", handler.GetFullPostHandler)

	return handler
}
