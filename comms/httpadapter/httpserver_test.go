package httpadapter_test

import (
	"net/http"
	"testing"

	"github.com/perocha/goadapters/comms/httpadapter"
	"github.com/stretchr/testify/assert"
)

func TestHTTPServerAdapterInit(t *testing.T) {
	// Mock context and dependencies
	ctx := initializeTelemetry()
	portNumber := "8080"
	path := "/test"

	// Initialize the HTTP adapter
	adapter, err := httpadapter.HTTPServerAdapterInit(ctx, portNumber, path)
	assert.NoError(t, err)
	assert.NotNil(t, adapter)

	// Get the endpoint URL and port number
	httpEndPoint, ok := adapter.GetEndPoint().(*httpadapter.HTTPEndPoint)
	if !ok {
		// Test error, should not reach here
		assert.Fail(t, "Failed to cast to HTTPEndPoint")
	}

	// Check the HTTP endpoint
	assert.NotNil(t, httpEndPoint)
	assert.Equal(t, portNumber, httpEndPoint.GetPortNumber())
	assert.Equal(t, path, httpEndPoint.GetPath())
}

func TestHttpAdapter_StartAndStop(t *testing.T) {
	// Mock context and dependencies
	ctx := initializeTelemetry()
	portNumber := "8080"
	path := "/test"

	// Initialize the HTTP adapter
	adapter, err := httpadapter.HTTPServerAdapterInit(ctx, portNumber, path)
	assert.NoError(t, err)
	assert.NotNil(t, adapter)

	// Start the HTTP adapter
	err = adapter.Start(ctx, adapter.GetEndPoint())
	assert.NoError(t, err)

	// Create a test request to ensure the server is running
	req, err := http.NewRequest("GET", "http://localhost:"+portNumber+path, nil)
	assert.NoError(t, err)

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Stop the HTTP adapter
	err = adapter.Stop(ctx)
	assert.NoError(t, err)
}

func TestHttpAdapter_RegisterEndPoint(t *testing.T) {
	// Mock context and dependencies
	ctx := initializeTelemetry()
	portNumber := "8080"
	path := "/test"

	// Initialize the HTTP adapter
	adapter, err := httpadapter.HTTPServerAdapterInit(ctx, portNumber, path)
	assert.NoError(t, err)
	assert.NotNil(t, adapter)

	// Register a test handler function
	testHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Register an endpoint
	err = adapter.RegisterEndPoint(ctx, adapter.GetEndPoint(), testHandler)
	assert.NoError(t, err)

	// Create a test request
	req, err := http.NewRequest("GET", "http://localhost:"+portNumber+path, nil)
	assert.NoError(t, err)

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
