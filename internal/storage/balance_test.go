package storage

import (
	"context"
	"diploma-1/internal/storage/mocks"
	"diploma-1/internal/types"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetBalanceByUser(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		user    *types.User
		want    *types.Balance
		wantErr error
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			user:    new(types.User),
			want:    new(types.Balance),
			wantErr: nil,
		},
		{
			name:    "user not found",
			ctx:     context.Background(),
			user:    new(types.User),
			want:    nil,
			wantErr: types.ErrUserNotFound,
		},
		{
			name:    "unexpected error",
			ctx:     context.Background(),
			user:    new(types.User),
			want:    nil,
			wantErr: errors.New("some unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(mocks.Storage)
			mockStorage.On("GetBalanceByUser", tt.ctx, tt.user).Return(tt.want, tt.wantErr)
			UsedStorage = mockStorage
			got, err := GetBalanceByUser(tt.ctx, tt.user)
			require.Equal(t, tt.want, got)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestWithdrawByUser(t *testing.T) {
	tests := []struct {
		name     string
		ctx      context.Context
		user     *types.User
		withdraw *types.Withdraw
		wantErr  error
	}{
		{
			name:     "success",
			ctx:      context.Background(),
			user:     new(types.User),
			withdraw: new(types.Withdraw),
			wantErr:  nil,
		},
		{
			name:     "user not found",
			ctx:      context.Background(),
			user:     new(types.User),
			withdraw: new(types.Withdraw),
			wantErr:  types.ErrUserNotFound,
		},
		{
			name:     "insufficient funds",
			ctx:      context.Background(),
			user:     new(types.User),
			withdraw: new(types.Withdraw),
			wantErr:  types.ErrInsufficientFunds,
		},
		{
			name:     "order already exists for this user",
			ctx:      context.Background(),
			user:     new(types.User),
			withdraw: new(types.Withdraw),
			wantErr:  types.ErrOrderAlreadyExistsForThisUser,
		},
		{
			name:     "order already exists for another user",
			ctx:      context.Background(),
			user:     new(types.User),
			withdraw: new(types.Withdraw),
			wantErr:  types.ErrOrderAlreadyExistsForAnotherUser,
		},
		{
			name:     "unexpected error",
			ctx:      context.Background(),
			user:     new(types.User),
			withdraw: new(types.Withdraw),
			wantErr:  errors.New("some unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(mocks.Storage)
			mockStorage.On("WithdrawByUser", tt.ctx, tt.user, tt.withdraw).Return(tt.wantErr)
			UsedStorage = mockStorage
			err := WithdrawByUser(tt.ctx, tt.user, tt.withdraw)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGetWithdrawalsByUser(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		user    *types.User
		want    []types.WithdrawWithTS
		wantErr error
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			user:    new(types.User),
			want:    []types.WithdrawWithTS{},
			wantErr: nil,
		},
		{
			name:    "user not found",
			ctx:     context.Background(),
			user:    new(types.User),
			want:    nil,
			wantErr: types.ErrUserNotFound,
		},
		{
			name:    "unexpected error",
			ctx:     context.Background(),
			user:    new(types.User),
			want:    nil,
			wantErr: errors.New("some unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(mocks.Storage)
			mockStorage.On("GetWithdrawalsByUser", tt.ctx, tt.user).Return(tt.want, tt.wantErr)
			UsedStorage = mockStorage
			got, err := GetWithdrawalsByUser(tt.ctx, tt.user)
			require.Equal(t, tt.want, got)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
