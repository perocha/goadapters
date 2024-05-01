package eventhub

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs"
	"github.com/perocha/goadapters/messaging"
	"github.com/perocha/goutils/pkg/telemetry"
)

// Publish an event to the EventHub
func (p *EventHubAdapterImpl) Publish(ctx context.Context, data messaging.Message) error {
	xTelemetry := telemetry.GetXTelemetryClient(ctx)

	// Add the operation ID to the context
	ctx = context.WithValue(context.Background(), telemetry.OperationIDKeyContextKey, data.GetOperationID())
	startTime := time.Now()

	// Check if EventHub is initialized
	if p == nil {
		err := errors.New("eventhub producer is not initialized")
		xTelemetry.Error(ctx, "EventHub::Publish::Failed", telemetry.String("Error", err.Error()))
		return err
	}

	// Create a new batch
	batch, err := p.ehProducerClient.NewEventDataBatch(ctx, nil)
	if err != nil {
		panic(err)
	}

	// Convert the message to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		// Failed to marshal message, log dependency failure to App Insights
		xTelemetry.Error(ctx, "EventHub::Publish::Failed", telemetry.String("Error", err.Error()))
		return err
	}

	// Can be called multiple times with new messages until you receive an azeventhubs.ErrMessageTooLarge
	err = batch.AddEventData(&azeventhubs.EventData{
		Body: []byte(jsonData),
	}, nil)

	if errors.Is(err, azeventhubs.ErrEventDataTooLarge) {
		// Message too large to fit into this batch.
		//
		// At this point you'd usually just send the batch (using ProducerClient.SendEventDataBatch),
		// create a new one, and start filling up the batch again.
		//
		// If this is the _only_ message being added to the batch then it's too big in general, and
		// will need to be split or shrunk to fit.
		xTelemetry.Error(ctx, "EventHub::Publish::Message too large to fit into this batch", telemetry.String("Error", err.Error()))
		return err
	} else if err != nil {
		// Some other error occurred
		xTelemetry.Error(ctx, "EventHub::Publish::Failed to add message to batch", telemetry.String("Error", err.Error()))
		return err
	}

	// Send the batch
	err = p.ehProducerClient.SendEventDataBatch(context.TODO(), batch, nil)

	if err != nil {
		xTelemetry.Error(ctx, "EventHub::Publish::Failed to send message", telemetry.String("Error", err.Error()))
		return err
	}

	// Track the dependency to App Insights, using event hub name as the target
	xTelemetry.Dependency(ctx, "EventHub", p.eventHubName, true, startTime, time.Now(), "Publish EventHub message success")

	return nil
}
