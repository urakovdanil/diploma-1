package storage

import (
	"context"
	"diploma-1/internal/types"
	"errors"
	"fmt"
)

func CreateOrder(ctx context.Context, order *types.Order) (*types.Order, error) {
	res, err := UsedStorage.CreateOrder(ctx, order)
	if err != nil {
		if errors.Is(err, types.ErrOrderAlreadyExistsForThisUser) || errors.Is(err, types.ErrOrderAlreadyExistsForAnotherUser) {
			return nil, err
		}
		return nil, fmt.Errorf("unexpected error on CreateOrder: %w", err)
	}
	return res, nil
}
