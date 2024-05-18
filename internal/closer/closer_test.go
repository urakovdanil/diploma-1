package closer

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNew(t *testing.T) {
	New()
	require.NotNil(t, c)
	require.Equal(t, len(c), 0)
	require.Equal(t, cap(c), 10)
}

func TestAdd(t *testing.T) {
	New()
	Add(func() error { return nil })
	require.Equal(t, len(c), 1)
}

func TestClose(t *testing.T) {
	New()
	Add(func() error { return nil })
	errorText := "TestClose"
	Add(func() error { return errors.New(errorText) })
	require.Len(t, c, 2)
	res := Close()
	require.Len(t, res, 1)
	require.Equal(t, res[0].Error(), errorText)
}
