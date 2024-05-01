package httpadapter

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/perocha/goadapters/messaging"
	"github.com/perocha/goutils/pkg/telemetry"
)

// HTTPAdapterImpl implements the MessagingSystem interface
type HTTPAdapter struct {
	httpClient   *http.Client
	httpEndPoint *HTTPEndPoint
}

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

// PublisherInitializer
func PublisherInitializer(ctx context.Context, endPointURL string, portNumber string) (*HTTPAdapter, error) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::PublisherInitializer", telemetry.String("endPointURL", endPointURL), telemetry.String("PortNumber", portNumber))

	// Create a new HTTP client
	httpClient := &http.Client{}

	// Create a new HTTP endpoint
	httpEndPoint := NewEndpoint(endPointURL, portNumber)

	return &HTTPAdapter{
		httpClient:   httpClient,
		httpEndPoint: httpEndPoint,
	}, nil
}

// Publish a message with an optional endpoint URL
func (a *HTTPAdapter) Publish(ctx context.Context, data messaging.Message) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::Publish", telemetry.String("Command", data.GetCommand()), telemetry.String("Status", data.GetStatus()), telemetry.String("Data", string(data.GetData())))

	// Convert the message to JSON
	jsonData, err := data.Serialize()
	if err != nil {
		xTelemetry.Error(ctx, "HTTPAdapter::Publish::Failed", telemetry.String("Error", err.Error()))
		return err
	}

	// Construct the full endpoint URL
	endpointURL := a.httpEndPoint.GetEndPointURL() + ":" + a.httpEndPoint.GetPortNumber()

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpointURL, bytes.NewBuffer(jsonData))
	if err != nil {
		xTelemetry.Error(ctx, "HTTPAdapter::Publish::Failed to create HTTP request", telemetry.String("Error", err.Error()))
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request
	resp, err := a.httpClient.Do(req)
	if err != nil {
		xTelemetry.Error(ctx, "HTTPAdapter::Publish::Failed to make HTTP request", telemetry.String("Error", err.Error()))
		return err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		// Read the response body to capture the error message or details
		respBody, _ := io.ReadAll(resp.Body)
		xTelemetry.Error(ctx, "HTTPAdapter::Publish::Server returned non-OK status code", telemetry.Int("StatusCode", resp.StatusCode), telemetry.String("Response", string(respBody)))
		return errors.New("server returned non-OK status code")
	}

	return nil
}

// Subscribe to messages
func (a *HTTPAdapter) Subscribe(ctx context.Context) (<-chan messaging.Message, context.CancelFunc, error) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)

	xTelemetry.Debug(ctx, "HTTPAdapter::Subscribe")
	return nil, nil, nil
}

// Close the adapter
func (a *HTTPAdapter) Close(ctx context.Context) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::Close http adapter")

	a.httpClient.CloseIdleConnections()
	return nil
}

// Set the endpoint URL and port number
func (a *HTTPAdapter) UpdateEndPoint(ctx context.Context, endPoint HTTPEndPoint) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::SetEndPoint", telemetry.String("endPointURL", a.httpEndPoint.GetEndPointURL()), telemetry.String("PortNumber", a.httpEndPoint.GetPortNumber()))

	a.httpEndPoint.SetEndPointURL(endPoint.GetEndPointURL())
	a.httpEndPoint.SetPortNumber(endPoint.GetPortNumber())
}

// Set the endpoint
func (a *HTTPAdapter) SetEndPoint(endPoint HTTPEndPoint) error {
	a.httpEndPoint = &endPoint
	return nil
}
