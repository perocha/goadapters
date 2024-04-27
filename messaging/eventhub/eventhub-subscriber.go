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
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	eventChannel := make(chan message.Message)

	// Run all partition clients
	go a.dispatchPartitionClients(ctx, eventChannel)

	processorCtx, processorCancel := context.WithCancel(context.TODO())

	go func() {
		if err := a.ehProcessor.Run(processorCtx); err != nil {
			xTelemetry.Error(ctx, "EventHubAdapter::Subscribe::Error processor run", telemetry.String("Error", err.Error()))
			processorCancel()
			a.ehConsumerClient.Close(context.TODO())
		}
	}()

	return eventChannel, processorCancel, nil
}

func (a *EventHubAdapterImpl) dispatchPartitionClients(ctx context.Context, eventChannel chan message.Message) {
	for {
		xTelemetry := telemetry.GetXTelemetryClient(ctx)

		// Get the next partition client
		partitionClient := a.ehProcessor.NextPartitionClient(context.TODO())

		if partitionClient == nil {
			// No more partition clients to process
			break
		}

		go func() {
			// Initialize the partition client
			xTelemetry.Info(ctx, "EventHubAdapter::dispatchPartitionClients::Client initialized", telemetry.String("PartitionID", partitionClient.PartitionID()))

			// Process events for the partition client
			if err := a.processEventsForPartition(ctx, partitionClient, eventChannel); err != nil {
				xTelemetry.Error(ctx, "EventHubAdapter::dispatchPartitionClients::Error processing events", telemetry.String("PartitionID", partitionClient.PartitionID()), telemetry.String("Error", err.Error()))
				//panic(err)
				return
			}
		}()
	}
}

// ProcessEvents implements the logic that is executed when events are received from the event hub
func (a *EventHubAdapterImpl) processEventsForPartition(ctx context.Context, partitionClient *azeventhubs.ProcessorPartitionClient, eventChannel chan message.Message) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)

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
			xTelemetry.Error(ctx, "EventHubAdapter::processEventsForPartition::Error receiving events", telemetry.String("PartitionID", partitionClient.PartitionID()), telemetry.String("Error", err.Error()))
			return err
		}

		// Uncomment the following line to verify that the consumer is trying to receive events
		xTelemetry.Debug(ctx, "EventHubAdapter::processEventsForPartition", telemetry.String("PartitionID", partitionClient.PartitionID()), telemetry.Int("Events", len(events)))

		for _, eventItem := range events {
			// Track the current time to log the telemetry and create a new operation uuid (add to the context)
			//startTime := time.Now()

			// eventItem.Body is a byte slice and needs to be unmarshalled into a message
			receivedMessage := message.NewMessage("", nil, "", "", nil)
			err := receivedMessage.Deserialize(eventItem.Body)

			if err != nil {
				// Error unmarshalling the event body, send an error event to the event channel
				xTelemetry.Error(ctx, "EventHubAdapter::processEventsForPartition::Error unmarshalling event body", telemetry.String("PartitionID", partitionClient.PartitionID()), telemetry.String("Error", err.Error()))
				errorMessage := message.NewMessage("", err, "", "", nil)
				eventChannel <- errorMessage
			} else {
				// If we reach this point, we have a message!! Get the operation ID from the message and add it to the context
				ctx := context.WithValue(context.Background(), telemetry.OperationIDKeyContextKey, receivedMessage.GetOperationID())
				xTelemetry.Debug(ctx, "EventHubAdapter::processEventsForPartition::Message received", telemetry.String("PartitionID", partitionClient.PartitionID()))
				eventChannel <- receivedMessage
				xTelemetry.Info(ctx, "EventHubAdapter::processEventsForPartition::Message received", telemetry.String("PartitionID", partitionClient.PartitionID()), telemetry.String("Command", receivedMessage.GetCommand()), telemetry.String("Status", receivedMessage.GetStatus()))
				/*
						telemetryProps := map[string]string{
							"OperationID": receivedMessage.GetOperationID(),
							"Command":     receivedMessage.GetCommand(),
							"Status":      receivedMessage.GetStatus(),
						}
					// TODO ---> telemetryClient.TrackDependency(ctx, "EventHubAdapter::processEventsForPartition", "Process message", "EventHub", a.eventHubName, true, startTime, time.Now(), telemetryProps, true)
				*/
			}

			//telemetryClient.TrackDependency(ctx, "EventHubAdapter::processEventsForPartition", "Process message", "EventHub", a.eventHubName, true, startTime, time.Now(), nil, true)
		}

		if len(events) != 0 {
			if err := partitionClient.UpdateCheckpoint(context.TODO(), events[len(events)-1], nil); err != nil {
				xTelemetry.Error(ctx, "EventHubAdapter::processEventsForPartition::Error updating checkpoint", telemetry.String("PartitionID", partitionClient.PartitionID()), telemetry.String("Error", err.Error()))
				return err
			}
		}
	}
}

// Closes the partition client
func shutdownPartitionResources(ctx context.Context, partitionClient *azeventhubs.ProcessorPartitionClient) {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)
	xTelemetry.Debug(ctx, "EventHubAdapter::shutdownPartitionResources", telemetry.String("PartitionID", partitionClient.PartitionID()))

	// Close the partition client
	defer partitionClient.Close(context.TODO())
}
