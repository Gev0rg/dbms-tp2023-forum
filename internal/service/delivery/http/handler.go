package http

import (
	"context"
	service "dbms/internal/service/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
)

type serviceHandler struct {
	serviceUsecase service.Usecase
}

func (h *serviceHandler) ClearHandler(ctx echo.Context) error {
	err := h.serviceUsecase.Clear(context.TODO())
	if err != nil {
		return err
	}

	return ctx.NoContent(http.StatusOK)
}

func (h *serviceHandler) GetStatusHandler(ctx echo.Context) error {
	status, err := h.serviceUsecase.GetStatus(context.TODO())
	if err!= nil {
        return err
    }

	return ctx.JSON(http.StatusOK, status)
}

func NewHandler(e *echo.Echo, serviceUsecase service.Usecase) serviceHandler {
	handler := serviceHandler{serviceUsecase: serviceUsecase}

	e.POST("/service/clear", handler.ClearHandler)
	e.GET("/service/status", handler.GetStatusHandler)

	return handler
}
