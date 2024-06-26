package httpadapter

import "net/http"

// HttpSender implements the sender part of comms interface
type HttpSender struct {
	httpClient *http.Client
}

// HttpReceiver implements the receiver part of comms interface
type HttpReceiver struct {
	httpServer *http.Server
	portNumber string
}

// HTTPStatus represents custom HTTP status codes
type HTTPStatus int

const (
	// Success status codes
	StatusOK HTTPStatus = 200

	// Redirection status codes
	StatusMovedPermanently HTTPStatus = 301
	StatusFound            HTTPStatus = 302

	// Client error status codes
	StatusBadRequest           HTTPStatus = 400
	StatusUnauthorized         HTTPStatus = 401
	StatusForbidden            HTTPStatus = 403
	StatusNotFound             HTTPStatus = 404
	StatusMethodNotAllowed     HTTPStatus = 405
	StatusRequestTimeout       HTTPStatus = 408
	StatusConflict             HTTPStatus = 409
	StatusPreconditionFailed   HTTPStatus = 412
	StatusUnsupportedMediaType HTTPStatus = 415
	StatusTooManyRequests      HTTPStatus = 429

	// Server error status codes
	StatusInternalServerError HTTPStatus = 500
	StatusNotImplemented      HTTPStatus = 501
	StatusBadGateway          HTTPStatus = 502
	StatusServiceUnavailable  HTTPStatus = 503
	StatusGatewayTimeout      HTTPStatus = 504
)
