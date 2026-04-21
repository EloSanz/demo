package circuitbreaker_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/elosanz/demo/pkg/circuitbreaker"
	"github.com/stretchr/testify/assert"
)

func TestCircuitBreaker_StateTransitions(t *testing.T) {
	cb := circuitbreaker.New(circuitbreaker.Config{
		MaxFailures:       2,
		HalfOpenSuccesses: 2,
		Timeout:           100 * time.Millisecond,
	})

	ctx := context.Background()
	expectedErr := errors.New("some specific error")
	failFunc := func(ctx context.Context) error { return expectedErr }
	successFunc := func(ctx context.Context) error { return nil }

	// 1. Initial state is Closed
	assert.Equal(t, circuitbreaker.StateClosed, cb.State())

	// 2. One failure (threshold is 2), still closed
	err := cb.Execute(ctx, failFunc)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, circuitbreaker.StateClosed, cb.State())

	// 3. Second failure, transitions to Open
	err = cb.Execute(ctx, failFunc)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, circuitbreaker.StateOpen, cb.State())

	// 4. Requests are blocked in Open state
	err = cb.Execute(ctx, successFunc)
	assert.ErrorIs(t, err, circuitbreaker.ErrCircuitOpen)

	// 5. Wait for timeout to enter HalfOpen
	time.Sleep(150 * time.Millisecond)

	// First request after timeout should be allowed (transitions to HalfOpen inside)
	err = cb.Execute(ctx, successFunc)
	assert.NoError(t, err)
	assert.Equal(t, circuitbreaker.StateHalfOpen, cb.State())

	// 6. Another success (threshold is 2), transitions to Closed
	err = cb.Execute(ctx, successFunc)
	assert.NoError(t, err)
	assert.Equal(t, circuitbreaker.StateClosed, cb.State())

	// 7. Test HalfOpen failing again
	_ = cb.Execute(ctx, failFunc)
	_ = cb.Execute(ctx, failFunc)
	assert.Equal(t, circuitbreaker.StateOpen, cb.State()) // Re-opened

	time.Sleep(150 * time.Millisecond)
	err = cb.Execute(ctx, failFunc) // Allowed, but fails -> opens again
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, circuitbreaker.StateOpen, cb.State())
}
