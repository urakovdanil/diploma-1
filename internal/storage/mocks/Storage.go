// Code generated by mockery v2.43.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	types "diploma-1/internal/types"
)

// Storage is an autogenerated mock type for the Storage type
type Storage struct {
	mock.Mock
}

// CreateOrder provides a mock function with given fields: ctx, order
func (_m *Storage) CreateOrder(ctx context.Context, order *types.Order) (*types.Order, error) {
	ret := _m.Called(ctx, order)

	if len(ret) == 0 {
		panic("no return value specified for CreateOrder")
	}

	var r0 *types.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.Order) (*types.Order, error)); ok {
		return rf(ctx, order)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.Order) *types.Order); ok {
		r0 = rf(ctx, order)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.Order) error); ok {
		r1 = rf(ctx, order)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateUser provides a mock function with given fields: ctx, user
func (_m *Storage) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 *types.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.User) (*types.User, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.User) *types.User); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBalanceByUser provides a mock function with given fields: ctx, user
func (_m *Storage) GetBalanceByUser(ctx context.Context, user *types.User) (*types.Balance, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for GetBalanceByUser")
	}

	var r0 *types.Balance
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.User) (*types.Balance, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.User) *types.Balance); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Balance)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrdersByUser provides a mock function with given fields: ctx, user
func (_m *Storage) GetOrdersByUser(ctx context.Context, user *types.User) ([]types.Order, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for GetOrdersByUser")
	}

	var r0 []types.Order
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.User) ([]types.Order, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.User) []types.Order); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.Order)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByLogin provides a mock function with given fields: ctx, login
func (_m *Storage) GetUserByLogin(ctx context.Context, login string) (*types.User, error) {
	ret := _m.Called(ctx, login)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByLogin")
	}

	var r0 *types.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*types.User, error)); ok {
		return rf(ctx, login)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *types.User); ok {
		r0 = rf(ctx, login)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, login)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWithdrawalsByUser provides a mock function with given fields: ctx, user
func (_m *Storage) GetWithdrawalsByUser(ctx context.Context, user *types.User) ([]types.WithdrawWithTS, error) {
	ret := _m.Called(ctx, user)

	if len(ret) == 0 {
		panic("no return value specified for GetWithdrawalsByUser")
	}

	var r0 []types.WithdrawWithTS
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.User) ([]types.WithdrawWithTS, error)); ok {
		return rf(ctx, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *types.User) []types.WithdrawWithTS); ok {
		r0 = rf(ctx, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]types.WithdrawWithTS)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *types.User) error); ok {
		r1 = rf(ctx, user)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOrderFromAccrual provides a mock function with given fields: ctx, order
func (_m *Storage) UpdateOrderFromAccrual(ctx context.Context, order *types.OrderFromAccrual) error {
	ret := _m.Called(ctx, order)

	if len(ret) == 0 {
		panic("no return value specified for UpdateOrderFromAccrual")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.OrderFromAccrual) error); ok {
		r0 = rf(ctx, order)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// WithdrawByUser provides a mock function with given fields: ctx, user, withdraw
func (_m *Storage) WithdrawByUser(ctx context.Context, user *types.User, withdraw *types.Withdraw) error {
	ret := _m.Called(ctx, user, withdraw)

	if len(ret) == 0 {
		panic("no return value specified for WithdrawByUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *types.User, *types.Withdraw) error); ok {
		r0 = rf(ctx, user, withdraw)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStorage creates a new instance of Storage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStorage(t interface {
	mock.TestingT
	Cleanup(func())
}) *Storage {
	mock := &Storage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}