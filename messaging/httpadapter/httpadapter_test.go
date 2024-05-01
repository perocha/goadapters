package httpadapter_test

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/perocha/goadapters/messaging"
	"github.com/perocha/goadapters/messaging/httpadapter"
	"github.com/perocha/goutils/pkg/telemetry"
	"github.com/stretchr/testify/assert"
)

type MockMessage struct {
	messaging.Message
}

func (m *MockMessage) Serialize() ([]byte, error) {
	return nil, errors.New("forced serialize error")
}

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
	msg := messaging.NewMessage("1234", nil, "success", "test", []byte("test"))
	err = adapter.Publish(ctx, msg)
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	ctx := initializeTelemetry()
	endPointURL := "http://localhost"
	portNumber := "8080"

	// Initialize the HTTP adapter
	adapter, err := httpadapter.ConsumerInitializer(ctx, endPointURL, portNumber)
	assert.NoError(t, err)

	// Get the endpoint URL and port number
	assert.Equal(t, endPointURL, adapter.GetEndPoint())
	assert.Equal(t, portNumber, adapter.GetPortNumber())
}

func TestSet(t *testing.T) {
	ctx := initializeTelemetry()
	endPointURL := "http://localhost"
	portNumber := "8080"

	// Initialize the HTTP adapter
	adapter, err := httpadapter.ConsumerInitializer(ctx, endPointURL, portNumber)
	assert.NoError(t, err)

	// Set the endpoint URL and port number
	newEndPointURL := "http://newhost"
	newPortNumber := "9090"
	adapter.SetEndPoint(ctx, newEndPointURL, newPortNumber)
	assert.Equal(t, newEndPointURL, adapter.GetEndPoint())
	assert.Equal(t, newPortNumber, adapter.GetPortNumber())
}

func TestSubscribe(t *testing.T) {
	ctx := initializeTelemetry()
	endPointURL := "http://localhost"
	portNumber := "8080"

	// Initialize the HTTP adapter
	adapter, err := httpadapter.ConsumerInitializer(ctx, endPointURL, portNumber)
	assert.NoError(t, err)

	// Subscribe to messages
	_, _, err = adapter.Subscribe(ctx)
	assert.NoError(t, err)
}

func TestPublish_ErrorMakingHttpRequest(t *testing.T) {
	ctx := initializeTelemetry()
	endPointURL := "http://localhost"
	portNumber := "8080"

	adapter, _ := httpadapter.PublisherInitializer(ctx, endPointURL, portNumber)

	// Create a message that will not cause an error when serialized
	msg := messaging.NewMessage("1234", nil, "success", "test", []byte("test"))

	// Force an error when making the HTTP request by providing an invalid URL
	adapter.SetEndPoint(ctx, "://invalid-url", "8080")

	err := adapter.Publish(ctx, msg)
	assert.Error(t, err)
}

func TestPublish_ErrorIncorrectPort(t *testing.T) {
	ctx := initializeTelemetry()
	endPointURL := "http://localhost"
	portNumber := "8081"

	adapter, _ := httpadapter.PublisherInitializer(ctx, endPointURL, portNumber)

	// Create a message that will not cause an error when serialized
	msg := messaging.NewMessage("1234", nil, "success", "test", []byte("test"))

	// Force an error when making the HTTP request by providing an invalid URL
	adapter.SetEndPoint(ctx, "http://localhost", "8080")

	err := adapter.Publish(ctx, msg)
	assert.Error(t, err)
}

func TestPublish_NonOKResponse(t *testing.T) {
	ctx := initializeTelemetry()
	// Create a mock HTTP server that responds with a 500 Internal Server Error status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	endPointURL := "http://localhost"
	portNumber := strings.Split(server.URL, ":")[2]
	adapter, _ := httpadapter.PublisherInitializer(ctx, endPointURL, portNumber)

	// Create a message that will not cause an error when serialized
	msg := messaging.NewMessage("1234", nil, "success", "test", []byte("test"))

	err := adapter.Publish(ctx, msg)
	assert.Error(t, err)
	assert.Equal(t, "server returned non-OK status code", err.Error())
}

func TestPublish_ErrorSerializing(t *testing.T) {
	ctx := initializeTelemetry()
	endPointURL := "http://localhost"
	portNumber := "8080"

	adapter, _ := httpadapter.PublisherInitializer(ctx, endPointURL, portNumber)

	// Create a mock message that will cause an error when serialized
	msg := &MockMessage{
		Message: messaging.NewMessage("1234", nil, "success", "test", []byte("test")),
	}

	err := adapter.Publish(ctx, msg)
	assert.Error(t, err)
	assert.Equal(t, "forced serialize error", err.Error())
}

func TestClose(t *testing.T) {
	ctx := initializeTelemetry()
	endPointURL := "http://localhost"
	portNumber := "8080"

	adapter, _ := httpadapter.PublisherInitializer(ctx, endPointURL, portNumber)

	err := adapter.Close(ctx)
	assert.NoError(t, err)
}
