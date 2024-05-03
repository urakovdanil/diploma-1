package postgres

import (
	"context"
	"diploma-1/internal/closer"
	"diploma-1/internal/config"
	"diploma-1/internal/logger"
	"diploma-1/internal/types"
	"errors"
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

func getFunctionName(f func(tx pgx.Tx) error) string {
	fullFuncName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	parts := strings.Split(fullFuncName, ".")
	return parts[len(parts)-2]
}

type Storage struct {
	pool *pgxpool.Pool
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

func (s *Storage) GetOrdersByUser(ctx context.Context, user *types.User) ([]types.Order, error) {
	orders := make([]types.Order, 0, 10)
	if err := s.withTx(ctx, readCommittedTXOptions, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, ordersListByUser, user.ID)
		if err != nil {
			return err
		}
		for rows.Next() {
			order := types.Order{}
			if err = rows.Scan(&order.ID, &order.Number, &order.Status, &order.Accrual, &order.UserID, &order.CreatedAt); err != nil {
				return err
			}
			orders = append(orders, order)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if len(orders) == 0 {
		return nil, types.ErrOrderNotFound
	}
	return orders, nil
}

func (s *Storage) UpdateOrderFromAccrual(ctx context.Context, order *types.OrderFromAccrual) error {
	if err := s.withTx(ctx, readCommittedTXOptions, func(tx pgx.Tx) error {
		ct, err := tx.Exec(ctx, orderUpdateFromAccrual, order.Status, order.Accrual, order.Number)
		if ct.RowsAffected() == 0 {
			return types.ErrOrderNotFound
		}
		return err
	}); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetBalanceByUser(ctx context.Context, user *types.User) (*types.Balance, error) {
	balance := &types.Balance{}
	if err := s.withTx(ctx, readCommittedTXOptions, func(tx pgx.Tx) error {
		if err := tx.QueryRow(ctx, balanceSelectByUser, user.ID).Scan(&balance.Current, &balance.Withdrawn); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return types.ErrUserNotFound
			}
			return err
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return balance, nil
}

func (s *Storage) WithdrawByUser(ctx context.Context, user *types.User, withdraw *types.Withdraw) error {
	if err := s.withTx(ctx, readCommittedTXOptions, func(tx pgx.Tx) error {
		existingOrder := &types.Order{}
		err := tx.QueryRow(ctx, orderSelect, withdraw.Order).Scan(&existingOrder.ID, &existingOrder.Number, &existingOrder.Status, &existingOrder.Accrual, &existingOrder.UserID)
		if err == nil {
			if existingOrder.UserID == user.ID {
				return types.ErrOrderAlreadyExistsForThisUser
			}
			return types.ErrOrderAlreadyExistsForAnotherUser
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		balance := &types.Balance{}
		if err := tx.QueryRow(ctx, balanceSelectByUser, user.ID).Scan(&balance.Current, &balance.Withdrawn); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return types.ErrUserNotFound
			}
			return err
		}
		if balance.Current < withdraw.Sum {
			return types.ErrInsufficientFunds
		}
		if _, err := tx.Exec(ctx, balanceWithdrawByUser, withdraw.Order, types.OrderStatusProcessed, -withdraw.Sum, user.ID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetWithdrawalsByUser(ctx context.Context, user *types.User) ([]types.WithdrawWithTS, error) {
	res := make([]types.WithdrawWithTS, 0, 10)
	if err := s.withTx(ctx, readCommittedTXOptions, func(tx pgx.Tx) error {
		rows, err := tx.Query(ctx, withdrawalsByUser, user.ID)
		if err != nil {
			return err
		}
		for rows.Next() {
			if rows.Err() != nil {
				return err
			}
			withdraw := types.WithdrawWithTS{Withdraw: &types.Withdraw{}}
			if err = rows.Scan(&withdraw.Order, &withdraw.Sum, &withdraw.ProcessedAt); err != nil {
				return err
			}
			res = append(res, withdraw)
		}
		return nil
	}); err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, types.ErrOrderNotFound
	}
	return res, nil
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
	closer.Add(func() error {
		db.Close()
		return nil
	})
	return &Storage{pool: db}, nil
}
