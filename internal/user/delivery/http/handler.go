package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"dbms/internal/models"
	user "dbms/internal/user/usecase"
)

type userHandler struct {
	userUsecase user.Usecase
}

func (h userHandler) CreateUserHandler(ctx echo.Context) error {
	var createUser model.User

	err := ctx.Bind(&createUser)
	if err != nil {
		return err
	}

	nickname := ctx.Param("nickname")
	createUser.Nickname = nickname

	user, err := h.userUsecase.CreateUser(context.TODO(), createUser)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, user)
}

func (u *userHandler) GetUserHandler(ctx echo.Context) error {
	nickname := ctx.Param("nickname")

	user, err := u.userUsecase.GetUserByNickname(context.TODO(), nickname)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, user)
}

func (u *userHandler) UpdateUserHandler(ctx echo.Context) error {
	var user model.User
	err := ctx.Bind(&user)

	if err != nil {
		return err
	}

	nickname := ctx.Param("nickname")
	user.Nickname = nickname

	user, err = u.userUsecase.UpdateUserByNickname(context.TODO(), user)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, user)
}

func NewHandler(e *echo.Echo, userUsecase user.Usecase) userHandler {
	handler := userHandler{userUsecase: userUsecase}

	getUserUrl := "/user/:nickname/profile"
	createUserUrl := "/user/:nickname/create"
	editUserUrl := "/user/:nickname/profile"

	e.GET(getUserUrl, handler.CreateUserHandler)
	e.POST(createUserUrl, handler.GetUserHandler)
	e.POST(editUserUrl, handler.UpdateUserHandler)

	return handler

}
