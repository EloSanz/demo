package notification

import "context"

// Notification represents a message to be processed asynchronously.
type Notification struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type NotificationService interface {
	Publish(ctx context.Context, n Notification) error
	StartWorker(ctx context.Context)
}
