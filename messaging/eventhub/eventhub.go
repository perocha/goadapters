package eventhub

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/checkpoints"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/perocha/goutils/pkg/telemetry"
)

type EventHubAdapterImpl struct {
	ehProcessor      *azeventhubs.Processor
	ehConsumerClient *azeventhubs.ConsumerClient
	checkpointStore  *checkpoints.BlobStore
	checkClient      *container.Client
	ehProducerClient *azeventhubs.ProducerClient
	eventHubName     string
}

// Initializes only the consumer client
func ConsumerInitializer(ctx context.Context, eventHubName, consumerConnectionString, containerName, checkpointStoreConnectionString string) (*EventHubAdapterImpl, error) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)

	// create a container client using a connection string and container name
	checkClient, err := container.NewClientFromConnectionString(checkpointStoreConnectionString, containerName, nil)
	if err != nil {
		xTelemetry.Error(ctx, "EventHubAdapter::Error creating container client", telemetry.String("Error", err.Error()))
		return nil, err
	}

	// create a checkpoint store that will be used by the event hub
	checkpointStore, err := checkpoints.NewBlobStore(checkClient, nil)
	if err != nil {
		xTelemetry.Error(ctx, "EventHubAdapter::Error creating checkpoint store", telemetry.String("Error", err.Error()))
		return nil, err
	}

	// create a consumer client using a connection string to the namespace and the event hub
	consumerClient, err := azeventhubs.NewConsumerClientFromConnectionString(consumerConnectionString, eventHubName, azeventhubs.DefaultConsumerGroup, nil)
	if err != nil {
		xTelemetry.Error(ctx, "EventHubAdapter::Error creating consumer client", telemetry.String("Error", err.Error()))
		return nil, err
	}

	// Create a processor to receive and process events
	processor, err := azeventhubs.NewProcessor(consumerClient, checkpointStore, nil)
	if err != nil {
		xTelemetry.Error(ctx, "EventHubAdapter::Error creating processor", telemetry.String("Error", err.Error()))
		return nil, err
	}

	// Obtain the eventHubName from the consumer client
	eventHubProperties, err := consumerClient.GetEventHubProperties(ctx, nil)
	if err != nil {
		xTelemetry.Error(ctx, "EventHubAdapter::Error getting event hub properties", telemetry.String("Error", err.Error()))
		return nil, err
	}

	adapter := &EventHubAdapterImpl{
		ehProcessor:      processor,
		ehConsumerClient: consumerClient,
		checkpointStore:  checkpointStore,
		checkClient:      checkClient,
		eventHubName:     eventHubProperties.Name,
	}

	return adapter, nil
}

// Initializes only the producer client
func ProducerInitializer(ctx context.Context, eventHubName, producerConnectionString string) (*EventHubAdapterImpl, error) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)

	// Create a new producer client
	producerClient, err := azeventhubs.NewProducerClientFromConnectionString(producerConnectionString, eventHubName, nil)
	if err != nil {
		xTelemetry.Error(ctx, "EventHubAdapter::Failed to initialize producer", telemetry.String("Error", err.Error()))
		return nil, err
	}

	// Obtain the eventHubName from the producer client
	eventHubProperties, err := producerClient.GetEventHubProperties(ctx, nil)
	if err != nil {
		xTelemetry.Error(ctx, "EventHubAdapter::Error getting event hub properties", telemetry.String("Error", err.Error()))
		return nil, err
	}

	adapter := &EventHubAdapterImpl{
		ehProducerClient: producerClient,
		eventHubName:     eventHubProperties.Name,
	}

	return adapter, nil
}

// Close the EventHub adapter, both the consumer and producer clients
func (a *EventHubAdapterImpl) Close(ctx context.Context) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Info(ctx, "EventHubAdapter::Close::Stopping event hub consumer and producer clients")

	// Close the consumer client
	if a.ehConsumerClient != nil {
		err := a.ehConsumerClient.Close(context.TODO())
		if err != nil {
			xTelemetry.Error(ctx, "EventHubAdapter::Error closing consumer client", telemetry.String("Error", err.Error()))
			return err
		}
	}

	// Close the producer client
	if a.ehProducerClient != nil {
		err := a.ehProducerClient.Close(ctx)
		if err != nil {
			xTelemetry.Error(ctx, "EventHubAdapter::Error closing producer client", telemetry.String("Error", err.Error()))
			return err
		}
	}

	xTelemetry.Info(ctx, "EventHubAdapter::Close::Event hub consumer and producer clients stopped")

	return nil
}
