package httpadapter

import (
	"net/url"
)

// EndPoint interface
type EndPoint interface {
	GetEndPoint() string
	SetEndPoint(endPoint string)
}

// HTTPEndPoint implements the EndPoint interface
type HTTPEndPoint struct {
	scheme     string
	host       string
	portNumber string
	path       string
}

// Create a new HTTPEndPoint
func NewEndpoint(host string, portNumber string, path string) *HTTPEndPoint {
	scheme := "http"

	return &HTTPEndPoint{
		scheme:     scheme,
		host:       host,
		portNumber: portNumber,
		path:       path,
	}
}

// Generic method to get the endpoint
func (e *HTTPEndPoint) GetEndPoint() string {
	// Construct the url
	url := e.scheme + "://" + e.host + ":" + e.portNumber + e.path

	return url
}

// Generic method to set the endpoint
func (e *HTTPEndPoint) SetHost(host string) {
	e.host = host
}

// Get the endpoint URL
func (e *HTTPEndPoint) GetHost() string {
	return e.host
}

// Get the port number
func (e *HTTPEndPoint) GetPortNumber() string {
	return e.portNumber
}

// Get the path
func (e *HTTPEndPoint) GetPath() string {
	return e.path
}

// Set the port number
func (e *HTTPEndPoint) SetPortNumber(portNumber string) {
	e.portNumber = portNumber
}

// Set the path
func (e *HTTPEndPoint) SetPath(path string) {
	e.path = path
}

// Set the endpoint URL
func (e *HTTPEndPoint) SetEndPoint(endPoint string) {
	// Parse the URL
	url, err := url.Parse(endPoint)
	if err != nil {
		return
	}

	// Set the host and port number
	e.host = url.Hostname()
	e.portNumber = url.Port()
	e.path = url.Path
}
