package messaging

import (
	"context"
)

// EventWithOperationID struct that contains an operation ID and an Event
type Message struct {
	OperationID string
	Error       error
	Data        interface{}
}

type MessagingSystem interface {
	Publish(ctx context.Context, data interface{}) error
	// TODO how to deal with "topic" concept?
	Subscribe(ctx context.Context) (<-chan Message, context.CancelFunc, error)
	Close(ctx context.Context) error
}
