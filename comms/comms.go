package comms

import (
	"context"

	"github.com/perocha/goadapters/messaging"
)

type CommsSender interface {
	SendRequest(ctx context.Context, endpoint EndPoint, data messaging.Message) error
}

type CommsReceiver interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	RegisterEndPoint(ctx context.Context, endpointPath string, handler HandlerFunc) error
}

// HandlerFunc defines the interface for the handler function
type HandlerFunc func(ResponseWriter, Request)

// ResponseWriter interface to abstract response writing
type ResponseWriter interface {
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

// Request interface to abstract incoming requests
type Request interface {
	Header(key string) string
	Body() []byte
}

// EndPoint interface
type EndPoint interface {
	GetEndPoint() string
	SetEndPoint(endPoint string)
}
