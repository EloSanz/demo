package notification

import "context"

// Notification represents a message to be processed asynchronously.
type Notification struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// NotificationService defines the contract for sending and consuming notifications.
type NotificationService interface {
	// Publish sends a notification to the queue.
	Publish(ctx context.Context, n Notification) error
	
	// StartWorker starts a background process to consume notifications.
	StartWorker(ctx context.Context)
}
