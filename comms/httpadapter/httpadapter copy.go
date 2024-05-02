package httpadapter

import (
	"net/http"
)

// HTTPAdapterImpl implements the MessagingSystem interface
type HTTPAdapter struct {
	httpClient   *http.Client
	httpEndPoint *HTTPEndPoint
}

/*
// ConsumerInitializer
func ConsumerInitializer(ctx context.Context, endPointURL string, portNumber string) (*HTTPAdapter, error) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)

	// Create a new HTTP client
	httpClient := &http.Client{}

	// Create a new HTTP endpoint
	httpEndPoint := NewEndpoint(endPointURL, portNumber)

	xTelemetry.Debug(ctx, "HTTPAdapter::ConsumerInitializer", telemetry.String("PortNumber", portNumber))

	return &HTTPAdapter{
		httpClient:   httpClient,
		httpEndPoint: httpEndPoint,
	}, nil
}

// Subscribe to messages
func (a *HTTPAdapter) Subscribe(ctx context.Context) (<-chan messaging.Message, context.CancelFunc, error) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::Subscribe")

	// Create a new channel
	messageChannel := make(chan messaging.Message)

	// Create a new context with a cancel function
	ctx, cancel := context.WithCancel(ctx)

	// Construct the full endpoint URL
	endpointURL := a.httpEndPoint.GetEndPointURL() + ":" + a.httpEndPoint.GetPortNumber()

	// Listen for incoming HTTP requests
	http.HandleFunc(endpointURL, func(w http.ResponseWriter, r *http.Request) {
		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			xTelemetry.Error(ctx, "HTTPAdapter::Subscribe::Failed to read request body", telemetry.String("Error", err.Error()))
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// Create a new message
		message := messaging.NewMessage("", nil, "", "", nil)
		err = message.Deserialize(body)
		if err != nil {
			xTelemetry.Error(ctx, "HTTPAdapter::Subscribe::Failed to deserialize message", telemetry.String("Error", err.Error()))
			http.Error(w, "Failed to deserialize message", http.StatusInternalServerError)
			return
		}

		// Send the message to the channel
		messageChannel <- message

		// Send an OK response
		w.WriteHeader(http.StatusOK)
	})

	// Start the HTTP server
	go func() {
		err := http.ListenAndServe(a.httpEndPoint.GetEndPoint(), nil)
		if err != nil {
			xTelemetry.Error(ctx, "HTTPAdapter::Subscribe::Failed to start HTTP server", telemetry.String("Error", err.Error()))
		}
	}()

	return messageChannel, cancel, nil
}

// Close the adapter
func (a *HTTPAdapter) Close(ctx context.Context) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::Close http adapter")

	a.httpClient.CloseIdleConnections()
	return nil
}

/*
// Generic endpoint update method
func (a *HTTPAdapter) UpdateEndPoint(ctx context.Context, endPoint messaging.EndPoint) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::SetEndPoint", telemetry.String("endPointURL", a.httpEndPoint.GetEndPointURL()), telemetry.String("PortNumber", a.httpEndPoint.GetPortNumber()))

	httpEndPoint, ok := endPoint.(*HTTPEndPoint)
	if !ok {
		return errors.New("endpoint type is not supported by HTTP adapter")
	}

	return a.UpdateHTTPAdapterEndPoint(ctx, httpEndPoint)
}

// Update the HTTP adapter endpoint
func (a *HTTPAdapter) UpdateHTTPAdapterEndPoint(ctx context.Context, endPoint *HTTPEndPoint) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::UpdateHTTPAdapterEndPoint", telemetry.String("endPointURL", endPoint.GetEndPointURL()), telemetry.String("PortNumber", endPoint.GetPortNumber()))

	a.httpEndPoint.SetEndPointURL(endPoint.GetEndPointURL())
	a.httpEndPoint.SetPortNumber(endPoint.GetPortNumber())
	return nil
}

// Generic method to get the endpoint
func (a *HTTPAdapter) GetEndPoint() string {
	return a.GetHTTPAdapterEndPoint()
}

// Get the HTTP adapter endpoint
func (a *HTTPAdapter) GetHTTPAdapterEndPoint() string {
	return a.httpEndPoint.GetEndPoint()
}

*/
