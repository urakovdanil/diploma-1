package types

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type Withdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}
