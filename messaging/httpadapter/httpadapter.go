package httpadapter

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/perocha/goadapters/messaging/message"
	"github.com/perocha/goutils/pkg/telemetry"
)

// HTTPAdapterImpl implements the MessagingSystem interface
type HTTPAdapterImpl struct {
	httpClient  *http.Client
	endPointURL string
	portNumber  string
}

// ConsumerInitializer
func ConsumerInitializer(ctx context.Context, endPointURL string, portNumber string) (*HTTPAdapterImpl, error) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)

	httpClient := &http.Client{}

	xTelemetry.Info(ctx, "HTTPAdapter::ConsumerInitializer", telemetry.String("PortNumber", portNumber))

	return &HTTPAdapterImpl{
		httpClient:  httpClient,
		endPointURL: endPointURL,
		portNumber:  portNumber,
	}, nil
}

// PublisherInitializer
func PublisherInitializer(ctx context.Context, endPointURL string, portNumber string) (*HTTPAdapterImpl, error) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Info(ctx, "HTTPAdapter::PublisherInitializer", telemetry.String("endPointURL", endPointURL), telemetry.String("PortNumber", portNumber))

	httpClient := &http.Client{}

	return &HTTPAdapterImpl{
		httpClient:  httpClient,
		endPointURL: endPointURL,
		portNumber:  portNumber,
	}, nil
}

// Publish a message with an optional endpoint URL
func (a *HTTPAdapterImpl) Publish(ctx context.Context, data message.Message) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Info(ctx, "HTTPAdapter::Publish", telemetry.String("Command", data.GetCommand()), telemetry.String("Status", data.GetStatus()), telemetry.String("Data", string(data.GetData())))

	// Convert the message to JSON
	jsonData, err := data.Serialize()
	if err != nil {
		xTelemetry.Error(ctx, "HTTPAdapter::Publish::Failed", telemetry.String("Error", err.Error()))
		return err
	}

	// Construct the full endpoint URL
	endpointURL := a.endPointURL + ":" + a.portNumber

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
func (a *HTTPAdapterImpl) Subscribe(ctx context.Context) (<-chan message.Message, context.CancelFunc, error) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)

	xTelemetry.Info(ctx, "HTTPAdapter::Subscribe")
	return nil, nil, nil
}

// Set the endpoint URL and port number
func (a *HTTPAdapterImpl) SetEndPoint(ctx context.Context, endPointURL string, portNumber string) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Info(ctx, "HTTPAdapter::SetEndPoint", telemetry.String("endPointURL", endPointURL), telemetry.String("PortNumber", portNumber))
	a.endPointURL = endPointURL
	a.portNumber = portNumber
}

// Get the endpoint URL
func (a *HTTPAdapterImpl) GetEndPoint() string {
	return a.endPointURL
}

// Get the port number
func (a *HTTPAdapterImpl) GetPortNumber() string {
	return a.portNumber
}
