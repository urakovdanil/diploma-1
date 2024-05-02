package accrual

import (
	"github.com/cenkalti/backoff"
	"time"
)

const (
	defaultBackoffInitialInterval = 500 * time.Millisecond
	defaultBackoffMaxInterval     = 20 * time.Second
	defaultBackoffMultiplier      = 1.2
	defaultBackoffMaxElapsedTime  = 30 * time.Minute
)

func newBackoff() *backoff.ExponentialBackOff {
	bo := backoff.NewExponentialBackOff()
	bo.InitialInterval = defaultBackoffInitialInterval
	bo.MaxInterval = defaultBackoffMaxInterval
	bo.Multiplier = defaultBackoffMultiplier
	bo.MaxElapsedTime = defaultBackoffMaxElapsedTime
	return bo
}
