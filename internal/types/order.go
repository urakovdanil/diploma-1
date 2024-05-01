package types

type Order struct {
	ID      int64
	Number  string
	Status  OrderStatus
	Accrual int64
	UserID  int64
}

type OrderStatus string

const (
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusProcessed  OrderStatus = "PROCESSED"
	OrderStatusInvalid    OrderStatus = "INVALID"
	OrderStatusRegistered OrderStatus = "REGISTERED"
)
