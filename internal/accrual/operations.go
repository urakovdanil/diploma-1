package accrual

import (
	"context"
	"diploma-1/internal/logger"
	"diploma-1/internal/types"
)

func Track(ctx context.Context, ord *types.Order) {
	internalOrder := &order{Order: ord}
	requestID := ctx.Value(types.CtxKeyRequestID)
	if requestID != nil {
		internalOrder.requestID = requestID.(string)
	}
	go func() {
		logger.Debugf(context.Background(), "adding order %s from request %s to queue", ord.Number, requestID)
		client.ordersToBeUpdated <- internalOrder
		logger.Debugf(context.Background(), "added order %s from request %s to queue", ord.Number, requestID)
	}()
}
