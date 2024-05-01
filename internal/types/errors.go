package types

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrWrongPassword     = errors.New("wrong password")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidAuthInput  = errors.New("empty login or password")

	ErrEmptyOrderNumber                 = errors.New("empty order number")
	ErrNonDigitalOrderNumber            = errors.New("non-digital order number")
	ErrInvalidOrderNumber               = errors.New("invalid order number")
	ErrOrderAlreadyExistsForThisUser    = errors.New("order already exists for provided user")
	ErrOrderAlreadyExistsForAnotherUser = errors.New("order already exists for another user")
	ErrOrderNotFound                    = errors.New("no orders found")
)
