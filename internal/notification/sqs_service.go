package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type sqsNotificationService struct {
	client   *sqs.Client
	queueURL string
}

// NewSQSNotificationService constructs the SQS implementation.
func NewSQSNotificationService(client *sqs.Client, queueURL string) NotificationService {
	return &sqsNotificationService{
		client:   client,
		queueURL: queueURL,
	}
}

func (s *sqsNotificationService) Publish(ctx context.Context, n Notification) error {
	body, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("marshalling notification: %w", err)
	}

	_, err = s.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.queueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return fmt.Errorf("sending SQS message: %w", err)
	}

	slog.Info("notification published to SQS", "type", n.Type)
	return nil
}

func (s *sqsNotificationService) StartWorker(ctx context.Context) {
	slog.Info("starting SQS worker", "queue", s.queueURL)

	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Info("stopping SQS worker")
				return
			default:
				// Long polling for messages (20 seconds max)
				output, err := s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
					QueueUrl:            aws.String(s.queueURL),
					MaxNumberOfMessages: 1,
					WaitTimeSeconds:     10, // Wait up to 10s for a message
				})

				if err != nil {
					slog.Error("error receiving SQS message (Queue dummy or unreachable)", "error", err)
					
					// IMPROVED: Context-aware sleep. This allows immediate shutdown.
					select {
					case <-ctx.Done():
						return
					case <-time.After(10 * time.Second):
						continue
					}
				}

				for _, msg := range output.Messages {
					s.processMessage(ctx, msg)
				}
			}
		}
	}()
}

func (s *sqsNotificationService) processMessage(ctx context.Context, msg types.Message) {
	var n Notification
	if err := json.Unmarshal([]byte(*msg.Body), &n); err != nil {
		slog.Error("error unmarshalling SQS message", "error", err)
		return
	}

	// ─── SIMULATION OF WORK ──────────────────────────────────────────────────
	slog.Info("RECEIVED ASYNC NOTIFICATION", "type", n.Type, "content", n.Content)
	time.Sleep(2 * time.Second) // Simulate processing time
	// ─────────────────────────────────────────────────────────────────────────

	// Delete message from queue after processing
	_, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.queueURL),
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		slog.Error("error deleting SQS message", "error", err)
	}
}
