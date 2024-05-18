package types

import "time"

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type WithdrawWithTS struct {
	*Withdraw
	ProcessedAt time.Time `json:"processed_at"`
}
