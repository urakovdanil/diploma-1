package storage

import (
	"context"
	"diploma-1/internal/config"
	"diploma-1/internal/storage/db/postgres"
	"diploma-1/internal/types"
)

type Storage interface {
	GetUser(ctx context.Context, login string) (*types.User, error)
}

var UsedStorage Storage

func New(ctx context.Context, su *config.StartUp) error {
	var err error
	UsedStorage, err = postgres.New(ctx, su)
	return err
}
