package types

import "time"

type Order struct {
	ID        int64       `json:"-"`
	Number    string      `json:"number"`
	Status    OrderStatus `json:"status"`
	Accrual   float64     `json:"accrual,omitempty"`
	UserID    int64       `json:"-"`
	CreatedAt time.Time   `json:"uploaded_at"`
}
type OrderFromAccrual struct {
	Number  string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual float64     `json:"accrual"`
}

type OrderStatus string

const (
	OrderStatusNew        OrderStatus = "NEW"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusRegistered OrderStatus = "REGISTERED"
)

var FinalOrderStatuses = map[OrderStatus]struct{}{
	OrderStatusProcessed: {},
	OrderStatusInvalid:   {},
}
