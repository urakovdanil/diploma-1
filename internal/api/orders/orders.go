package orders

import (
	"bytes"
	"diploma-1/internal/accrual"
	"diploma-1/internal/storage"
	"diploma-1/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

const (
	Path = "/api/user/orders"
)

type Orders struct{}

func New() *Orders {
	return &Orders{}
}

func (o *Orders) validateInput(order string) error {
	if order == "" {
		return types.ErrEmptyOrderNumber
	}
	sum := 0
	alternate := false
	for i := len(order) - 1; i >= 0; i-- {
		digit, err := strconv.Atoi(string(order[i]))
		if err != nil {
			return types.ErrNonDigitalOrderNumber
		}
		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
		alternate = !alternate
	}
	if sum%10 != 0 {
		return types.ErrInvalidOrderNumber
	}
	return nil
}

func (o *Orders) CreateOrderHandlerFunc(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		http.Error(w, fmt.Sprintf("unable to read request body: %v", err), http.StatusBadRequest)
		return
	}
	orderNumber := string(buf.Bytes())
	if err := o.validateInput(orderNumber); err != nil {
		if errors.Is(err, types.ErrInvalidOrderNumber) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := r.Context().Value(types.CtxKeyUser).(*types.User)
	if user == nil {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}
	order := &types.Order{
		Number: orderNumber,
		UserID: user.ID,
		Status: types.OrderStatusNew,
	}
	order, err := storage.CreateOrder(ctx, order)
	if err != nil {
		if errors.Is(err, types.ErrOrderAlreadyExistsForThisUser) {
			w.WriteHeader(http.StatusOK)
			return
		}
		if errors.Is(err, types.ErrOrderAlreadyExistsForAnotherUser) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("unable to create order: %v", err), http.StatusInternalServerError)
		return
	}
	accrual.Track(ctx, order)
	w.WriteHeader(http.StatusAccepted)
}

func (o *Orders) GetOrdersHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	user := ctx.Value(types.CtxKeyUser).(*types.User)
	if user == nil {
		http.Error(w, "user not found in context", http.StatusInternalServerError)
		return
	}

	orders, err := storage.GetOrdersByUser(ctx, user)
	if err != nil {
		if errors.Is(err, types.ErrOrderNotFound) {
			http.Error(w, err.Error(), http.StatusNoContent)
			return
		}
		http.Error(w, fmt.Sprintf("unable to get orders: %v", err), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, fmt.Sprintf("unable to encode orders: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
