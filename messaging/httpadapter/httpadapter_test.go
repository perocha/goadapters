package httpadapter_test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/perocha/goadapters/messaging/httpadapter"
	"github.com/perocha/goadapters/messaging/message"
	"github.com/perocha/goutils/pkg/telemetry"
	"github.com/stretchr/testify/assert"
)

func initializeTelemetry() context.Context {
	// Initialize telemetry package
	serviceName := "httpadapter"
	telemetryConfig := telemetry.NewXTelemetryConfig("", serviceName, "info", 1)
	xTelemetry, err := telemetry.NewXTelemetry(telemetryConfig)
	if err != nil {
		log.Fatalf("Main::Fatal error::Failed to initialize XTelemetry %s\n", err.Error())
	}
	// Add telemetry object to the context, so that it can be reused across the application
	ctx := context.WithValue(context.Background(), telemetry.TelemetryContextKey, xTelemetry)
	return ctx
}

func TestConsumerInitializer(t *testing.T) {
	ctx := initializeTelemetry()
	endPointURL := "http://localhost"
	portNumber := "8080"

	adapter, err := httpadapter.ConsumerInitializer(ctx, endPointURL, portNumber)

	assert.NoError(t, err)
	assert.NotNil(t, adapter)
	assert.IsType(t, &httpadapter.HTTPAdapterImpl{}, adapter)
	assert.Equal(t, endPointURL, adapter.GetEndPoint())
	assert.Equal(t, portNumber, adapter.GetPortNumber())
}

func TestPublisherInitializer(t *testing.T) {
	ctx := initializeTelemetry()
	endPointURL := "http://localhost"
	portNumber := "8080"

	adapter, err := httpadapter.PublisherInitializer(ctx, endPointURL, portNumber)

	assert.NoError(t, err)
	assert.NotNil(t, adapter)
	assert.IsType(t, &httpadapter.HTTPAdapterImpl{}, adapter)
	assert.Equal(t, endPointURL, adapter.GetEndPoint())
	assert.Equal(t, portNumber, adapter.GetPortNumber())
}

func TestPublish(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx := initializeTelemetry()
	endPointURL := "http://localhost"
	portNumber := strings.Split(server.URL, ":")[2]

	// Set the URL of the mock server
	server.URL = endPointURL + ":" + portNumber

	// Initialize the HTTP adapter
	adapter, err := httpadapter.PublisherInitializer(ctx, endPointURL, portNumber)
	assert.NoError(t, err)

	// Build mock Message
	msg := message.NewMessage("1234", nil, "success", "test", []byte("test"))
	err = adapter.Publish(ctx, msg)
	assert.NoError(t, err)
}
