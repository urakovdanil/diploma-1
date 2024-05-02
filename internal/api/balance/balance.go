package balance

import (
	"diploma-1/internal/storage"
	"diploma-1/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	Path = "/api/user/balance"
)

type Balance struct{}

func New() *Balance {
	return &Balance{}
}

func (b *Balance) GetBalanceHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()
	user := ctx.Value(types.CtxKeyUser).(*types.User)
	if user == nil {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}
	balance, err := storage.GetBalanceByUser(ctx, user)
	if err != nil {
		if errors.Is(err, types.ErrUserNotFound) {
			balance = &types.Balance{}
		} else {
			http.Error(w, fmt.Sprintf("unable to get balance: %v", err), http.StatusInternalServerError)
		}
	}
	if err = json.NewEncoder(w).Encode(balance); err != nil {
		http.Error(w, fmt.Sprintf("unable to encode balance: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
