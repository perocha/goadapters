package httpadapter

// HTTPEndPoint implements the EndPoint interface
type HTTPEndPoint struct {
	endPointURL string
	portNumber  string
}

// Create a new HTTPEndPoint
func NewEndpoint(endPointURL string, portNumber string) *HTTPEndPoint {
	return &HTTPEndPoint{
		endPointURL: endPointURL,
		portNumber:  portNumber,
	}
}

// Generic method to get the endpoint
func (e *HTTPEndPoint) GetEndPoint() string {
	return e.endPointURL + ":" + e.portNumber
}

// Generic method to set the endpoint
func (e *HTTPEndPoint) SetEndPoint(endPoint string) {
	e.endPointURL = endPoint
}

// Get the endpoint URL
func (e *HTTPEndPoint) GetEndPointURL() string {
	return e.endPointURL
}

// Get the port number
func (e *HTTPEndPoint) GetPortNumber() string {
	return e.portNumber
}

// Set the endpoint URL
func (e *HTTPEndPoint) SetEndPointURL(endPointURL string) {
	e.endPointURL = endPointURL
}

// Set the port number
func (e *HTTPEndPoint) SetPortNumber(portNumber string) {
	e.portNumber = portNumber
}
