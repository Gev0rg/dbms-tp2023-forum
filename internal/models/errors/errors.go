package errors

import "errors"

var (
	ErrInvalidUsername = errors.New("Invalid username")
	ErrInvalidEmail    = errors.New("Invalid email")
	ErrInvalidName     = errors.New("Invalid name")
	ErrInvalidPassword = errors.New("Invalid password")

	ErrUserIsAlreadyCreated       = errors.New("The user is already created")
	ErrSessionIsAlreadyCreated    = errors.New("The session is already created")
	ErrEmailIsAlreadyRegistered    = errors.New("The email is already registered")
	ErrUsernameIsAlreadyRegistered = errors.New("The username is already registered")

	ErrCookieNotFound = errors.New("Cookie not found")

	ErrSessionNotFound   = errors.New("Session not found")
	ErrUserNotFound      = errors.New("User not found")
	ErrIncorrectPassword = errors.New("Incorrect password")

	ErrInternal = errors.New("Internal error")
)
