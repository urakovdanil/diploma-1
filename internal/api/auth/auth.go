package auth

import (
	"bytes"
	"context"
	"diploma-1/internal/config"
	"diploma-1/internal/logger"
	"diploma-1/internal/storage"
	"diploma-1/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

const (
	LoginPath    = "/api/user/login"
	RegisterPath = "/api/user/register"
)

func authenticate(ctx context.Context, user *types.User) error {
	fromStorage, err := storage.GetUserByLogin(ctx, user.Login)
	if err != nil {
		return err
	}

	if fromStorage.Password != user.Password {
		return types.ErrWrongPassword
	}

	return nil
}

func generateToken(su *config.StartUp, user *types.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user.Login
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(su.GetJWTTokenTTLMinutes())).Unix()

	tokenString, err := token.SignedString(su.GetJWTSecretKey())
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

type Auth struct{}

func New() *Auth {
	return &Auth{}
}

func (a *Auth) validateInput(user *types.User) error {
	if user.Login == "" || user.Password == "" {
		return types.ErrInvalidAuthInput
	}
	return nil
}

func (a *Auth) parseBody(r *http.Request) (*types.User, error) {
	user := &types.User{}
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		return nil, fmt.Errorf("unable to read request body: %v", err)
	}
	if err := json.Unmarshal(buf.Bytes(), user); err != nil {
		return nil, fmt.Errorf("unable to unmarshal request body: %v", err)
	}
	return user, nil
}

func (a *Auth) RegisterHandlerFunc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	w.Header().Set("Content-Type", "application/json")

	user, err := a.parseBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.validateInput(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = storage.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, types.ErrUserAlreadyExists) {
			http.Error(w, fmt.Sprintf("user with login %s already exists", user.Login), http.StatusConflict)
			return
		}
		logger.Errorf(ctx, "unexpected error on user creation: %v", err)
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	token, err := generateToken(config.Applied, user)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	if _, err = w.Write([]byte(token)); err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}

func (a *Auth) LoginHandlerFunc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := a.parseBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.validateInput(user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := authenticate(ctx, user); err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			http.Error(w, fmt.Sprintf("user with login %s does not exist", user.Login), http.StatusUnauthorized)
			return
		}
		if errors.Is(err, types.ErrWrongPassword) {
			http.Error(w, fmt.Sprintf("wrong password for user %s", user.Login), http.StatusUnauthorized)
			return
		}
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	token, err := generateToken(config.Applied, user)
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
	if _, err = w.Write([]byte(token)); err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}
}
