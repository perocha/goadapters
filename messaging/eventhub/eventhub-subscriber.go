package eventhub

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/perocha/goadapters/messaging/message"
	"github.com/perocha/goutils/pkg/telemetry"
)

func (a *EventHubAdapterImpl) Subscribe(ctx context.Context) (<-chan message.Message, context.CancelFunc, error) {
	telemetryClient := telemetry.GetTelemetryClient(ctx)
	eventChannel := make(chan message.Message)

	// Run all partition clients
	go a.dispatchPartitionClients(ctx, eventChannel)

	processorCtx, processorCancel := context.WithCancel(context.TODO())

	go func() {
		if err := a.ehProcessor.Run(processorCtx); err != nil {
			telemetryClient.TrackException(ctx, "EventHubAdapter::Subscribe::Error processor run", err, telemetry.Critical, nil, true)
			processorCancel()
			a.ehConsumerClient.Close(context.TODO())
		}
	}()

	return eventChannel, processorCancel, nil
}

func (a *EventHubAdapterImpl) dispatchPartitionClients(ctx context.Context, eventChannel chan message.Message) {
	for {
		telemetryClient := telemetry.GetTelemetryClient(ctx)

		// Get the next partition client
		partitionClient := a.ehProcessor.NextPartitionClient(context.TODO())

		if partitionClient == nil {
			// No more partition clients to process
			break
		}

		go func() {
			telemetryClient.TrackTrace(ctx, "EventHubAdapter::dispatchPartitionClients::Partition ID "+partitionClient.PartitionID()+"::Client initialized", telemetry.Information, nil, true)

			// Process events for the partition client
			if err := a.processEventsForPartition(ctx, partitionClient, eventChannel); err != nil {
				properties := map[string]string{
					"PartitionID": partitionClient.PartitionID(),
				}
				telemetryClient.TrackException(ctx, "EventHubAdapter::dispatchPartitionClients::Error processing events", err, telemetry.Error, properties, true)
				//panic(err)
				return
			}
		}()
	}
}

// ProcessEvents implements the logic that is executed when events are received from the event hub
func (a *EventHubAdapterImpl) processEventsForPartition(ctx context.Context, partitionClient *azeventhubs.ProcessorPartitionClient, eventChannel chan message.Message) error {
	telemetryClient := telemetry.GetTelemetryClient(ctx)

	// Defer the shutdown of the partition resources
	defer func() {
		shutdownPartitionResources(ctx, partitionClient)
	}()

	for {
		// Receive events from the partition client with a timeout of 20 seconds
		timeout := time.Second * 20
		receiveCtx, receiveCtxCancel := context.WithTimeout(context.TODO(), timeout)

		// Limit the wait for a number of events to receive
		limitEvents := 10
		events, err := partitionClient.ReceiveEvents(receiveCtx, limitEvents, nil)
		receiveCtxCancel()

		if err != nil && !errors.Is(err, context.DeadlineExceeded) {
			telemetryClient.TrackException(ctx, "EventHubAdapter::processEventsForPartition::Error receiving events", err, telemetry.Error, nil, true)
			return err
		}

		// Uncomment the following line to verify that the consumer is trying to receive events
		// log.Printf("EventHubAdapter::processEventsForPartition::PartitionID=%s::Processing %d event(s)\n", partitionClient.PartitionID(), len(events))

		for _, eventItem := range events {
			// Track the current time to log the telemetry and create a new operation uuid (add to the context)
			startTime := time.Now()

			// eventItem.Body is a byte slice and needs to be unmarshalled into a message
			receivedMessage := message.NewMessage("", nil, "", "", nil)
			err := receivedMessage.Deserialize(eventItem.Body)

			if err != nil {
				// Error unmarshalling the event body, send an error event to the event channel
				telemetryClient.TrackException(ctx, "EventHubAdapter::processEventsForPartition::Error unmarshalling event body", err, telemetry.Error, nil, true)
				errorMessage := message.NewMessage("", err, "", "", nil)
				eventChannel <- errorMessage
			} else {
				// If we reach this point, we have a message!! Get the operation ID from the message and add it to the context
				ctx := context.WithValue(context.Background(), telemetry.OperationIDKeyContextKey, receivedMessage.GetOperationID())
				telemetryClient.TrackTrace(ctx, "EventHubAdapter::processEventsForPartition::Bonzai!!", telemetry.Information, nil, true)
				eventChannel <- receivedMessage
			}

			telemetryClient.TrackDependency(ctx, "EventHubAdapter::processEventsForPartition", "Process message", "EventHub", a.eventHubName, true, startTime, time.Now(), nil, true)
		}

		if len(events) != 0 {
			if err := partitionClient.UpdateCheckpoint(context.TODO(), events[len(events)-1], nil); err != nil {
				telemetryClient.TrackException(ctx, "EventHubAdapter::processEventsForPartition::Error updating checkpoint", err, telemetry.Error, nil, true)
				return err
			}
		}
	}
}

// Closes the partition client
func shutdownPartitionResources(ctx context.Context, partitionClient *azeventhubs.ProcessorPartitionClient) {
	telemetryClient := telemetry.GetTelemetryClient(ctx)
	telemetryClient.TrackTrace(ctx, "EventHubAdapter::shutdownPartitionResources::Closing partition client", telemetry.Information, nil, true)

	// Close the partition client
	defer partitionClient.Close(context.TODO())
}
