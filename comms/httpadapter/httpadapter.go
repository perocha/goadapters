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
type HttpSendAdapter struct {
	httpClient   *http.Client
	httpEndPoint *HTTPEndPoint
}

// Initialize the HTTP adapter
func Initializer(ctx context.Context, host string, portNumber string, path string) (*HttpSendAdapter, error) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::PublisherInitializer", telemetry.String("host", host), telemetry.String("PortNumber", portNumber), telemetry.String("Path", path))

	// Create a new HTTP client
	httpClient := &http.Client{}

	// Create a new HTTP endpoint
	httpEndPoint := NewEndpoint(host, portNumber, path)

	return &HttpSendAdapter{
		httpClient:   httpClient,
		httpEndPoint: httpEndPoint,
	}, nil
}

// Send a request
func (a *HttpSendAdapter) SendRequest(ctx context.Context, data messaging.Message) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::Publish", telemetry.String("Command", data.GetCommand()), telemetry.String("Status", data.GetStatus()), telemetry.String("Data", string(data.GetData())))

	// Convert the message to JSON
	jsonData, err := data.Serialize()
	if err != nil {
		xTelemetry.Error(ctx, "HTTPAdapter::Publish::Failed", telemetry.String("Error", err.Error()))
		return err
	}

	// Construct the full endpoint URL
	endpointURL := a.httpEndPoint.GetEndPoint()

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

// Generic endpoint update method
func (a *HttpSendAdapter) SetEndPoint(ctx context.Context, endPoint HTTPEndPoint) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::SetEndPoint", telemetry.String("endPoint", endPoint.GetEndPoint()))

	a.httpEndPoint = &endPoint

	return nil
}

// Get endpoint object
func (a *HttpSendAdapter) GetEndPoint() *HTTPEndPoint {
	return a.httpEndPoint
}
