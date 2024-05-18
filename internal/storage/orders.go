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

func GetOrdersByUser(ctx context.Context, user *types.User) ([]types.Order, error) {
	res, err := UsedStorage.GetOrdersByUser(ctx, user)
	if err != nil {
		if errors.Is(err, types.ErrOrderNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("unexpected error on GetOrdersByUser: %w", err)
	}
	return res, nil
}

func UpdateOrderFromAccrual(ctx context.Context, order *types.OrderFromAccrual) error {
	if err := UsedStorage.UpdateOrderFromAccrual(ctx, order); err != nil {
		if errors.Is(err, types.ErrOrderNotFound) {
			return err
		}
		return fmt.Errorf("unexpected error on UpdateOrderFromAccrual: %w", err)
	}
	return nil
}
