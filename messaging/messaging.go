package messaging

import "context"

// Interface for messaging systems
type MessagingSystem interface {
	Publish(ctx context.Context, data Message) error
	Subscribe(ctx context.Context) (<-chan Message, context.CancelFunc, error)
	Close(ctx context.Context) error
}
