package accrual

import "testing"

func TestNewBackoff(t *testing.T) {
	bo := newBackoff()
	if bo.InitialInterval != defaultBackoffInitialInterval {
		t.Errorf("expected %v got %v", defaultBackoffInitialInterval, bo.InitialInterval)
	}
	if bo.MaxInterval != defaultBackoffMaxInterval {
		t.Errorf("expected %v got %v", defaultBackoffMaxInterval, bo.MaxInterval)
	}
	if bo.Multiplier != defaultBackoffMultiplier {
		t.Errorf("expected %v got %v", defaultBackoffMultiplier, bo.Multiplier)
	}
	if bo.MaxElapsedTime != defaultBackoffMaxElapsedTime {
		t.Errorf("expected %v got %v", defaultBackoffMaxElapsedTime, bo.MaxElapsedTime)
	}
}
