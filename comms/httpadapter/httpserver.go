package httpadapter

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/perocha/goadapters/comms"
	"github.com/perocha/goutils/pkg/telemetry"
)

func HTTPServerAdapterInit(ctx context.Context, port string) (*HttpReceiver, error) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::HTTPServerAdapterInit", telemetry.String("Port", port))

	// Create a new server
	httpServer := &http.Server{Addr: ":" + port}

	return &HttpReceiver{
		httpServer: httpServer,
		portNumber: port,
	}, nil
}

// Start the HTTP server
func (a *HttpReceiver) Start(ctx context.Context) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::Start")

	// Validate server is not nil and port number is not empty
	if a.httpServer == nil {
		xTelemetry.Error(ctx, "HTTPAdapter::Start::HTTP server is nil")
		err := errors.New("HTTP server is nil")
		return err
	}
	if a.portNumber == "" {
		xTelemetry.Error(ctx, "HTTPAdapter::Start::Port number is empty")
		err := errors.New("port number is empty")
		return err
	}

	// Start the server
	go func() {
		err := a.httpServer.ListenAndServe()
		if err != nil {
			xTelemetry.Error(ctx, "HTTPAdapter::Start::Failed to start HTTP server", telemetry.String("Error", err.Error()))
		}
	}()

	return nil
}

// Stop the HTTP server
func (a *HttpReceiver) Stop(ctx context.Context) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::Stop")

	// Shutdown the server
	err := a.httpServer.Shutdown(ctx)
	if err != nil {
		xTelemetry.Error(ctx, "HTTPAdapter::Stop::Failed to shutdown HTTP server", telemetry.String("Error", err.Error()))
		return err
	}

	return nil
}

// Register a new endpoint
func (a *HttpReceiver) RegisterEndPoint(ctx context.Context, endpointPath string, handler comms.HandlerFunc) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::RegisterEndPoint", telemetry.String("endpointPath", endpointPath))

	// Register the endpoint with the adapter function
	http.HandleFunc(endpointPath, func(w http.ResponseWriter, r *http.Request) {
		// Convert http.ResponseWriter to comms.ResponseWriter
		commsWriter := &responseWriterAdapter{
			w,
			http.StatusOK,
		}

		// Convert *http.Request to comms.Request
		commsReq := &requestAdapter{
			r,
		}

		// Call the handler function, wrapping the handler with telemetry logging
		wrappedHandler := func(ctx context.Context, w comms.ResponseWriter, r comms.Request) {
			// Generate a new OperationID and append to the context
			operationID := uuid.New().String()
			ctx = telemetry.SetOperationID(ctx, operationID)

			// Get service name from context
			serviceName := telemetry.GetServiceName(ctx)

			startTime := time.Now()
			// Call the original handler
			newCtx, err := handler(ctx, w, r)

			// Decide on the message based on the error
			message := ""
			if err == nil {
				message = "Request processed successfully"
			} else {
				message = err.Error()
			}

			// Log telemetry after calling the original handler
			statusCode := w.Status()
			success := isSuccess(statusCode)
			xTelemetry.Request(newCtx, http.MethodPost, endpointPath, startTime, time.Now(), strconv.Itoa(statusCode), success, serviceName, message)
		}

		// Call the wrapped handler
		wrappedHandler(ctx, commsWriter, commsReq)

	})

	return nil
}

// Check if the status code is a success status code
func isSuccess(statusCode int) bool {
	switch statusCode {
	case http.StatusOK:
		return true
	case http.StatusCreated:
		return true
	case http.StatusAccepted:
		return true
	case http.StatusNoContent:
		return true
	default:
		return false
	}
}
