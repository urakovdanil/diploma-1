package balance

import (
	"bytes"
	"context"
	"diploma-1/internal/api/auth"
	"diploma-1/internal/api/middleware"
	"diploma-1/internal/config"
	"diploma-1/internal/storage"
	"diploma-1/internal/storage/mocks"
	"diploma-1/internal/types"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBalance_GetBalanceHandlerFunc(t *testing.T) {
	tests := []struct {
		name              string
		user              *types.User
		storageUserErr    error
		storageUser       *types.User
		storageBalance    *types.Balance
		storageBalanceErr error
		statusCode        int
	}{
		{
			name:              "no data for user",
			user:              &types.User{Login: "login", Password: "password"},
			storageUser:       &types.User{Login: "login", Password: "password"},
			storageUserErr:    nil,
			storageBalance:    nil,
			storageBalanceErr: types.ErrUserNotFound,
			statusCode:        http.StatusOK,
		},
		{
			name:              "unexpected storage error",
			user:              &types.User{Login: "login", Password: "password"},
			storageUser:       &types.User{Login: "login", Password: "password"},
			storageUserErr:    nil,
			storageBalance:    nil,
			storageBalanceErr: errors.New("some unexpected error"),
			statusCode:        http.StatusInternalServerError,
		},
		{
			name:              "success",
			user:              &types.User{Login: "login", Password: "password"},
			storageUser:       &types.User{Login: "login", Password: "password"},
			storageUserErr:    nil,
			storageBalance:    &types.Balance{Current: 100, Withdrawn: 0},
			storageBalanceErr: nil,
			statusCode:        http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, config.New(context.Background()))
			mockStorage := new(mocks.Storage)
			mockStorage.On("GetUserByLogin", mock.Anything, tt.user.Login).Return(tt.storageUser, tt.storageUserErr)
			mockStorage.On("GetBalanceByUser", mock.Anything, tt.user).Return(tt.storageBalance, tt.storageBalanceErr)
			storage.UsedStorage = mockStorage

			router := chi.NewRouter()
			a := New()
			router.Use(middleware.IsAuthenticated)
			router.Get(Path, a.GetBalanceHandlerFunc)
			router.Post(auth.LoginPath, auth.New().LoginHandlerFunc)
			ts := httptest.NewServer(router)
			defer ts.Close()

			pre, err := json.Marshal(tt.user)
			require.NoError(t, err)
			body := bytes.NewBuffer(pre)
			request, err := http.NewRequest(http.MethodPost, ts.URL+auth.LoginPath, body)
			require.NoError(t, err)

			res, err := ts.Client().Do(request)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
			require.Contains(t, res.Header.Get("Authorization"), "Bearer ")

			request, err = http.NewRequest(http.MethodGet, ts.URL+Path, nil)
			require.NoError(t, err)
			request.Header.Set("Authorization", res.Header.Get("Authorization"))

			res, err = ts.Client().Do(request)
			require.NoError(t, err)
			require.Equal(t, tt.statusCode, res.StatusCode)
		})
	}
}

func TestBalance_WithdrawHandlerFunc(t *testing.T) {
	tests := []struct {
		name               string
		user               *types.User
		withdrawal         *types.Withdraw
		storageUserErr     error
		storageUser        *types.User
		storageWithdrawErr error
		statusCode         int
	}{
		{
			name:               "order already exists",
			withdrawal:         &types.Withdraw{Order: "123", Sum: 10},
			user:               &types.User{Login: "login", Password: "password"},
			storageUser:        &types.User{Login: "login", Password: "password"},
			storageUserErr:     nil,
			storageWithdrawErr: types.ErrOrderAlreadyExistsForThisUser,
			statusCode:         http.StatusUnprocessableEntity,
		},
		{
			name:               "invalid order number",
			withdrawal:         &types.Withdraw{Order: "123", Sum: 10},
			user:               &types.User{Login: "login", Password: "password"},
			storageUser:        &types.User{Login: "login", Password: "password"},
			storageUserErr:     nil,
			storageWithdrawErr: nil,
			statusCode:         http.StatusUnprocessableEntity,
		},
		{
			name:               "insufficient funds",
			withdrawal:         &types.Withdraw{Order: "371449635398431", Sum: 10},
			user:               &types.User{Login: "login", Password: "password"},
			storageUser:        &types.User{Login: "login", Password: "password"},
			storageUserErr:     nil,
			storageWithdrawErr: types.ErrInsufficientFunds,
			statusCode:         http.StatusPaymentRequired,
		},
		{
			name:               "unexpected storage error",
			withdrawal:         &types.Withdraw{Order: "371449635398431", Sum: 10},
			user:               &types.User{Login: "login", Password: "password"},
			storageUser:        &types.User{Login: "login", Password: "password"},
			storageUserErr:     nil,
			storageWithdrawErr: errors.New("some unexpected error"),
			statusCode:         http.StatusInternalServerError,
		},
		{
			name:               "success",
			withdrawal:         &types.Withdraw{Order: "371449635398431", Sum: 10},
			user:               &types.User{Login: "login", Password: "password"},
			storageUser:        &types.User{Login: "login", Password: "password"},
			storageUserErr:     nil,
			storageWithdrawErr: nil,
			statusCode:         http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, config.New(context.Background()))
			mockStorage := new(mocks.Storage)
			mockStorage.On("GetUserByLogin", mock.Anything, tt.user.Login).Return(tt.storageUser, tt.storageUserErr)
			mockStorage.On("WithdrawByUser", mock.Anything, tt.user, tt.withdrawal).Return(tt.storageWithdrawErr)
			storage.UsedStorage = mockStorage

			router := chi.NewRouter()
			a := New()
			router.Use(middleware.IsAuthenticated)
			router.Post(Path, a.WithdrawHandlerFunc)
			router.Post(auth.LoginPath, auth.New().LoginHandlerFunc)
			ts := httptest.NewServer(router)
			defer ts.Close()

			pre, err := json.Marshal(tt.user)
			require.NoError(t, err)
			body := bytes.NewBuffer(pre)
			request, err := http.NewRequest(http.MethodPost, ts.URL+auth.LoginPath, body)
			require.NoError(t, err)

			res, err := ts.Client().Do(request)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, res.StatusCode)
			require.Contains(t, res.Header.Get("Authorization"), "Bearer ")

			pre, err = json.Marshal(tt.withdrawal)
			require.NoError(t, err)
			body = bytes.NewBuffer(pre)
			request, err = http.NewRequest(http.MethodPost, ts.URL+Path, body)
			require.NoError(t, err)
			request.Header.Set("Authorization", res.Header.Get("Authorization"))

			res, err = ts.Client().Do(request)
			require.NoError(t, err)
			require.Equal(t, tt.statusCode, res.StatusCode)
		})
	}
}
