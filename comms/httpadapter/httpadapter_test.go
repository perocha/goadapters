package httpadapter_test

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/perocha/goadapters/comms/httpadapter"
	"github.com/perocha/goadapters/messaging"
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

func TestInitializer(t *testing.T) {
	ctx := initializeTelemetry()
	host := "http://localhost"
	portNumber := "8080"
	path := "/test"

	adapter, err := httpadapter.Initializer(ctx, host, portNumber, path)

	assert.NoError(t, err)
	assert.NotNil(t, adapter)
	assert.IsType(t, &httpadapter.HttpSendAdapter{}, adapter)
	assert.Equal(t, host, adapter.GetEndPoint().GetHost())
	assert.Equal(t, portNumber, adapter.GetEndPoint().GetPortNumber())
	assert.Equal(t, path, adapter.GetEndPoint().GetPath())
}

func TestPublish(t *testing.T) {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx := initializeTelemetry()
	host := "localhost"
	portNumber := strings.Split(server.URL, ":")[2]
	path := "/test"

	// Set the URL of the mock server
	server.URL = host + ":" + portNumber + path

	// Initialize the HTTP adapter
	adapter, err := httpadapter.Initializer(ctx, host, portNumber, path)
	assert.NoError(t, err)

	// Build mock Message
	msg := messaging.NewMessage("1234", nil, "success", "test", []byte("test"))
	err = adapter.SendRequest(ctx, msg)
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	ctx := initializeTelemetry()
	host := "http://localhost"
	portNumber := "8080"
	path := "/test"

	// Initialize the HTTP adapter
	adapter, err := httpadapter.Initializer(ctx, host, portNumber, path)
	assert.NoError(t, err)

	// Get host, port number, and path
	assert.Equal(t, host, adapter.GetEndPoint().GetHost())
	assert.Equal(t, portNumber, adapter.GetEndPoint().GetPortNumber())
	assert.Equal(t, path, adapter.GetEndPoint().GetPath())
}

func TestSet(t *testing.T) {
	ctx := initializeTelemetry()
	host := "http://localhost"
	portNumber := "8080"
	path := "/test"

	// Initialize the HTTP adapter
	adapter, err := httpadapter.Initializer(ctx, host, portNumber, path)
	assert.NoError(t, err)

	// Set the endpoint URL and port number
	newHost := "http://newhost"
	newPortNumber := "9090"
	newPath := "/newtest"
	newEndpoint := httpadapter.NewEndpoint(newHost, newPortNumber, newPath)
	adapter.SetEndPoint(ctx, *newEndpoint)
	assert.Equal(t, newHost, adapter.GetEndPoint().GetHost())
	assert.Equal(t, newPortNumber, adapter.GetEndPoint().GetPortNumber())
	assert.Equal(t, newPath, adapter.GetEndPoint().GetPath())
}

func TestPublish_ErrorMakingHttpRequest(t *testing.T) {
	ctx := initializeTelemetry()
	host := "http://localhost"
	portNumber := "8080"
	path := "/test"

	adapter, _ := httpadapter.Initializer(ctx, host, portNumber, path)

	// Create a message that will not cause an error when serialized
	msg := messaging.NewMessage("1234", nil, "success", "test", []byte("test"))

	// Force an error when making the HTTP request by providing an invalid URL
	newEndpoint := httpadapter.NewEndpoint("://invalid-url", "8080", "/test")
	adapter.SetEndPoint(ctx, *newEndpoint)

	err := adapter.SendRequest(ctx, msg)
	assert.Error(t, err)
}

func TestPublish_ErrorIncorrectPort(t *testing.T) {
	ctx := initializeTelemetry()
	host := "http://localhost"
	portNumber := "8081"
	path := "/test"

	adapter, _ := httpadapter.Initializer(ctx, host, portNumber, path)

	// Create a message that will not cause an error when serialized
	msg := messaging.NewMessage("1234", nil, "success", "test", []byte("test"))

	// Force an error when making the HTTP request by providing an invalid URL
	newEndpoint := httpadapter.NewEndpoint("http://localhost", "8080", "/test")
	adapter.SetEndPoint(ctx, *newEndpoint)

	err := adapter.SendRequest(ctx, msg)
	assert.Error(t, err)
}

func TestPublish_NonOKResponse(t *testing.T) {
	ctx := initializeTelemetry()
	// Create a mock HTTP server that responds with a 500 Internal Server Error status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	host := "localhost"
	portNumber := strings.Split(server.URL, ":")[2]
	path := "/test"
	adapter, _ := httpadapter.Initializer(ctx, host, portNumber, path)

	// Create a message that will not cause an error when serialized
	msg := messaging.NewMessage("1234", nil, "success", "test", []byte("test"))

	err := adapter.SendRequest(ctx, msg)
	assert.Error(t, err)
	assert.Equal(t, "server returned non-OK status code", err.Error())
}

func TestPublish_ErrorSerializing(t *testing.T) {
	ctx := initializeTelemetry()
	host := "http://localhost"
	portNumber := "8080"
	path := "/test"

	adapter, _ := httpadapter.Initializer(ctx, host, portNumber, path)

	// Create a mock message that will cause an error when serialized
	msg := &MockMessage{
		Message: messaging.NewMessage("1234", nil, "success", "test", []byte("test")),
	}

	err := adapter.SendRequest(ctx, msg)
	assert.Error(t, err)
	assert.Equal(t, "forced serialize error", err.Error())
}
