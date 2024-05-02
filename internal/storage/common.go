package storage

import (
	"context"
	"diploma-1/internal/config"
	"diploma-1/internal/storage/db/postgres"
	"diploma-1/internal/types"
)

type Storage interface {
	CreateUser(ctx context.Context, user *types.User) (*types.User, error)
	GetUserByLogin(ctx context.Context, login string) (*types.User, error)

	CreateOrder(ctx context.Context, order *types.Order) (*types.Order, error)
	GetOrdersByUser(ctx context.Context, user *types.User) ([]types.Order, error)
	UpdateOrderFromAccrual(ctx context.Context, order *types.OrderFromAccrual) error
}

var UsedStorage Storage

func New(ctx context.Context, su *config.StartUp) error {
	var err error
	UsedStorage, err = postgres.New(ctx, su)
	return err
}
