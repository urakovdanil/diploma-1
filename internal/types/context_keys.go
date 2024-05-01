package types

type CtxKey string

var (
	CtxKeyUser      CtxKey = "user"
	CtxKeyRequestID CtxKey = "internal_request_id"
)
