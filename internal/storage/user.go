package storage

import (
	"context"
	"diploma-1/internal/types"
	"errors"
	"fmt"
)

func CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	user, err := UsedStorage.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, types.ErrUserAlreadyExists) {
			return nil, types.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("unexpected error on CreateUser: %w", err)
	}
	return user, nil
}

func GetUserByLogin(ctx context.Context, login string) (*types.User, error) {
	user, err := UsedStorage.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("unexpected error on GetUserByLogin: %w", err)
	}
	return user, nil
}
