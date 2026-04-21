package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

// State represents the state of the circuit breaker.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

// ErrCircuitOpen is returned when the circuit breaker is open and not allowing requests.
var ErrCircuitOpen = errors.New("circuit breaker is open")

// CircuitBreaker is a simple state machine that monitors failures and prevents
// execution when a failure threshold is reached.
type CircuitBreaker struct {
	mu           sync.RWMutex
	state        State
	failureCount int
	successCount int

	maxFailures       int
	halfOpenSuccesses int
	timeout           time.Duration
	openedAt          time.Time
}

// Config configures the CircuitBreaker.
type Config struct {
	// MaxFailures is the number of consecutive failures before the circuit opens.
	MaxFailures int
	// HalfOpenSuccesses is the number of consecutive successes required to close a half-open circuit.
	HalfOpenSuccesses int
	// Timeout is the duration the circuit remains open before transitioning to half-open.
	Timeout time.Duration
}

// New creates a new CircuitBreaker with the given configuration.
func New(cfg Config) *CircuitBreaker {
	if cfg.MaxFailures <= 0 {
		cfg.MaxFailures = 5
	}
	if cfg.HalfOpenSuccesses <= 0 {
		cfg.HalfOpenSuccesses = 1
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = 10 * time.Second
	}

	return &CircuitBreaker{
		state:             StateClosed,
		maxFailures:       cfg.MaxFailures,
		halfOpenSuccesses: cfg.HalfOpenSuccesses,
		timeout:           cfg.Timeout,
	}
}

// Execute runs the given function if the circuit is closed or half-open.
// If the circuit is open, it returns ErrCircuitOpen immediately.
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	if !cb.allowRequest() {
		return ErrCircuitOpen
	}

	err := fn(ctx)
	cb.recordResult(err == nil)
	return err
}

// State returns the current state of the circuit breaker.
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

func (cb *CircuitBreaker) allowRequest() bool {
	cb.mu.RLock()
	state := cb.state
	openedAt := cb.openedAt
	timeout := cb.timeout
	cb.mu.RUnlock()

	switch state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(openedAt) > timeout {
			cb.mu.Lock()
			defer cb.mu.Unlock()
			// Another goroutine might have already changed the state
			if cb.state == StateOpen {
				cb.state = StateHalfOpen
				cb.successCount = 0
			}
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

func (cb *CircuitBreaker) recordResult(success bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case StateClosed:
		if success {
			cb.failureCount = 0
		} else {
			cb.failureCount++
			if cb.failureCount >= cb.maxFailures {
				cb.transitionToOpen()
			}
		}
	case StateHalfOpen:
		if success {
			cb.successCount++
			if cb.successCount >= cb.halfOpenSuccesses {
				cb.transitionToClosed()
			}
		} else {
			// Any failure in half-open state immediately re-opens the circuit
			cb.transitionToOpen()
		}
	}
}

func (cb *CircuitBreaker) transitionToOpen() {
	cb.state = StateOpen
	cb.openedAt = time.Now()
}

func (cb *CircuitBreaker) transitionToClosed() {
	cb.state = StateClosed
	cb.failureCount = 0
	cb.successCount = 0
}
