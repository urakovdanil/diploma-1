package storage

import (
	"context"
	"diploma-1/internal/storage/mocks"
	"diploma-1/internal/types"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestCreateOrder(t *testing.T) {

	tests := []struct {
		name    string
		ctx     context.Context
		order   *types.Order
		want    *types.Order
		wantErr error
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			order:   new(types.Order),
			want:    new(types.Order),
			wantErr: nil,
		},
		{
			name:    "order already exists for current user",
			ctx:     context.Background(),
			order:   new(types.Order),
			want:    nil,
			wantErr: types.ErrOrderAlreadyExistsForThisUser,
		},
		{
			name:    "order already exists for another user",
			ctx:     context.Background(),
			order:   new(types.Order),
			want:    nil,
			wantErr: types.ErrOrderAlreadyExistsForAnotherUser,
		},
		{
			name:    "unexpected error",
			ctx:     context.Background(),
			order:   new(types.Order),
			want:    nil,
			wantErr: errors.New("some unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(mocks.Storage)
			mockStorage.On("CreateOrder", tt.ctx, tt.order).Return(tt.want, tt.wantErr)
			UsedStorage = mockStorage
			got, err := CreateOrder(tt.ctx, tt.order)
			require.Equal(t, tt.want, got)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGetOrdersByUser(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		user    *types.User
		want    []types.Order
		wantErr error
	}{
		{
			name: "success",
			ctx:  context.Background(),
			user: new(types.User),
			want: []types.Order{
				{ID: 1, UserID: 1, Status: types.OrderStatusNew, Accrual: 0, Number: "1234", CreatedAt: time.Now().UTC()},
				{ID: 2, UserID: 1, Status: types.OrderStatusProcessed, Accrual: 10, Number: "4321", CreatedAt: time.Now().UTC()},
			},
			wantErr: nil,
		},
		{
			name:    "orders not found",
			ctx:     context.Background(),
			user:    new(types.User),
			want:    nil,
			wantErr: types.ErrOrderNotFound,
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
			mockStorage.On("GetOrdersByUser", tt.ctx, tt.user).Return(tt.want, tt.wantErr)
			UsedStorage = mockStorage
			got, err := GetOrdersByUser(tt.ctx, tt.user)
			require.Equal(t, tt.want, got)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestUpdateOrderFromAccrual(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		order   *types.OrderFromAccrual
		wantErr error
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			order:   new(types.OrderFromAccrual),
			wantErr: nil,
		},
		{
			name:    "order not found",
			ctx:     context.Background(),
			order:   new(types.OrderFromAccrual),
			wantErr: types.ErrOrderNotFound,
		},
		{
			name:    "unexpected error",
			ctx:     context.Background(),
			order:   new(types.OrderFromAccrual),
			wantErr: errors.New("some unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(mocks.Storage)
			mockStorage.On("UpdateOrderFromAccrual", tt.ctx, tt.order).Return(tt.wantErr)
			UsedStorage = mockStorage
			err := UpdateOrderFromAccrual(tt.ctx, tt.order)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
