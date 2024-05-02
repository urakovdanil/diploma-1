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

func WithdrawByUser(ctx context.Context, user *types.User, withdraw *types.Withdraw) error {
	if err := UsedStorage.WithdrawByUser(ctx, user, withdraw); err != nil {
		if errors.Is(err, types.ErrInsufficientFunds) || errors.Is(err, types.ErrOrderAlreadyExistsForThisUser) || errors.Is(err, types.ErrOrderAlreadyExistsForAnotherUser) {
			return err
		}
		return fmt.Errorf("unexpected error on WithdrawByUser: %w", err)
	}
	return nil
}

func GetWithdrawalsByUser(ctx context.Context, user *types.User) ([]types.WithdrawWithTS, error) {
	res, err := UsedStorage.GetWithdrawalsByUser(ctx, user)
	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("unexpected error on GetWithdrawalsByUser: %w", err)
	}
	return res, nil
}
