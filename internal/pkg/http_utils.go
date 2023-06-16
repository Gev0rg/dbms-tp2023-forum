package http_utils

import (
	"errors"
	"net/http"
	myErrors "dbms/internal/models/errors"
)

func StatusCode(err error) int {
	switch {
	case errors.Is(err, myErrors.ErrInvalidUsername):
		return http.StatusBadRequest
	case errors.Is(err, myErrors.ErrInvalidEmail):
		return http.StatusBadRequest
	case errors.Is(err, myErrors.ErrInvalidPassword):
		return http.StatusBadRequest
	case errors.Is(err, myErrors.ErrEmailIsAlreadyRegistered):
		return http.StatusConflict
	case errors.Is(err, myErrors.ErrUsernameIsAlreadyRegistered):
		return http.StatusConflict
	case errors.Is(err, myErrors.ErrSessionIsAlreadyCreated):
		return http.StatusConflict
	case errors.Is(err, myErrors.ErrCookieNotFound):
		return http.StatusUnauthorized
	case errors.Is(err, myErrors.ErrSessionNotFound):
		return http.StatusNotFound
	case errors.Is(err, myErrors.ErrUserNotFound):
		return http.StatusNotFound
	case errors.Is(err, myErrors.ErrIncorrectPassword):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
