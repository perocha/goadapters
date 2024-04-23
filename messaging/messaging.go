package messaging

import (
	"context"
	"encoding/json"
	"errors"
)

// EventWithOperationID struct that contains an operation ID and an Event
type Message struct {
	OperationID string
	Error       error
	Data        interface{}
}

// ToMap tries to convert the Data field of the Message to a map[string]interface{}
func (m *Message) ToMap() (map[string]interface{}, error) {
	if m.Data == nil {
		return nil, errors.New("message::ToMap::Data is nil")
	}

	// Try to assert the Data field to a map[string]interface{}
	dataMap, ok := m.Data.(map[string]interface{})
	if ok {
		return dataMap, nil
	}

	// If the above assertion failed, try to assert the Data field to []byte
	dataBytes, ok := m.Data.([]byte)
	if !ok {
		return nil, errors.New("message::ToMap::Error Data is not a valid JSON")
	}

	// Unmarshal the data to a map[string]interface{}
	var result map[string]interface{}
	if err := json.Unmarshal(dataBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Interface for messaging systems
type MessagingSystem interface {
	Publish(ctx context.Context, data interface{}) error
	Subscribe(ctx context.Context) (<-chan Message, context.CancelFunc, error)
	Close(ctx context.Context) error
}
