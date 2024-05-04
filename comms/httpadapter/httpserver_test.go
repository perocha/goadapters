package httpadapter_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/perocha/goadapters/comms"
	"github.com/perocha/goadapters/comms/httpadapter"
	"github.com/stretchr/testify/assert"
)

func TestHTTPServerAdapterInit(t *testing.T) {
	// Mock context and dependencies
	ctx := initializeTelemetry()
	portNumber := "8080"

	// Initialize the HTTP adapter
	adapter, err := httpadapter.HTTPServerAdapterInit(ctx, portNumber)
	assert.NoError(t, err)
	assert.NotNil(t, adapter)
}

func TestHttpAdapter_StartAndStop(t *testing.T) {
	// Mock context and dependencies
	ctx := initializeTelemetry()
	portNumber := "8080"

	// Initialize the HTTP adapter
	adapter, err := httpadapter.HTTPServerAdapterInit(ctx, portNumber)
	assert.NoError(t, err)
	assert.NotNil(t, adapter)

	// Start the HTTP adapter
	err = adapter.Start(ctx)
	assert.NoError(t, err)

	// Stop the HTTP adapter
	err = adapter.Stop(ctx)
	assert.NoError(t, err)
}

func TestHttpAdapter_RegisterEndPoint(t *testing.T) {
	// Mock context and dependencies
	ctx := initializeTelemetry()
	host := "localhost"
	portNumber := "8080"
	path := "/test"
	endpoint := httpadapter.NewEndpoint(host, portNumber, path)

	// Initialize the HTTP adapter
	adapter, err := httpadapter.HTTPServerAdapterInit(ctx, portNumber)
	assert.NoError(t, err)
	assert.NotNil(t, adapter)

	// Register a test handler function
	testHandler := func(w comms.ResponseWriter, r comms.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// Register an endpoint
	err = adapter.RegisterEndPoint(ctx, endpoint, testHandler)
	assert.NoError(t, err)

	// Start the HTTP adapter
	err = adapter.Start(ctx)
	assert.NoError(t, err)

	// Wait for the server to start
	time.Sleep(5 * time.Second)

	// Create a test request
	req, err := http.NewRequest("POST", "http://localhost:"+portNumber+path, nil)
	assert.NoError(t, err)

	// Send the request
	resp, err := http.DefaultClient.Do(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Stop the HTTP adapter
	err = adapter.Stop(ctx)
	assert.NoError(t, err)
}
