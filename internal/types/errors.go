package types

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")

	ErrWrongPassword = errors.New("wrong password")

	ErrUserAlreadyExists = errors.New("user already exists")

	ErrInvalidAuthInput = errors.New("empty login or password")
)
