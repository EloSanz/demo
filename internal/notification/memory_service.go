package notification

import (
	"context"
	"errors"
	"log/slog"
	"time"
)

type memoryNotificationService struct {
	queue chan Notification
}

// NewMemoryNotificationService creates a notification service that uses Go channels.
// bufferSize allows the queue to hold messages in memory.
func NewMemoryNotificationService(bufferSize int) NotificationService {
	return &memoryNotificationService{
		queue: make(chan Notification, bufferSize),
	}
}

func (s *memoryNotificationService) Publish(ctx context.Context, n Notification) error {
	select {
	case s.queue <- n:
		slog.Info("notification published to internal memory queue", "type", n.Type)
		return nil
	default:
		slog.Warn("internal queue is full")
		return errors.New("internal queue is full")
	}
}

func (s *memoryNotificationService) StartWorker(ctx context.Context) {
	slog.Info("starting in-memory notification worker")

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Info("stopping in-memory worker")
				return
			case n := <-s.queue:
				// Procesamos la notificación
				slog.Info("WORKER: processing notification from memory", "type", n.Type, "content", n.Content)
				
				// Simulamos trabajo
				time.Sleep(2 * time.Second)
				
				slog.Info("WORKER: notification processed successfully")
			}
		}
	}()
}
