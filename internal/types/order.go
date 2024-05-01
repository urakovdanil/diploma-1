package types

import "time"

type Order struct {
	ID        int64       `json:"-"`
	Number    string      `json:"number"`
	Status    OrderStatus `json:"status"`
	Accrual   int64       `json:"accrual,omitempty"`
	UserID    int64       `json:"-"`
	CreatedAt time.Time   `json:"uploaded_at"`
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusRegistered OrderStatus = "REGISTERED"
)
