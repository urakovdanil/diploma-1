package postgres

import (
	"context"
	"diploma-1/internal/config"
	"diploma-1/internal/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	pool *pgxpool.Pool
}

func (s *Storage) GetUser(ctx context.Context, login string) (*types.User, error) {
	// TODO: implement
	return nil, nil
}

func New(ctx context.Context, su *config.StartUp) (*Storage, error) {
	if err := migrateUp(su); err != nil {
		return nil, err
	}
	db, err := pgxpool.New(context.Background(), su.DatabaseURI)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &Storage{pool: db}, nil
}
