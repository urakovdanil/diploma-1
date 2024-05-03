package balance

import (
	"diploma-1/internal/api/orders"
	"diploma-1/internal/storage"
	"diploma-1/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	Path             = "/api/user/balance"
	WithdrawPath     = "/api/user/balance/withdraw"
	WithdrawListPath = "/api/user/withdrawals"
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
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(balance); err != nil {
		http.Error(w, fmt.Sprintf("unable to encode balance: %v", err), http.StatusInternalServerError)
		return
	}
}

func (b *Balance) WithdrawHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	user := ctx.Value(types.CtxKeyUser).(*types.User)
	if user == nil {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}
	var withdraw types.Withdraw
	if err := json.NewDecoder(r.Body).Decode(&withdraw); err != nil {
		http.Error(w, fmt.Sprintf("unable to decode withdraw: %v", err), http.StatusBadRequest)
		return
	}
	if err := orders.New().ValidateInput(withdraw.Order); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	if err := storage.WithdrawByUser(ctx, user, &withdraw); err != nil {
		if errors.Is(err, types.ErrInsufficientFunds) {
			http.Error(w, err.Error(), http.StatusPaymentRequired)
			return
		}
		if errors.Is(err, types.ErrOrderAlreadyExistsForThisUser) || errors.Is(err, types.ErrOrderAlreadyExistsForAnotherUser) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, fmt.Sprintf("unable to withdraw balance: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (b *Balance) GetWithdrawalsHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	user := ctx.Value(types.CtxKeyUser).(*types.User)
	if user == nil {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}
	withdrawals, err := storage.GetWithdrawalsByUser(ctx, user)
	if err != nil {
		if errors.Is(err, types.ErrOrderNotFound) {
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}
		http.Error(w, fmt.Sprintf("unable to get withdrawals: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(withdrawals); err != nil {
		http.Error(w, fmt.Sprintf("unable to encode withdrawals: %v", err), http.StatusInternalServerError)
		return
	}
}
