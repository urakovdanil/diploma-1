package postgres

import (
	"context"
	"diploma-1/internal/types"
)

type Storage struct {
	// TODO: implement
}

func (s *Storage) GetUser(ctx context.Context, login string) (*types.User, error) {
	// TODO: implement
	return nil, nil
}

func New(ctx context.Context) (*Storage, error) {
	// TODO: implement
	return nil, nil
}
