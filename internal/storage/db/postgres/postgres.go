package postgres

import (
	"context"
	"diploma-1/internal/config"
	"diploma-1/internal/logger"
	"diploma-1/internal/types"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"reflect"
	"runtime"
	"strings"
)

const (
	uniqueConstraintViolationCode = "23505"
)

var readCommittedTXOptions = pgx.TxOptions{
	IsoLevel: pgx.ReadCommitted,
}

type Storage struct {
	pool *pgxpool.Pool
}

func getFunctionName(f func(tx pgx.Tx) error) string {
	fullFuncName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	parts := strings.Split(fullFuncName, ".")
	return parts[len(parts)-2]
}

func (s *Storage) withTx(ctx context.Context, txOptions pgx.TxOptions, f func(tx pgx.Tx) error) error {
	fName := getFunctionName(f)
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		logger.Errorf(ctx, "unable to acquire connection for decorated %s: %v", fName, err)
		return err
	}
	defer conn.Release()
	tx, err := conn.BeginTx(ctx, txOptions)
	if err != nil {
		logger.Errorf(ctx, "unable to begin tx for decorated %s: %v", fName, err)
		return err
	}
	logger.Debugf(ctx, "tx started successfully for decorated %s", fName)
	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			logger.Errorf(ctx, "error on tx rollback in decorated %s: %v", fName, err)
			return
		}
		logger.Debugf(ctx, "tx exited successfully in decorated %s", fName)
	}()
	if err = f(tx); err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		logger.Errorf(ctx, "unable to commit tx for decorated %s: %v", fName, err)
		return err
	}
	return nil
}

func (s *Storage) GetUserByLogin(ctx context.Context, login string) (*types.User, error) {
	var user types.User
	if err := s.withTx(ctx, readCommittedTXOptions, func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx, userSelectByLogin, login)
		if err := row.Scan(&user.ID, &user.Login, &user.Password); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return types.ErrUserNotFound
			}
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Storage) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	if err := s.withTx(ctx, readCommittedTXOptions, func(tx pgx.Tx) error {
		row := tx.QueryRow(ctx, userInsert, user.Login, user.Password)
		if err := row.Scan(&user.ID); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == uniqueConstraintViolationCode {
				return types.ErrUserAlreadyExists
			}
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Storage) CreateOrder(ctx context.Context, order *types.Order) (*types.Order, error) {
	if err := s.withTx(ctx, readCommittedTXOptions, func(tx pgx.Tx) error {
		existingOrder := &types.Order{}
		err := tx.QueryRow(ctx, orderSelect, order.Number).Scan(&existingOrder.ID, &existingOrder.Number, &existingOrder.Status, &existingOrder.Accrual, &existingOrder.UserID)
		if err == nil {
			if existingOrder.UserID == order.UserID {
				order = existingOrder
				return types.ErrOrderAlreadyExistsForThisUser
			}
			return types.ErrOrderAlreadyExistsForAnotherUser
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			fmt.Println("HERE")
			return err
		}
		if err = tx.QueryRow(ctx, orderInsert, order.Number, order.Status, order.Accrual, order.UserID).Scan(&order.ID, &order.Number, &order.Status, &order.Accrual, &order.UserID); err != nil {
			// код выше должен гарантировать отсутствие дубликатов, проверка ниже для перестраховки
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == uniqueConstraintViolationCode {
				return types.ErrOrderAlreadyExistsForThisUser
			}
			return err
		}
		return nil
	}); err != nil {
		return order, err
	}
	return order, nil
}

func New(ctx context.Context, su *config.StartUp) (*Storage, error) {
	if err := migrateUp(su); err != nil {
		return nil, err
	}
	db, err := pgxpool.New(context.Background(), su.GetDatabaseURI())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &Storage{pool: db}, nil
}
