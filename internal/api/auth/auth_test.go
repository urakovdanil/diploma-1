package auth

import (
	"bytes"
	"context"
	"diploma-1/internal/config"
	"diploma-1/internal/storage"
	"diploma-1/internal/storage/mocks"
	"diploma-1/internal/types"
	"encoding/json"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestValidateInput(t *testing.T) {
	tests := []struct {
		name    string
		user    *types.User
		wantErr error
	}{
		{
			name:    "success",
			user:    &types.User{Login: "login", Password: "password"},
			wantErr: nil,
		},
		{
			name:    "empty login",
			user:    &types.User{Login: "", Password: "password"},
			wantErr: types.ErrInvalidAuthInput,
		},
		{
			name:    "empty password",
			user:    &types.User{Login: "login", Password: ""},
			wantErr: types.ErrInvalidAuthInput,
		},
		{
			name:    "empty login and password",
			user:    &types.User{Login: "", Password: ""},
			wantErr: types.ErrInvalidAuthInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New()
			err := a.validateInput(tt.user)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestAuthenticate(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		user           *types.User
		fromStorage    *types.User
		fromStorageErr error
		wantErr        error
	}{
		{
			name:           "success",
			ctx:            context.Background(),
			user:           &types.User{Login: "login", Password: "password"},
			fromStorage:    &types.User{Login: "login", Password: "password"},
			fromStorageErr: nil,
			wantErr:        nil,
		},
		{
			name:           "user not found",
			ctx:            context.Background(),
			user:           &types.User{Login: "login", Password: "password"},
			fromStorage:    nil,
			fromStorageErr: types.ErrUserNotFound,
			wantErr:        types.ErrUserNotFound,
		},
		{
			name:           "wrong password",
			ctx:            context.Background(),
			user:           &types.User{Login: "login", Password: "wrong password"},
			fromStorage:    &types.User{Login: "login", Password: "password"},
			fromStorageErr: nil,
			wantErr:        types.ErrWrongPassword,
		},
		{
			name:           "unexpected error",
			ctx:            context.Background(),
			user:           &types.User{Login: "login", Password: "password"},
			fromStorage:    nil,
			fromStorageErr: errors.New("some unexpected error"),
			wantErr:        errors.New("some unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(mocks.Storage)
			mockStorage.On("GetUserByLogin", tt.ctx, tt.user.Login).Return(tt.fromStorage, tt.fromStorageErr)
			storage.UsedStorage = mockStorage
			err := authenticate(tt.ctx, tt.user)
			if tt.name == "unexpected error" {
				require.ErrorContains(t, err, tt.fromStorageErr.Error())
			} else {
				require.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name    string
		user    *types.User
		wantErr error
	}{
		{
			name:    "success",
			user:    &types.User{Login: "login", Password: "password"},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.New(context.Background())
			require.NoError(t, err)
			got, err := generateToken(config.Applied, tt.user)
			require.NoError(t, err)
			require.NotEmpty(t, got)

			token, err := jwt.Parse(got, func(token *jwt.Token) (interface{}, error) {
				return config.Applied.GetJWTSecretKey(), nil
			})
			require.NoError(t, err)
			require.True(t, token.Valid)
			claims := token.Claims.(jwt.MapClaims)
			require.LessOrEqual(t, float64(time.Now().Unix()), claims["exp"].(float64))
			require.Equal(t, claims["username"], tt.user.Login)
		})
	}
}

func TestAuth_RegisterHandlerFunc(t *testing.T) {
	tests := []struct {
		name       string
		user       *types.User
		statusCode int
		storageErr error
	}{
		{
			name:       "success",
			user:       &types.User{Login: "login", Password: "password"},
			statusCode: http.StatusOK,
		},
		{
			name:       "invalid body",
			user:       &types.User{Login: "login", Password: ""},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "empty password",
			user:       &types.User{Login: "login", Password: ""},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "empty login",
			user:       &types.User{Login: "", Password: "password"},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "empty login and password",
			user:       &types.User{Login: "", Password: ""},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "user already exists",
			user:       &types.User{Login: "login", Password: "password"},
			statusCode: http.StatusConflict,
			storageErr: types.ErrUserAlreadyExists,
		},
		{
			name:       "unexpected error",
			user:       &types.User{Login: "login", Password: "password"},
			statusCode: http.StatusInternalServerError,
			storageErr: errors.New("some unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, config.New(context.Background()))
			mockStorage := new(mocks.Storage)
			mockStorage.On("CreateUser", mock.Anything, tt.user).Return(tt.user, tt.storageErr)
			storage.UsedStorage = mockStorage

			router := chi.NewRouter()
			a := New()
			router.Post(RegisterPath, a.RegisterHandlerFunc)
			ts := httptest.NewServer(router)
			defer ts.Close()

			pre, err := json.Marshal(tt.user)
			require.NoError(t, err)
			body := bytes.NewBuffer(pre)
			if tt.name == "invalid body" {
				body.WriteString("{")
			}
			request, err := http.NewRequest(http.MethodPost, ts.URL+RegisterPath, body)

			res, err := ts.Client().Do(request)
			require.NoError(t, err)
			require.Equal(t, tt.statusCode, res.StatusCode)
			if tt.statusCode == http.StatusOK {
				require.Equal(t, "application/json", res.Header.Get("Content-Type"))
				require.Contains(t, res.Header.Get("Authorization"), "Bearer ")
			} else {
				require.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
			}
			_ = res.Body.Close()
		})
	}
}

func TestAuth_LoginHandlerFunc(t *testing.T) {
	tests := []struct {
		name        string
		user        *types.User
		statusCode  int
		storageUser *types.User
		storageErr  error
	}{
		{
			name:        "success",
			user:        &types.User{Login: "login", Password: "password"},
			statusCode:  http.StatusOK,
			storageUser: &types.User{Login: "login", Password: "password"},
		},
		{
			name:       "invalid body",
			user:       &types.User{Login: "login", Password: ""},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "empty password",
			user:       &types.User{Login: "login", Password: ""},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "empty login",
			user:       &types.User{Login: "", Password: "password"},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "empty login and password",
			user:       &types.User{Login: "", Password: ""},
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "user not found",
			user:       &types.User{Login: "login", Password: "password"},
			statusCode: http.StatusUnauthorized,
			storageErr: types.ErrUserNotFound,
		},
		{
			name:        "wrong password",
			user:        &types.User{Login: "login", Password: "worng_password"},
			statusCode:  http.StatusUnauthorized,
			storageUser: &types.User{Login: "login", Password: "password"},
		},
		{
			name:       "unexpected error",
			user:       &types.User{Login: "login", Password: "password"},
			statusCode: http.StatusInternalServerError,
			storageErr: errors.New("some unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, config.New(context.Background()))
			mockStorage := new(mocks.Storage)
			mockStorage.On("GetUserByLogin", mock.Anything, tt.user.Login).Return(tt.storageUser, tt.storageErr)
			storage.UsedStorage = mockStorage

			router := chi.NewRouter()
			a := New()
			router.Post(LoginPath, a.LoginHandlerFunc)
			ts := httptest.NewServer(router)
			defer ts.Close()

			pre, err := json.Marshal(tt.user)
			require.NoError(t, err)
			body := bytes.NewBuffer(pre)
			if tt.name == "invalid body" {
				body.WriteString("{")
			}
			request, err := http.NewRequest(http.MethodPost, ts.URL+LoginPath, body)
			require.NoError(t, err)

			res, err := ts.Client().Do(request)
			require.NoError(t, err)
			require.Equal(t, tt.statusCode, res.StatusCode)
			if tt.statusCode == http.StatusOK {
				require.Contains(t, res.Header.Get("Authorization"), "Bearer ")
			} else {
				require.Equal(t, "text/plain; charset=utf-8", res.Header.Get("Content-Type"))
			}
			_ = res.Body.Close()
		})
	}
}
