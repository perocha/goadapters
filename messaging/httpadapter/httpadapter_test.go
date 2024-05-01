package httpadapter_test

import (
	"context"
	"log"
	"testing"

	"github.com/perocha/goadapters/messaging/httpadapter"
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
