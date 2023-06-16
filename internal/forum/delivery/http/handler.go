package http

import (
	"context"
	forum "dbms/internal/forum/usecase"
	"dbms/internal/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type forumHandler struct {
	forumUsecase forum.Usecase
}

func (h *forumHandler) CreateForumHandler(ctx echo.Context) error {
	var createForum models.CreateForum

	err := ctx.Bind(&createForum)
	if err != nil {
		return err
	}

	forum, err := h.forumUsecase.CreateForum(context.TODO(), models.Forum{
		Title: createForum.Title,
		User: createForum.User,
		Slug: createForum.Slug,
	})
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, forum)
}

func (h *forumHandler) CreateForumThreadHandler(ctx echo.Context) error {
	slug := ctx.Param("slug")
	var createForumThread models.CreateThread

	err := ctx.Bind(&createForumThread)
	if err!= nil {
        return err
    }

	thread, err := h.forumUsecase.CreateForumThread(context.TODO(), slug, models.Thread{
		Title: createForumThread.Title,
        Author: createForumThread.Author,
        Message: createForumThread.Message,
		Created: createForumThread.Created,
	})
	if err!= nil {
        return err
    }

	return ctx.JSON(http.StatusCreated, thread)
}

func (h *forumHandler) GetForumHandler(ctx echo.Context) error {
	slug := ctx.Param("slug")

	forum, err := h.forumUsecase.GetForumBySlug(context.TODO(), slug)
	if err!= nil {
        return err
    }

	return ctx.JSON(http.StatusOK, forum)
}

func (h *forumHandler) GetForumUsersHandler(ctx echo.Context) error {
	slug := ctx.Param("slug")
	limitStr := ctx.Param("limit")
	since := ctx.Param("since")
	descStr := ctx.Param("desc")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		return err
	}

	desc, err := strconv.ParseBool(descStr)
	if err != nil {
		return err
	}

	users, err := h.forumUsecase.GetForumUsersBySlug(context.TODO(), slug, limit, since, desc)
	if err != nil {
        return err
    }

	return ctx.JSON(http.StatusOK, users)
}

func (h *forumHandler) GetForumThreadsHandler(ctx echo.Context) error {
	slug := ctx.Param("slug")
    limitStr := ctx.Param("limit")
    since := ctx.Param("since")
    descStr := ctx.Param("desc")

    limit, err := strconv.ParseInt(limitStr, 10, 64)
    if err!= nil {
        return err
    }

    desc, err := strconv.ParseBool(descStr)
    if err!= nil {
        return err
    }

    threads, err := h.forumUsecase.GetForumThreadsBySlug(context.TODO(), slug, limit, since, desc)
    if err!= nil {
        return err
    }

    return ctx.JSON(http.StatusOK, threads)
}

func NewHandler(e *echo.Echo, forumUsecase forum.Usecase) forumHandler {
	handler := forumHandler{forumUsecase: forumUsecase}
	
	e.POST("/forum/create", handler.CreateForumHandler)
	e.POST("/forum/:slug/create", handler.CreateForumThreadHandler)

	e.GET("/forum/:slug/details", handler.GetForumHandler)
	e.GET("/forum/:slug/users", handler.GetForumUsersHandler)
	e.GET("/forum/:slug/threads", handler.GetForumThreadsHandler)
	
	return handler
}
