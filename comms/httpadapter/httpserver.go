package httpadapter

import (
	"context"
	"errors"
	"net/http"

	"github.com/perocha/goadapters/comms"
	"github.com/perocha/goutils/pkg/telemetry"
)

func HTTPServerAdapterInit(ctx context.Context, endPoint comms.EndPoint) (*HttpAdapter, error) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)

	// Get the HTTP endpoint
	httpEndPoint, ok := endPoint.(*HTTPEndPoint)
	if !ok {
		xTelemetry.Error(ctx, "HTTPAdapter::Start::Failed to cast to HTTPEndPoint")
		err := errors.New("failed to cast to HTTPEndPoint")
		return nil, err
	}

	xTelemetry.Debug(ctx, "HTTPAdapter::HTTPServerAdapterInit", telemetry.String("PortNumber", httpEndPoint.GetPortNumber()), telemetry.String("Path", httpEndPoint.GetPath()))

	// Create a new server
	httpServer := &http.Server{Addr: ":" + httpEndPoint.GetPortNumber()}

	return &HttpAdapter{
		httpClient:   nil,
		httpEndPoint: httpEndPoint,
		httpServer:   httpServer,
	}, nil
}

// Start the HTTP server
func (a *HttpAdapter) Start(ctx context.Context) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::Start")

	// Validate if there's an endpoint set
	if a.httpEndPoint == nil {
		xTelemetry.Error(ctx, "HTTPAdapter::Start::No endpoint set")
		err := errors.New("no endpoint set")
		return err
	}

	// Register the endpoint
	a.httpServer = &http.Server{Addr: ":" + a.httpEndPoint.GetPortNumber()}

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
func (a *HttpAdapter) RegisterEndPoint(ctx context.Context, endPoint comms.EndPoint, handler comms.HandlerFunc) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "HTTPAdapter::RegisterEndPoint", telemetry.String("EndPoint", endPoint.GetEndPoint()))

	// Get the HTTP endpoint
	httpEndPoint, ok := endPoint.(*HTTPEndPoint)
	if !ok {
		xTelemetry.Error(ctx, "HTTPAdapter::RegisterEndPoint::Failed to cast to HTTPEndPoint")
		err := errors.New("failed to cast to HTTPEndPoint")
		return err
	}

	// Register the endpoint with the adapter function
	http.HandleFunc(httpEndPoint.GetPath(), func(w http.ResponseWriter, r *http.Request) {
		// Convert http.ResponseWriter to comms.ResponseWriter
		commsWriter := &responseWriterAdapter{w}

		// Convert *http.Request to comms.Request
		commsReq := &requestAdapter{r}

		// Call the handler function
		handler(commsWriter, commsReq)
	})

	return nil
}

// Adapter functions to convert http.ResponseWriter and *http.Request to comms.ResponseWriter and comms.Request respectively

type responseWriterAdapter struct {
	http.ResponseWriter
}

func (r *responseWriterAdapter) Write(data []byte) (int, error) {
	return r.ResponseWriter.Write(data)
}

func (r *responseWriterAdapter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
}

type requestAdapter struct {
	*http.Request
}

func (r *requestAdapter) Header(key string) string {
	return r.Request.Header.Get(key)
}

func (r *requestAdapter) Body() []byte {
	// Implement your logic to read the request body if needed
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
