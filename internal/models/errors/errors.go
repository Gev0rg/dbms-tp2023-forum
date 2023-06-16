package errors

import "errors"

var (
	ErrInvalidUsername = errors.New("Invalid username")
	ErrInvalidEmail    = errors.New("Invalid email")
	ErrInvalidName     = errors.New("Invalid name")
	ErrInvalidPassword = errors.New("Invalid password")

	ErrUserIsAlreadyCreated       = errors.New("The user is already created")
	ErrSessionIsAlreadyCreated    = errors.New("The session is already created")
	ErrEmailIsAlreadyRegistred    = errors.New("The email is already registered")
	ErrUsernameIsAlreadyRegistred = errors.New("The username is already registered")

	ErrCookieNotFound = errors.New("Cookie not found")

	ErrSessionNotFound   = errors.New("Session not found")
	ErrThreadNotFound    = errors.New("Thread not found")
	ErrUserNotFound      = errors.New("User not found")
	ErrForumNotFound     = errors.New("Forum not found")
	ErrIncorrectPassword = errors.New("Incorrect password")

	ErrThreadIsAlreadyExisted = errors.New("Thread is already existed")
	ErrForumIsAlreadyExisted  = errors.New("Forum is already existed")

	ErrInternal = errors.New("Internal error")
)
