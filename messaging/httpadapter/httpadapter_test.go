package httpadapter_test

import (
	"context"
	"testing"

	"github.com/perocha/goadapters/messaging/httpadapter"
	"github.com/stretchr/testify/assert"
)

func TestConsumerInitializer(t *testing.T) {
	ctx := context.Background()
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
	ctx := context.Background()
	endPointURL := "http://localhost"
	portNumber := "8080"

	adapter, err := httpadapter.PublisherInitializer(ctx, endPointURL, portNumber)

	assert.NoError(t, err)
	assert.NotNil(t, adapter)
	assert.IsType(t, &httpadapter.HTTPAdapterImpl{}, adapter)
	assert.Equal(t, endPointURL, adapter.GetEndPoint())
	assert.Equal(t, portNumber, adapter.GetPortNumber())
}
