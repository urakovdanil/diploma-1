package accrual

import (
	"context"
	"diploma-1/internal/config"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
	client = nil
}

func TestNew(t *testing.T) {
	require.Nil(t, client)
	ctx := context.Background()
	err := config.New(ctx)
	require.NoError(t, err)

	err = New(ctx, config.Applied)
	require.NoError(t, err)

	require.NotNil(t, client)
	require.NotNil(t, client.ordersToBeUpdated)
	require.IsType(t, make(chan *order, 100), client.ordersToBeUpdated)
	require.True(t, client.canSend)
	require.NotNil(t, client.notificationChan)
	require.IsType(t, make(chan struct{}), client.notificationChan)
	require.NotNil(t, client.client)
	require.IsType(t, client.client, &resty.Client{})
}
