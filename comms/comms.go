package comms

import (
	"context"

	"github.com/perocha/goadapters/messaging"
)

type CommsSystem interface {
	SendRequest(ctx context.Context, data messaging.Message) error
	SetEndPoint(ctx context.Context, endPoint EndPoint) error
	GetEndPoint() EndPoint
}

// EndPoint interface
type EndPoint interface {
	GetEndPoint() string
	SetEndPoint(endPoint string)
}
