package httpadapter

import (
	"context"
	"errors"
	"net/http"

	"github.com/perocha/goadapters/comms"
	"github.com/perocha/goutils/pkg/telemetry"
)

func HTTPServerAdapterInit(ctx context.Context, portNumber string, path string) (*HttpAdapter, error) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::HTTPServerAdapterInit", telemetry.String("PortNumber", portNumber), telemetry.String("Path", path))

	// Create a new server
	httpServer := &http.Server{Addr: ":" + portNumber}

	// Create a new HTTP endpoint
	httpEndPoint := NewEndpoint("localhost", portNumber, path)

	return &HttpAdapter{
		httpClient:   nil,
		httpEndPoint: httpEndPoint,
		httpServer:   httpServer,
	}, nil
}

// Start the HTTP server
func (a *HttpAdapter) Start(ctx context.Context, endPoint comms.EndPoint) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::Start", telemetry.String("EndPoint", endPoint.GetEndPoint()))

	// Get the HTTP endpoint
	httpEndPoint, ok := endPoint.(*HTTPEndPoint)
	if !ok {
		xTelemetry.Error(ctx, "HTTPAdapter::Start::Failed to cast to HTTPEndPoint")
		err := errors.New("failed to cast to HTTPEndPoint")
		return err
	}

	// Register the endpoint
	a.httpServer = &http.Server{Addr: ":" + httpEndPoint.GetPortNumber()}

	// Start the server
	go func() {
		err := a.httpServer.ListenAndServe()
		if err != nil {
			xTelemetry.Error(ctx, "HTTPAdapter::Start::Failed to start HTTP server", telemetry.String("Error", err.Error()))
		}
	}()

	return nil
}

// Register a new endpoint
func (a *HttpAdapter) RegisterEndPoint(ctx context.Context, endPoint comms.EndPoint, handler http.HandlerFunc) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::RegisterEndPoint", telemetry.String("EndPoint", endPoint.GetEndPoint()))

	// Get the HTTP endpoint
	httpEndPoint, ok := endPoint.(*HTTPEndPoint)
	if !ok {
		xTelemetry.Error(ctx, "HTTPAdapter::RegisterEndPoint::Failed to cast to HTTPEndPoint")
		err := errors.New("failed to cast to HTTPEndPoint")
		return err
	}

	// Register the endpoint
	http.HandleFunc(httpEndPoint.GetPath(), handler)

	return nil
}

// Stop the HTTP server
func (a *HttpAdapter) Stop(ctx context.Context) error {
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
