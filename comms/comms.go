package comms

import (
	"context"
	"net/http"

	"github.com/perocha/goadapters/messaging"
)

type CommsSystem interface {
	SendRequest(ctx context.Context, data messaging.Message) error
	SetEndPoint(ctx context.Context, endPoint EndPoint) error
	GetEndPoint() EndPoint
	Start(ctx context.Context, endPoint EndPoint) error
	Stop(ctx context.Context) error
	RegisterEndPoint(ctx context.Context, endPoint EndPoint, handler http.HandlerFunc) error
}

// EndPoint interface
type EndPoint interface {
	GetEndPoint() string
	SetEndPoint(endPoint string)
}
