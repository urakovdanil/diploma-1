package storage

import (
	"context"
	"diploma-1/internal/types"
	"errors"
	"fmt"
)

func GetBalanceByUser(ctx context.Context, user *types.User) (*types.Balance, error) {
	res, err := UsedStorage.GetBalanceByUser(ctx, user)
	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("unexpected error on GetBalanceByUser: %w", err)
	}
	return res, err
}
