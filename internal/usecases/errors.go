package usecases

import "errors"

var (
	ErrCoffeeEntryNotFound = errors.New("coffee entry not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidInput        = errors.New("invalid input")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInternalError       = errors.New("internal error")
	ErrNotFound           = errors.New("not found")
	ErrEntryAlreadyExists  = errors.New("entry already exists")
)
