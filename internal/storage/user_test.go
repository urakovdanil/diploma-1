package storage

import (
	"context"
	"diploma-1/internal/storage/mocks"
	"diploma-1/internal/types"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetUserByLogin(t *testing.T) {

	tests := []struct {
		name    string
		ctx     context.Context
		login   string
		want    *types.User
		wantErr error
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			login:   "test",
			want:    new(types.User),
			wantErr: nil,
		},
		{
			name:    "user not found",
			ctx:     context.Background(),
			login:   "test",
			want:    nil,
			wantErr: types.ErrUserNotFound,
		},
		{
			name:    "unexpected error",
			ctx:     context.Background(),
			login:   "test",
			want:    nil,
			wantErr: errors.New("some unexpected error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := new(mocks.Storage)
			mockStorage.On("GetUserByLogin", tt.ctx, tt.login).Return(tt.want, tt.wantErr)
			UsedStorage = mockStorage
			got, err := GetUserByLogin(tt.ctx, tt.login)
			require.Equal(t, tt.want, got)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCreateUser(t *testing.T) {

	tests := []struct {
		name    string
		ctx     context.Context
		user    *types.User
		want    *types.User
		wantErr error
	}{
		{
			name:    "success",
			ctx:     context.Background(),
			user:    new(types.User),
			want:    new(types.User),
			wantErr: nil,
		},
		{
			name:    "user already exists",
			ctx:     context.Background(),
			user:    new(types.User),
			want:    nil,
			wantErr: types.ErrUserAlreadyExists,
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
			mockStorage.On("CreateUser", tt.ctx, tt.user).Return(tt.want, tt.wantErr)
			UsedStorage = mockStorage
			got, err := CreateUser(tt.ctx, tt.user)
			require.Equal(t, tt.want, got)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
