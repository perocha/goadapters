package httpadapter

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/perocha/goadapters/comms"
	"github.com/perocha/goadapters/messaging"
	"github.com/perocha/goutils/pkg/telemetry"
)

// Initialize the HTTP adapter
func HttpSenderInit(ctx context.Context) (*HttpSender, error) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::HttpSenderInit")

	// Create a new HTTP client
	httpClient := &http.Client{}

	return &HttpSender{
		httpClient: httpClient,
	}, nil
}

// Send a request
func (a *HttpSender) SendRequest(ctx context.Context, endpoint comms.EndPoint, data messaging.Message) error {
	startTime := time.Now()
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::Publish", telemetry.String("Command", data.GetCommand()), telemetry.String("Status", data.GetStatus()), telemetry.String("Data", string(data.GetData())))

	// Create a new operation id (new uuid) and add it to the context
	operationID := uuid.New().String()
	ctx = context.WithValue(context.Background(), telemetry.OperationIDKeyContextKey, operationID)
	data.SetOperationID(operationID)

	// Convert the message to JSON
	jsonData, err := data.Serialize()
	if err != nil {
		xTelemetry.Error(ctx, "HTTPAdapter::Publish::Failed", telemetry.String("Error", err.Error()))
		return err
	}

	// Get the endpoint URL
	httpEndPoint, ok := endpoint.(*HTTPEndPoint)
	if !ok {
		return errors.New("endpoint is not of type HTTPEndPoint")
	}

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, httpEndPoint.GetEndPoint(), bytes.NewBuffer(jsonData))
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

	// Log the telemetry request
	duration := time.Since(startTime)
	xTelemetry.Request(ctx, http.MethodPost, httpEndPoint.GetEndPoint(), duration, strconv.Itoa(http.StatusOK), true, httpEndPoint.GetHost(), "HTTPAdapter::Publish::Success")

	return nil
}
