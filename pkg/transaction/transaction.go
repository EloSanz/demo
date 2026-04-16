package transaction

import (
	"fmt"
	"log/slog"
)

// Action define una unidad de trabajo con su respectiva compensación (rollback).
type Action struct {
	Name     string
	Execute  func() error
	Rollback func() error
}

// Helper maneja la lógica de ejecución secuencial y rollbacks en cascada.
type Helper struct {
	Retries int
}

func NewHelper(retries int) *Helper {
	return &Helper{Retries: retries}
}

// Execute corre las acciones en orden. Si una falla, ejecuta los rollbacks en orden inverso.
func (h *Helper) Execute(actions []Action) error {
	executedCount := 0

	for i, act := range actions {
		slog.Debug("executing transaction step", "step", act.Name, "index", i)
		if err := act.Execute(); err != nil {
			slog.Error("transaction step failed", 
				"step", act.Name, 
				"index", i, 
				"error", err,
			)
			h.rollbackAll(actions, executedCount)
			return fmt.Errorf("transaction failed at step %q (index %d): %w", act.Name, i, err)
		}
		executedCount++
	}

	return nil
}

func (h *Helper) rollbackAll(actions []Action, executedCount int) {
	slog.Info("starting rollback sequence", "count", executedCount)
	
	for i := executedCount - 1; i >= 0; i-- {
		act := actions[i]
		if act.Rollback == nil {
			continue
		}

		slog.Warn("rolling back step", "step", act.Name, "index", i)
		err := h.executeWithRetries(act.Name, act.Rollback)
		if err != nil {
			slog.Error("CRITICAL: rollback failed after retries", 
				"step", act.Name, 
				"index", i, 
				"error", err,
			)
		}
	}
}

func (h *Helper) executeWithRetries(name string, action func() error) error {
	var lastErr error
	for attempt := 0; attempt <= h.Retries; attempt++ {
		lastErr = action()
		if lastErr == nil {
			return nil
		}
		slog.Warn("rollback attempt failed, retrying...", 
			"step", name, 
			"attempt", attempt+1, 
			"error", lastErr,
		)
	}
	return lastErr
}
