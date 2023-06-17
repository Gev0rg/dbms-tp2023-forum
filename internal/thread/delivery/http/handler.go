package http

import (
	"context"
	"dbms/internal/models"
	thread "dbms/internal/thread/usecase"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type threadHandler struct {
	threadUsecase thread.Usecase
}

func (h *threadHandler) CreatePostsInThreadHandler(ctx echo.Context) error {
	var posts []models.Post

	var createPost []models.CreatePost
	err := ctx.Bind(&createPost)
	if err != nil {
		return err
	}

	slugOrId := ctx.Param("slug_or_id")
	id, err := strconv.ParseInt(slugOrId, 10, 64)
	
	if err != nil {
		posts, err = h.threadUsecase.CreatePostsBySlug(context.TODO(), slugOrId, createPost)
	} else {
		posts, err = h.threadUsecase.CreatePostsById(context.TODO(), id, createPost)
	}

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, posts)
}

func (h *threadHandler) GetThreadHandler(ctx echo.Context) error {
	var thread models.Thread

	slugOrId := ctx.Param("slug_or_id")
	id, err := strconv.ParseInt(slugOrId, 10, 64)
	
	if err != nil {
		thread, err = h.threadUsecase.GetThreadBySlug(context.TODO(), slugOrId)
	} else {
		thread, err = h.threadUsecase.GetThreadById(context.TODO(), id)
	}

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, thread)
}

func (h *threadHandler) UpdateThreadHandler(ctx echo.Context) error {
	var thread models.Thread

	var updateThread models.UpdateThread
	err := ctx.Bind(&updateThread)
	if err != nil {
		return err
	}

	slugOrId := ctx.Param("slug_or_id")
	id, err := strconv.ParseInt(slugOrId, 10, 64)
	
	if err != nil {
		thread, err = h.threadUsecase.EditThreadBySlug(context.TODO(), slugOrId, updateThread)
	} else {
		thread, err = h.threadUsecase.EditThreadById(context.TODO(), id, updateThread)
	}

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, thread)
}

func (h *threadHandler) GetThreadPostsHandler(ctx echo.Context) error {
	var posts []models.Post

	slugOrId := ctx.Param("slug_or_id")
    limitStr := ctx.Param("limit")
    since := ctx.Param("since")
	sort := ctx.Param("sort")
    descStr := ctx.Param("desc")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
    if err!= nil {
        return err
    }

    desc, err := strconv.ParseBool(descStr)
    if err!= nil {
        return err
    }

	id, err := strconv.ParseInt(slugOrId, 10, 64)
	
	if err != nil {
		posts, err = h.threadUsecase.GetThreadPostsBySlug(context.TODO(), slugOrId, models.GetThreadPostsBySlug{
			Slug: slugOrId,
			Limit: limit,
            Since: since,
            Sort: sort,
            Desc: desc,
		})
	} else {
		posts, err = h.threadUsecase.GetThreadPostsById(context.TODO(), id, models.GetThreadPostsById{
			Id: id,
			Limit: limit,
            Since: since,
            Sort: sort,
            Desc: desc,
		})
	}

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, posts)
}

func (h *threadHandler) CreateVoteHandler(ctx echo.Context) error {
	var thread models.Thread

	var vote models.Vote
	err := ctx.Bind(&vote)
	if err != nil {
		return err
	}

	slugOrId := ctx.Param("slug_or_id")
	id, err := strconv.ParseInt(slugOrId, 10, 64)
	
	if err != nil {
		thread, err = h.threadUsecase.VoteBySlug(context.TODO(), slugOrId, vote)
	} else {
		thread, err = h.threadUsecase.VoteById(context.TODO(), id, vote)
	}

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, thread)
}

func NewHandler(e *echo.Echo, threadUsecase thread.Usecase) threadHandler {
	handler := threadHandler{threadUsecase: threadUsecase}

	return handler
}
