package accrual

import "diploma-1/internal/types"

type order struct {
	*types.Order
	requestID string
}
