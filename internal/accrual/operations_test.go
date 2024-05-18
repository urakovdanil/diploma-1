package accrual

import (
	"context"
	"diploma-1/internal/config"
	"diploma-1/internal/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTrack(t *testing.T) {
	err := newMock(context.Background(), config.Applied)
	require.NoError(t, err)

	ord := new(types.Order)

	Track(context.Background(), ord)
	fromChan, ok := <-client.ordersToBeUpdated
	require.True(t, ok)
	require.Equal(t, ord, fromChan.Order)
	require.Equal(t, "", fromChan.requestID)
}
